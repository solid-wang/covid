package options

import (
	"github.com/solid-wang/covid/pkg/controller/devops/continuousdeployment/config"
	"github.com/spf13/pflag"
)

type ContinuousDeploymentControllerOptions struct {
}

// AddFlags adds flags related to ServerController for controller manager to the specified FlagSet.
func (o *ContinuousDeploymentControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
}

// ApplyTo fills up ServerController config with options.
func (o *ContinuousDeploymentControllerOptions) ApplyTo(cfg *config.ContinuousDeploymentControllerConfiguration) error {
	if o == nil {
		return nil
	}
	return nil
}

// Validate checks validation of ServerControllerOptions.
func (o *ContinuousDeploymentControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}
	errs := make([]error, 0)
	return errs
}
