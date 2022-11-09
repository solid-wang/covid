package app

import (
	"context"
	"fmt"
	"github.com/solid-wang/covid/cmd/devops-controller-manager/app/config"
	"github.com/solid-wang/covid/pkg/controller/devops/application"
	"github.com/solid-wang/covid/pkg/controller/devops/continuousdeployment"
	"github.com/solid-wang/covid/pkg/controller/devops/continuousintegration"
)

func startContinuousIntegrationController(ctx context.Context, c *config.CompletedConfig) error {
	ContinuousIntegrationController, err := continuousintegration.NewContinuousIntegrationController(
		c.InformerFactory.Batch().V1().ContinuousIntegrations(),
		c.InformerFactory.Batch().V1().ContinuousDeployments(),
		c.Client)
	if err != nil {
		return fmt.Errorf("failed to start continuous integration controller: %v", err)
	}
	go ContinuousIntegrationController.Run(ctx, 1)
	return nil
}

func startContinuousDeploymentController(ctx context.Context, c *config.CompletedConfig) error {
	ContinuousDeploymentController, err := continuousdeployment.NewContinuousDeploymentController(
		c.InformerFactory.Batch().V1().ContinuousDeployments(),
		c.Client)
	if err != nil {
		return fmt.Errorf("failed to start continuous deployment controller: %v", err)
	}
	go ContinuousDeploymentController.Run(ctx, 1)
	return nil
}

func startApplicationController(ctx context.Context, c *config.CompletedConfig) error {
	ApplicationController, err := application.NewApplicationController(
		c.InformerFactory.App().V1().Applications(),
		c.InformerFactory.Service().V1().Gitlabs(),
		c.Client,
		c.ComponentConfig.Application.ExternalUrl)
	if err != nil {
		return fmt.Errorf("failed to start application controller: %v", err)
	}
	go ApplicationController.Run(ctx, 1)
	return nil
}
