#!/usr/bin/env bash
#
# Build and push Docker images to quay.io.
#

cd "$(dirname "$0")" || exit 1

export IMAGE_PREFIX=deisci
docker login -e="$QUAY_EMAIL" -u="$QUAY_USERNAME" -p="$QUAY_PASSWORD" quay.io
DEIS_REGISTRY=quay.io/ make -C .. docker-build docker-push
