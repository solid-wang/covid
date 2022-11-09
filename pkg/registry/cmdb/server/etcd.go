package kubernetes

import (
	"context"
	"fmt"
	"github.com/solid-wang/covid/pkg/apis/cmdb"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"k8s.io/apiserver/pkg/util/dryrun"
)

// REST implements a RESTStorage for API services against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &cmdb.Server{} },
		NewListFunc:              func() runtime.Object { return &cmdb.ServerList{} },
		PredicateFunc:            Match,
		DefaultQualifiedResource: cmdb.Resource("servers"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(cmdb.Resource("servers")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}

// Delete enforces life-cycle rules for ContinuousIntegration termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	Obj, err := r.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	server := Obj.(*cmdb.Server)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &server.UID
	} else if *options.Preconditions.UID != server.UID {
		err = apierrors.NewConflict(
			cmdb.Resource("servers"),
			name,
			fmt.Errorf("Precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, server.UID),
		)
		return nil, false, err
	}
	if options.Preconditions.ResourceVersion != nil && *options.Preconditions.ResourceVersion != server.ResourceVersion {
		err = apierrors.NewConflict(
			cmdb.Resource("servers"),
			name,
			fmt.Errorf("Precondition failed: ResourceVersion in precondition: %v, ResourceVersion in object meta: %v", *options.Preconditions.ResourceVersion, server.ResourceVersion),
		)
		return nil, false, err
	}

	// upon first request to delete, we switch the phase to start namespace termination
	// TODO: enhance graceful deletion's calls to DeleteStrategy to allow phase change and finalizer patterns
	if server.DeletionTimestamp.IsZero() {

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
				existingServer, ok := existing.(*cmdb.Server)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *cmdb.Server, got %v", existing)
				}
				if err := deleteValidation(ctx, existingServer); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingServer.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingServer.DeletionTimestamp = &now
				}

				// the current finalizers which are on ci
				currentFinalizers := map[string]bool{}
				for _, f := range existingServer.Finalizers {
					currentFinalizers[f] = true
				}
				// the finalizers we should ensure on ci
				shouldHaveFinalizers := map[string]bool{
					cmdb.FinalizerServer: true,
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
					existingServer.Finalizers = newFinalizers
				}
				return existingServer, nil
			}),
			dryrun.IsDryRun(options.DryRun),
			nil,
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, cmdb.Resource("servers"), name)
			err = storageerr.InterpretUpdateError(err, cmdb.Resource("servers"), name)
			if _, ok := err.(*apierrors.StatusError); !ok {
				err = apierrors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(server.Finalizers) != 0 {
		return server, false, nil
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}
