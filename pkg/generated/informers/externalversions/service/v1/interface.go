// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/solid-wang/covid/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// DockerRepositories returns a DockerRepositoryInformer.
	DockerRepositories() DockerRepositoryInformer
	// Gitlabs returns a GitlabInformer.
	Gitlabs() GitlabInformer
	// Kuberneteses returns a KubernetesInformer.
	Kuberneteses() KubernetesInformer
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

// DockerRepositories returns a DockerRepositoryInformer.
func (v *version) DockerRepositories() DockerRepositoryInformer {
	return &dockerRepositoryInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Gitlabs returns a GitlabInformer.
func (v *version) Gitlabs() GitlabInformer {
	return &gitlabInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Kuberneteses returns a KubernetesInformer.
func (v *version) Kuberneteses() KubernetesInformer {
	return &kubernetesInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
