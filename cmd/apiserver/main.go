package main

import (
	"github.com/solid-wang/covid/cmd/apiserver/app"
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()
	options := app.NewCovidServerOptions(os.Stdout, os.Stderr)
	cmd := app.NewCommandStartCovidServer(options, stopCh)
	code := cli.Run(cmd)
	os.Exit(code)
}
