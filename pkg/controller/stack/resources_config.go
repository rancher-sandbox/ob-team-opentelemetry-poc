package controller

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/internal/encoder"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	"github.com/samber/lo"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
)

func (g *ClusterStackGenerator) constructOpenTelemetryConfig() error {
	if err := g.constructExtensions(); err != nil {
		return fmt.Errorf("failed to construct extensions : %w", err)
	}
	if err := g.constructReceivers(); err != nil {
		return fmt.Errorf("failed to construct receivers : %w", err)
	}

	if err := g.constructExporters(); err != nil {
		return fmt.Errorf("failed to construct exporters : %w", err)
	}

	if err := g.constructPipelines(); err != nil {
		return fmt.Errorf("failed to construct pipelines : %w", err)
	}
	return nil
}

func (g *ClusterStackGenerator) constructPipelines() error {
	extensionMap, ok := g.managedConfig[configExtensions].(map[string]any)
	if !ok {
		return fmt.Errorf("internal error, couldn't extract extensions from managed config")
	}

	receiverMap, ok := g.managedConfig[configReceivers].(map[string]any)
	if !ok {
		return fmt.Errorf("internal error, couldn't extract receivers from managed config")
	}

	exporterMap, ok := g.managedConfig[configExporters].(map[string]any)
	if !ok {
		return fmt.Errorf("internal error, couldn't extract expoters from managed config")
	}

	registeredExtensions := lo.Keys(extensionMap)
	registeredReceivers := lo.Keys(receiverMap)
	registeredExporters := lo.Keys(exporterMap)
	g.managedConfig[configServices] = map[string]any{
		configExtensions: registeredExtensions,
		configPipelines: map[string]any{
			"logs": map[string]any{
				"receivers": registeredReceivers,
				"exporters": registeredExporters,
			},
		},
	}
	return nil
}

func (g *ClusterStackGenerator) constructReceivers() error {
	g.encoder()
	receiverMap := map[string]any{}
	filelogId := filelogreceiver.NewFactory().Type().String()
	journaldId := journaldreceiver.NewFactory().Type().String()
	if g.clusterstack.Spec.CollectPodLogs {
		podKey := fmt.Sprintf("%s/k8slog", filelogId)
		podConfig, err := g.constructPodLogReceiver()
		if err != nil {
			return err
		}
		receiverMap[podKey] = podConfig
	}
	if g.clusterstack.Spec.CollectAuditLogs {
		auditKey := fmt.Sprintf("%s/k8saudit", filelogId)
		auditConfig, err := g.constructAuditLogReceiver()
		if err != nil {
			return err
		}
		receiverMap[auditKey] = auditConfig
	}
	if g.clusterstack.Spec.CollectRKE2 {
		journaldKey := fmt.Sprintf("%s/rke2", journaldId)
		journaldConfig, err := g.constructRke2JournaldReceiver()
		if err != nil {
			return err
		}

		agentKey := fmt.Sprintf("%s/rke2agent", filelogId)
		agentConfig, err := g.constructRke2AgentLogReceiver()
		if err != nil {
			return err
		}
		receiverMap[journaldKey] = journaldConfig
		receiverMap[agentKey] = agentConfig
	}

	if g.clusterstack.Spec.CollectK3s {
		journaldKey := fmt.Sprintf("%s/k3s", journaldId)
		journaldConfig, err := g.constructK3sJournaldReceiver()
		if err != nil {
			return err
		}
		receiverMap[journaldKey] = journaldConfig
	}
	if len(receiverMap) == 0 {
		return fmt.Errorf("no receivers configured")
	}
	g.managedConfig[configReceivers] = receiverMap
	// filelog receivers
	return nil
}

func (g *ClusterStackGenerator) encoder() *encoder.Encoder {
	return encoder.New(&encoder.EncoderConfig{
		EncodeHook: mapstructure.ComposeDecodeHookFunc(
			encoder.YamlMarshalerHookFunc(),
			encoder.TextMarshalerHookFunc(),
		),
	})
}

func (g *ClusterStackGenerator) constructExtensions() error {
	enc := g.encoder()
	state := filestorage.Config{
		Directory: g.stateStoragePath(),
		Timeout:   time.Second,
		FSync:     false,
		//TODO figure out how defaults are constructed in upstream
		Compaction:      &filestorage.CompactionConfig{},
		CreateDirectory: false,
	}

	if err := state.Validate(); err != nil {
		return err
	}

	stateConstructed, err := enc.Encode(state)
	if err != nil {
		return err
	}
	pprof := pprofextension.Config{
		TCPAddr: confignet.NewDefaultTCPAddrConfig(),
	}

	if err := pprof.Validate(); err != nil {
		return err
	}

	pprofConstructed, err := enc.Encode(pprof)
	if err != nil {
		return err
	}

	g.managedConfig[configExtensions] = map[string]any{
		pprofextension.NewFactory().Type().String(): pprofConstructed,
		filestorage.NewFactory().Type().String():    stateConstructed,
	}
	return nil
}

func (g *ClusterStackGenerator) stateStoragePath() string {
	return "/var/otel/storage"
}

