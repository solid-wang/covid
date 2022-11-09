package gitlab

import (
	"context"
	"fmt"
	cmdbv1 "github.com/solid-wang/covid/pkg/apis/cmdb/v1"
	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	typedcorev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/core/v1"
	cmdbinformers "github.com/solid-wang/covid/pkg/generated/informers/externalversions/cmdb/v1"
	cmdblister "github.com/solid-wang/covid/pkg/generated/listers/cmdb/v1"
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
	"reflect"
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

	// gitlabLister can list/get server from the shared informer's store
	gitlabLister cmdblister.GitlabLister

	// gitlabListerSynced returns true if the gitlab store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	gitlabListerSynced cache.InformerSynced

	eventRecorder record.EventRecorder
	// ci that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewGitlabController creates a new ServerController.
func NewGitlabController(gitlabInformer cmdbinformers.GitlabInformer, client clientset.Interface) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events("")})

	c := &Controller{
		client:        client,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "server-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "server"),
	}

	gitlabInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addGitlab,
		UpdateFunc: c.updateGitlab,
		DeleteFunc: c.deleteGitlab,
	})

	c.gitlabLister = gitlabInformer.Lister()
	c.gitlabListerSynced = gitlabInformer.Informer().HasSynced

	return c, nil
}

// Run begins watching and syncing.
func (c *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.InfoS("Starting controller", "controller", "gitlab")
	defer klog.InfoS("Shutting down controller", "controller", "gitlab")

	if !cache.WaitForNamedCacheSync("gitlab", ctx.Done(), c.gitlabListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.worker, time.Second)
	}

	<-ctx.Done()
}

func (c *Controller) addGitlab(obj interface{}) {
	s := obj.(*cmdbv1.Gitlab)
	klog.V(4).InfoS("Adding gitlab", "gitlab", klog.KObj(s))
	c.enqueue(s)
}

func (c *Controller) updateGitlab(old, cur interface{}) {
	oldCI := old.(*cmdbv1.Gitlab)
	curCI := cur.(*cmdbv1.Gitlab)
	klog.V(4).InfoS("Updating gitlab", "gitlab", klog.KObj(oldCI))
	c.enqueue(curCI)
}

func (c *Controller) deleteGitlab(obj interface{}) {
	s, ok := obj.(*cmdbv1.Gitlab)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		s, ok = tombstone.Obj.(*cmdbv1.Gitlab)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a gitlab %#v", obj))
			return
		}
	}
	klog.V(4).InfoS("Deleting gitlab", "gitlab", klog.KObj(s))
	c.enqueue(s)
}

func (c *Controller) enqueue(s *cmdbv1.Gitlab) {
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
		klog.V(2).InfoS("Error syncing gitlab", "gitlab", klog.KRef(ns, name), "err", err)
		c.queue.AddRateLimited(key)
		return
	}

	utilruntime.HandleError(err)
	klog.V(2).InfoS("Dropping ci out of the queue", "gitlab", klog.KRef(ns, name), "err", err)

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
	klog.V(4).InfoS("Started syncing gitlab", "gitlab", klog.KRef(namespace, name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing gitlab", "gitlab", klog.KRef(namespace, name), "duration", time.Since(startTime))
	}()

	g, err := c.gitlabLister.Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("gitlab has been deleted", "gitlab", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return err
	}

	gl := g.DeepCopy()

	gitlabClient := gitlab.NewForGitlab(gl)

	for projectStringID, project := range gl.Spec.ProjectIndex {
		projectID, _ := strconv.Atoi(string(projectStringID))
		for eventType, hookID := range project.HooksMap {

			if len(project.ServersMap) == 0 && hookID != nil {
				return gitlabClient.Project().Webhook(projectID).Delete(*hookID)
			}
			if hookID == nil {
				hook := &webhooks.Webhook{
					URL:                      strings.Join([]string{gl.Spec.URL, gl.Name, "webhook", string(eventType)}, "/"),
					PushEvents:               false,
					PushEventsBranchFilter:   "",
					IssuesEvents:             false,
					ConfidentialIssuesEvents: false,
					MergeRequestsEvents:      false,
					TagPushEvents:            eventType == cmdbv1.GitlabWebhookEventTagPush,
					NoteEvents:               false,
					ConfidentialNoteEvents:   false,
					JobEvents:                false,
					PipelineEvents:           eventType == cmdbv1.GitlabWebhookEventPipeline,
					WikiPageEvents:           false,
					DeploymentEvents:         false,
					ReleasesEvents:           false,
					EnableSSLVerification:    false,
				}

				add, err := gitlabClient.Project().Webhook(projectID).Add(hook)
				if err != nil {
					return err
				}
				gl.Spec.ProjectIndex[projectStringID].HooksMap[eventType] = &add.ID
			}
		}
	}

	if reflect.DeepEqual(g, gl) {
		_, err := c.client.CmdbV1().Gitlabs().Update(ctx, gl, metav1.UpdateOptions{})
		return err
	}
	return nil
}
