package app

import (
	"context"
	"fmt"
	"github.com/solid-wang/covid/cmd/cmdb-controller-manager/app/config"
	"github.com/solid-wang/covid/cmd/cmdb-controller-manager/app/options"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"os"
)

// NewControllerManagerCommand creates a *cobra.Command object with default parameters
func NewControllerManagerCommand() *cobra.Command {
	s := options.NewCMDBControllerOptions()

	cmd := &cobra.Command{
		Use: "devops-controller-manager",
		Long: `The devops controller manager is a daemon that embeds
the core control loops shipped with Kubernetes. In applications of robotics and
automation, a control loop is a non-terminating loop that regulates the state of
the system. In Kubernetes, a controller is a control loop that watches the shared
state of the cluster through the apiserver and makes changes attempting to move the
current state towards the desired state. Examples of controllers that ship with
Kubernetes today are the replication controller, endpoints controller, namespace
controller, and serviceaccounts controller.`,
		Run: func(cmd *cobra.Command, args []string) {
			cliflag.PrintFlags(cmd.Flags())

			c, err := s.Config()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := Run(c.Complete(), wait.NeverStop); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	s.Flags(fs)

	return cmd
}

// InitFunc is used to launch a particular controller. It returns a controller
// that can optionally implement other interfaces so that the controller manager
// can support the requested features.
// The returned controller may be nil, which will be considered an anonymous controller
// that requests no additional features from the controller manager.
// Any error returned will cause the controller process to `Fatal`
// The bool indicates whether the controller was enabled.
type InitFunc func(ctx context.Context, c *config.CompletedConfig) error

// NewControllerInitializers is a public map of named controller groups (you can start more than one in an init func)
// paired to their InitFunc.  This allows for structured downstream composition and subdivision.
func NewControllerInitializers() map[string]InitFunc {
	controllers := map[string]InitFunc{}
	controllers["server"] = startServerController

	return controllers
}

// Run runs the KubeControllerManagerOptions.  This should never exit.
func Run(c *config.CompletedConfig, stopCh <-chan struct{}) error {

	if err := StartControllers(context.TODO(), c, NewControllerInitializers()); err != nil {
		klog.Fatalf("error starting controllers: %v", err)
	}

	c.InformerFactory.Start(stopCh)

	select {}
}

// StartControllers starts a set of controllers with a specified ControllerContext
func StartControllers(ctx context.Context, c *config.CompletedConfig, controllers map[string]InitFunc) error {
	for controllerName, initFn := range controllers {
		klog.V(1).Infof("Starting %q", controllerName)
		err := initFn(ctx, c)
		if err != nil {
			klog.Errorf("Error starting %q", controllerName)
			return err
		}
		klog.Infof("Started %q", controllerName)
	}

	return nil
}
