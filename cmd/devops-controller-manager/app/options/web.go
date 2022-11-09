package options

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/solid-wang/covid/pkg/controller/devops/web"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	"github.com/spf13/pflag"
)

type WebServerOptions struct {
	BindPort int32
}

// AddFlags adds flags related to ServerController for controller manager to the specified FlagSet.
func (o *WebServerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.Int32Var(&o.BindPort, "web-bind-port", 8080, "DevOps controller web listen port.")
}

// ApplyTo fills up ServerController config with options.
func (o *WebServerOptions) ApplyTo(cfg *web.WebConfiguration, feishu *lark.Client, client *clientset.Clientset) error {
	if o == nil {
		return nil
	}
	cfg.Port = o.BindPort
	cfg.FeiShu = feishu
	cfg.Client = client
	return nil
}

// Validate checks validation of ServerControllerOptions.
func (o *WebServerOptions) Validate() []error {
	if o == nil {
		return nil
	}
	errs := make([]error, 0)
	return errs
}
