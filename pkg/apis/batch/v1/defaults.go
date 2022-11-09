package v1

import "k8s.io/apimachinery/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ContinuousIntegration(obj *ContinuousIntegration) {
	if obj.Status.Phase == "" {
		obj.Status.Phase = DevOpsPendingStatus
	}
}

func SetDefaults_ContinuousDeployment(obj *ContinuousDeployment) {
	if obj.Status.Phase == "" {
		obj.Status.Phase = DevOpsPendingStatus
	}
}
