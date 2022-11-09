package options

import (
	"github.com/solid-wang/covid/pkg/controller/devops/continuousintegration/config"
	"github.com/spf13/pflag"
)

type ContinuousIntegrationControllerOptions struct {
}

// AddFlags adds flags related to ServerController for controller manager to the specified FlagSet.
func (o *ContinuousIntegrationControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
}

// ApplyTo fills up ServerController config with options.
func (o *ContinuousIntegrationControllerOptions) ApplyTo(cfg *config.ContinuousIntegrationControllerConfiguration) error {
	if o == nil {
		return nil
	}
	return nil
}

// Validate checks validation of ServerControllerOptions.
func (o *ContinuousIntegrationControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}
	errs := make([]error, 0)
	return errs
}
