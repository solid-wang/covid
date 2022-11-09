package rest

import (
	"github.com/solid-wang/covid/pkg/apis/core"
	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	eventstore "github.com/solid-wang/covid/pkg/registry/core/event"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// StorageProvider is a struct for apps REST storage.
type StorageProvider struct{}

// NewRESTStorage returns APIGroupInfo object.
func (p StorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(core.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)

	if storageMap, err := p.v1Storage(scheme.Scheme, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[corev1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil
}

func (p StorageProvider) v1Storage(scheme *runtime.Scheme, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// events
	eventStorage, err := eventstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["events"] = eventStorage

	return storage, nil
}

// GroupName returns name of the group
func (p StorageProvider) GroupName() string {
	return core.GroupName
}
