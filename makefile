SHELL := /bin/bash

# ==============================================================================
# Testing running system

# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3030/v1/users/token
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3030/v1/users/1/2
#
# For testing load on the service.
# go install github.com/rakyll/hey@latest
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3030/v1/users/1/2
#
# Access metrics directly (5050) or through the sidecar (3031)
# go install github.com/divan/expvarmon@latest
# expvarmon -ports=":5050" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# expvarmon -ports=":3031" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#
# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./admin genkey
#
# Testing coverage.
# go test -coverprofile p.out
# go tool cover -html p.out
#
# Test debug endpoints.
# curl http://localhost:5050/debug/liveness
# curl http://localhost:5050/debug/readiness
#
# Running pgcli client for database.
# brew install pgcli
# pgcli postgresql://postgres:postgres@localhost
#
# Launch zipkin.
# http://localhost:9411/zipkin/
#
# Database access
# dblab --host 0.0.0.0 --user postgres --db postgres --pass postgres --ssl disable --port 5432 --driver postgres

# ==============================================================================
# Local commands

build:
	go build -o bin/moneyflow-api ./app/services/moneyflow-api

run:
	go run app/services/moneyflow-api/main.go | go run app/tooling/logfmt/main.go

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

# ==============================================================================
# Install dependencies

dev.setup.mac:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
APP_NAME := moneyflow-api-amd64
METRICS_NAME := metrics-amd64
VERSION := 1.0

all: moneyflow metrics

moneyflow:
	docker build \
		-f zarf/docker/dockerfile.moneyflow-api \
		-t $(APP_NAME):$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t $(METRICS_NAME):$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := gloompi-moneyflow-cluster

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
	cd zarf/k8s/kind/moneyflow-pod; kustomize edit set image metrics-image=$(METRICS_NAME):$(VERSION)
	kind load docker-image $(APP_NAME):$(VERSION) --name $(KIND_CLUSTER)
	kind load docker-image $(METRICS_NAME):$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/zipkin-pod | kubectl apply -f -
	kubectl wait --namespace=zipkin-system --timeout=120s --for=condition=Available deployment/zipkin-pod
	kustomize build zarf/k8s/kind/moneyflow-pod | kubectl apply -f -

kind-services-delete:
	kustomize build zarf/k8s/kind/moneyflow-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/zipkin-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/database-pod | kubectl delete -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-moneyflow:
	kubectl get pods -o wide --watch

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-status-zipkin:
	kubectl get pods -o wide --watch --namespace=zipkin-system

kind-logs:
	kubectl logs -l app=moneyflow --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-logs-moneyflow:
	kubectl logs -l app=moneyflow --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=MONEYFLOW-API

kind-logs-metrics:
	kubectl logs -l app=moneyflow --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=METRICS

kind-logs-db:
	kubectl logs -l app=database --namespace=database-system --all-containers=true -f --tail=100

kind-logs-zipkin:
	kubectl logs -l app=zipkin --namespace=zipkin-system --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment moneyflow-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pods -l app=moneyflow

kind-describe-deployment:
	kubectl describe deployment moneyflow-pod

kind-describe-replicaset:
	kubectl get rs
	kubectl describe rs -l app=moneyflow

kind-events:
	kubectl get ev --sort-by metadata.creationTimestamp

kind-events-warn:
	kubectl get ev --field-selector type=Warning --sort-by metadata.creationTimestamp

kind-context-moneyflow:
	kubectl config set-context --current --namespace=moneyflow-system

kind-shell:
	kubectl exec -it $(shell kubectl get pods | cut -c1-26) --container moneyflow-api -- /bin/sh

kind-database:
	# ./admin --db-disable-tls=1 migrate
	# ./admin --db-disable-tls=1 seed
