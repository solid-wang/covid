package event

import (
	"fmt"
	"github.com/solid-wang/covid/pkg/apis/core"
	"github.com/solid-wang/covid/pkg/apis/core/v1"
	"github.com/solid-wang/covid/pkg/printers"
	"github.com/solid-wang/covid/pkg/printers/duration"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
	"time"
)

// AddHandlers adds print handlers for default Kubernetes types dealing with internal versions.
// reference: https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go
// TODO: handle errors from Handler
func AddHandlers(h printers.PrintHandler) {
	eventColumnDefinitions := []metav1.TableColumnDefinition{
		{Name: "Last Seen", Type: "string", Description: v1.Event{}.SwaggerDoc()["lastTimestamp"]},
		{Name: "Type", Type: "string", Description: v1.Event{}.SwaggerDoc()["type"]},
		{Name: "Reason", Type: "string", Description: v1.Event{}.SwaggerDoc()["reason"]},
		{Name: "Object", Type: "string", Description: v1.Event{}.SwaggerDoc()["involvedObject"]},
		{Name: "Subobject", Type: "string", Priority: 1, Description: v1.Event{}.InvolvedObject.SwaggerDoc()["fieldPath"]},
		{Name: "Source", Type: "string", Priority: 1, Description: v1.Event{}.SwaggerDoc()["source"]},
		{Name: "Message", Type: "string", Description: v1.Event{}.SwaggerDoc()["message"]},
		{Name: "First Seen", Type: "string", Priority: 1, Description: v1.Event{}.SwaggerDoc()["firstTimestamp"]},
		{Name: "Count", Type: "string", Priority: 1, Description: v1.Event{}.SwaggerDoc()["count"]},
		{Name: "Name", Type: "string", Priority: 1, Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
	}

	h.TableHandler(eventColumnDefinitions, printEvent)
	h.TableHandler(eventColumnDefinitions, printEventList)
}

func printEvent(obj *core.Event, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	row := metav1.TableRow{
		Object: runtime.RawExtension{Object: obj},
	}

	firstTimestamp := translateTimestampSince(obj.FirstTimestamp)
	if obj.FirstTimestamp.IsZero() {
		firstTimestamp = translateMicroTimestampSince(obj.EventTime)
	}

	lastTimestamp := translateTimestampSince(obj.LastTimestamp)
	if obj.LastTimestamp.IsZero() {
		lastTimestamp = firstTimestamp
	}

	count := obj.Count
	if obj.Series != nil {
		lastTimestamp = translateMicroTimestampSince(obj.Series.LastObservedTime)
		count = obj.Series.Count
	} else if count == 0 {
		// Singleton events don't have a count set in the new API.
		count = 1
	}

	var target string
	if len(obj.InvolvedObject.Name) > 0 {
		target = fmt.Sprintf("%s/%s", strings.ToLower(obj.InvolvedObject.Kind), obj.InvolvedObject.Name)
	} else {
		target = strings.ToLower(obj.InvolvedObject.Kind)
	}
	if options.Wide {
		row.Cells = append(row.Cells,
			lastTimestamp,
			obj.Type,
			obj.Reason,
			target,
			obj.InvolvedObject.FieldPath,
			formatEventSource(obj.Source, obj.ReportingController, obj.ReportingInstance),
			strings.TrimSpace(obj.Message),
			firstTimestamp,
			int64(count),
			obj.Name,
		)
	} else {
		row.Cells = append(row.Cells,
			lastTimestamp,
			obj.Type,
			obj.Reason,
			target,
			strings.TrimSpace(obj.Message),
		)
	}

	return []metav1.TableRow{row}, nil
}

// Sorts and prints the EventList in a human-friendly format.
func printEventList(list *core.EventList, options printers.GenerateOptions) ([]metav1.TableRow, error) {
	rows := make([]metav1.TableRow, 0, len(list.Items))
	for i := range list.Items {
		r, err := printEvent(&list.Items[i], options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// translateTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestampSince(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

// translateMicroTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func translateMicroTimestampSince(timestamp metav1.MicroTime) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

func formatEventSourceComponentInstance(component, instance string) string {
	if len(instance) == 0 {
		return component
	}
	return component + ", " + instance
}

// formatEventSource formats EventSource as a comma separated string excluding Host when empty.
// It uses reportingController when Source.Component is empty and reportingInstance when Source.Host is empty
func formatEventSource(es core.EventSource, reportingController, reportingInstance string) string {
	return formatEventSourceComponentInstance(
		firstNonEmpty(es.Component, reportingController),
		firstNonEmpty(es.Host, reportingInstance),
	)
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}
