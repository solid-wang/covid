package continuousintegration

import (
	"context"
	"fmt"
	"github.com/solid-wang/covid/pkg/apis/devops"
	devopsvalidation "github.com/solid-wang/covid/pkg/apis/devops/validation"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// NewStrategy creates and returns a Strategy instance
func NewStrategy(typer runtime.ObjectTyper) (Strategy, StatusStrategy) {
	s := Strategy{typer, names.SimpleNameGenerator}
	ss := StatusStrategy{s}
	return s, ss
}

// GetAttrs returns labels.Set, fields.Set, the presence of Initializers if any
// and error in case the given runtime.Object is not a example
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*devops.ContinuousIntegration)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Example")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// Match is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func Match(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *devops.ContinuousIntegration) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
	specificFieldsSet := fields.Set{
		"metadata.namespace": obj.Namespace,
		"metadata.name":      obj.Name,
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

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (Strategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	f := map[fieldpath.APIVersion]*fieldpath.Set{
		"devops/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}
	return f
}

func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	ci := obj.(*devops.ContinuousIntegration)
	ci.Status = devops.ContinuousIntegrationStatus{}
	ci.Generation = 1
}

func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newCI := obj.(*devops.ContinuousIntegration)
	oldCI := old.(*devops.ContinuousIntegration)
	newCI.Status = oldCI.Status

	// Spec updates bump the generation so that we can distinguish between
	// scaling events and template changes, annotation updates bump the generation
	// because annotations are copied from deployments to their replica sets.
	if !apiequality.Semantic.DeepEqual(newCI.Spec, oldCI.Spec) ||
		!apiequality.Semantic.DeepEqual(newCI.Annotations, oldCI.Annotations) {
		newCI.Generation = oldCI.Generation + 1
	}
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
	updateCI := obj.(*devops.ContinuousIntegration)
	oldCI := old.(*devops.ContinuousIntegration)
	return devopsvalidation.ContinuousIntegrationUpdate(updateCI, oldCI)
}

type StatusStrategy struct {
	Strategy
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (StatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"devops/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
			fieldpath.MakePathOrDie("metadata", "labels"),
		),
	}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status
func (StatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newDeployment := obj.(*devops.ContinuousIntegration)
	oldDeployment := old.(*devops.ContinuousIntegration)
	newDeployment.Spec = oldDeployment.Spec
	newDeployment.Labels = oldDeployment.Labels
}

func (StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (s StatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
