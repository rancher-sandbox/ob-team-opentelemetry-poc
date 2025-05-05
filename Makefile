.PHONY: collector
.PHONY: images

VERSION?=dev
TARGET?=gateway
TYPE?=minimal

collector:
	ocb --config=./collector/$(TARGET)-$(TYPE).yaml --skip-compilation

build:
	go build -o ./bin/operator ./cmd/operator/main.go

images:
	@echo "Building gateway"
	docker build ./collector/images/gateway/minimal/ -t rancher-sandbox/gateway:$(VERSION) -f ./collector/images/gateway/minimal/Dockerfile
	@echo "Building logging collector"
	docker build ./collector/images/logging/minimal/ -t rancher-sandbox/node:$(VERSION) -f ./collector/images/logging/minimal/Dockerfile