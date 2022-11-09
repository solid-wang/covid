package web

import (
	appv1 "github.com/solid-wang/covid/pkg/apis/app/v1"
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	"github.com/solid-wang/covid/pkg/controller/devops/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newContinuousIntegrationsFromApp(app *appv1.Application, gitlabHost, gitlabToken, image, version, ref, registry, registryUser, registryPassword string) *batchv1.ContinuousIntegration {
	labels := map[string]string{
		util.LabelApp:     app.Name,
		util.LabelProduct: app.Namespace,
		util.LabelVersion: version,
	}
	return &batchv1.ContinuousIntegration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name + "-" + version,
			Namespace: app.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.ContinuousIntegrationSpec{
			ConfigPath:       app.Spec.ConfigPath,
			BuildImage:       app.Spec.BuildImage,
			FromImage:        app.Spec.FromImage,
			GitlabHost:       gitlabHost,
			GitlabToken:      gitlabToken,
			BuildDir:         app.Spec.BuildDir,
			BuildCommand:     app.Spec.BuildCommand,
			ArtifactPath:     app.Spec.ArtifactPath,
			Image:            image,
			ProjectID:        app.Spec.ProjectID,
			Version:          version,
			Ref:              ref,
			Registry:         registry,
			RegistryUser:     registryUser,
			RegistryPassword: registryPassword,
		},
	}
}

func newContinuousDeploymentFromApp(app *appv1.Application, k8sConfig []byte, env, image, version string) *batchv1.ContinuousDeployment {
	labels := map[string]string{
		util.LabelApp:     app.Name,
		util.LabelProduct: app.Namespace,
		util.LabelVersion: version,
		util.LabelEnv:     env,
	}
	dPorts := []batchv1.DeploymentPort{}
	sPorts := []batchv1.ServicePort{}
	for _, port := range app.Spec.ContinuousDeploymentTemplate.Ports {
		dPorts = append(dPorts, batchv1.DeploymentPort{
			Name:          port.Name,
			ContainerPort: port.TargetPort,
		})
		sPorts = append(sPorts, batchv1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: int(port.TargetPort),
		})
	}

	vars := []batchv1.Variable{}
	for _, v := range app.Spec.ContinuousDeploymentTemplate.EnvMap[env].Variables {
		vars = append(vars, batchv1.Variable{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return &batchv1.ContinuousDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name + "-" + version + "-" + env,
			Namespace: app.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.ContinuousDeploymentSpec{
			Env:              env,
			KubernetesConfig: k8sConfig,
			Deployment: batchv1.Deployment{
				Name:            app.Name,
				Namespace:       app.Spec.ContinuousDeploymentTemplate.EnvMap[env].KubernetesNamespace,
				Labels:          labels,
				Replicas:        app.Spec.ContinuousDeploymentTemplate.EnvMap[env].Replicas,
				Image:           image + ":" + version,
				Command:         app.Spec.ContinuousDeploymentTemplate.Command,
				Ports:           dPorts,
				Variables:       vars,
				ConfigMapName:   &app.Name,
				ConfigMountPath: &app.Spec.ContinuousDeploymentTemplate.EnvMap[env].ConfigMountPath,
			},
			Service: batchv1.Service{
				Name:      app.Name,
				Namespace: app.Spec.ContinuousDeploymentTemplate.EnvMap[env].KubernetesNamespace,
				Selector:  labels,
				Ports:     sPorts,
			},
			ConfigMap: batchv1.ConfigMap{
				Name:      app.Name,
				Namespace: app.Spec.ContinuousDeploymentTemplate.EnvMap[env].KubernetesNamespace,
				Labels:    labels,
				Data:      app.Spec.ContinuousDeploymentTemplate.EnvMap[env].ConfigMap,
			},
		},
	}
}
