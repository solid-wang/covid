package validation

import (
	corevalidation "github.com/solid-wang/covid/pkg/apis/core/validation"
	"github.com/solid-wang/covid/pkg/apis/devops"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ContinuousIntegrationUpdate tests if an update to a Deployment is valid.
func ContinuousIntegrationUpdate(update, old *devops.ContinuousIntegration) field.ErrorList {
	allErrs := corevalidation.ValidateObjectMetaUpdate(&update.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	if update.Spec.CIConfigPath != old.Spec.CIConfigPath {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "ciConfigPath"), update.Spec.CIConfigPath, "modifying project id is not allowed"))
	}
	return allErrs
}
