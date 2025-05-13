.PHONY: collector
.PHONY: images

VERSION?=dev
TARGET?=gateway
TYPE?=minimal
REPO?=rancher-sandbox

ALL_TARGETS := gateway logging

collector-all: 
	@$(foreach item,$(ALL_TARGETS), $(MAKE) collector TARGET=$(item);)

collector:
	ocb --config=./collector/$(TARGET)-$(TYPE).yaml --skip-compilation

compile-all:
	@$(foreach item,$(ALL_TARGETS), $(MAKE) compile TARGET=$(item);)


compile:
	ocb --config=./collector/$(TARGET)-$(TYPE).yaml
	cp ./collector/images/$(TARGET)/$(TYPE)/rancher-$(TARGET)-$(TYPE) ./bin

build:
	go build -o ./bin/operator ./cmd/operator/main.go

images:
	@$(foreach item,$(ALL_TARGETS), $(MAKE) image TARGET=$(item);)

image:
	@echo "Building $(TARGET)"
	docker build ./collector/ -t $(REPO)/otel-$(TARGET):$(VERSION) -f ./package/Dockerfile.$(TARGET)

pushall:
	@$(foreach item,$(ALL_TARGETS), $(MAKE) push TARGET=$(item);)

push:
	@echo "Pushing $(TARGET)"
	docker push $(REPO)/otel-$(TARGET):$(VERSION)

imagepushall:
	@$(foreach item,$(ALL_TARGETS), $(MAKE) imagepush TARGET=$(item);)

imagepush: image push