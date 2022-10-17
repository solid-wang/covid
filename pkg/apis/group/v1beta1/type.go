package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DemoList is a list of Server objects.
type DemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Demo `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// DemoSpec is the specification of a Demo.
type DemoSpec struct {
	V1 string `json:"v1,omitempty" protobuf:"bytes,1,opt,name=v1beta1"`
}

// DemoStatus is the status of a Demo.
type DemoStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Demo is a demo type with a spec and a status.
type Demo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   DemoSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status DemoStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
