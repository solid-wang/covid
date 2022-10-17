/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ConditionStatus defines conditions of resources
type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition;
// "ConditionFalse" means a resource is not in the condition; "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// NamespaceSpec describes the attributes on a Namespace
type NamespaceSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage
	Finalizers []FinalizerName
}

// FinalizerName is the name identifying a finalizer during namespace lifecycle.
type FinalizerName string

// These are internal finalizer values to Kubernetes, must be qualified name unless defined here or
// in metav1.
const (
	FinalizerKubernetes FinalizerName = "kubernetes"
)

// NamespaceStatus is information about the current status of a Namespace.
type NamespaceStatus struct {
	// Phase is the current lifecycle phase of the namespace.
	// +optional
	Phase NamespacePhase
	// +optional
	Conditions []NamespaceCondition
}

// NamespacePhase defines the phase in which the namespace is
type NamespacePhase string

// These are the valid phases of a namespace.
const (
	// NamespaceActive means the namespace is available for use in the system
	NamespaceActive NamespacePhase = "Active"
	// NamespaceTerminating means the namespace is undergoing graceful termination
	NamespaceTerminating NamespacePhase = "Terminating"
)

// NamespaceConditionType defines constants reporting on status during namespace lifetime and deletion progress
type NamespaceConditionType string

// These are valid conditions of a namespace.
const (
	NamespaceDeletionDiscoveryFailure NamespaceConditionType = "NamespaceDeletionDiscoveryFailure"
	NamespaceDeletionContentFailure   NamespaceConditionType = "NamespaceDeletionContentFailure"
	NamespaceDeletionGVParsingFailure NamespaceConditionType = "NamespaceDeletionGroupVersionParsingFailure"
)

// NamespaceCondition contains details about state of namespace.
type NamespaceCondition struct {
	// Type of namespace controller condition.
	Type NamespaceConditionType
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus
	// +optional
	LastTransitionTime metav1.Time
	// +optional
	Reason string
	// +optional
	Message string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Namespace provides a scope for Names.
// Use of multiple namespaces is optional
type Namespace struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the behavior of the Namespace.
	// +optional
	Spec NamespaceSpec

	// Status describes the current status of a Namespace
	// +optional
	Status NamespaceStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceList is a list of Namespaces.
type NamespaceList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	Items []Namespace
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ObjectReference struct {
	// +optional
	Kind string
	// +optional
	Namespace string
	// +optional
	Name string
	// +optional
	UID types.UID
	// +optional
	APIVersion string
	// +optional
	ResourceVersion string

	// Optional. If referring to a piece of an object instead of an entire object, this string
	// should contain information to identify the sub-object. For example, if the object
	// reference is to a container within a pod, this would take on a value like:
	// "spec.containers{name}" (where "name" refers to the name of the container that triggered
	// the event) or if no container name is specified "spec.containers[2]" (container with
	// index 2 in this pod). This syntax is chosen only to have some well-defined way of
	// referencing a part of an object.
	// TODO: this design is not final and this field is subject to change in the future.
	// +optional
	FieldPath string
}

// EventSource represents the source from which an event is generated
type EventSource struct {
	// Component from which the event is generated.
	// +optional
	Component string
	// Node name on which the event is generated.
	// +optional
	Host string
}

// Valid values for event types (new types could be added in future)
const (
	// Information only and will not cause any problems
	EventTypeNormal string = "Normal"
	// These events are to warn that something might go wrong
	EventTypeWarning string = "Warning"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Event is a report of an event somewhere in the cluster.  Events
// have a limited retention time and triggers and messages may evolve
// with time.  Event consumers should not rely on the timing of an event
// with a given Reason reflecting a consistent underlying trigger, or the
// continued existence of events with that Reason.  Events should be
// treated as informative, best-effort, supplemental data.
// TODO: Decide whether to store these separately or with the object they apply to.
type Event struct {
	metav1.TypeMeta

	metav1.ObjectMeta

	// The object that this event is about. Mapped to events.Event.regarding
	// +optional
	InvolvedObject ObjectReference

	// Optional; this should be a short, machine understandable string that gives the reason
	// for this event being generated. For example, if the event is reporting that a container
	// can't start, the Reason might be "ImageNotFound".
	// TODO: provide exact specification for format.
	// +optional
	Reason string

	// Optional. A human-readable description of the status of this operation.
	// TODO: decide on maximum length. Mapped to events.Event.note
	// +optional
	Message string

	// Optional. The component reporting this event. Should be a short machine understandable string.
	// +optional
	Source EventSource

	// The time at which the event was first recorded. (Time of server receipt is in TypeMeta.)
	// +optional
	FirstTimestamp metav1.Time

	// The time at which the most recent occurrence of this event was recorded.
	// +optional
	LastTimestamp metav1.Time

	// The number of times this event has occurred.
	// +optional
	Count int32

	// Type of this event (Normal, Warning), new types could be added in the future.
	// +optional
	Type string

	// Time when this Event was first observed.
	// +optional
	EventTime metav1.MicroTime

	// Data about the Event series this event represents or nil if it's a singleton Event.
	// +optional
	Series *EventSeries

	// What action was taken/failed regarding to the Regarding object.
	// +optional
	Action string

	// Optional secondary object for more complex actions.
	// +optional
	Related *ObjectReference

	// Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
	// +optional
	ReportingController string

	// ID of the controller instance, e.g. `kubelet-xyzf`.
	// +optional
	ReportingInstance string
}

// EventSeries represents a series ov events
type EventSeries struct {
	// Number of occurrences in this series up to the last heartbeat time
	Count int32
	// Time of the last occurrence observed
	LastObservedTime metav1.MicroTime
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventList is a list of events.
type EventList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	Items []Event
}
