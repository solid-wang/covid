package v1

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
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ContinuousIntegration `json:"items"`
}

// ContinuousIntegrationSpec is the specification of a ContinuousIntegration.
type ContinuousIntegrationSpec struct {
	GitlabHost       string `json:"gitlabHost"`
	GitlabToken      string `json:"gitlabToken"`
	ProjectID        int    `json:"projectID"`
	Ref              string `json:"ref"`
	ConfigPath       string `json:"configPath"`
	BuildImage       string `json:"buildImage"`
	FromImage        string `json:"fromImage"`
	BuildDir         string `json:"buildDir"`
	BuildCommand     string `json:"buildCommand"`
	ArtifactPath     string `json:"artifactPath"`
	Image            string `json:"image"`
	Version          string `json:"version"`
	Registry         string `json:"registry"`
	RegistryUser     string `json:"registryUser"`
	RegistryPassword string `json:"registryPassword"`
}

// ContinuousIntegrationStatus is the status of a ContinuousIntegration.
type ContinuousIntegrationStatus struct {
	Phase                       DevOpsStatus `json:"phase"`
	ContinuousDeploymentTrigger *string      `json:"continuousDeploymentTrigger"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousIntegration is a ContinuousIntegration type with a spec and a status.
type ContinuousIntegration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContinuousIntegrationSpec   `json:"spec,omitempty"`
	Status ContinuousIntegrationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousDeploymentList is a list of ContinuousDeployment objects.
type ContinuousDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ContinuousDeployment `json:"items"`
}

// ContinuousDeploymentSpec is the specification of a ContinuousDeployment.
type ContinuousDeploymentSpec struct {
	Env              string `json:"env"`
	KubernetesConfig []byte `json:"kubernetesConfig"`
	Deployment       `json:"deployment"`
	Service          `json:"service"`
	ConfigMap        `json:"configMap"`
}

type Deployment struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Labels          map[string]string `json:"labels"`
	Replicas        int32             `json:"replicas"`
	Image           string            `json:"image"`
	Command         []string          `json:"command"`
	Ports           []DeploymentPort  `json:"ports"`
	Variables       []Variable        `json:"variables"`
	ConfigMapName   *string           `json:"configMapName"`
	ConfigMountPath *string           `json:"configMountPath"`
}

type DeploymentPort struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Service struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Selector  map[string]string `json:"selector"`
	Ports     []ServicePort     `json:"ports"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort int    `json:"targetPort"`
}

type ConfigMap struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Data      map[string]string `json:"data"`
}

// ContinuousDeploymentStatus is the status of a ContinuousDeployment.
type ContinuousDeploymentStatus struct {
	Phase DevOpsStatus `json:"phase"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousDeployment is a ContinuousDeployment type with a spec and a status.
type ContinuousDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContinuousDeploymentSpec   `json:"spec,omitempty"`
	Status ContinuousDeploymentStatus `json:"status,omitempty"`
}
