package config

import (
	"github.com/larksuite/oapi-sdk-go/v3"
	applicationconfig "github.com/solid-wang/covid/pkg/controller/devops/application/config"
	continuousdeploymentconfig "github.com/solid-wang/covid/pkg/controller/devops/continuousdeployment/config"
	continuousintegrationconfig "github.com/solid-wang/covid/pkg/controller/devops/continuousintegration/config"
	"github.com/solid-wang/covid/pkg/controller/devops/web"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	informers "github.com/solid-wang/covid/pkg/generated/informers/externalversions"
	"github.com/solid-wang/covid/pkg/tools/record"
)

// Config is the main context object for the controller manager.
type Config struct {
	WebServer web.WebConfiguration

	ComponentConfig DevOpsControllerManagerConfiguration

	FeiShu *lark.Client

	// the general covid client
	Client *clientset.Clientset

	// InformerFactory gives access to informers for the controller.
	InformerFactory informers.SharedInformerFactory

	// the event sink
	EventRecorder record.EventRecorder
}

// DevOpsControllerManagerConfiguration contains devops describing devops-controller manager.
type DevOpsControllerManagerConfiguration struct {
	ContinuousIntegration continuousintegrationconfig.ContinuousIntegrationControllerConfiguration
	ContinuousDeployment  continuousdeploymentconfig.ContinuousDeploymentControllerConfiguration
	Application           applicationconfig.ApplicationControllerConfiguration
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	cc := completedConfig{c}
	return &CompletedConfig{&cc}
}
