.PHONY: collector
.PHONY: images

VERSION?=dev
TARGET?=gateway
TYPE?=minimal
REPO?=rancher-sandbox

collector:
	ocb --config=./collector/$(TARGET)-$(TYPE).yaml --skip-compilation

compile:
	ocb --config=./collector/$(TARGET)-$(TYPE).yaml

build:
	go build -o ./bin/operator ./cmd/operator/main.go

images:
	@echo "Building gateway"
	docker build ./collector/images/gateway/minimal/ -t $(REPO)/gateway:$(VERSION) -f ./collector/images/gateway/minimal/Dockerfile
	@echo "Building logging collector"
	docker build ./collector/images/logging/minimal/ -t $(REPO)/node:$(VERSION) -f ./collector/images/logging/minimal/Dockerfile