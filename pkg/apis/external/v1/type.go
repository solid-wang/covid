package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList is a list of Application objects.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Application `json:"items"`
}

// ApplicationSpec is the specification of a Application.
type ApplicationSpec struct {
	GitlabReference `json:"gitlabReference"`
	ProductName     string `json:"productName"`
	Owner           string `json:"owner"`
}

type GitlabReference struct {
	Name      string `json:"name"`
	ProjectID int    `json:"projectID"`
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
