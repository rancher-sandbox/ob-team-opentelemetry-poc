// Code generated by codegen. DO NOT EDIT.

// +k8s:deepcopy-gen=package
// +groupName=otel.stack.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenTelemetryStackList is a list of OpenTelemetryStack resources
type OpenTelemetryStackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []OpenTelemetryStack `json:"items"`
}

func NewOpenTelemetryStack(namespace, name string, obj OpenTelemetryStack) *OpenTelemetryStack {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("OpenTelemetryStack").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenTelemetryClusterStackList is a list of OpenTelemetryClusterStack resources
type OpenTelemetryClusterStackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []OpenTelemetryClusterStack `json:"items"`
}

func NewOpenTelemetryClusterStack(namespace, name string, obj OpenTelemetryClusterStack) *OpenTelemetryClusterStack {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("OpenTelemetryClusterStack").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}
