package main

import (
	"github.com/solid-wang/covid/cmd/server"
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()
	options := server.NewCovidServerOptions(os.Stdout, os.Stderr)
	cmd := server.NewCommandStartCovidServer(options, stopCh)
	code := cli.Run(cmd)
	os.Exit(code)
}
