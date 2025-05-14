package controller

import (
	"testing"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const testNs = "test"

func TestClusterStackGeneratorExpoter(t *testing.T) {
	clusterstack := &v1alpha1.OpenTelemetryClusterStack{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.ClusterStackSpec{
			GatewayRefs: []v1alpha1.GatewayRef{
				{
					Name:      "todo",
					Namespace: "todo",
				},
			},
		},
	}
	emptyRefs := &v1alpha1.OpenTelemetryClusterStack{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.ClusterStackSpec{
			GatewayRefs: []v1alpha1.GatewayRef{},
		},
	}

	tc := []*v1alpha1.OpenTelemetryClusterStack{
		clusterstack,
		emptyRefs,
	}

	for _, clusterstack := range tc {

		g := NewClusterStackGenerator(testNs, clusterstack)
		err := g.constructExporters()
		require.NoError(t, err)

		require.NotEmpty(t, g.managedConfig)
		assert.NotNil(t, g.managedConfig[configExporters])

		rawBytes, err := yaml.Marshal(g.managedConfig)
		require.NoError(t, err)
		require.NotEmpty(t, rawBytes)

		workingMap := map[string]any{}
		err = yaml.Unmarshal(rawBytes, workingMap)
		require.NoError(t, err)
		val, ok := workingMap[configExporters]
		require.True(t, ok)

		exporterMap, ok := val.(map[string]any)
		require.True(t, ok)
		require.NotEmpty(t, exporterMap)
	}
}

func TestClusterStackGeneratorExtension(t *testing.T) {
	stack := &v1alpha1.OpenTelemetryClusterStack{}

	g := NewClusterStackGenerator(testNs, stack)
	err := g.constructExtensions()
	require.NoError(t, err)
	require.NotEmpty(t, g.managedConfig)

	outBytes, err := yaml.Marshal(g.managedConfig)
	require.NoError(t, err)
	require.NotEmpty(t, outBytes)

}

func TestClusterStackGeneratorReceiver(t *testing.T) {
	all := &v1alpha1.OpenTelemetryClusterStack{
		Spec: v1alpha1.ClusterStackSpec{
			CollectPodLogs:   true,
			CollectAuditLogs: true,
			CollectK3s:       true,
			CollectRKE2:      true,
			RKE2LogPath:      "/var/log/journal",
			K3sLogPath:       "/var/log/journal",
			AuditLogPath:     "/var/log/kube-audit/audit.log",
		},
	}
	tc := []*v1alpha1.OpenTelemetryClusterStack{
		all,
	}

	for _, stack := range tc {
		g := NewClusterStackGenerator(testNs, stack)
		err := g.constructReceivers()
		require.NoError(t, err)

		rawConfig, err := yaml.Marshal(g.managedConfig)
		require.NoError(t, err)
		require.NotEmpty(t, rawConfig)
	}
}

func TestClusterStackGeneratorOpenTelemetryConfig(t *testing.T) {
	all := &v1alpha1.OpenTelemetryClusterStack{
		Spec: v1alpha1.ClusterStackSpec{
			GatewayRefs: []v1alpha1.GatewayRef{
				{
					Name:      "test",
					Namespace: "test",
				},
			},
			CollectPodLogs:   true,
			CollectAuditLogs: true,
			CollectK3s:       true,
			CollectRKE2:      true,
			RKE2LogPath:      "/var/log/journald",
			K3sLogPath:       "/var/log/journald",
			AuditLogPath:     "/var/log/kube-audit/audit.log",
		},
	}
	tc := []*v1alpha1.OpenTelemetryClusterStack{
		all,
	}

	for _, stack := range tc {
		g := NewClusterStackGenerator(testNs, stack)
		err := g.constructOpenTelemetryConfig()
		require.NoError(t, err)

		rawBytes, err := yaml.Marshal(g.managedConfig)
		require.NoError(t, err)
		require.NotEmpty(t, rawBytes)
	}
}

func TestStackGeneratorOpenTelemetryConfig(t *testing.T) {
	stack := &v1alpha1.OpenTelemetryStack{}

	g := NewStackGenerator(stack)

	err := g.constructOpenTelemetryConfig()
	require.NoError(t, err)

}
