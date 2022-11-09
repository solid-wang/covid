package app

import (
	"context"
	appv1 "github.com/solid-wang/covid/pkg/apis/app/v1"
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	servicev1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	"github.com/solid-wang/covid/pkg/apiserver"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	informers "github.com/solid-wang/covid/pkg/generated/informers/externalversions"
	covidopenapi "github.com/solid-wang/covid/pkg/generated/openapi"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"net/http"
	"os"
)

const defaultEtcdPathPrefix = "/registry/covid"

// CovidServerOptions contains state for master/api server
type CovidServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions

	SharedInformerFactory informers.SharedInformerFactory
	StdOut                io.Writer
	StdErr                io.Writer
}

// NewCovidServerOptions returns a new CovidServerOptions
func NewCovidServerOptions(out, errOut io.Writer) *CovidServerOptions {

	groupVersioners := schema.GroupVersions{
		corev1.SchemeGroupVersion,
		appv1.SchemeGroupVersion,
		batchv1.SchemeGroupVersion,
		servicev1.SchemeGroupVersion,
	}

	o := &CovidServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			scheme.Codecs.LegacyCodec(groupVersioners...),
		),

		StdOut: out,
		StdErr: errOut,
	}

	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = groupVersioners
	return o
}

// NewCommandStartCovidServer provides a CLI handler for 'start master' command
// with a default CovidServerOptions.
func NewCommandStartCovidServer(defaults *CovidServerOptions, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch a Covid API server",
		Long:  "Launch a Covid API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunCovidServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)
	//utilfeature.DefaultMutableFeatureGate.AddFlag(flags)

	return cmd
}

// Validate validates CovidServerOptions
func (o *CovidServerOptions) Validate(args []string) error {
	var errors []error
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *CovidServerOptions) Complete() error {
	return nil
}

// Config returns config for the api server given CovidServerOptions
func (o *CovidServerOptions) Config() (*apiserver.Config, error) {

	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)

	o.RecommendedOptions.Authentication = nil
	o.RecommendedOptions.Authorization = nil
	o.RecommendedOptions.CoreAPI = nil
	o.RecommendedOptions.Admission = nil

	serverConfig := genericapiserver.NewRecommendedConfig(scheme.Codecs)

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(covidopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(scheme.Scheme))
	serverConfig.OpenAPIConfig.Info.Title = "Covid"
	serverConfig.OpenAPIConfig.Info.Version = "1.0"

	if utilfeature.DefaultFeatureGate.Enabled(features.OpenAPIV3) {
		serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(covidopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(scheme.Scheme))
		serverConfig.OpenAPIV3Config.Info.Title = "Covid"
		serverConfig.OpenAPIV3Config.Info.Version = "1.0"
	}

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
	}
	return config, nil
}

// RunCovidServer starts a new CovidServer given CovidServerOptions
func (o *CovidServerOptions) RunCovidServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server.GenericAPIServer.PrepareRun().Handler,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			klog.Error(err)
			os.Exit(0)
		}
	}()

	<-stopCh
	return httpServer.Shutdown(context.Background())
}
