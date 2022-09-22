package example

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Demo1List is a list of Server objects.
type Demo1List struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Demo1
}

// Demo1Spec is the specification of a Demo1.
type Demo1Spec struct {
	V1 string
}

// Demo1Status is the status of a Demo1.
type Demo1Status struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Demo1 is a demo type with a spec and a status.
type Demo1 struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   Demo1Spec
	Status Demo1Status
}
