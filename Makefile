CMD_DIR    := ./cmd
DIST_DIR   := ./dist
DOCKER_DIR := ./docker
DEPLOY_DIR := ./deploy

EXTERNAL_SERVICE_PKG := $(CMD_DIR)/externalservice
INTERNAL_SERVICE_PKG := $(CMD_DIR)/internalservice

KUBECONFIG := $(DEPLOY_DIR)/kubeconfig.yaml

ifneq ("$(wildcard $(KUBECONFIG))","")
    CONTROL_PLANE_IP := $(shell cd deploy && terraform show | grep -E -o "([0-9]{1,3}[\.]){3}[0-9]{1,3}" | head -n 1)
endif

DOCKER_IMG_PREFIX := charliekenney23/lke-in-cluster-hostname-test

GO_BUILD_ARGS ?= GOOS=linux GOARCH=amd64 CGO_ENABLED=0

push-docker-images:
	docker push $(DOCKER_IMG_PREFIX)-external
	docker push $(DOCKER_IMG_PREFIX)-internal

deploy-infra:
	cd deploy && terraform init && terraform apply --auto-approve

deploy-k8s:
	kubectl --kubeconfig $(KUBECONFIG) apply -f k8s/

deploy: deploy-infra deploy-k8s

build-binaries:
	$(GO_BUILD_ARGS) go build -o $(DIST_DIR)/externalservice $(EXTERNAL_SERVICE_PKG)/main.go
	$(GO_BUILD_ARGS) go build -o $(DIST_DIR)/internalservice $(INTERNAL_SERVICE_PKG)/main.go

build-images:
	docker build -f $(DOCKER_DIR)/Dockerfile.external -t $(DOCKER_IMG_PREFIX)-external .
	docker build -f $(DOCKER_DIR)/Dockerfile.internal -t $(DOCKER_IMG_PREFIX)-internal .

build: build-binaries build-images push-docker-images

control-plane-ssh:
	ssh root@$(CONTROL_PLANE_IP)
