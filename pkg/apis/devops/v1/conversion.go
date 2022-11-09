package v1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	if err := AddFieldLabelConversionsForContinuousIntegration(scheme); err != nil {
		return err
	}
	return nil
}

func AddFieldLabelConversionsForContinuousIntegration(scheme *runtime.Scheme) error {
	return scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("ContinuousIntegration"),
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.namespace",
				"metadata.name":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
}
