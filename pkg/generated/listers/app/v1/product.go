// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/solid-wang/covid/pkg/apis/app/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ProductLister helps list Products.
// All objects returned here must be treated as read-only.
type ProductLister interface {
	// List lists all Products in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Product, err error)
	// Get retrieves the Product from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Product, error)
	ProductListerExpansion
}

// productLister implements the ProductLister interface.
type productLister struct {
	indexer cache.Indexer
}

// NewProductLister returns a new ProductLister.
func NewProductLister(indexer cache.Indexer) ProductLister {
	return &productLister{indexer: indexer}
}

// List lists all Products in the indexer.
func (s *productLister) List(selector labels.Selector) (ret []*v1.Product, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Product))
	})
	return ret, err
}

// Get retrieves the Product from the index for a given name.
func (s *productLister) Get(name string) (*v1.Product, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("product"), name)
	}
	return obj.(*v1.Product), nil
}