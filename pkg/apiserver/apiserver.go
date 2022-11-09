package apiserver

import (
	appinstall "github.com/solid-wang/covid/pkg/apis/app/install"
	batchinstall "github.com/solid-wang/covid/pkg/apis/batch/install"
	coreinstall "github.com/solid-wang/covid/pkg/apis/core/install"
	serviceinstall "github.com/solid-wang/covid/pkg/apis/service/install"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/registry/generic"

	apprest "github.com/solid-wang/covid/pkg/registry/app/rest"
	batchrest "github.com/solid-wang/covid/pkg/registry/batch/rest"
	corerest "github.com/solid-wang/covid/pkg/registry/core/rest"
	servicerest "github.com/solid-wang/covid/pkg/registry/service/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func init() {
	coreinstall.Install(scheme.Scheme)
	batchinstall.Install(scheme.Scheme)
	appinstall.Install(scheme.Scheme)
	serviceinstall.Install(scheme.Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(scheme.Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	scheme.Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

type RESTStorageProvider interface {
	GroupName() string
	NewRESTStorage(scheme *runtime.Scheme, codes serializer.CodecFactory, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error)
}

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	// Place you custom config here.
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

// CovidServer contains state for a Kubernetes cluster master/api server.
type CovidServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of CovidServer from the given config.
func (c completedConfig) New() (*CovidServer, error) {
	genericServer, err := c.GenericConfig.New("Covid-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &CovidServer{
		GenericAPIServer: genericServer,
	}

	coreStoreProvider := corerest.StorageProvider{}
	coreGroupInfo, err := coreStoreProvider.NewRESTStorage(c.GenericConfig.RESTOptionsGetter)

	if err := s.GenericAPIServer.InstallLegacyAPIGroup("/api", &coreGroupInfo); err != nil {
		return nil, err
	}

	batchStoreProvider := batchrest.StorageProvider{}
	batchGroupInfo, err := batchStoreProvider.NewRESTStorage(c.GenericConfig.RESTOptionsGetter)

	appStoreProvider := apprest.StorageProvider{}
	appGroupInfo, err := appStoreProvider.NewRESTStorage(c.GenericConfig.RESTOptionsGetter)

	serviceStoreProvider := servicerest.StorageProvider{}
	serviceGroupInfo, err := serviceStoreProvider.NewRESTStorage(c.GenericConfig.RESTOptionsGetter)

	if err := s.GenericAPIServer.InstallAPIGroups(&batchGroupInfo, &appGroupInfo, &serviceGroupInfo); err != nil {
		return nil, err
	}

	return s, nil
}
