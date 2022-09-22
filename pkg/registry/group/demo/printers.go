package demo

import (
	"github.com/solid-wang/covid/pkg/apis/group"
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
	demoColumnDefinitions := []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Version", Type: "string", Description: "version."},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
		{Name: "Other", Type: "string", Priority: 1, Description: "test filed."},
	}
	h.TableHandler(demoColumnDefinitions, printDemoList)
	h.TableHandler(demoColumnDefinitions, printDemo)
}

func printDemoList(demoList *group.DemoList, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0, len(demoList.Items))
	for i := range demoList.Items {
		r, err := printDemo(&demoList.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

func printDemo(demo *group.Demo, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: demo},
	}

	version := demo.Spec.V1

	row.Cells = append(row.Cells, demo.Name, version, translateTimestampSince(demo.CreationTimestamp))
	if options.Wide {

		otherCell := "test"

		row.Cells = append(row.Cells, otherCell)
	}

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
