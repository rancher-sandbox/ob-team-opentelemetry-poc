package setup

import (
	"context"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/generated/controllers/otel.stack.io"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/generated/controllers/otel.stack.io/v1alpha1"
	"github.com/rancher/wrangler/v3/pkg/apply"
	"github.com/rancher/wrangler/v3/pkg/generated/controllers/apps"
	appcontroller "github.com/rancher/wrangler/v3/pkg/generated/controllers/apps/v1"
	"github.com/rancher/wrangler/v3/pkg/generated/controllers/core"
	corecontroller "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/v3/pkg/start"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

type AppContext struct {
	context.Context
	threadiness int

	Core         corecontroller.Interface
	Apps         appcontroller.Interface
	Stack        v1alpha1.OpenTelemetryStackController
	ClusterStack v1alpha1.OpenTelemetryClusterStackController
	Apply        apply.Apply
	starters     []start.Starter
}

func (a *AppContext) Start(ctx context.Context) error {
	return start.All(ctx, a.threadiness, a.starters...)
}

func NewAppContext(ctx context.Context, threadiness int, cfg clientcmd.ClientConfig) (*AppContext, error) {
	client, err := cfg.ClientConfig()
	if err != nil {
		return nil, err
	}

	discovery, err := discovery.NewDiscoveryClientForConfig(client)
	if err != nil {
		return nil, err
	}

	core, err := core.NewFactoryFromConfig(client)
	if err != nil {
		return nil, err
	}

	apps, err := apps.NewFactoryFromConfig(client)
	if err != nil {
		return nil, err
	}
	stack, err := otel.NewFactoryFromConfig(client)
	if err != nil {
		return nil, err
	}

	applier := apply.New(discovery, apply.NewClientFactory(client))
	appsv1 := apps.Apps().V1()
	corev1 := core.Core().V1()
	return &AppContext{
		Core:         corev1,
		Apps:         appsv1,
		Stack:        stack.Otel().V1alpha1().OpenTelemetryStack(),
		ClusterStack: stack.Otel().V1alpha1().OpenTelemetryClusterStack(),
		Context:      ctx,
		threadiness:  threadiness,
		Apply:        applier,
		starters: []start.Starter{
			core,
			apps,
			stack,
		},
	}, nil
}
