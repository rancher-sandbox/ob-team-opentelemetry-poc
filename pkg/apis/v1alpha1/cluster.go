package v1alpha1

import (
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/generic"
	"github.com/rancher/wrangler/v3/pkg/genericcondition"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OpenTelemetryClusterStack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterStackSpec   `json:"spec"`
	Status            ClusterStackStatus `json:"status"`
}

type ClusterStackStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`
}

type ClusterStackSpec struct {
	GatewayRefs []GatewayRef         `json:"gatewayRefs"`
	Image       generic.GenericImage `json:"image"`

	// FIXME: these are just example flags for the POC
	CollectPodLogs   bool   `json:"collectPodLogs"`
	CollectAuditLogs bool   `json:"collectAuditLogs"`
	AuditLogPath     string `json:"auditLogPath"`
	CollectK3s       bool   `json:"collectK3s"`
	CollectRKE2      bool   `json:"collectRKE2"`
	RKE2LogPath      string `json:"rke2LogPath"` // should be /var/log/journald by default
	K3sLogPath       string `json:"k3sLogPath"`  // should be /var/log/journald by default
}
