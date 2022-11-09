package v1

import "k8s.io/apimachinery/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ContinuousIntegration(obj *ContinuousIntegration) {
	if obj.Status.Phase == "" {
		obj.Status.Phase = ContinuousIntegrationPending
	}
}
