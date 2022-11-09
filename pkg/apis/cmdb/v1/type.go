package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	FinalizerServer = "server"
	FinalizerGitlab = "gitlab"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesList is a list of Kubernetes objects.
type KubernetesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Kubernetes `json:"items"`
}

// KubernetesSpec is the specification of a Kubernetes.
type KubernetesSpec struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

// KubernetesStatus is the status of a Kubernetes.
type KubernetesStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Kubernetes is a demo type with a spec and a status.
type Kubernetes struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesSpec   `json:"spec,omitempty"`
	Status KubernetesStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProductList is a list of Product objects.
type ProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Product `json:"items"`
}

// ProductSpec is the specification of a Product.
type ProductSpec struct {
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

// GitlabList is a list of Gitlab objects.
type GitlabList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Gitlab `json:"items"`
}

type ProjectStringID string

// GitlabSpec is the specification of a Gitlab.
type GitlabSpec struct {
	URL          string                      `json:"url"`
	Token        string                      `json:"token"`
	ProjectIndex map[ProjectStringID]Project `json:"projectIndex"`
}

type GitlabWebhookEventType string
type ServerName string
type ServerProduct string

const (
	GitlabWebhookEventTagPush  GitlabWebhookEventType = "tag_push"
	GitlabWebhookEventPipeline GitlabWebhookEventType = "pipeline"
)

type Project struct {
	ServersMap map[ServerName]ServerProduct    `json:"serversMap"`
	HooksMap   map[GitlabWebhookEventType]*int `json:"hooksMap"`
}

// GitlabStatus is the status of a Gitlab.
type GitlabStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Gitlab is a demo type with a spec and a status.
type Gitlab struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitlabSpec   `json:"spec,omitempty"`
	Status GitlabStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FeiShuAppList is a list of FeiShuApp objects.
type FeiShuAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []FeiShuApp `json:"items"`
}

// FeiShuAppSpec is the specification of a FeiShuApp.
type FeiShuAppSpec struct {
	ID           string `json:"id"`
	Secret       string `json:"secret"`
	ApprovalCode string `json:"approvalCode"`
}

// FeiShuAppStatus is the status of a FeiShuApp.
type FeiShuAppStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// FeiShuApp is a demo type with a spec and a status.
type FeiShuApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeiShuAppSpec   `json:"spec,omitempty"`
	Status FeiShuAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServerList is a list of Server objects.
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Server `json:"items"`
}

// ServerSpec is the specification of a Server.
type ServerSpec struct {
	GitlabInfo                `json:"gitlabInfo"`
	ContinuousIntegrationInfo `json:"continuousIntegrationInfo"`
	ContinuousDeploymentInfo  `json:"continuousDeploymentInfo"`
}

type GitlabInfo struct {
	Name      string `json:"name"`
	ProjectID int    `json:"projectID"`
}

type ContinuousIntegrationInfo struct {
	ConfigPath     string `json:"configPath"`
	BuildImage     string `json:"buildImage"`
	FromImage      string `json:"fromImage"`
	BuildDir       string `json:"buildDir"`
	BuildCommand   string `json:"buildCommand"`
	ArtifactPath   string `json:"artifactPath"`
	PushRepository string `json:"pushRepository"`
}

type ContinuousDeploymentInfo struct {
	EnvMap map[string]EnvInfo `json:"envMap"`
}

type EnvInfo struct {
	Approval       bool `json:"approval"`
	KubernetesInfo `json:"kubernetesInfo"`
}

type KubernetesInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// ServerStatus is the status of a Server.
type ServerStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Server is a demo type with a spec and a status.
type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerSpec   `json:"spec,omitempty"`
	Status ServerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DockerRepositoryList is a list of DockerRepository objects.
type DockerRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DockerRepository `json:"items"`
}

// DockerRepositorySpec is the specification of a DockerRepository.
type DockerRepositorySpec struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// DockerRepositoryStatus is the status of a DockerRepository.
type DockerRepositoryStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// DockerRepository is a demo type with a spec and a status.
type DockerRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DockerRepositorySpec   `json:"spec,omitempty"`
	Status DockerRepositoryStatus `json:"status,omitempty"`
}
