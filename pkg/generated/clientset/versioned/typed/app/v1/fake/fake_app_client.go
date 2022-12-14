// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/app/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAppV1 struct {
	*testing.Fake
}

func (c *FakeAppV1) Applications(namespace string) v1.ApplicationInterface {
	return &FakeApplications{c, namespace}
}

func (c *FakeAppV1) Products() v1.ProductInterface {
	return &FakeProducts{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAppV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
