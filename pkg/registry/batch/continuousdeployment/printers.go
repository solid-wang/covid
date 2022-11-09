package continuousdeployment

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
	continuousDeploymentColumnDefinitions := []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Status", Type: "string", Description: "status."},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
	}
	h.TableHandler(continuousDeploymentColumnDefinitions, printContinuousDeploymentList)
	h.TableHandler(continuousDeploymentColumnDefinitions, printContinuousDeployment)
}

func printContinuousDeploymentList(continuousDeploymentList *batch.ContinuousDeploymentList, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0, len(continuousDeploymentList.Items))
	for i := range continuousDeploymentList.Items {
		r, err := printContinuousDeployment(&continuousDeploymentList.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func printContinuousDeployment(continuousDeployment *batch.ContinuousDeployment, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: continuousDeployment},
	}

	status := continuousDeployment.Status.Phase

	row.Cells = append(row.Cells, continuousDeployment.Name, status, translateTimestampSince(continuousDeployment.CreationTimestamp))

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
