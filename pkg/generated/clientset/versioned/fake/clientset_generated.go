// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	appv1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/app/v1"
	fakeappv1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/app/v1/fake"
	batchv1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/batch/v1"
	fakebatchv1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/batch/v1/fake"
	corev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/core/v1"
	fakecorev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/core/v1/fake"
	servicev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/service/v1"
	fakeservicev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/service/v1/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &Clientset{tracker: o}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return cs
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
	tracker   testing.ObjectTracker
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() testing.ObjectTracker {
	return c.tracker
}

var (
	_ clientset.Interface = &Clientset{}
	_ testing.FakeClient  = &Clientset{}
)

// AppV1 retrieves the AppV1Client
func (c *Clientset) AppV1() appv1.AppV1Interface {
	return &fakeappv1.FakeAppV1{Fake: &c.Fake}
}

// BatchV1 retrieves the BatchV1Client
func (c *Clientset) BatchV1() batchv1.BatchV1Interface {
	return &fakebatchv1.FakeBatchV1{Fake: &c.Fake}
}

// CoreV1 retrieves the CoreV1Client
func (c *Clientset) CoreV1() corev1.CoreV1Interface {
	return &fakecorev1.FakeCoreV1{Fake: &c.Fake}
}

// ServiceV1 retrieves the ServiceV1Client
func (c *Clientset) ServiceV1() servicev1.ServiceV1Interface {
	return &fakeservicev1.FakeServiceV1{Fake: &c.Fake}
}
