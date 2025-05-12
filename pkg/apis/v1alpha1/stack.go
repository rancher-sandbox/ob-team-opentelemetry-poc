package v1alpha1

import (
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/generic"
	"github.com/rancher/wrangler/v3/pkg/genericcondition"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OpenTelemetryStack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StackSpec   `json:"spec"`
	Status            StackStatus `json:"status"`
}

type StackSpec struct {
	Gateway GatewayStack `json:"gateway,omitempty"`
	Node    NodeStack    `json:"node,omitempty"`
}

type StackStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`
}

type GatewayStack struct {
	Enabled bool                 `json:"enabled"`
	Image   generic.GenericImage `json:"image"`
	// GrpcDebugLogging flag enables the gRPC debug logging from the collector pods
	GrpcDebugLogging bool `json:"grpcDebugLogging"`
}

type NodeStack struct {
	Enabled bool `json:"enabled"`
	// GatewayRef is a reference to another gateway stack, if none is specified assume the one
	// deployed by the current OpenTelemetryStack to be the gateway ref
	GatewayRef GatewayRef           `json:"gatewayRef"`
	Image      generic.GenericImage `json:"image"`
	// GrpcDebugLogging flag enables the gRPC debug logging from the collector pods
	GrpcDebugLogging bool `json:"grpcDebugLogging"`
}

type GatewayRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}
