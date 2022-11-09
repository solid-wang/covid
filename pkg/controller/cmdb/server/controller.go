package server

import (
	"context"
	"fmt"
	cmdbv1 "github.com/solid-wang/covid/pkg/apis/cmdb/v1"
	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	devopsv1 "github.com/solid-wang/covid/pkg/apis/devops/v1"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	typedcorev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/core/v1"
	cmdbinformers "github.com/solid-wang/covid/pkg/generated/informers/externalversions/cmdb/v1"
	devopsinformers "github.com/solid-wang/covid/pkg/generated/informers/externalversions/devops/v1"
	cmdblister "github.com/solid-wang/covid/pkg/generated/listers/cmdb/v1"
	devopslister "github.com/solid-wang/covid/pkg/generated/listers/devops/v1"
	"github.com/solid-wang/covid/pkg/tools/record"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"time"
)

const (
	// maxRetries is the number of times a ci will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// a ci is going to be requeued:
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15
)

type Controller struct {
	client clientset.Interface

	// serverLister can list/get server from the shared informer's store
	serverLister cmdblister.ServerLister

	// gitlabLister can list/get server from the shared informer's store
	gitlabLister cmdblister.GitlabLister

	// ProductLister can list/get product from the shared informer's store
	productLister cmdblister.ProductLister

	// continuousIntegrationLister can list/get continuousIntegration from the shared informer's store
	continuousIntegrationLister devopslister.ContinuousIntegrationLister

	// continuousDeploymentLister can list/get continuousDeploymentLister from the shared informer's store
	continuousDeploymentLister devopslister.ContinuousDeploymentLister

	// serverListerSynced returns true if the server store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	serverListerSynced cache.InformerSynced

	// gitlabListerSynced returns true if the gitlab store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	gitlabListerSynced cache.InformerSynced

	// ProductListerSynced returns true if the product store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	productListerSynced cache.InformerSynced

	// continuousIntegrationListerSynced returns true if the continuousIntegration store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	continuousIntegrationListerSynced cache.InformerSynced

	// continuousDeploymentListerSynced returns true if the continuousDeployment store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	continuousDeploymentListerSynced cache.InformerSynced

	eventRecorder record.EventRecorder
	// ci that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewServerController creates a new ServerController.
func NewServerController(serverInformer cmdbinformers.ServerInformer, gitlabInformer cmdbinformers.GitlabInformer, productInformer cmdbinformers.ProductInformer, continuousIntegrationInformer devopsinformers.ContinuousIntegrationInformer, continuousDeploymentInformer devopsinformers.ContinuousDeploymentInformer, client clientset.Interface) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events("")})

	c := &Controller{
		client:        client,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "server-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "server"),
	}

	serverInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addServer,
		UpdateFunc: c.updateServer,
		DeleteFunc: c.deleteServer,
	})

	continuousIntegrationInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: c.updateContinuousIntegration,
		DeleteFunc: c.deleteContinuousIntegration,
	})

	productInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteProduct,
	})

	continuousDeploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: c.updateContinuousDeployment,
		DeleteFunc: c.deleteContinuousDeployment,
	})

	c.serverLister = serverInformer.Lister()
	c.gitlabLister = gitlabInformer.Lister()
	c.continuousIntegrationLister = continuousIntegrationInformer.Lister()
	c.productLister = productInformer.Lister()
	c.continuousDeploymentLister = continuousDeploymentInformer.Lister()
	c.serverListerSynced = serverInformer.Informer().HasSynced
	c.gitlabListerSynced = gitlabInformer.Informer().HasSynced
	c.continuousIntegrationListerSynced = continuousIntegrationInformer.Informer().HasSynced
	c.productListerSynced = productInformer.Informer().HasSynced
	c.continuousDeploymentListerSynced = continuousDeploymentInformer.Informer().HasSynced

	return c, nil
}

// Run begins watching and syncing.
func (c *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.InfoS("Starting controller", "controller", "server")
	defer klog.InfoS("Shutting down controller", "controller", "server")

	if !cache.WaitForNamedCacheSync("server", ctx.Done(), c.serverListerSynced, c.productListerSynced, c.gitlabListerSynced, c.continuousIntegrationListerSynced, c.continuousDeploymentListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.worker, time.Second)
	}

	<-ctx.Done()
}

func (c *Controller) addServer(obj interface{}) {
	s := obj.(*cmdbv1.Server)
	klog.V(4).InfoS("Adding server", "server", klog.KObj(s))
	c.enqueue(s)
}

func (c *Controller) updateServer(old, cur interface{}) {
	oldCI := old.(*cmdbv1.Server)
	curCI := cur.(*cmdbv1.Server)
	klog.V(4).InfoS("Updating server", "server", klog.KObj(oldCI))
	c.enqueue(curCI)
}

func (c *Controller) deleteServer(obj interface{}) {
	s, ok := obj.(*cmdbv1.Server)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		s, ok = tombstone.Obj.(*cmdbv1.Server)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a server %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("Deleting server", "server", klog.KObj(s))
	c.enqueue(s)
}

