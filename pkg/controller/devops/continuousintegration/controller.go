package continuousintegration

import (
	"context"
	"fmt"
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	typedcorev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/core/v1"
	batchinformers "github.com/solid-wang/covid/pkg/generated/informers/externalversions/batch/v1"
	batchlister "github.com/solid-wang/covid/pkg/generated/listers/batch/v1"
	"github.com/solid-wang/covid/pkg/tools/record"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	//gitlab *gitlab.Gitlab
	client clientset.Interface

	// continuousIntegrationLister can list/get ci from the shared informer's store
	continuousIntegrationLister batchlister.ContinuousIntegrationLister

	// continuousIntegrationListerSynced returns true if the ci store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	continuousIntegrationListerSynced cache.InformerSynced

	// ContinuousDeploymentLister can list/get cd from the shared informer's store
	continuousDeploymentLister batchlister.ContinuousDeploymentLister

	// serverVersionListerSynced returns true if the ci store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	continuousDeploymentListerSynced cache.InformerSynced

	eventRecorder record.EventRecorder
	// ci that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewContinuousIntegrationController creates a new ContinuousIntegrationController.
func NewContinuousIntegrationController(continuousIntegrationInformer batchinformers.ContinuousIntegrationInformer, continuousDeploymentInformer batchinformers.ContinuousDeploymentInformer, client clientset.Interface) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events("")})

	c := &Controller{
		client:        client,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "continuous-integration-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "continuous-integration"),
	}

	continuousIntegrationInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addContinuousIntegration,
		UpdateFunc: c.updateContinuousIntegration,
		DeleteFunc: c.deleteContinuousIntegration,
	})

	c.continuousIntegrationLister = continuousIntegrationInformer.Lister()
	c.continuousDeploymentLister = continuousDeploymentInformer.Lister()
	c.continuousIntegrationListerSynced = continuousIntegrationInformer.Informer().HasSynced
	c.continuousDeploymentListerSynced = continuousDeploymentInformer.Informer().HasSynced

	return c, nil
}

// Run begins watching and syncing.
func (c *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.InfoS("Starting controller", "controller", "continuous-integration")
	defer klog.InfoS("Shutting down controller", "controller", "continuous-integration")

	if !cache.WaitForNamedCacheSync("continuous-integration", ctx.Done(), c.continuousIntegrationListerSynced, c.continuousDeploymentListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.worker, time.Second)
	}

	<-ctx.Done()
}

func (c *Controller) addContinuousIntegration(obj interface{}) {
	ci := obj.(*batchv1.ContinuousIntegration)
	klog.V(4).InfoS("Adding ci", "ci", klog.KObj(ci))
	c.enqueue(ci)
}

func (c *Controller) updateContinuousIntegration(old, cur interface{}) {
	oldCI := old.(*batchv1.ContinuousIntegration)
	curCI := cur.(*batchv1.ContinuousIntegration)
	klog.V(4).InfoS("Updating ci", "ci", klog.KObj(oldCI))
	c.enqueue(curCI)
}

func (c *Controller) deleteContinuousIntegration(obj interface{}) {
	ci, ok := obj.(*batchv1.ContinuousIntegration)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		ci, ok = tombstone.Obj.(*batchv1.ContinuousIntegration)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a ci %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("Deleting ci", "ci", klog.KObj(ci))
	c.enqueue(ci)
}

func (c *Controller) enqueue(ci *batchv1.ContinuousIntegration) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(ci)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", ci, err))
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
		klog.V(2).InfoS("Error syncing ci", "ci", klog.KRef(ns, name), "err", err)
		c.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	klog.V(2).InfoS("Dropping ci out of the queue", "ci", klog.KRef(ns, name), "err", err)

	continuousIntegration, err := c.continuousIntegrationLister.ContinuousIntegrations(ns).Get(name)
	if err != nil {
		klog.V(2).InfoS("ci has been deleted", "ci", klog.KRef(ns, name))
		c.queue.Forget(key)
		return
	}
	ci := continuousIntegration.DeepCopy()
	ci.Status.Phase = batchv1.DevOpsFailedStatus
	c.client.BatchV1().ContinuousIntegrations(ci.Namespace).UpdateStatus(context.Background(), ci, metav1.UpdateOptions{})

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
	klog.V(4).InfoS("Started syncing ci", "ci", klog.KRef(namespace, name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing ci", "ci", klog.KRef(namespace, name), "duration", time.Since(startTime))
	}()

	continuousIntegration, err := c.continuousIntegrationLister.ContinuousIntegrations(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("ci has been deleted", "ci", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return err
	}

	ci := continuousIntegration.DeepCopy()

	if ci.Status.Phase == batchv1.DevOpsPendingStatus {
		ci.Status.Phase = batchv1.DevOpsRunningStatus
		_, err := c.client.BatchV1().ContinuousIntegrations(ci.Namespace).UpdateStatus(ctx, ci, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		if err := createPipeline(ci); err != nil {
			c.eventRecorder.Eventf(ci, v1.EventTypeWarning, "Create Pipeline Err", "%s", err)
			return err
		}
	}

	if ci.Status.ContinuousDeploymentTrigger == nil {
		return nil
	}

	cd, err := c.continuousDeploymentLister.ContinuousDeployments(ci.Namespace).Get(*ci.Status.ContinuousDeploymentTrigger)
	if err != nil {
		c.eventRecorder.Eventf(ci, v1.EventTypeWarning, "Get CD Err", "%s", err)
		return err
	}
	continuousDeployment := cd.DeepCopy()
	switch ci.Status.Phase {
	case batchv1.DevOpsSuccessStatus:
		continuousDeployment.Status.Phase = batchv1.DevOpsRunningStatus
	case batchv1.DevOpsFailedStatus, batchv1.DevOpsCancelStatus:
		continuousDeployment.Status.Phase = batchv1.DevOpsCancelStatus
	default:
	}
	_, err = c.client.BatchV1().ContinuousDeployments(cd.Namespace).UpdateStatus(ctx, continuousDeployment, metav1.UpdateOptions{})
	if err != nil {
		c.eventRecorder.Eventf(ci, v1.EventTypeWarning, "Sync CD Err", "%s", err)
		return err
	}

	return nil
}
