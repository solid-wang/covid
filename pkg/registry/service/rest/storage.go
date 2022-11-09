package rest

import (
	"github.com/solid-wang/covid/pkg/apis/service"
	servicev1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	dockerrepositorystore "github.com/solid-wang/covid/pkg/registry/service/dockerrepository"
	gitlabstore "github.com/solid-wang/covid/pkg/registry/service/gitlab"
	kubernetesstore "github.com/solid-wang/covid/pkg/registry/service/kubernetes"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// StorageProvider is a struct for apps REST storage.
type StorageProvider struct{}

// NewRESTStorage returns APIGroupInfo object.
func (p StorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(service.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)

	if storageMap, err := p.v1Storage(scheme.Scheme, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[servicev1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil
}

func (p StorageProvider) v1Storage(scheme *runtime.Scheme, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// gitlabs
	gitlabStorage, err := gitlabstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["gitlabs"] = gitlabStorage

	// kuberneteses
	kubernetesStorage, err := kubernetesstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["kuberneteses"] = kubernetesStorage

	// dockerrepositorys
	dockerrepositoryStorage, err := dockerrepositorystore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["dockerrepositories"] = dockerrepositoryStorage

	return storage, nil
}

// GroupName returns name of the group
func (p StorageProvider) GroupName() string {
	return service.GroupName
}
