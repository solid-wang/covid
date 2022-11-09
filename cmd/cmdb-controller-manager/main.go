package main

import (
	"github.com/solid-wang/covid/cmd/cmdb-controller-manager/app"
	"k8s.io/component-base/cli"
	"os"
)

func main() {
	command := app.NewControllerManagerCommand()
	code := cli.Run(command)
	os.Exit(code)
}
