package application

import (
	"context"
	"fmt"
	appv1 "github.com/solid-wang/covid/pkg/apis/app/v1"
	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	servicev1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	typedcorev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/core/v1"
	appinformers "github.com/solid-wang/covid/pkg/generated/informers/externalversions/app/v1"
	serviceinformers "github.com/solid-wang/covid/pkg/generated/informers/externalversions/service/v1"
	applister "github.com/solid-wang/covid/pkg/generated/listers/app/v1"
	servicelister "github.com/solid-wang/covid/pkg/generated/listers/service/v1"
	"github.com/solid-wang/covid/pkg/gitlab"
	"github.com/solid-wang/covid/pkg/gitlab/projects/webhooks"
	"github.com/solid-wang/covid/pkg/tools/record"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
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

	// appLister can list/get application from the shared informer's store
	appLister applister.ApplicationLister

	// appListerSynced returns true if the application store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	appListerSynced cache.InformerSynced

	// gitlabLister can list/get gitlab from the shared informer's store
	gitlabLister servicelister.GitlabLister

	// gitlabListerSynced returns true if the gitlab store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	gitlabListerSynced cache.InformerSynced

	eventRecorder record.EventRecorder
	// ci that need to be synced
	queue workqueue.RateLimitingInterface

	externalUrl string
}

// NewApplicationController creates a new ApplicationController.
func NewApplicationController(applicationInformer appinformers.ApplicationInformer, gitlabInformer serviceinformers.GitlabInformer, client clientset.Interface, externalUrl string) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events("")})

	c := &Controller{
		client:        client,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "application-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "application"),
		externalUrl:   externalUrl,
	}

	applicationInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addApplication,
		UpdateFunc: c.updateApplication,
		DeleteFunc: c.deleteApplication,
	})

	c.appLister = applicationInformer.Lister()
	c.gitlabLister = gitlabInformer.Lister()
	c.appListerSynced = applicationInformer.Informer().HasSynced
	c.gitlabListerSynced = gitlabInformer.Informer().HasSynced

	return c, nil
}

// Run begins watching and syncing.
func (c *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.InfoS("Starting controller", "controller", "application")
	defer klog.InfoS("Shutting down controller", "controller", "application")

	if !cache.WaitForNamedCacheSync("application", ctx.Done(), c.appListerSynced, c.gitlabListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.worker, time.Second)
	}

	<-ctx.Done()
}

func (c *Controller) addApplication(obj interface{}) {
	s := obj.(*appv1.Application)
	klog.V(4).InfoS("Adding application", "application", klog.KObj(s))
	c.enqueue(s)
}

func (c *Controller) updateApplication(old, cur interface{}) {
	oldCI := old.(*appv1.Application)
	curCI := cur.(*appv1.Application)
	klog.V(4).InfoS("Updating application", "application", klog.KObj(oldCI))
	c.enqueue(curCI)
}

func (c *Controller) deleteApplication(obj interface{}) {
	s, ok := obj.(*appv1.Application)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		s, ok = tombstone.Obj.(*appv1.Application)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a application %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("Deleting application", "application", klog.KObj(s))
	c.enqueue(s)
}

func (c *Controller) enqueue(s *appv1.Application) {
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
		klog.V(2).InfoS("Error syncing application", "application", klog.KRef(ns, name), "err", err)
		c.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	klog.V(2).InfoS("Dropping application out of the queue", "application", klog.KRef(ns, name), "err", err)

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
	klog.V(4).InfoS("Started syncing application", "application", klog.KRef(namespace, name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing application", "application", klog.KRef(namespace, name), "duration", time.Since(startTime))
	}()

	app, err := c.appLister.Applications(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("application has been deleted", "application", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return err
	}

	application := app.DeepCopy()

	gl, err := c.gitlabLister.Get(application.Spec.GitlabName)
	if err != nil {
		klog.V(2).ErrorS(err, "gitlab", klog.KRef("", application.Spec.GitlabName))
		return err
	}
	g := gl.DeepCopy()

	project, ok := g.Spec.ProjectIndex[strconv.Itoa(application.Spec.ProjectID)]
	if !ok {
		project = servicev1.Project{
			ApplicationProductMap: map[string]*string{app.Name: &app.Namespace},
			HooksMap:              make(map[servicev1.GitlabWebhookEventType]*int),
		}
		g.Spec.ProjectIndex[strconv.Itoa(application.Spec.ProjectID)] = project
	}

	if project.ApplicationProductMap[app.Name] == nil {
		g.Spec.ProjectIndex[strconv.Itoa(application.Spec.ProjectID)].ApplicationProductMap[app.Name] = &app.Namespace
	}

	gitlabClient := gitlab.NewForGitlab(g)

	if project.HooksMap[servicev1.GitlabWebhookEventTagMergeRequest] == nil {
		hook := &webhooks.Webhook{
			URL:                 strings.Join([]string{c.externalUrl, gl.Name, "webhook", string(servicev1.GitlabWebhookEventTagMergeRequest)}, "/"),
			MergeRequestsEvents: true,
		}
		add, err := gitlabClient.Project().Webhook(app.Spec.ProjectID).Add(hook)
		if err != nil {
			return err
		}
		g.Spec.ProjectIndex[strconv.Itoa(application.Spec.ProjectID)].HooksMap[servicev1.GitlabWebhookEventTagMergeRequest] = &add.ID
	}

	if project.HooksMap[servicev1.GitlabWebhookEventPipeline] == nil {
		hook := &webhooks.Webhook{
			URL:            strings.Join([]string{c.externalUrl, gl.Name, "webhook", string(servicev1.GitlabWebhookEventPipeline)}, "/"),
			PipelineEvents: true,
		}
		add, err := gitlabClient.Project().Webhook(app.Spec.ProjectID).Add(hook)
		if err != nil {
			return err
		}
		g.Spec.ProjectIndex[strconv.Itoa(application.Spec.ProjectID)].HooksMap[servicev1.GitlabWebhookEventPipeline] = &add.ID
	}

	_, err = c.client.ServiceV1().Gitlabs().Update(ctx, g, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
