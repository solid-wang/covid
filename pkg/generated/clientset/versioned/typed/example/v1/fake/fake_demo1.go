// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	examplev1 "github.com/solid-wang/covid/pkg/apis/example/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDemo1s implements Demo1Interface
type FakeDemo1s struct {
	Fake *FakeExampleV1
	ns   string
}

var demo1sResource = schema.GroupVersionResource{Group: "example", Version: "v1", Resource: "demo1s"}

var demo1sKind = schema.GroupVersionKind{Group: "example", Version: "v1", Kind: "Demo1"}

// Get takes name of the demo1, and returns the corresponding demo1 object, and an error if there is any.
func (c *FakeDemo1s) Get(ctx context.Context, name string, options v1.GetOptions) (result *examplev1.Demo1, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(demo1sResource, c.ns, name), &examplev1.Demo1{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Demo1), err
}

// List takes label and field selectors, and returns the list of Demo1s that match those selectors.
func (c *FakeDemo1s) List(ctx context.Context, opts v1.ListOptions) (result *examplev1.Demo1List, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(demo1sResource, demo1sKind, c.ns, opts), &examplev1.Demo1List{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &examplev1.Demo1List{ListMeta: obj.(*examplev1.Demo1List).ListMeta}
	for _, item := range obj.(*examplev1.Demo1List).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested demo1s.
func (c *FakeDemo1s) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(demo1sResource, c.ns, opts))

}

// Create takes the representation of a demo1 and creates it.  Returns the server's representation of the demo1, and an error, if there is any.
func (c *FakeDemo1s) Create(ctx context.Context, demo1 *examplev1.Demo1, opts v1.CreateOptions) (result *examplev1.Demo1, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(demo1sResource, c.ns, demo1), &examplev1.Demo1{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Demo1), err
}

// Update takes the representation of a demo1 and updates it. Returns the server's representation of the demo1, and an error, if there is any.
func (c *FakeDemo1s) Update(ctx context.Context, demo1 *examplev1.Demo1, opts v1.UpdateOptions) (result *examplev1.Demo1, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(demo1sResource, c.ns, demo1), &examplev1.Demo1{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Demo1), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDemo1s) UpdateStatus(ctx context.Context, demo1 *examplev1.Demo1, opts v1.UpdateOptions) (*examplev1.Demo1, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(demo1sResource, "status", c.ns, demo1), &examplev1.Demo1{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Demo1), err
}

// Delete takes name of the demo1 and deletes it. Returns an error if one occurs.
func (c *FakeDemo1s) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(demo1sResource, c.ns, name, opts), &examplev1.Demo1{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDemo1s) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(demo1sResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &examplev1.Demo1List{})
	return err
}

// Patch applies the patch and returns the patched demo1.
func (c *FakeDemo1s) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *examplev1.Demo1, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(demo1sResource, c.ns, name, pt, data, subresources...), &examplev1.Demo1{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Demo1), err
}
