#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

# shellcheck disable=SC2046
pkgs=$(go list $(glide novendor))

echo "" > coverage.txt
for p in $pkgs; do
    go test -covermode=atomic -coverprofile=profile.out -tags integration "$p"
    if [ -s profile.out ]; then
        cat profile.out >> coverage.txt
    fi
    rm -f profile.out
done
