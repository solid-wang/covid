package continuousintegration

import (
	"github.com/solid-wang/covid/pkg/apis/batch"
	"github.com/solid-wang/covid/pkg/printers"
	"github.com/solid-wang/covid/pkg/printers/duration"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

// AddHandlers adds print handlers for default Kubernetes types dealing with internal versions.
// reference: https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go
// TODO: handle errors from Handler
func AddHandlers(h printers.PrintHandler) {
	continuousIntegrationColumnDefinitions := []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Status", Type: "string", Description: "status."},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
	}
	h.TableHandler(continuousIntegrationColumnDefinitions, printContinuousIntegrationList)
	h.TableHandler(continuousIntegrationColumnDefinitions, printContinuousIntegration)
}

func printContinuousIntegrationList(continuousIntegrationList *batch.ContinuousIntegrationList, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0, len(continuousIntegrationList.Items))
	for i := range continuousIntegrationList.Items {
		r, err := printContinuousIntegration(&continuousIntegrationList.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func printContinuousIntegration(continuousIntegration *batch.ContinuousIntegration, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: continuousIntegration},
	}

	status := continuousIntegration.Status.Phase

	row.Cells = append(row.Cells, continuousIntegration.Name, status, translateTimestampSince(continuousIntegration.CreationTimestamp))

	return []metav1.TableRow{row}, nil
}

// translateTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestampSince(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}
