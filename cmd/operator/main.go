package main

import (
	"context"
	"os"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	controller "github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/controller/stack"
	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/setup"
	"github.com/rancher/wrangler/v3/pkg/crd"
	"github.com/rancher/wrangler/v3/pkg/signals"
	"k8s.io/client-go/rest"

	"github.com/rancher/wrangler/v3/pkg/kubeconfig"
)

func main() {
	k := os.Getenv("KUBECONFIG")
	ctx := signals.SetupSignalContext()
	cfg := kubeconfig.GetNonInteractiveClientConfig(k)

	client, err := cfg.ClientConfig()
	if err != nil {
		panic(err)
	}

	stackCrd := crd.NamespacedType("OpenTelemetryStack.otel.stack.io/v1alpha1").WithSchemaFromStruct(&v1alpha1.OpenTelemetryStack{})
	clusterStackCrd := crd.NonNamespacedType("OpenTelemetryClusterStack.otel.stack.io/v1alpha1").WithSchemaFromStruct(&v1alpha1.OpenTelemetryClusterStack{})

	if err := createCrd(ctx, client, []crd.CRD{
		stackCrd,
		clusterStackCrd,
	}); err != nil {
		panic(err)
	}

	app, err := setup.NewAppContext(ctx, 2, cfg)
	if err != nil {
		panic(err)
	}

	controller.Register(ctx, app)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	<-ctx.Done()
}

func createCrd(ctx context.Context, cfg *rest.Config, crds []crd.CRD) error {
	factory, err := crd.NewFactoryFromClient(cfg)
	if err != nil {
		return err
	}
	return factory.BatchCreateCRDs(ctx, crds...).BatchWait()
}
