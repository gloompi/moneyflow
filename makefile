SHELL := /bin/bash

# ==============================================================================
# Testing running system
#
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#
# ==============================================================================

run:
	go run app/services/moneyflow-api/main.go | go run app/tooling/logfmt/main.go

# ==============================================================================
# Building containers

APP_NAME := moneyflow-api-amd64
VERSION := 1.0

all: moneyflow

moneyflow:
	docker build \
		-f zarf/docker/dockerfile.moneyflow-api \
		-t $(APP_NAME):$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := gloompi-starter-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.0@sha256:406fd86d48eaf4c04c7280cd1d2ca1d61e7d0d61ddef0125cb097bc7b82ed6a1 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=moneyflow-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/moneyflow-pod; kustomize edit set image moneyflow-api-image=$(APP_NAME):$(VERSION)
	kind load docker-image $(APP_NAME):$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/moneyflow-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-moneyflow:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=moneyflow --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment moneyflow-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pods -l app=moneyflow

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor