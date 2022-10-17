package event

import (
	"github.com/solid-wang/covid/pkg/apis/core"
	"github.com/solid-wang/covid/pkg/printers"
	"github.com/solid-wang/covid/pkg/printers/storage"
	"github.com/solid-wang/covid/pkg/registry"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
)

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:       func() runtime.Object { return &core.Event{} },
		NewListFunc:   func() runtime.Object { return &core.EventList{} },
		PredicateFunc: MatchExample,
		TTLFunc: func(runtime.Object, uint64, bool) (uint64, error) {
			return uint64(3600), nil
		},
		DefaultQualifiedResource: core.Resource("events"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: storage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &registry.REST{store}, nil
}
