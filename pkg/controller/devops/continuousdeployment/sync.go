package continuousdeployment

import (
	"context"
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"time"
)

const deployTimeout = time.Minute

func (c *Controller) deploy(ctx context.Context, cd *batchv1.ContinuousDeployment) {
	startTime := time.Now()
	klog.V(4).InfoS("Started deploy", "deploy", klog.KRef(cd.Namespace, cd.Name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished deploy", "deploy", klog.KRef(cd.Namespace, cd.Name), "duration", time.Since(startTime))
	}()

	kubernetesClient, err := newKubernetesClient(cd.Spec.KubernetesConfig)
	if err != nil {
		cd.Status.Phase = batchv1.DevOpsFailedStatus
	}
	configMap := getConfigMap(&cd.Spec.ConfigMap)
	service := getService(&cd.Spec.Service)
	deployment := getDeployment(&cd.Spec.Deployment)

	if err := applyConfigMap(ctx, kubernetesClient, configMap); err != nil {
		cd.Status.Phase = batchv1.DevOpsFailedStatus
	}

	if err := applyService(ctx, kubernetesClient, service); err != nil {
		cd.Status.Phase = batchv1.DevOpsFailedStatus
	}

	if err := applyDeployment(ctx, kubernetesClient, deployment); err != nil {
		cd.Status.Phase = batchv1.DevOpsFailedStatus
	}

	timer := time.NewTimer(deployTimeout)
	c.waitingDeploymentReady(cd, kubernetesClient, timer.C)
	c.client.BatchV1().ContinuousDeployments(cd.Namespace).UpdateStatus(ctx, cd, metav1.UpdateOptions{})

}

func (c *Controller) waitingDeploymentReady(cd *batchv1.ContinuousDeployment, client *kubernetes.Clientset, timeout <-chan time.Time) {
	klog.V(2).InfoS("waiting deployment ready", "deployment", klog.KRef(cd.Spec.Deployment.Namespace, cd.Spec.Deployment.Name))
	for {
		time.Sleep(3 * time.Second)
		select {
		case <-timeout:
			cd.Status.Phase = batchv1.DevOpsFailedStatus
			return
		default:
			check, err := client.AppsV1().Deployments(cd.Spec.Deployment.Namespace).Get(context.Background(), cd.Spec.Deployment.Name, metav1.GetOptions{})
			if err != nil {
				cd.Status.Phase = batchv1.DevOpsFailedStatus
				return
			}
			if check.Status.AvailableReplicas == check.Status.ReadyReplicas && check.Status.AvailableReplicas == check.Status.Replicas {
				cd.Status.Phase = batchv1.DevOpsSuccessStatus
				return
			}
		}
	}
}

func newKubernetesClient(c []byte) (*kubernetes.Clientset, error) {
	k8sCLient, err := clientcmd.NewClientConfigFromBytes(c)
	if err != nil {
		return nil, err
	}
	config, err := k8sCLient.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func applyConfigMap(ctx context.Context, clientSet *kubernetes.Clientset, configMap *corev1.ConfigMap) error {
	_, err := clientSet.CoreV1().ConfigMaps(configMap.Namespace).Get(ctx, configMap.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		_, err := clientSet.CoreV1().ConfigMaps(configMap.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	_, err = clientSet.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func applyService(ctx context.Context, clientSet *kubernetes.Clientset, svc *corev1.Service) error {
	_, err := clientSet.CoreV1().Services(svc.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		_, err := clientSet.CoreV1().Services(svc.Namespace).Create(ctx, svc, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	_, err = clientSet.CoreV1().Services(svc.Namespace).Update(ctx, svc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func applyDeployment(ctx context.Context, clientSet *kubernetes.Clientset, deployment *appv1.Deployment) error {
	_, err := clientSet.AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		_, err := clientSet.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	_, err = clientSet.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
