package v1alpha1

import (
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/generic"
	"github.com/rancher/wrangler/v3/pkg/genericcondition"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
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
}

type StackStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`
}

type GatewayStack struct {
	Enabled bool                 `json:"enabled"`
	Image   generic.GenericImage `json:"image"`
	// GrpcDebugLogging flag enables the gRPC debug logging from the collector pods
	GrpcDebugLogging bool       `json:"grpcDebugLogging"`
	Exporters        GenericMap `json:"exporters"`
}

// +kubebuilder:pruning:PreserveUnknownFields
// +kubebuilder:validation:EmbeddedResource

// GenericMap is a wrapper on arbitrary JSON / YAML resources
type GenericMap map[string]interface{}

func (in *GenericMap) DeepCopy() *GenericMap {
	if in == nil {
		return nil
	}
	out := new(GenericMap)
	*out = runtime.DeepCopyJSON(*in)
	return out
}

type GatewayRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}
