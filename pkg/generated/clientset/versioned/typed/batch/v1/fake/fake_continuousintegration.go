// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeContinuousIntegrations implements ContinuousIntegrationInterface
type FakeContinuousIntegrations struct {
	Fake *FakeBatchV1
	ns   string
}

var continuousintegrationsResource = schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "continuousintegrations"}

var continuousintegrationsKind = schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "ContinuousIntegration"}

// Get takes name of the continuousIntegration, and returns the corresponding continuousIntegration object, and an error if there is any.
func (c *FakeContinuousIntegrations) Get(ctx context.Context, name string, options v1.GetOptions) (result *batchv1.ContinuousIntegration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(continuousintegrationsResource, c.ns, name), &batchv1.ContinuousIntegration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*batchv1.ContinuousIntegration), err
}

// List takes label and field selectors, and returns the list of ContinuousIntegrations that match those selectors.
func (c *FakeContinuousIntegrations) List(ctx context.Context, opts v1.ListOptions) (result *batchv1.ContinuousIntegrationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(continuousintegrationsResource, continuousintegrationsKind, c.ns, opts), &batchv1.ContinuousIntegrationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &batchv1.ContinuousIntegrationList{ListMeta: obj.(*batchv1.ContinuousIntegrationList).ListMeta}
	for _, item := range obj.(*batchv1.ContinuousIntegrationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested continuousIntegrations.
func (c *FakeContinuousIntegrations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(continuousintegrationsResource, c.ns, opts))

}

// Create takes the representation of a continuousIntegration and creates it.  Returns the server's representation of the continuousIntegration, and an error, if there is any.
func (c *FakeContinuousIntegrations) Create(ctx context.Context, continuousIntegration *batchv1.ContinuousIntegration, opts v1.CreateOptions) (result *batchv1.ContinuousIntegration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(continuousintegrationsResource, c.ns, continuousIntegration), &batchv1.ContinuousIntegration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*batchv1.ContinuousIntegration), err
}

// Update takes the representation of a continuousIntegration and updates it. Returns the server's representation of the continuousIntegration, and an error, if there is any.
func (c *FakeContinuousIntegrations) Update(ctx context.Context, continuousIntegration *batchv1.ContinuousIntegration, opts v1.UpdateOptions) (result *batchv1.ContinuousIntegration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(continuousintegrationsResource, c.ns, continuousIntegration), &batchv1.ContinuousIntegration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*batchv1.ContinuousIntegration), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeContinuousIntegrations) UpdateStatus(ctx context.Context, continuousIntegration *batchv1.ContinuousIntegration, opts v1.UpdateOptions) (*batchv1.ContinuousIntegration, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(continuousintegrationsResource, "status", c.ns, continuousIntegration), &batchv1.ContinuousIntegration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*batchv1.ContinuousIntegration), err
}

// Delete takes name of the continuousIntegration and deletes it. Returns an error if one occurs.
func (c *FakeContinuousIntegrations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(continuousintegrationsResource, c.ns, name, opts), &batchv1.ContinuousIntegration{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeContinuousIntegrations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(continuousintegrationsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &batchv1.ContinuousIntegrationList{})
	return err
}

// Patch applies the patch and returns the patched continuousIntegration.
func (c *FakeContinuousIntegrations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *batchv1.ContinuousIntegration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(continuousintegrationsResource, c.ns, name, pt, data, subresources...), &batchv1.ContinuousIntegration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*batchv1.ContinuousIntegration), err
}