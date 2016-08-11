#!/usr/bin/env bash

docker login -u "${QUAY_USERNAME}" -p "${QUAY_PASSWORD}" https://quay.io
make docker-build docker-push
