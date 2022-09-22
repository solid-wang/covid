package install

import (
	"github.com/solid-wang/covid/pkg/apis/group"
	"github.com/solid-wang/covid/pkg/apis/group/v1"
	"github.com/solid-wang/covid/pkg/apis/group/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(group.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(v1beta1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion, v1beta1.SchemeGroupVersion))
}
