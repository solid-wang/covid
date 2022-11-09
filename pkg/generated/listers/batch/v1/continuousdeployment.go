// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ContinuousDeploymentLister helps list ContinuousDeployments.
// All objects returned here must be treated as read-only.
type ContinuousDeploymentLister interface {
	// List lists all ContinuousDeployments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ContinuousDeployment, err error)
	// ContinuousDeployments returns an object that can list and get ContinuousDeployments.
	ContinuousDeployments(namespace string) ContinuousDeploymentNamespaceLister
	ContinuousDeploymentListerExpansion
}

// continuousDeploymentLister implements the ContinuousDeploymentLister interface.
type continuousDeploymentLister struct {
	indexer cache.Indexer
}

// NewContinuousDeploymentLister returns a new ContinuousDeploymentLister.
func NewContinuousDeploymentLister(indexer cache.Indexer) ContinuousDeploymentLister {
	return &continuousDeploymentLister{indexer: indexer}
}

// List lists all ContinuousDeployments in the indexer.
func (s *continuousDeploymentLister) List(selector labels.Selector) (ret []*v1.ContinuousDeployment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ContinuousDeployment))
	})
	return ret, err
}

// ContinuousDeployments returns an object that can list and get ContinuousDeployments.
func (s *continuousDeploymentLister) ContinuousDeployments(namespace string) ContinuousDeploymentNamespaceLister {
	return continuousDeploymentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ContinuousDeploymentNamespaceLister helps list and get ContinuousDeployments.
// All objects returned here must be treated as read-only.
type ContinuousDeploymentNamespaceLister interface {
	// List lists all ContinuousDeployments in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ContinuousDeployment, err error)
	// Get retrieves the ContinuousDeployment from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ContinuousDeployment, error)
	ContinuousDeploymentNamespaceListerExpansion
}

// continuousDeploymentNamespaceLister implements the ContinuousDeploymentNamespaceLister
// interface.
type continuousDeploymentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ContinuousDeployments in the indexer for a given namespace.
func (s continuousDeploymentNamespaceLister) List(selector labels.Selector) (ret []*v1.ContinuousDeployment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ContinuousDeployment))
	})
	return ret, err
}

// Get retrieves the ContinuousDeployment from the indexer for a given namespace and name.
func (s continuousDeploymentNamespaceLister) Get(name string) (*v1.ContinuousDeployment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("continuousdeployment"), name)
	}
	return obj.(*v1.ContinuousDeployment), nil
}