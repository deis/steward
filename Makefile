SHORT_NAME := steward

include versioning.mk

REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.19.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
DEV_ENV_CMD_INT := ${DEV_ENV_PREFIX} -it ${DEV_ENV_IMAGE}

VERSION ?= "dev"
BINARY_DEST_DIR := rootfs/bin

all:
	@echo "Use a Makefile to control top-level building of the project."

# Allow developers to step into the containerized development environment
dev:
	${DEV_ENV_CMD_INT} bash

bootstrap:
	${DEV_ENV_CMD} glide install

glideup:
	${DEV_ENV_CMD} glide up

build:
	${DEV_ENV_CMD} sh -c "GOOS=linux GOARCH=amd64 go build -o ${BINARY_DEST_DIR}/steward ."

test:
	${DEV_ENV_CMD} sh -c 'go test $$(glide nv)'

test-cover:
	${DEV_ENV_CMD} test-cover.sh

docker-build: build
	${DEV_ENV_CMD} upx -9 ${BINARY_DEST_DIR}/steward
	docker build --rm -t ${IMAGE} rootfs
	docker tag ${IMAGE} ${MUTABLE_IMAGE}

install-namespace:
	kubectl get ns steward || kubectl create -f manifests/steward-namespace.yaml


STEWARD_IMAGE ?= quay.io/deisci/steward:devel

install-cf-steward:
ifndef CF_BROKER_NAME
	$(error CF_BROKER_NAME is undefined)
endif
ifndef CF_BROKER_SCHEME
	$(error CF_BROKER_SCHEME is undefined)
endif
ifndef CF_BROKER_HOSTNAME
	$(error CF_BROKER_HOSTNAME is undefined)
endif
ifndef CF_BROKER_PORT
	$(error CF_BROKER_PORT is undefined)
endif
ifndef CF_BROKER_USERNAME
	$(error CF_BROKER_USERNAME is undefined)
endif
ifndef CF_BROKER_PASSWORD
	$(error CF_BROKER_PASSWORD is undefined)
endif
	sed "s/#cf_broker_name#/${CF_BROKER_NAME}/g" manifests/steward-template-cf.yaml > manifests/${CF_BROKER_NAME}-cf-steward.yaml
	sed -i.bak "s/#cf_broker_scheme#/${CF_BROKER_SCHEME}/g" manifests/${CF_BROKER_NAME}-cf-steward.yaml
	sed -i.bak "s/#cf_broker_hostname#/${CF_BROKER_HOSTNAME}/g" manifests/${CF_BROKER_NAME}-cf-steward.yaml
	sed -i.bak "s/#cf_broker_port#/${CF_BROKER_PORT}/g" manifests/${CF_BROKER_NAME}-cf-steward.yaml
	sed -i.bak "s/#cf_broker_username#/${CF_BROKER_USERNAME}/g" manifests/${CF_BROKER_NAME}-cf-steward.yaml
	sed -i.bak "s/#cf_broker_password#/${CF_BROKER_PASSWORD}/g" manifests/${CF_BROKER_NAME}-cf-steward.yaml
	sed -i.bak "s#\#steward_image\##${STEWARD_IMAGE}#g" manifests/${CF_BROKER_NAME}-cf-steward.yaml
	rm manifests/${CF_BROKER_NAME}-cf-steward.yaml.bak
	kubectl get deployment ${CF_BROKER_NAME}-steward --namespace=steward && \
	kubectl apply -f manifests/${CF_BROKER_NAME}-cf-steward.yaml || \
	kubectl create -f manifests/${CF_BROKER_NAME}-cf-steward.yaml

install-helm-steward:
ifndef HELM_CHART_NAME
	$(error HELM_CHART_NAME is undefined)
endif
ifndef HELM_TILLER_IP
	$(error HELM_TILLER_IP is undefined)
endif
ifndef HELM_TILLER_PORT
	$(error HELM_TILLER_PORT is undefined)
endif
ifndef HELM_CHART_URL
	$(error HELM_CHART_URL is undefined)
endif
ifndef HELM_CHART_INSTALL_NAMESPACE
	$(error HELM_CHART_INSTALL_NAMESPACE is undefined)
endif
ifndef HELM_PROVISION_BEHAVIOR
	$(error HELM_PROVISION_BEHAVIOR is undefined)
endif
ifndef HELM_SERVICE_ID
	$(error HELM_SERVICE_ID is undefined)
endif
ifndef HELM_SERVICE_NAME
	$(error HELM_SERVICE_NAME is undefined)
