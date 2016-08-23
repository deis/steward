SHORT_NAME := steward

include versioning.mk

REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.17.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD ?= ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

VERSION ?= "dev"
LDFLAGS := "-s -w -X main.version=${VERSION}"
BINARY_DEST_DIR := rootfs/bin

all:
	@echo "Use a Makefile to control top-level building of the project."

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

install-3prs:
	kubectl get thirdpartyresource service-catalog-entry.steward.deis.io || \
	kubectl create -f manifests/service-catalog-entry.yaml

STEWARD_IMAGE ?= quay.io/deisci/steward:devel

install-steward:
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
	sed "s/#cf_broker_name#/${CF_BROKER_NAME}/g" manifests/steward-template.yaml > manifests/${CF_BROKER_NAME}-steward.yaml
	sed -i.bak "s/#cf_broker_scheme#/${CF_BROKER_SCHEME}/g" manifests/${CF_BROKER_NAME}-steward.yaml
	sed -i.bak "s/#cf_broker_hostname#/${CF_BROKER_HOSTNAME}/g" manifests/${CF_BROKER_NAME}-steward.yaml
	sed -i.bak "s/#cf_broker_port#/${CF_BROKER_PORT}/g" manifests/${CF_BROKER_NAME}-steward.yaml
	sed -i.bak "s/#cf_broker_username#/${CF_BROKER_USERNAME}/g" manifests/${CF_BROKER_NAME}-steward.yaml
	sed -i.bak "s/#cf_broker_password#/${CF_BROKER_PASSWORD}/g" manifests/${CF_BROKER_NAME}-steward.yaml
	sed -i.bak "s#\#steward_image\##${STEWARD_IMAGE}#g" manifests/${CF_BROKER_NAME}-steward.yaml
	rm manifests/${CF_BROKER_NAME}-steward.yaml.bak
	kubectl get deployment ${CF_BROKER_NAME}-steward --namespace=steward && \
	kubectl apply -f manifests/${CF_BROKER_NAME}-steward.yaml || \
	kubectl create -f manifests/${CF_BROKER_NAME}-steward.yaml

deploy: install-namespace install-3prs install-steward

dev-deploy: docker-build docker-push
	STEWARD_IMAGE=${IMAGE} $(MAKE) deploy
