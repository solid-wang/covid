package rest

import (
	"github.com/solid-wang/covid/pkg/apis/batch"
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	continuousdeploymentstore "github.com/solid-wang/covid/pkg/registry/batch/continuousdeployment"
	continuousintegrationstore "github.com/solid-wang/covid/pkg/registry/batch/continuousintegration"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// StorageProvider is a struct for apps REST storage.
type StorageProvider struct{}

// NewRESTStorage returns APIGroupInfo object.
func (p StorageProvider) NewRESTStorage(restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(batch.GroupName, scheme.Scheme, scheme.ParameterCodec, scheme.Codecs)

	if storageMap, err := p.v1Storage(scheme.Scheme, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[batchv1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil
}

func (p StorageProvider) v1Storage(scheme *runtime.Scheme, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// continuousintegrations
	continuousintegrationStorage, err := continuousintegrationstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["continuousintegrations"] = continuousintegrationStorage.ContinuousIntegration
	storage["continuousintegrations/status"] = continuousintegrationStorage.Status

	// continuousdeployments
	continuousdeploymentStorage, err := continuousdeploymentstore.NewREST(scheme, restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["continuousdeployments"] = continuousdeploymentStorage.ContinuousDeployment
	storage["continuousdeployments/status"] = continuousdeploymentStorage.Status

	return storage, nil
}

// GroupName returns name of the group
func (p StorageProvider) GroupName() string {
	return batch.GroupName
}
