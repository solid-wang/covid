package v1

import "k8s.io/apimachinery/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_Gitlab(obj *Gitlab) {
	if obj.Spec.ProjectIndex == nil {
		obj.Spec.ProjectIndex = make(map[string]Project)
	}
}
