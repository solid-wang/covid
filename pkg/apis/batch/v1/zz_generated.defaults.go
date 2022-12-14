//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by defaulter-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&ContinuousDeployment{}, func(obj interface{}) { SetObjectDefaults_ContinuousDeployment(obj.(*ContinuousDeployment)) })
	scheme.AddTypeDefaultingFunc(&ContinuousDeploymentList{}, func(obj interface{}) { SetObjectDefaults_ContinuousDeploymentList(obj.(*ContinuousDeploymentList)) })
	scheme.AddTypeDefaultingFunc(&ContinuousIntegration{}, func(obj interface{}) { SetObjectDefaults_ContinuousIntegration(obj.(*ContinuousIntegration)) })
	scheme.AddTypeDefaultingFunc(&ContinuousIntegrationList{}, func(obj interface{}) { SetObjectDefaults_ContinuousIntegrationList(obj.(*ContinuousIntegrationList)) })
	return nil
}

func SetObjectDefaults_ContinuousDeployment(in *ContinuousDeployment) {
	SetDefaults_ContinuousDeployment(in)
}

func SetObjectDefaults_ContinuousDeploymentList(in *ContinuousDeploymentList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_ContinuousDeployment(a)
	}
}

func SetObjectDefaults_ContinuousIntegration(in *ContinuousIntegration) {
	SetDefaults_ContinuousIntegration(in)
}

func SetObjectDefaults_ContinuousIntegrationList(in *ContinuousIntegrationList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_ContinuousIntegration(a)
	}
}
