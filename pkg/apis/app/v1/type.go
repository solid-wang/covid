package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProductList is a list of Product objects.
type ProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Product `json:"items"`
}

// ProductSpec is the specification of a Product.
type ProductSpec struct {
	Owner string `json:"owner"`
}

// ProductStatus is the status of a Product.
type ProductStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Product is a demo type with a spec and a status.
type Product struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductSpec   `json:"spec,omitempty"`
	Status ProductStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList is a list of Application objects.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Application `json:"items"`
}

// ApplicationSpec is the specification of a Application.
type ApplicationSpec struct {
	GitlabName                    string            `json:"gitlabName"`
	ProjectID                     int               `json:"projectID"`
	DockerRepositoryName          string            `json:"dockerRepositoryName"`
	BranchEnvMap                  map[string]string `json:"branchEnvMap"`
	ContinuousIntegrationTemplate `json:"continuousIntegrationTemplate"`
	ContinuousDeploymentTemplate  `json:"continuousDeploymentTemplate"`
	Owner                         string `json:"owner"`
}

type ContinuousIntegrationTemplate struct {
	ConfigPath   string `json:"configPath"`
	BuildImage   string `json:"buildImage"`
	FromImage    string `json:"fromImage"`
	BuildDir     string `json:"buildDir"`
	BuildCommand string `json:"buildCommand"`
	ArtifactPath string `json:"artifactPath"`
}

type ContinuousDeploymentTemplate struct {
	Command []string        `json:"command"`
	Ports   []Port          `json:"ports"`
	EnvMap  map[string]*Env `json:"envMap"`
}

type Port struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Env struct {
	Approval            bool              `json:"approval"`
	Replicas            int32             `json:"replicas"`
	Variables           []Variable        `json:"variables"`
	KubernetesName      string            `json:"kubernetesName"`
	KubernetesNamespace string            `json:"kubernetesNamespace"`
	ConfigMap           map[string]string `json:"configMap"`
	ConfigMountPath     string            `json:"configMountPath"`
}

// ApplicationStatus is the status of a Application.
type ApplicationStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Application is a Application type with a spec and a status.
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}
