package continuousdeployment

import (
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getService(svc *batchv1.Service) *corev1.Service {
	ports := []corev1.ServicePort{}
	for _, port := range svc.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:       port.Name,
			Protocol:   corev1.ProtocolTCP,
			Port:       port.Port,
			TargetPort: intstr.FromInt(port.TargetPort),
		})
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Labels:    svc.Selector,
		},
		Spec: corev1.ServiceSpec{
			Selector: svc.Selector,
			Ports:    ports,
			Type:     corev1.ServiceTypeNodePort,
		},
		Status: corev1.ServiceStatus{},
	}
}

func getDeployment(deployment *batchv1.Deployment) *appv1.Deployment {
	ports := []corev1.ContainerPort{}
	for _, port := range deployment.Ports {
		ports = append(ports, corev1.ContainerPort{
			Name:          port.Name,
			ContainerPort: port.ContainerPort,
			Protocol:      corev1.ProtocolTCP,
		})
	}
	vars := []corev1.EnvVar{}
	for _, variable := range deployment.Variables {
		vars = append(vars, corev1.EnvVar{
			Name:  variable.Key,
			Value: variable.Value,
		})
	}
	d := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Labels:    deployment.Labels,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &deployment.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: deployment.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: deployment.Labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deployment.Name,
							Image:           deployment.Image,
							Command:         deployment.Command,
							Ports:           ports,
							Env:             vars,
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}
	if deployment.ConfigMapName != nil && deployment.ConfigMountPath != nil {
		volume := corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: *deployment.ConfigMapName},
				},
			},
		}
		volumeMount := corev1.VolumeMount{
			Name:      "config",
			MountPath: *deployment.ConfigMountPath,
		}
		d.Spec.Template.Spec.Volumes = []corev1.Volume{volume}
		d.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{volumeMount}
	}
	return d
}

func getConfigMap(cm *batchv1.ConfigMap) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cm.Name,
			Namespace: cm.Namespace,
			Labels:    cm.Labels,
		},
		Data: cm.Data,
	}
}
