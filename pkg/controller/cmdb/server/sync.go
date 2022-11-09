package server

import (
	"context"
	cmdbv1 "github.com/solid-wang/covid/pkg/apis/cmdb/v1"
	devopsv1 "github.com/solid-wang/covid/pkg/apis/devops/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"strconv"
	"strings"
)

func (c *Controller) ensureProductExists(ctx context.Context, server *cmdbv1.Server) error {
	_, err := c.productLister.Get(server.Namespace)
	if errors.IsNotFound(err) {
		product := &cmdbv1.Product{
			ObjectMeta: metav1.ObjectMeta{
				Name: server.Namespace,
			},
		}
		_, err = c.client.CmdbV1().Products().Create(ctx, product, metav1.CreateOptions{})
	}
	return err
}

func (c *Controller) syncGitlabIndex(ctx context.Context, server *cmdbv1.Server) error {
	gitlabName := server.Spec.GitlabInfo.Name
	gitlabProjectID := strconv.Itoa(server.Spec.GitlabInfo.ProjectID)
	g, err := c.gitlabLister.Get(gitlabName)
	if err != nil {
		return err
	}

	gitlab := g.DeepCopy()
	if gitlab.Spec.ProjectIndex == nil {
		gitlab.Spec.ProjectIndex = make(map[cmdbv1.ProjectStringID]cmdbv1.Project)
	}
	p := gitlab.Spec.ProjectIndex[cmdbv1.ProjectStringID(gitlabProjectID)]
	if p.ServersMap == nil {
		p.ServersMap = make(map[cmdbv1.ServerName]cmdbv1.ServerProduct)
	}
	if p.HooksMap == nil {
		p.HooksMap = make(map[cmdbv1.GitlabWebhookEventType]*int)
	}
	p.ServersMap[cmdbv1.ServerName(server.Name)] = cmdbv1.ServerProduct(server.Namespace)
	if server.DeletionTimestamp != nil {
		delete(p.ServersMap, cmdbv1.ServerName(server.Name))
	}
	_, ok := p.HooksMap[cmdbv1.GitlabWebhookEventTagPush]
	if !ok {
		p.HooksMap[cmdbv1.GitlabWebhookEventTagPush] = nil
	}
	_, ok = p.HooksMap[cmdbv1.GitlabWebhookEventPipeline]
	if !ok {
		p.HooksMap[cmdbv1.GitlabWebhookEventPipeline] = nil
	}
	gitlab.Spec.ProjectIndex[cmdbv1.ProjectStringID(gitlabProjectID)] = p

	if !reflect.DeepEqual(g, gitlab) {
		_, err := c.client.CmdbV1().Gitlabs().Update(ctx, gitlab, metav1.UpdateOptions{})
		return err
	}
	return nil
}

func (c *Controller) syncContinuousIntegration(ctx context.Context, server *cmdbv1.Server) error {
	ci, err := c.continuousIntegrationLister.ContinuousIntegrations(server.Namespace).Get(server.Name)
	continuousIntegration := ci.DeepCopy()
	if errors.IsNotFound(err) {
		if server.DeletionTimestamp != nil {
			return nil
		}
		continuousIntegration = &devopsv1.ContinuousIntegration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      server.Name,
				Namespace: server.Namespace,
			},
			Spec: devopsv1.ContinuousIntegrationSpec{},
		}
	}

	if server.DeletionTimestamp != nil {
		return c.client.DevopsV1().ContinuousIntegrations(continuousIntegration.Namespace).Delete(ctx, continuousIntegration.Name, metav1.DeleteOptions{})
	}

	continuousIntegration.Spec.CIConfigPath = server.Spec.ConfigPath
	continuousIntegration.Spec.BuildImage = server.Spec.BuildImage
	continuousIntegration.Spec.FromImage = server.Spec.FromImage
	continuousIntegration.Spec.BuildDir = server.Spec.BuildDir
	continuousIntegration.Spec.BuildCommand = server.Spec.BuildCommand
	continuousIntegration.Spec.ArtifactPath = server.Spec.ArtifactPath
	continuousIntegration.Spec.Image = strings.Join([]string{server.Spec.PushRepository, server.Namespace, server.Name}, "/")
	if continuousIntegration.Spec.BuiltCommitHistory == nil {
		continuousIntegration.Spec.BuiltCommitHistory = make([]string, 0, 25)
	}

	if !reflect.DeepEqual(ci, continuousIntegration) {
		_, err := c.client.DevopsV1().ContinuousIntegrations(continuousIntegration.Namespace).Update(ctx, continuousIntegration, metav1.UpdateOptions{})
		return err
	}
	return nil
}

func (c *Controller) syncContinuousDeployment(ctx context.Context, server *cmdbv1.Server) error {
	cd, err := c.continuousDeploymentLister.ContinuousDeployments(server.Namespace).Get(server.Name)
	continuousDeployment := cd.DeepCopy()
	if errors.IsNotFound(err) {
		if server.DeletionTimestamp != nil {
			return nil
		}
		continuousDeployment = &devopsv1.ContinuousDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      server.Name,
				Namespace: server.Namespace,
			},
			Spec: devopsv1.ContinuousDeploymentSpec{},
		}
	}

	if server.DeletionTimestamp != nil {
		return c.client.DevopsV1().ContinuousDeployments(server.Namespace).Delete(ctx, continuousDeployment.Name, metav1.DeleteOptions{})
	}

	if continuousDeployment.Spec.EnvDeployMap == nil {
		continuousDeployment.Spec.EnvDeployMap = make(map[string]devopsv1.Deploy)
	}

	for env, info := range server.Spec.EnvMap {
		continuousDeployment.Spec.EnvDeployMap[env] = devopsv1.Deploy{
			Approval:               info.Approval,
			KubernetesName:         info.KubernetesInfo.Name,
			KubernetesNamespace:    info.KubernetesInfo.Namespace,
			DeployedManifestCommit: "",
			Deploying:              devopsv1.Deploying{},
		}
	}

	if !reflect.DeepEqual(cd, continuousDeployment) {
		_, err := c.client.DevopsV1().ContinuousDeployments(continuousDeployment.Namespace).Update(ctx, continuousDeployment, metav1.UpdateOptions{})
		return err
	}
	return nil
}
