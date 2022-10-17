// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEvents implements EventInterface
type FakeEvents struct {
	Fake *FakeCoreV1
	ns   string
}

var eventsResource = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}

var eventsKind = schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Event"}

// Get takes name of the event, and returns the corresponding event object, and an error if there is any.
func (c *FakeEvents) Get(ctx context.Context, name string, options v1.GetOptions) (result *corev1.Event, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(eventsResource, c.ns, name), &corev1.Event{})

	if obj == nil {
		return nil, err
	}
	return obj.(*corev1.Event), err
}

// List takes label and field selectors, and returns the list of Events that match those selectors.
func (c *FakeEvents) List(ctx context.Context, opts v1.ListOptions) (result *corev1.EventList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(eventsResource, eventsKind, c.ns, opts), &corev1.EventList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &corev1.EventList{ListMeta: obj.(*corev1.EventList).ListMeta}
	for _, item := range obj.(*corev1.EventList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested events.
func (c *FakeEvents) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(eventsResource, c.ns, opts))

}

// Create takes the representation of a event and creates it.  Returns the server's representation of the event, and an error, if there is any.
func (c *FakeEvents) Create(ctx context.Context, event *corev1.Event, opts v1.CreateOptions) (result *corev1.Event, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(eventsResource, c.ns, event), &corev1.Event{})

	if obj == nil {
		return nil, err
	}
	return obj.(*corev1.Event), err
}

// Update takes the representation of a event and updates it. Returns the server's representation of the event, and an error, if there is any.
func (c *FakeEvents) Update(ctx context.Context, event *corev1.Event, opts v1.UpdateOptions) (result *corev1.Event, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(eventsResource, c.ns, event), &corev1.Event{})

	if obj == nil {
		return nil, err
	}
	return obj.(*corev1.Event), err
}

// Delete takes name of the event and deletes it. Returns an error if one occurs.
func (c *FakeEvents) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(eventsResource, c.ns, name, opts), &corev1.Event{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEvents) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(eventsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &corev1.EventList{})
	return err
}

// Patch applies the patch and returns the patched event.
func (c *FakeEvents) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Event, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(eventsResource, c.ns, name, pt, data, subresources...), &corev1.Event{})

	if obj == nil {
		return nil, err
	}
	return obj.(*corev1.Event), err
}