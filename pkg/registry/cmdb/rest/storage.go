package rest

import (
	"github.com/solid-wang/covid/pkg/apis/cmdb"
	cmdbv1 "github.com/solid-wang/covid/pkg/apis/cmdb/v1"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	dockerrepositorystore "github.com/solid-wang/covid/pkg/registry/cmdb/dockerrepository"
	feishuappstore "github.com/solid-wang/covid/pkg/registry/cmdb/feishuapp"
	gitlabstore "github.com/solid-wang/covid/pkg/registry/cmdb/gitlab"
	kubernetesstore "github.com/solid-wang/covid/pkg/registry/cmdb/kubernetes"
	productstore "github.com/solid-wang/covid/pkg/registry/cmdb/product"
	serverstore "github.com/solid-wang/covid/pkg/registry/cmdb/server"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// StorageProvider is a struct for apps REST storage.
type StorageProvider struct{}

// NewRESTStorage returns APIGroupInfo object.
func (p StorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(cmdb.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)

	if storageMap, err := p.v1Storage(scheme.Scheme, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[cmdbv1.SchemeGroupVersion.Version] = storageMap
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

	// products
	productStorage, err := productstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["products"] = productStorage

	// projects
	feiShuAppStorage, err := feishuappstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["feishuapps"] = feiShuAppStorage

	// servers
	serverStorage, err := serverstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["servers"] = serverStorage

	// dockerrepositorys
	dockerrepositoryStorage, err := dockerrepositorystore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["dockerrepositorys"] = dockerrepositoryStorage

	return storage, nil
}

// GroupName returns name of the group
func (p StorageProvider) GroupName() string {
	return cmdb.GroupName
}
