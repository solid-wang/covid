// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	scheme "github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// GitlabsGetter has a method to return a GitlabInterface.
// A group's client should implement this interface.
type GitlabsGetter interface {
	Gitlabs() GitlabInterface
}

// GitlabInterface has methods to work with Gitlab resources.
type GitlabInterface interface {
	Create(ctx context.Context, gitlab *v1.Gitlab, opts metav1.CreateOptions) (*v1.Gitlab, error)
	Update(ctx context.Context, gitlab *v1.Gitlab, opts metav1.UpdateOptions) (*v1.Gitlab, error)
	UpdateStatus(ctx context.Context, gitlab *v1.Gitlab, opts metav1.UpdateOptions) (*v1.Gitlab, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Gitlab, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.GitlabList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Gitlab, err error)
	GitlabExpansion
}

// gitlabs implements GitlabInterface
type gitlabs struct {
	client rest.Interface
}

// newGitlabs returns a Gitlabs
func newGitlabs(c *ServiceV1Client) *gitlabs {
	return &gitlabs{
		client: c.RESTClient(),
	}
}

// Get takes name of the gitlab, and returns the corresponding gitlab object, and an error if there is any.
func (c *gitlabs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Gitlab, err error) {
	result = &v1.Gitlab{}
	err = c.client.Get().
		Resource("gitlabs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Gitlabs that match those selectors.
func (c *gitlabs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.GitlabList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.GitlabList{}
	err = c.client.Get().
		Resource("gitlabs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested gitlabs.
func (c *gitlabs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("gitlabs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a gitlab and creates it.  Returns the server's representation of the gitlab, and an error, if there is any.
func (c *gitlabs) Create(ctx context.Context, gitlab *v1.Gitlab, opts metav1.CreateOptions) (result *v1.Gitlab, err error) {
	result = &v1.Gitlab{}
	err = c.client.Post().
		Resource("gitlabs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(gitlab).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a gitlab and updates it. Returns the server's representation of the gitlab, and an error, if there is any.
func (c *gitlabs) Update(ctx context.Context, gitlab *v1.Gitlab, opts metav1.UpdateOptions) (result *v1.Gitlab, err error) {
	result = &v1.Gitlab{}
	err = c.client.Put().
		Resource("gitlabs").
		Name(gitlab.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(gitlab).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *gitlabs) UpdateStatus(ctx context.Context, gitlab *v1.Gitlab, opts metav1.UpdateOptions) (result *v1.Gitlab, err error) {
	result = &v1.Gitlab{}
	err = c.client.Put().
		Resource("gitlabs").
		Name(gitlab.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(gitlab).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the gitlab and deletes it. Returns an error if one occurs.
func (c *gitlabs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("gitlabs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *gitlabs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("gitlabs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched gitlab.
func (c *gitlabs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Gitlab, err error) {
	result = &v1.Gitlab{}
	err = c.client.Patch(pt).
		Resource("gitlabs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
