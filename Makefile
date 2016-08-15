SHORT_NAME := steward

include versioning.mk

REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.14.0
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

build-integration:
	go build -o integration/integration ./integration
