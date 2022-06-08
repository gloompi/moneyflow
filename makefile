SHELL := /bin/bash

# ==============================================================================
# Testing running system

# For testing a simple query on the system. Don't forget to `make seed` first.
# curl -il http://localhost:3030/v1/testauth
# curl --user "admin@example.com:gophers" http://localhost:3030/v1/users/token
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -ilH "Authorization: Bearer ${TOKEN}" http://localhost:3030/v1/testauth
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3030/v1/users/1/2
#
# For testing load on the service.
# hey -m GET -c 100 -n 10000 http://localhost:3030/v1/test
#
# Access metrics directly (4040) or through the sidecar (3031)
# expvarmon -ports=":4040" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# expvarmon -ports=":3031" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#
# Used to install expvarmon program for metrics dashboard.
# go install github.com/divan/expvarmon@latest
#
# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./admin genkey
#
# Database access
# dblab --host 0.0.0.0 --user postgres --db postgres --pass postgres --ssl disable --port 5432 --driver postgres
#
# ==============================================================================

build:
	go build -o bin/moneyflow-api ./app/services/moneyflow-api

run:
	go run app/services/moneyflow-api/main.go | go run app/tooling/logfmt/main.go

admin:
	go run app/tooling/admin/main.go

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
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
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/moneyflow-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-moneyflow:
	kubectl get pods -o wide --watch

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

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

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...
