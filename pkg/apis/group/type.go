package group

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DemoList is a list of Server objects.
type DemoList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Demo
}

// DemoSpec is the specification of a Demo.
type DemoSpec struct {
	V1 string
}

// DemoStatus is the status of a Demo.
type DemoStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Demo is a demo type with a spec and a status.
type Demo struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   DemoSpec
	Status DemoStatus
}
