package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Demo1List is a list of Server objects.
type Demo1List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Demo1 `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// Demo1Spec is the specification of a Demo1.
type Demo1Spec struct {
	V1 string `json:"v1,omitempty" protobuf:"bytes,1,opt,name=v1"`
}

// Demo1Status is the status of a Demo1.
type Demo1Status struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Demo1 is a demo type with a spec and a status.
type Demo1 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   Demo1Spec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status Demo1Status `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
