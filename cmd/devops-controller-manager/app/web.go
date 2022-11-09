package app

import (
	"github.com/solid-wang/covid/pkg/controller/devops/web"
	"github.com/solid-wang/covid/pkg/ginserver"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	server = ginserver.NewServer()
)

func StartWebServer(c *web.WebConfiguration) error {
	err := server.InstallAPIs(web.GetAPIGroupInfos(c))
	if err != nil {
		return err
	}
	go server.Run(c.Port, wait.NeverStop)
	return nil
}
