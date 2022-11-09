package devops

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	FinalizerContinuousIntegration = "continuousIntegration"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousIntegrationList is a list of ContinuousIntegration objects.
type ContinuousIntegrationList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []ContinuousIntegration
}

// ContinuousIntegrationSpec is the specification of a ContinuousIntegration.
type ContinuousIntegrationSpec struct {
	CIConfigPath       string
	BuildImage         string
	FromImage          string
	BuildDir           string
	BuildCommand       string
	ArtifactPath       string
	Image              string
	BuiltCommitHistory []string
}

// ContinuousIntegrationStatus is the status of a ContinuousIntegration.
type ContinuousIntegrationStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousIntegration is a ContinuousIntegration type with a spec and a status.
type ContinuousIntegration struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ContinuousIntegrationSpec
	Status ContinuousIntegrationStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousDeploymentList is a list of ContinuousDeployment objects.
type ContinuousDeploymentList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []ContinuousDeployment
}

// ContinuousDeploymentSpec is the specification of a ContinuousDeployment.
type ContinuousDeploymentSpec struct {
	ManifestRepo string
	EnvDeployMap map[string]Deploy
}

type Deploy struct {
	Approval               bool
	KubernetesName         string
	KubernetesNamespace    string
	DeployedManifestCommit string
	Deploying
}

type Deploying struct {
	ManifestCommit       string
	ApprovalInstanceCode string
}

// ContinuousDeploymentStatus is the status of a ContinuousDeployment.
type ContinuousDeploymentStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousDeployment is a ContinuousDeployment type with a spec and a status.
type ContinuousDeployment struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ContinuousDeploymentSpec
	Status ContinuousDeploymentStatus
}
