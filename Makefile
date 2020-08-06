SHELL := /bin/bash
.DEFAULT_GOAL := all

# Image URL to use all building/pushing image targets
IMG ?= pizzaorder-controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

BIN := $(CURDIR)/.bin
GOBIN := $(BIN)
PATH := $(abspath $(BIN)):$(PATH)
UNAME_OS := $(shell uname -s)

$(BIN):
	@mkdir -p $(BIN)

CONTROLLER_GEN := $(BIN)/controller-gen
$(CONTROLLER_GEN): | $(BIN)
	cd tools; \
		go build -o $(CONTROLLER_GEN) sigs.k8s.io/controller-tools/cmd/controller-gen

KIND := $(BIN)/kind
KIND_VERSION := v0.8.1
$(KIND): | $(BIN)
	@curl -sSLo $(KIND) "https://kind.sigs.k8s.io/dl/$(KIND_VERSION)/kind-$(UNAME_OS)-amd64"
	@chmod +x $(KIND)

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

TEST_CLUSTER_NAME ?= pizzaorder-controller-test
KUBECONFIG ?= $(CURDIR)/.kube/config
export KUBECONFIG
test-cluster: $(KIND)
	kind create cluster --name $(TEST_CLUSTER_NAME)
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.16.0/cert-manager.yaml

clean:
	kind delete cluster --name $(TEST_CLUSTER_NAME)

load-image: docker-build $(KIND)
	kind load docker-image ${IMG} --name $(TEST_CLUSTER_NAME)
