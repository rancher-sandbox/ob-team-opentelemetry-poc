package main

import (
	"log"
	"os"

	"github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/apis/v1alpha1"
	"github.com/rancher/wrangler/v3/pkg/cleanup"
	controllergen "github.com/rancher/wrangler/v3/pkg/controller-gen"
	"github.com/rancher/wrangler/v3/pkg/controller-gen/args"
)

func main() {
	if err := cleanup.Cleanup("./pkg/apis/"); err != nil {
		log.Fatalf("cleanup failed: %v", err)
	}

	if err := os.RemoveAll("./pkg/generated"); err != nil {
		log.Fatalf("cleanup failed: %v", err)
	}

	controllergen.Run(args.Options{
		OutputPackage: "github.com/rancher-sandbox/ob-team-opentelemetry-poc/pkg/generated",
		Boilerplate:   "./gen/boilerplate.go.txt",
		Groups: map[string]args.Group{
			// TODO : probably should be cattle.otel.io or something
			"otel.stack.io": {
				Types: []interface{}{
					v1alpha1.OpenTelemetryStack{},
					v1alpha1.OpenTelemetryClusterStack{},
				},
				GenerateTypes: true,
			},
		},
	})
}
