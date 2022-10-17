// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/solid-wang/covid/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Demo1s returns a Demo1Informer.
	Demo1s() Demo1Informer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Demo1s returns a Demo1Informer.
func (v *version) Demo1s() Demo1Informer {
	return &demo1Informer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
