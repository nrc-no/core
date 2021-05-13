package v1

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/api/patch"
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/v1"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/client/rest"
	"github.com/nrc-no/core/apps/api/pkg/server/extensions-apiserver/client/clientset/scheme"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	"time"
)

// CustomResourceDefinitionsGetter has a method to return a CustomResourceDefinitionInterface.
// A group's client should implement this interface.
type CustomResourceDefinitionsGetter interface {
	CustomResourceDefinitions() CustomResourceDefinitionInterface
}

// CustomResourceDefinitionInterface has methods to work with CustomResourceDefinition resources.
type CustomResourceDefinitionInterface interface {
	Create(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.CreateOptions) (*v1.CustomResourceDefinition, error)
	Update(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.UpdateOptions) (*v1.CustomResourceDefinition, error)
	UpdateStatus(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.UpdateOptions) (*v1.CustomResourceDefinition, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.CustomResourceDefinition, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.CustomResourceDefinitionList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	//Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.CustomResourceDefinition, err error)
}

// customResourceDefinitions implements CustomResourceDefinitionInterface
type customResourceDefinitions struct {
	client rest.Interface
}

// newCustomResourceDefinitions returns a CustomResourceDefinitions
func newCustomResourceDefinitions(c *ApiextensionsV1Client) *customResourceDefinitions {
	return &customResourceDefinitions{
		client: c.RESTClient(),
	}
}

// Get takes name of the customResourceDefinition, and returns the corresponding customResourceDefinition object, and an error if there is any.
func (c *customResourceDefinitions) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Get().
		Resource("customresourcedefinitions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CustomResourceDefinitions that match those selectors.
func (c *customResourceDefinitions) List(ctx context.Context, opts metav1.ListOptions) (result *v1.CustomResourceDefinitionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.CustomResourceDefinitionList{}
	err = c.client.Get().
		Resource("customresourcedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested customResourceDefinitions.
func (c *customResourceDefinitions) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("customresourcedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a customResourceDefinition and creates it.  Returns the server's representation of the customResourceDefinition, and an error, if there is any.
func (c *customResourceDefinitions) Create(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.CreateOptions) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Post().
		Resource("customresourcedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(customResourceDefinition).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a customResourceDefinition and updates it. Returns the server's representation of the customResourceDefinition, and an error, if there is any.
func (c *customResourceDefinitions) Update(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.UpdateOptions) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Put().
		Resource("customresourcedefinitions").
		Name(customResourceDefinition.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(customResourceDefinition).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *customResourceDefinitions) UpdateStatus(ctx context.Context, customResourceDefinition *v1.CustomResourceDefinition, opts metav1.UpdateOptions) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Put().
		Resource("customresourcedefinitions").
		Name(customResourceDefinition.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(customResourceDefinition).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the customResourceDefinition and deletes it. Returns an error if one occurs.
func (c *customResourceDefinitions) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("customresourcedefinitions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *customResourceDefinitions) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("customresourcedefinitions").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched customResourceDefinition.
func (c *customResourceDefinitions) Patch(ctx context.Context, name string, pt patch.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.CustomResourceDefinition, err error) {
	result = &v1.CustomResourceDefinition{}
	err = c.client.Patch(pt).
		Resource("customresourcedefinitions").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
