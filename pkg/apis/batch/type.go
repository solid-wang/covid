package batch

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DevOpsStatus string

const (
	DevOpsPendingStatus DevOpsStatus = "Pending"
	DevOpsRunningStatus DevOpsStatus = "Running"
	DevOpsFailedStatus  DevOpsStatus = "Failed"
	DevOpsSuccessStatus DevOpsStatus = "Success"
	DevOpsCancelStatus  DevOpsStatus = "Cancel"
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
	GitlabHost       string
	GitlabToken      string
	ProjectID        int
	Ref              string
	ConfigPath       string
	BuildImage       string
	FromImage        string
	BuildDir         string
	BuildCommand     string
	ArtifactPath     string
	Image            string
	Version          string
	Registry         string
	RegistryUser     string
	RegistryPassword string
}

// ContinuousIntegrationStatus is the status of a ContinuousIntegration.
type ContinuousIntegrationStatus struct {
	Phase                       DevOpsStatus
	ContinuousDeploymentTrigger *string
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
	Env              string
	KubernetesConfig []byte
	Deployment
	Service
	ConfigMap
}

type Deployment struct {
	Name            string
	Namespace       string
	Labels          map[string]string
	Replicas        int32
	Image           string
	Command         []string
	Ports           []DeploymentPort
	Variables       []Variable
	ConfigMapName   *string
	ConfigMountPath *string
}

type DeploymentPort struct {
	Name          string
	ContainerPort int32
}

type Variable struct {
	Key   string
	Value string
}

type Service struct {
	Name      string
	Namespace string
	Selector  map[string]string
	Ports     []ServicePort
}

type ServicePort struct {
	Name       string
	Port       int32
	TargetPort int
}

type ConfigMap struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Data      map[string]string
}

// ContinuousDeploymentStatus is the status of a ContinuousDeployment.
type ContinuousDeploymentStatus struct {
	Phase DevOpsStatus
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
