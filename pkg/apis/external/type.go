package external

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList is a list of Application objects.
type ApplicationList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Application
}

// ApplicationSpec is the specification of a Application.
type ApplicationSpec struct {
	GitlabReference
	ProductName string
	Owner       string
}

type GitlabReference struct {
	Name      string
	ProjectID int
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
