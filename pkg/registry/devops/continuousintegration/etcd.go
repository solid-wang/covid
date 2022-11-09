package continuousintegration

import (
	"context"
	"fmt"
	"github.com/solid-wang/covid/pkg/apis/devops"
	"github.com/solid-wang/covid/pkg/printers"
	printerstorage "github.com/solid-wang/covid/pkg/printers/storage"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"k8s.io/apiserver/pkg/util/dryrun"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// Storage includes dummy storage for Deployments and for Scale subresource.
type Storage struct {
	ContinuousIntegration *REST
	Status                *StatusREST
}

// REST implements a RESTStorage for events.
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (Storage, error) {
	strategy, statusStrategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &devops.ContinuousIntegration{} },
		NewListFunc:              func() runtime.Object { return &devops.ContinuousIntegrationList{} },
		PredicateFunc:            Match,
		DefaultQualifiedResource: devops.Resource("continuousintegrations"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return Storage{}, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = statusStrategy

	return Storage{
		ContinuousIntegration: &REST{store},
		Status:                &StatusREST{&statusStore},
	}, nil
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"ci"}
}

// Delete enforces life-cycle rules for ContinuousIntegration termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	Obj, err := r.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	continuousintegration := Obj.(*devops.ContinuousIntegration)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &continuousintegration.UID
	} else if *options.Preconditions.UID != continuousintegration.UID {
		err = apierrors.NewConflict(
			devops.Resource("continuousintegrations"),
			name,
			fmt.Errorf("Precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, continuousintegration.UID),
		)
		return nil, false, err
	}
	if options.Preconditions.ResourceVersion != nil && *options.Preconditions.ResourceVersion != continuousintegration.ResourceVersion {
		err = apierrors.NewConflict(
			devops.Resource("continuousintegrations"),
			name,
			fmt.Errorf("Precondition failed: ResourceVersion in precondition: %v, ResourceVersion in object meta: %v", *options.Preconditions.ResourceVersion, continuousintegration.ResourceVersion),
		)
		return nil, false, err
	}

	// upon first request to delete, we switch the phase to start namespace termination
	// TODO: enhance graceful deletion's calls to DeleteStrategy to allow phase change and finalizer patterns
	if continuousintegration.DeletionTimestamp.IsZero() {

		//key, err := r.store.KeyFunc(ctx, name)
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID, ResourceVersion: options.Preconditions.ResourceVersion}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingContinuousIntegration, ok := existing.(*devops.ContinuousIntegration)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *devops.ContinuousIntegration, got %v", existing)
				}
				if err := deleteValidation(ctx, existingContinuousIntegration); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingContinuousIntegration.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingContinuousIntegration.DeletionTimestamp = &now
				}
				// Set the ci phase to terminating, if needed
				if existingContinuousIntegration.Status.Phase != devops.ContinuousIntegrationTerminating {
					existingContinuousIntegration.Status.Phase = devops.ContinuousIntegrationTerminating
				}

				// the current finalizers which are on ci
				currentFinalizers := map[string]bool{}
				for _, f := range existingContinuousIntegration.Finalizers {
					currentFinalizers[f] = true
				}
				// the finalizers we should ensure on ci
				shouldHaveFinalizers := map[string]bool{
					devops.FinalizerContinuousIntegration: true,
				}
				// determine whether there are changes
				changeNeeded := false
				for finalizer, shouldHave := range shouldHaveFinalizers {
					changeNeeded = currentFinalizers[finalizer] != shouldHave || changeNeeded
					if shouldHave {
						currentFinalizers[finalizer] = true
					} else {
						delete(currentFinalizers, finalizer)
					}
				}
				// make the changes if needed
				if changeNeeded {
					newFinalizers := []string{}
					for f := range currentFinalizers {
						newFinalizers = append(newFinalizers, f)
					}
					existingContinuousIntegration.Finalizers = newFinalizers
				}
				return existingContinuousIntegration, nil
			}),
			dryrun.IsDryRun(options.DryRun),
			nil,
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, devops.Resource("continuousintegrations"), name)
			err = storageerr.InterpretUpdateError(err, devops.Resource("continuousintegrations"), name)
			if _, ok := err.(*apierrors.StatusError); !ok {
				err = apierrors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(continuousintegration.Finalizers) != 0 {
		return continuousintegration, false, nil
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the REST endpoint for changing the status of a deployment
type StatusREST struct {
	store *genericregistry.Store
}

func (r *StatusREST) Destroy() {
	r.store.Destroy()
}

// New returns empty Deployment object.
func (r *StatusREST) New() runtime.Object {
	return &devops.ContinuousIntegration{}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// GetResetFields implements rest.ResetFieldsStrategy
func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.store.GetResetFields()
}

func (r *StatusREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.store.ConvertToTable(ctx, object, tableOptions)
}
