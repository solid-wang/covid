package app

import (
	"context"
	"fmt"
	"github.com/solid-wang/covid/cmd/cmdb-controller-manager/app/config"
	"github.com/solid-wang/covid/pkg/controller/cmdb/server"
)

func startServerController(ctx context.Context, c *config.CompletedConfig) error {
	ServerController, err := server.NewServerController(
		c.InformerFactory.Cmdb().V1().Servers(),
		c.InformerFactory.Cmdb().V1().Gitlabs(),
		c.InformerFactory.Devops().V1().ContinuousIntegrations(),
		c.InformerFactory.Devops().V1().ServerVersions(),
		c.InformerFactory.Devops().V1().ContinuousDeployments(),
		c.Client,
	)
	if err != nil {
		return fmt.Errorf("failed to start server controller: %v", err)
	}
	go ServerController.Run(ctx, 1)
	return nil
}
