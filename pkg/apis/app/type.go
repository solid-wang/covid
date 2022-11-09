package app

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProductList is a list of Product objects.
type ProductList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Product
}

// ProductSpec is the specification of a Product.
type ProductSpec struct {
	Owner string
}

// ProductStatus is the status of a Product.
type ProductStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Product is a demo type with a spec and a status.
type Product struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ProductSpec
	Status ProductStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList is a list of Application objects.
type ApplicationList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Application
}

// ApplicationSpec is the specification of a Application.
type ApplicationSpec struct {
	GitlabName           string
	ProjectID            int
	DockerRepositoryName string
	BranchEnvMap         map[string]string
	ContinuousIntegrationTemplate
	ContinuousDeploymentTemplate
	Owner string
}

type ContinuousIntegrationTemplate struct {
	ConfigPath   string
	BuildImage   string
	FromImage    string
	BuildDir     string
	BuildCommand string
	ArtifactPath string
}

type ContinuousDeploymentTemplate struct {
	Command []string
	Ports   []Port
	EnvMap  map[string]*Env
}

type Port struct {
	Name       string
	Port       int32
	TargetPort int32
}

type Variable struct {
	Key   string
	Value string
}

type Env struct {
	Approval            bool
	Replicas            int32
	Variables           []Variable
	KubernetesName      string
	KubernetesNamespace string
	ConfigMap           map[string]string
	ConfigMountPath     string
}

// ApplicationStatus is the status of a Application.
type ApplicationStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Application is a Application type with a spec and a status.
type Application struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ApplicationSpec
	Status ApplicationStatus
}
