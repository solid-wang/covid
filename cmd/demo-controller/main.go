package main

import (
	"github.com/solid-wang/covid/cmd/demo-controller/app"
	"github.com/solid-wang/covid/pkg/util/signals"
	"k8s.io/component-base/cli"
	"os"
)

func main() {
	stopCh := signals.SetupSignalHandler()
	options := app.NewControllerOptions()
	cmd := app.NewControllerManagerCommand(options, stopCh)
	code := cli.Run(cmd)
	os.Exit(code)
}
