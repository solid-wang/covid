// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	servicev1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	versioned "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/solid-wang/covid/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/solid-wang/covid/pkg/generated/listers/service/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GitlabInformer provides access to a shared informer and lister for
// Gitlabs.
type GitlabInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.GitlabLister
}

type gitlabInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewGitlabInformer constructs a new informer for Gitlab type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGitlabInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGitlabInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredGitlabInformer constructs a new informer for Gitlab type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGitlabInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ServiceV1().Gitlabs().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ServiceV1().Gitlabs().Watch(context.TODO(), options)
			},
		},
		&servicev1.Gitlab{},
		resyncPeriod,
		indexers,
	)
}

func (f *gitlabInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGitlabInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gitlabInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&servicev1.Gitlab{}, f.defaultInformer)
}

func (f *gitlabInformer) Lister() v1.GitlabLister {
	return v1.NewGitlabLister(f.Informer().GetIndexer())
}
