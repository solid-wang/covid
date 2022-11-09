// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	appv1 "github.com/solid-wang/covid/pkg/apis/app/v1"
	versioned "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/solid-wang/covid/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/solid-wang/covid/pkg/generated/listers/app/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ProductInformer provides access to a shared informer and lister for
// Products.
type ProductInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ProductLister
}

type productInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewProductInformer constructs a new informer for Product type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewProductInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredProductInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredProductInformer constructs a new informer for Product type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredProductInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppV1().Products().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppV1().Products().Watch(context.TODO(), options)
			},
		},
		&appv1.Product{},
		resyncPeriod,
		indexers,
	)
}

func (f *productInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredProductInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *productInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&appv1.Product{}, f.defaultInformer)
}

func (f *productInformer) Lister() v1.ProductLister {
	return v1.NewProductLister(f.Informer().GetIndexer())
}
