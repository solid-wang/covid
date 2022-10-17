package event

import (
	"context"
	"fmt"
	"github.com/solid-wang/covid/pkg/apis/core"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

// NewStrategy creates and returns a Strategy instance
func NewStrategy(typer runtime.ObjectTyper) Strategy {
	return Strategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, the presence of Initializers if any
// and error in case the given runtime.Object is not a example
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*core.Event)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Example")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchExample is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchExample(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *core.Event) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
	source := obj.Source.Component
	if source == "" {
		source = obj.ReportingController
	}
	specificFieldsSet := fields.Set{
		"involvedObject.kind":            obj.InvolvedObject.Kind,
		"involvedObject.namespace":       obj.InvolvedObject.Namespace,
		"involvedObject.name":            obj.InvolvedObject.Name,
		"involvedObject.uid":             string(obj.InvolvedObject.UID),
		"involvedObject.apiVersion":      obj.InvolvedObject.APIVersion,
		"involvedObject.resourceVersion": obj.InvolvedObject.ResourceVersion,
		"involvedObject.fieldPath":       obj.InvolvedObject.FieldPath,
		"reason":                         obj.Reason,
		"reportingComponent":             obj.ReportingController, // use the core/v1 field name
		"source":                         source,
		"type":                           obj.Type,
	}
	return generic.MergeFieldsSets(specificFieldsSet, objectMetaFieldsSet)
}

type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (s Strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (s Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (Strategy) NamespaceScoped() bool {
	return true
}

func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (Strategy) AllowCreateOnUpdate() bool {
	return false
}

func (Strategy) AllowUnconditionalUpdate() bool {
	return false
}

func (Strategy) Canonicalize(obj runtime.Object) {
}

func (Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
