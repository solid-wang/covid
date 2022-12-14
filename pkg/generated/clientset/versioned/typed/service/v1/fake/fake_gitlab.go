// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	servicev1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGitlabs implements GitlabInterface
type FakeGitlabs struct {
	Fake *FakeServiceV1
}

var gitlabsResource = schema.GroupVersionResource{Group: "service", Version: "v1", Resource: "gitlabs"}

var gitlabsKind = schema.GroupVersionKind{Group: "service", Version: "v1", Kind: "Gitlab"}

// Get takes name of the gitlab, and returns the corresponding gitlab object, and an error if there is any.
func (c *FakeGitlabs) Get(ctx context.Context, name string, options v1.GetOptions) (result *servicev1.Gitlab, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(gitlabsResource, name), &servicev1.Gitlab{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicev1.Gitlab), err
}

// List takes label and field selectors, and returns the list of Gitlabs that match those selectors.
func (c *FakeGitlabs) List(ctx context.Context, opts v1.ListOptions) (result *servicev1.GitlabList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(gitlabsResource, gitlabsKind, opts), &servicev1.GitlabList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &servicev1.GitlabList{ListMeta: obj.(*servicev1.GitlabList).ListMeta}
	for _, item := range obj.(*servicev1.GitlabList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gitlabs.
func (c *FakeGitlabs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(gitlabsResource, opts))
}

// Create takes the representation of a gitlab and creates it.  Returns the server's representation of the gitlab, and an error, if there is any.
func (c *FakeGitlabs) Create(ctx context.Context, gitlab *servicev1.Gitlab, opts v1.CreateOptions) (result *servicev1.Gitlab, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(gitlabsResource, gitlab), &servicev1.Gitlab{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicev1.Gitlab), err
}

// Update takes the representation of a gitlab and updates it. Returns the server's representation of the gitlab, and an error, if there is any.
func (c *FakeGitlabs) Update(ctx context.Context, gitlab *servicev1.Gitlab, opts v1.UpdateOptions) (result *servicev1.Gitlab, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(gitlabsResource, gitlab), &servicev1.Gitlab{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicev1.Gitlab), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeGitlabs) UpdateStatus(ctx context.Context, gitlab *servicev1.Gitlab, opts v1.UpdateOptions) (*servicev1.Gitlab, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(gitlabsResource, "status", gitlab), &servicev1.Gitlab{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicev1.Gitlab), err
}

// Delete takes name of the gitlab and deletes it. Returns an error if one occurs.
func (c *FakeGitlabs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(gitlabsResource, name, opts), &servicev1.Gitlab{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGitlabs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(gitlabsResource, listOpts)

	_, err := c.Fake.Invokes(action, &servicev1.GitlabList{})
	return err
}

// Patch applies the patch and returns the patched gitlab.
func (c *FakeGitlabs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *servicev1.Gitlab, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(gitlabsResource, name, pt, data, subresources...), &servicev1.Gitlab{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicev1.Gitlab), err
}