func (c *Controller) updateContinuousIntegration(old, cur interface{}) {
	ci := cur.(*devopsv1.ContinuousIntegration)
	klog.V(4).InfoS("ContinuousIntegration updated", "ContinuousIntegration", klog.KObj(ci))
	c.queue.Add(ci.GetNamespace() + "/" + ci.GetName())
}

func (c *Controller) deleteContinuousIntegration(obj interface{}) {
	ci, ok := obj.(*devopsv1.ContinuousIntegration)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		ci, ok = tombstone.Obj.(*devopsv1.ContinuousIntegration)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a server %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("ContinuousIntegration deleted", "ContinuousIntegration", klog.KObj(ci))
	c.queue.Add(ci.GetNamespace() + "/" + ci.GetName())
}

func (c *Controller) deleteProduct(obj interface{}) {
	sv, ok := obj.(*cmdbv1.Product)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		sv, ok = tombstone.Obj.(*cmdbv1.Product)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a server %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("Product deleted", "Product", klog.KObj(sv))
	servers, _ := c.serverLister.Servers(sv.GetName()).List(labels.Everything())
	if len(servers) == 0 {
		return
	}
	c.queue.Add(sv.GetName() + "/" + servers[0].GetName())
}

func (c *Controller) updateContinuousDeployment(old, cur interface{}) {
	cd := cur.(*devopsv1.ContinuousDeployment)
	klog.V(4).InfoS("ContinuousDeployment updated", "ContinuousDeployment", klog.KObj(cd))
	c.queue.Add(cd.GetNamespace() + "/" + cd.GetName())
}

func (c *Controller) deleteContinuousDeployment(obj interface{}) {
	cd, ok := obj.(*devopsv1.ContinuousDeployment)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		cd, ok = tombstone.Obj.(*devopsv1.ContinuousDeployment)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a server %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("ContinuousDeployment deleted", "ContinuousDeployment", klog.KObj(cd))
	c.queue.Add(cd.GetNamespace() + "/" + cd.GetName())
}

func (c *Controller) enqueue(s *cmdbv1.Server) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(s)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", s, err))
		return
	}

	c.queue.Add(key)
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (c *Controller) worker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncHandler(ctx, key.(string))
	c.handleErr(err, key)

	return true
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}
	ns, name, keyErr := cache.SplitMetaNamespaceKey(key.(string))
	if keyErr != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
	}

	if c.queue.NumRequeues(key) < maxRetries {
		klog.V(2).InfoS("Error syncing server", "server", klog.KRef(ns, name), "err", err)
		c.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	klog.V(2).InfoS("Dropping ci out of the queue", "server", klog.KRef(ns, name), "err", err)

	c.queue.Forget(key)
}

// syncHandler will sync the ci with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (c *Controller) syncHandler(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
		return err
	}

	startTime := time.Now()
	klog.V(4).InfoS("Started syncing server", "server", klog.KRef(namespace, name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing server", "server", klog.KRef(namespace, name), "duration", time.Since(startTime))
	}()

	server, err := c.serverLister.Servers(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("server has been deleted", "server", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return err
	}

	s := server.DeepCopy()

	// Ensure that the product to which the service belongs exists to support the normal operation of devops.
	if err := c.ensureProductExists(ctx, s); err != nil {
		c.eventRecorder.Eventf(s, v1.EventTypeWarning, "Sync", "Sync Product err: %s", err)
		return err
	}
	c.eventRecorder.Event(s, v1.EventTypeNormal, "Sync", "Synced Product")

	// Make sure the server is in gitlab's index so that devops events can find the server exactly
	if err := c.syncGitlabIndex(ctx, s); err != nil {
		c.eventRecorder.Eventf(s, v1.EventTypeWarning, "Sync", "Sync Gitlab err: %s", err)
		return err
	}
	c.eventRecorder.Event(s, v1.EventTypeNormal, "Sync", "Synced Gitlab")

	// sync ci
	if err := c.syncContinuousIntegration(ctx, s); err != nil {
		c.eventRecorder.Eventf(s, v1.EventTypeWarning, "Sync", "Sync ContinuousIntegration err: %s", err)
		return err
	}
	c.eventRecorder.Event(s, v1.EventTypeNormal, "Sync", "Synced ContinuousIntegration")

	// sync cd
	if err := c.syncContinuousDeployment(ctx, s); err != nil {
		c.eventRecorder.Eventf(s, v1.EventTypeWarning, "Sync", "Sync ContinuousDeployment err: %s", err)
		return err
	}
	c.eventRecorder.Event(s, v1.EventTypeNormal, "Sync", "Synced ContinuousDeployment")

	if s.DeletionTimestamp != nil {
		s.Finalizers = nil
		c.client.CmdbV1().Servers(s.Namespace).Update(ctx, s, metav1.UpdateOptions{})
		c.eventRecorder.Eventf(s, v1.EventTypeNormal, "Terminating", "Remove Server %s Finalizers", s.Name)
	}

	return nil
}
