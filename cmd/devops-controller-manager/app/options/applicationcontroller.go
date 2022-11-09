package options

import (
	"github.com/solid-wang/covid/pkg/controller/devops/application/config"
	"github.com/spf13/pflag"
)

type ApplicationControllerOptions struct {
	ExternalUrl string
}

// AddFlags adds flags related to ServerController for controller manager to the specified FlagSet.
func (o *ApplicationControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.ExternalUrl, "external-url", "", "Gitlab Webhook calling host address.")
}

// ApplyTo fills up ServerController config with options.
func (o *ApplicationControllerOptions) ApplyTo(cfg *config.ApplicationControllerConfiguration) error {
	if o == nil {
		return nil
	}
	cfg.ExternalUrl = o.ExternalUrl
	return nil
}

// Validate checks validation of ServerControllerOptions.
func (o *ApplicationControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}
	errs := make([]error, 0)
	return errs
}
