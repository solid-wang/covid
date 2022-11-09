package continuousdeployment

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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appslisters "k8s.io/client-go/listers/apps/v1"
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

	// ContinuousDeploymentLister can list/get cd from the shared informer's store
	continuousDeploymentLister batchlister.ContinuousDeploymentLister

	// serverVersionListerSynced returns true if the ci store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	continuousDeploymentListerSynced cache.InformerSynced

	// DeploymentLister can list/get deployments from the shared informer's store
	DeploymentLister appslisters.DeploymentLister

	// DeploymentListerSynced returns true if the ci store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	DeploymentListerSynced cache.InformerSynced

	eventRecorder record.EventRecorder
	// ci that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewContinuousDeploymentController creates a new ContinuousDeploymentController.
func NewContinuousDeploymentController(continuousDeploymentInformer batchinformers.ContinuousDeploymentInformer, client clientset.Interface) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events("")})

	c := &Controller{
		client:        client,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "continuous-deployment-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "continuous-deployment"),
	}

	continuousDeploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addContinuousDeployment,
		UpdateFunc: c.updateContinuousDeployment,
		DeleteFunc: c.deleteContinuousDeployment,
	})

	c.continuousDeploymentLister = continuousDeploymentInformer.Lister()
	c.continuousDeploymentListerSynced = continuousDeploymentInformer.Informer().HasSynced

	return c, nil
}

// Run begins watching and syncing.
func (c *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.InfoS("Starting controller", "controller", "continuous-deployment")
	defer klog.InfoS("Shutting down controller", "controller", "continuous-deployment")

	if !cache.WaitForNamedCacheSync("continuous-deployment", ctx.Done(), c.continuousDeploymentListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.worker, time.Second)
	}

	<-ctx.Done()
}

func (c *Controller) addContinuousDeployment(obj interface{}) {
	cd := obj.(*batchv1.ContinuousDeployment)
	klog.V(4).InfoS("Adding cd", "cd", klog.KObj(cd))
	c.enqueue(cd)
}

func (c *Controller) updateContinuousDeployment(old, cur interface{}) {
	oldCD := old.(*batchv1.ContinuousDeployment)
	curCD := cur.(*batchv1.ContinuousDeployment)
	klog.V(4).InfoS("Updating cd", "cd", klog.KObj(oldCD))
	c.enqueue(curCD)
}

func (c *Controller) deleteContinuousDeployment(obj interface{}) {
	cd, ok := obj.(*batchv1.ContinuousDeployment)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		cd, ok = tombstone.Obj.(*batchv1.ContinuousDeployment)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a ci %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("Deleting cd", "cd", klog.KObj(cd))
	c.enqueue(cd)
}

//func (c *Controller) updateDeployment(old, cur interface{}) {
//	curD := cur.(*appsv1.Deployment)
//	if curD.Status.AvailableReplicas == curD.Status.ReadyReplicas && curD.Status.AvailableReplicas == curD.Status.Replicas {
//		//name := curD.Labels["app"]
//		//namespace := curD.Labels["product"]
//		//continuousDeployment, err := c.continuousDeploymentLister.ContinuousDeployments(namespace).Get(name)
//		//if errors.IsNotFound(err) {
//		//	klog.V(2).InfoS("cd has been deleted", "cd", klog.KRef(namespace, name))
//		//	return
//		//}
//		//if err != nil {
//		//	return
//		//}
//		//cd := continuousDeployment.DeepCopy()
//		//cd.Status.Phase = batchv1.DevOpsSuccessStatus
//		//c.client.BatchV1().ContinuousDeployments(namespace).UpdateStatus(context.Background(), cd, metav1.UpdateOptions{})
//		c.queue.Add(curD.Labels[util.LabelProduct] + "/" + curD.Labels[util.LabelApp])
//	}
//}

func (c *Controller) enqueue(cd *batchv1.ContinuousDeployment) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cd)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", cd, err))
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
		klog.V(2).InfoS("Error syncing cd", "cd", klog.KRef(ns, name), "err", err)
		c.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	klog.V(2).InfoS("Dropping cd out of the queue", "cd", klog.KRef(ns, name), "err", err)

	continuousDeployment, err := c.continuousDeploymentLister.ContinuousDeployments(ns).Get(name)
	if err != nil {
		klog.V(2).InfoS("cd has been deleted", "cd", klog.KRef(ns, name))
		c.queue.Forget(key)
		return
	}
	cd := continuousDeployment.DeepCopy()
	cd.Status.Phase = batchv1.DevOpsFailedStatus
	c.client.BatchV1().ContinuousDeployments(cd.Namespace).UpdateStatus(context.Background(), cd, metav1.UpdateOptions{})

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
	klog.V(4).InfoS("Started syncing cd", "cd", klog.KRef(namespace, name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing cd", "cd", klog.KRef(namespace, name), "duration", time.Since(startTime))
	}()

	continuousDeployment, err := c.continuousDeploymentLister.ContinuousDeployments(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("cd has been deleted", "cd", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return err
	}

	cd := continuousDeployment.DeepCopy()

	if cd.Status.Phase == batchv1.DevOpsRunningStatus {
		go c.deploy(ctx, cd)
	}
	return nil
}
