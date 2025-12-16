#!/usr/bin/env bash

set -euo pipefail

VERSION=${1-`git describe --tags --always`}
OS=$(go env GOOS)
ARCH=$(go env GOARCH)

echo "building version: ${VERSION} for ${OS}/${ARCH}"

go build -ldflags="-s -w -X main._version=${VERSION} -X main._os=${OS} -X main._arch=${ARCH}" -o bin/memory-calculator