func (g *ClusterStackGenerator) storageId() *component.ID {
	return lo.ToPtr(component.NewID(filestorage.NewFactory().Type()))
}

func (g *ClusterStackGenerator) constructPodLogReceiver() (any, error) {
	return yamlMarshalString(PodConfig)

}

// FIXME: probably want to replace this with our own receiver
func (g *ClusterStackGenerator) constructAuditLogReceiver() (any, error) {
	var b bytes.Buffer
	auditPath := g.clusterstack.Spec.AuditLogPath
	if auditPath == "" {
		return "", fmt.Errorf("empty audit log path")
	}
	if err := AuditConfigTemplate.Execute(&b, AuditConfigParams{
		AuditLogPath: auditPath,
		StorageId:    g.storageId().String(),
	}); err != nil {
		return "", err
	}
	return yamlMarshalString(b.String())
}

// FIXME: perhaps we also want to wrap this in our own receiver
func (g *ClusterStackGenerator) constructRke2JournaldReceiver() (any, error) {
	rke2Path := g.clusterstack.Spec.RKE2LogPath
	if rke2Path == "" {
		return nil, fmt.Errorf("empty rke2 journald path")
	}
	var b bytes.Buffer
	if err := RK2JournalConfigTemplate.Execute(&b, RKE2JournalParams{
		JournalPath: rke2Path,
		StorageId:   g.storageId().String(),
	}); err != nil {
		return nil, err
	}

	return yamlMarshalString(b.String())
}

// FIXME: perhaps we also want to wrap this in our own receiver
func (g *ClusterStackGenerator) constructK3sJournaldReceiver() (any, error) {
	k3sPath := g.clusterstack.Spec.K3sLogPath
	if k3sPath == "" {
		return nil, fmt.Errorf("empty k3s journald path")
	}
	var b bytes.Buffer

	if err := K3sJournalConfigTemplate.Execute(&b, K3sJournalParams{
		JournalPath: k3sPath,
		StorageId:   g.storageId().String(),
	}); err != nil {
		return nil, err
	}

	return yamlMarshalString(b.String())

}

// FIXME: perhaps wrap this / find ways to validate this better
func (g *ClusterStackGenerator) constructRke2AgentLogReceiver() (any, error) {
	return yamlMarshalString(Rk2LogConfig)
}

func (g *ClusterStackGenerator) constructExporters() error {
	gatewayRefs := g.clusterstack.Spec.GatewayRefs
	enc := encoder.New(&encoder.EncoderConfig{
		EncodeHook: mapstructure.ComposeDecodeHookFunc(
			encoder.YamlMarshalerHookFunc(),
			encoder.TextMarshalerHookFunc(),
		),
	})
	if len(gatewayRefs) == 0 {
		// build debug exporter and return
		globalExporter := debugexporter.Config{
			Verbosity:          configtelemetry.LevelBasic,
			UseInternalLogger:  false,
			SamplingInitial:    2,
			SamplingThereafter: 60,
		}
		if err := globalExporter.Validate(); err != nil {
			return err
		}
		constructed, err := enc.Encode(globalExporter)
		if err != nil {
			return err
		}
		g.managedConfig[configExporters] = map[string]any{
			fmt.Sprintf("%s/global", debugexporter.NewFactory().Type().String()): constructed,
		}
		return nil
	}

	// check duplicate gateway refs
	tmap := lo.Associate(gatewayRefs, func(ref v1alpha1.GatewayRef) (string, struct{}) {
		return ref.Namespace + "-" + ref.Name, struct{}{}
	})

	if len(tmap) != len(gatewayRefs) {
		return fmt.Errorf("duplicate gateway ref")
	}

	gatewayConfigs := map[string]any{}

	// TODO : note we probably want to extend cluster stack CRD to encapsulate some otlp options & processors
	for _, gwRef := range gatewayRefs {
		targetSvc := gatewaySvcMeta(gwRef.Name, gwRef.Namespace)

		dns := fmt.Sprintf(serviceFmt, targetSvc.Name, targetSvc.Namespace)

		config := otlpexporter.Config{
			ClientConfig: configgrpc.ClientConfig{
				Endpoint: fmt.Sprintf("%s:%d", dns, 4317),
				TLSSetting: configtls.ClientConfig{
					Insecure:           true,
					InsecureSkipVerify: true,
				},
			},
			TimeoutConfig: exporterhelper.NewDefaultTimeoutConfig(),
			QueueConfig:   exporterhelper.NewDefaultQueueConfig(),
			BatcherConfig: exporterhelper.NewDefaultBatcherConfig(),
			RetryConfig:   configretry.NewDefaultBackOffConfig(),
		}

		if err := config.Validate(); err != nil {
			return err
		}

		constructed, err := enc.Encode(config)
		if err != nil {
			return err
		}

		exporterKey := fmt.Sprintf("%s/gateway-%s-%s", otlpexporter.NewFactory().Type().String(), gwRef.Namespace, gwRef.Name)
		gatewayConfigs[exporterKey] = constructed
	}
	g.managedConfig[configExporters] = gatewayConfigs
	return nil
}
