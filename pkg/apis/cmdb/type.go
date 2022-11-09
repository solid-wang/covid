package cmdb

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	FinalizerServer = "server"
	FinalizerGitlab = "gitlab"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesList is a list of Kubernetes objects.
type KubernetesList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Kubernetes
}

// KubernetesSpec is the specification of a Kubernetes.
type KubernetesSpec struct {
	Name   string
	Config string
}

// KubernetesStatus is the status of a Kubernetes.
type KubernetesStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Kubernetes is a demo type with a spec and a status.
type Kubernetes struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   KubernetesSpec
	Status KubernetesStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProductList is a list of Product objects.
type ProductList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Product
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
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ProductSpec
	Status ProductStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GitlabList is a list of Gitlab objects.
type GitlabList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Gitlab
}

type ProjectStringID string

// GitlabSpec is the specification of a Gitlab.
type GitlabSpec struct {
	URL          string
	Token        string
	ProjectIndex map[ProjectStringID]Project
}

type GitlabWebhookEventType string
type ServerName string
type ServerProduct string

const (
	GitlabWebhookEventTagPush  GitlabWebhookEventType = "tag_push"
	GitlabWebhookEventPipeline GitlabWebhookEventType = "pipeline"
)

type Project struct {
	ServersMap map[ServerName]ServerProduct
	HooksMap   map[GitlabWebhookEventType]int
}

// GitlabStatus is the status of a Gitlab.
type GitlabStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Gitlab is a demo type with a spec and a status.
type Gitlab struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   GitlabSpec
	Status GitlabStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FeiShuAppList is a list of FeiShuApp objects.
type FeiShuAppList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []FeiShuApp
}

// FeiShuAppSpec is the specification of a FeiShuApp.
type FeiShuAppSpec struct {
	ID           string
	Secret       string
	ApprovalCode string
}

// FeiShuAppStatus is the status of a FeiShuApp.
type FeiShuAppStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// FeiShuApp is a demo type with a spec and a status.
type FeiShuApp struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   FeiShuAppSpec
	Status FeiShuAppStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServerList is a list of Server objects.
type ServerList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Server
}

// ServerSpec is the specification of a Server.
type ServerSpec struct {
	GitlabInfo
	ContinuousIntegrationInfo
	ContinuousDeploymentInfo
}

type GitlabInfo struct {
	Name      string
	ProjectID int
}

type ContinuousIntegrationInfo struct {
	ConfigPath     string
	BuildImage     string
	FromImage      string
	BuildDir       string
	BuildCommand   string
	ArtifactPath   string
	PushRepository string
}

type ContinuousDeploymentInfo struct {
	EnvMap map[string]EnvInfo
}

type EnvInfo struct {
	Approval bool
	KubernetesInfo
}

type KubernetesInfo struct {
	Name      string
	Namespace string
}

// ServerStatus is the status of a Server.
type ServerStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Server is a demo type with a spec and a status.
type Server struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ServerSpec
	Status ServerStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DockerRepositoryList is a list of DockerRepository objects.
type DockerRepositoryList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []DockerRepository
}

// DockerRepositorySpec is the specification of a DockerRepository.
type DockerRepositorySpec struct {
	User     string
	Password string
}

// DockerRepositoryStatus is the status of a DockerRepository.
type DockerRepositoryStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// DockerRepository is a demo type with a spec and a status.
type DockerRepository struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   DockerRepositorySpec
	Status DockerRepositoryStatus
}