endif
ifndef HELM_SERVICE_DESCRIPTION
	$(error HELM_SERVICE_DESCRIPTION is undefined)
endif
ifndef HELM_PLAN_ID
	$(error HELM_PLAN_ID is undefined)
endif
ifndef HELM_PLAN_NAME
	$(error HELM_PLAN_NAME is undefined)
endif
ifndef HELM_PLAN_DESCRIPTION
	$(error HELM_PLAN_DESCRIPTION is undefined)
endif
	sed "s/#helm_name#/${HELM_CHART_NAME}/g" manifests/steward-template-helm.yaml > manifests/${HELM_CHART_NAME}-steward.yaml
	sed "s/#helm_tiller_ip#/${HELM_TILLER_IP}/g" manifests/steward-template-helm.yaml > manifests/${HELM_CHART_NAME}-steward.yaml
	sed -i.bak "s/#helm_tiller_port#/${HELM_TILLER_PORT}/g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s/#helm_chart_url#/${HELM_CHART_URL}/g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s/#helm_chart_install_namespace#/${HELM_CHART_INSTALL_NAMESPACE}/g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s/#helm_provision_behavior#/${HELM_PROVISION_BEHAVIOR}/g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s/#helm_service_id#/${HELM_SERVICE_ID}/g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s#\#helm_service_name\##${HELM_SERVICE_NAME}#g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s#\#helm_service_description\##${HELM_SERVICE_DESCRIPTION}#g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s#\#helm_plan_id\##${HELM_PLAN_ID}#g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s#\#helm_plan_name\##${HELM_PLAN_NAME}#g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	sed -i.bak "s#\#helm_plan_description\##${HELM_PLAN_DESCRIPTION}#g" manifests/${HELM_CHART_NAME}-helm-steward.yaml
	rm manifests/${HELM_CHART_NAME}-helm-steward.yaml.bak
	kubectl get deployment ${HELM_CHART_NAME}-steward --namespace=steward && \
	kubectl apply -f manifests/${HELM_CHART_NAME}-helm-steward.yaml || \
	kubectl create -f manifests/${HELM_CHART_NAME}-helm-steward.yaml

install-cmd-steward:
ifndef CMD_BROKER_NAME
	$(error CMD_BROKER_NAME is undefined)
endif
ifndef CMD_BROKER_IMAGE
	$(error CMD_BROKER_IMAGE is undefined)
endif
ifndef CMD_BROKER_CONFIG_MAP
	$(error CMD_BROKER_CONFIG_MAP is undefined)
endif
ifndef CMD_BROKER_SECRET
	$(error CMD_BROKER_SECRET is undefined)
endif
	sed "s/#cmd_broker_name#/${CMD_BROKER_NAME}/g" manifests/steward-template-cmd.yaml > manifests/${CMD_BROKER_NAME}-cmd-steward.yaml
	sed -i.bak "s#\#cmd_broker_image\##${CMD_BROKER_IMAGE}#g" manifests/${CMD_BROKER_NAME}-cmd-steward.yaml
	sed -i.bak "s#\#cmd_broker_config_map\##${CMD_BROKER_CONFIG_MAP}#g" manifests/${CMD_BROKER_NAME}-cmd-steward.yaml
	sed -i.bak "s#\#cmd_broker_secret\##${CMD_BROKER_SECRET}#g" manifests/${CMD_BROKER_NAME}-cmd-steward.yaml
	sed -i.bak "s#\#steward_image\##${STEWARD_IMAGE}#g" manifests/${CMD_BROKER_NAME}-cmd-steward.yaml
	rm manifests/${CMD_BROKER_NAME}-cmd-steward.yaml.bak
	kubectl get deployment ${CMD_BROKER_NAME}-steward --namespace=steward && \
	kubectl apply -f manifests/${CMD_BROKER_NAME}-cmd-steward.yaml || \
	kubectl create -f manifests/${CMD_BROKER_NAME}-cmd-steward.yaml

deploy-cf: install-namespace install-cf-steward

deploy-helm: install-namespace install-helm-steward

deploy-cmd: install-namespace install-cmd-steward

dev-deploy-cf: docker-build docker-push
	STEWARD_IMAGE=${IMAGE} $(MAKE) deploy-cf

dev-deploy-helm: docker-build docker-push
	STEWARD_IMAGE=${IMAGE} $(MAKE) deploy-helm

dev-deploy-cmd: docker-build docker-push
	STEWARD_IMAGE=${IMAGE} $(MAKE) deploy-cmd
