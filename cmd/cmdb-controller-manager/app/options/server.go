package options

import (
	"github.com/solid-wang/covid/pkg/controller/cmdb/server/config"
	"github.com/spf13/pflag"
)

type ServerControllerOptions struct {
}

// AddFlags adds flags related to ServerController for controller manager to the specified FlagSet.
func (o *ServerControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
}

// ApplyTo fills up ServerController config with options.
func (o *ServerControllerOptions) ApplyTo(cfg *config.ServerControllerConfiguration) error {
	if o == nil {
		return nil
	}
	return nil
}

// Validate checks validation of ServerControllerOptions.
func (o *ServerControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}
	errs := make([]error, 0)
	return errs
}
