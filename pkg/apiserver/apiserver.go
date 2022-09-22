package apiserver

import (
	"github.com/solid-wang/covid/pkg/apis/example"
	exampleInstall "github.com/solid-wang/covid/pkg/apis/example/install"
	"github.com/solid-wang/covid/pkg/apis/group"
	groupInstall "github.com/solid-wang/covid/pkg/apis/group/install"

	"github.com/solid-wang/covid/pkg/registry"
	demo1storage "github.com/solid-wang/covid/pkg/registry/example/demo1"
	demostorage "github.com/solid-wang/covid/pkg/registry/group/demo"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	exampleInstall.Install(Scheme)
	groupInstall.Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
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

	exampleGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(example.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	groupGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(group.GroupName, Scheme, metav1.ParameterCodec, Codecs)

	exampleGroupInfo.VersionedResourcesStorageMap["v1"] = map[string]rest.Storage{
		"demo1s": registry.RESTInPeace(demo1storage.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)),
	}

	groupGroupInfo.VersionedResourcesStorageMap["v1"] = map[string]rest.Storage{
		"demos": registry.RESTInPeace(demostorage.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)),
	}

	groupGroupInfo.VersionedResourcesStorageMap["v1beta1"] = map[string]rest.Storage{
		"demos": registry.RESTInPeace(demostorage.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)),
	}

	if err := s.GenericAPIServer.InstallAPIGroups(&exampleGroupInfo, &groupGroupInfo); err != nil {
		return nil, err
	}

	return s, nil
}
