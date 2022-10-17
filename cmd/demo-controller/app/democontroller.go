package app

import (
	"github.com/solid-wang/covid/pkg/democontroller"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	informers "github.com/solid-wang/covid/pkg/generated/informers/externalversions"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"time"
)

//type DemoController struct {
//}

func NewControllerManagerCommand(option *DemoControllerOptions, stopCh <-chan struct{}) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Launch a Demo Controller",
		Long:  "Launch a Demo Controller",
		RunE: func(c *cobra.Command, args []string) error {
			if err := option.Complete(); err != nil {
				return err
			}
			if err := option.Validate(args); err != nil {
				return err
			}
			if err := option.RunCovidServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	option.AddFlags(flags)
	return cmd
}

type DemoControllerOptions struct {
	CovidServer string
}

// NewControllerOptions returns a new CovidServerOptions
func NewControllerOptions() *DemoControllerOptions {
	return &DemoControllerOptions{}
}

func (o *DemoControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.CovidServer, "covid-server", o.CovidServer, "Covid server address.")
}

func (o *DemoControllerOptions) Complete() error {
	return nil
}

func (o *DemoControllerOptions) Validate(args []string) error {
	return nil
}

func (o *DemoControllerOptions) RunCovidServer(stopCh <-chan struct{}) error {
	config := &rest.Config{Host: o.CovidServer}
	clientSet, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)
	controller := democontroller.NewController(clientSet, informerFactory.Example().V1().Demo1s(), informerFactory.Group().V1().Demos())
	informerFactory.Start(stopCh)
	if err = controller.Run(1, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
	<-stopCh
	return nil
}
