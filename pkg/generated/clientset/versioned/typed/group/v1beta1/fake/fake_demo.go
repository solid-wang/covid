// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/solid-wang/covid/pkg/apis/group/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDemos implements DemoInterface
type FakeDemos struct {
	Fake *FakeGroupV1beta1
	ns   string
}

var demosResource = schema.GroupVersionResource{Group: "group", Version: "v1beta1", Resource: "demos"}

var demosKind = schema.GroupVersionKind{Group: "group", Version: "v1beta1", Kind: "Demo"}

// Get takes name of the demo, and returns the corresponding demo object, and an error if there is any.
func (c *FakeDemos) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Demo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(demosResource, c.ns, name), &v1beta1.Demo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Demo), err
}

// List takes label and field selectors, and returns the list of Demos that match those selectors.
func (c *FakeDemos) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.DemoList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(demosResource, demosKind, c.ns, opts), &v1beta1.DemoList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.DemoList{ListMeta: obj.(*v1beta1.DemoList).ListMeta}
	for _, item := range obj.(*v1beta1.DemoList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested demos.
func (c *FakeDemos) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(demosResource, c.ns, opts))

}

// Create takes the representation of a demo and creates it.  Returns the server's representation of the demo, and an error, if there is any.
func (c *FakeDemos) Create(ctx context.Context, demo *v1beta1.Demo, opts v1.CreateOptions) (result *v1beta1.Demo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(demosResource, c.ns, demo), &v1beta1.Demo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Demo), err
}

// Update takes the representation of a demo and updates it. Returns the server's representation of the demo, and an error, if there is any.
func (c *FakeDemos) Update(ctx context.Context, demo *v1beta1.Demo, opts v1.UpdateOptions) (result *v1beta1.Demo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(demosResource, c.ns, demo), &v1beta1.Demo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Demo), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDemos) UpdateStatus(ctx context.Context, demo *v1beta1.Demo, opts v1.UpdateOptions) (*v1beta1.Demo, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(demosResource, "status", c.ns, demo), &v1beta1.Demo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Demo), err
}

// Delete takes name of the demo and deletes it. Returns an error if one occurs.
func (c *FakeDemos) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(demosResource, c.ns, name, opts), &v1beta1.Demo{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDemos) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(demosResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.DemoList{})
	return err
}

// Patch applies the patch and returns the patched demo.
func (c *FakeDemos) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Demo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(demosResource, c.ns, name, pt, data, subresources...), &v1beta1.Demo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Demo), err
}
