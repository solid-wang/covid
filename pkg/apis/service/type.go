package service

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// GitlabList is a list of Gitlab objects.
type GitlabList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Gitlab
}

// GitlabSpec is the specification of a Gitlab.
type GitlabSpec struct {
	Host  string
	Token string
	// ProjectIndex index is string projectID
	ProjectIndex map[string]Project
}

type GitlabWebhookEventType string
type ApplicationProductMap map[string]*string

const (
	GitlabWebhookEventTagMergeRequest GitlabWebhookEventType = "merge_request"
	GitlabWebhookEventPipeline        GitlabWebhookEventType = "pipeline"
)

type Project struct {
	ApplicationProductMap
	HooksMap map[GitlabWebhookEventType]*int
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

// DockerRepositoryList is a list of DockerRepository objects.
type DockerRepositoryList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []DockerRepository
}

// DockerRepositorySpec is the specification of a DockerRepository.
type DockerRepositorySpec struct {
	Host     string
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
