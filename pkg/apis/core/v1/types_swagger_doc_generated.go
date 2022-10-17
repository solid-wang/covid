/*
Copyright The Kubernetes Authors.

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

package v1

var map_ObjectReference = map[string]string{
	"":                "ObjectReference contains enough information to let you inspect or modify the referred object.",
	"kind":            "Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
	"namespace":       "Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",
	"name":            "Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names",
	"uid":             "UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids",
	"apiVersion":      "API version of the referent.",
	"resourceVersion": "Specific resourceVersion to which this reference is made, if any. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency",
	"fieldPath":       "If referring to a piece of an object instead of an entire object, this string should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2]. For example, if the object reference is to a container within a pod, this would take on a value like: \"spec.containers{name}\" (where \"name\" refers to the name of the container that triggered the event) or if no container name is specified \"spec.containers[2]\" (container with index 2 in this pod). This syntax is chosen only to have some well-defined way of referencing a part of an object.",
}

func (ObjectReference) SwaggerDoc() map[string]string {
	return map_ObjectReference
}

var map_Event = map[string]string{
	"":                   "Event is a report of an event somewhere in the cluster.  Events have a limited retention time and triggers and messages may evolve with time.  Event consumers should not rely on the timing of an event with a given Reason reflecting a consistent underlying trigger, or the continued existence of events with that Reason.  Events should be treated as informative, best-effort, supplemental data.",
	"metadata":           "Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"involvedObject":     "The object that this event is about.",
	"reason":             "This should be a short, machine understandable string that gives the reason for the transition into the object's current status.",
	"message":            "A human-readable description of the status of this operation.",
	"source":             "The component reporting this event. Should be a short machine understandable string.",
	"firstTimestamp":     "The time at which the event was first recorded. (Time of server receipt is in TypeMeta.)",
	"lastTimestamp":      "The time at which the most recent occurrence of this event was recorded.",
	"count":              "The number of times this event has occurred.",
	"type":               "Type of this event (Normal, Warning), new types could be added in the future",
	"eventTime":          "Time when this Event was first observed.",
	"series":             "Data about the Event series this event represents or nil if it's a singleton Event.",
	"action":             "What action was taken/failed regarding to the Regarding object.",
	"related":            "Optional secondary object for more complex actions.",
	"reportingComponent": "Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.",
	"reportingInstance":  "ID of the controller instance, e.g. `kubelet-xyzf`.",
}

func (Event) SwaggerDoc() map[string]string {
	return map_Event
}

var map_EventList = map[string]string{
	"":         "EventList is a list of events.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
	"items":    "List of events",
}

func (EventList) SwaggerDoc() map[string]string {
	return map_EventList
}

var map_EventSeries = map[string]string{
	"":                 "EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time.",
	"count":            "Number of occurrences in this series up to the last heartbeat time",
	"lastObservedTime": "Time of the last occurrence observed",
}

func (EventSeries) SwaggerDoc() map[string]string {
	return map_EventSeries
}

var map_EventSource = map[string]string{
	"":          "EventSource contains information for an event.",
	"component": "Component from which the event is generated.",
	"host":      "Node name on which the event is generated.",
}

func (EventSource) SwaggerDoc() map[string]string {
	return map_EventSource
}
