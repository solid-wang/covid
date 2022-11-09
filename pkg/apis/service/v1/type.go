package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// GitlabList is a list of Gitlab objects.
type GitlabList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Gitlab `json:"items"`
}

// GitlabSpec is the specification of a Gitlab.
type GitlabSpec struct {
	Host  string `json:"host"`
	Token string `json:"token"`
	// ProjectIndex index is string projectID
	ProjectIndex map[string]Project `json:"projectIndex"`
}

type GitlabWebhookEventType string
type ApplicationProductMap map[string]*string

const (
	GitlabWebhookEventTagMergeRequest GitlabWebhookEventType = "merge_request"
	GitlabWebhookEventPipeline        GitlabWebhookEventType = "pipeline"
)

type Project struct {
	ApplicationProductMap `json:"applicationProductMap"`
	HooksMap              map[GitlabWebhookEventType]*int `json:"hooksMap"`
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

// DockerRepositoryList is a list of DockerRepository objects.
type DockerRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DockerRepository `json:"items"`
}

// DockerRepositorySpec is the specification of a DockerRepository.
type DockerRepositorySpec struct {
	Host     string `json:"host"`
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
