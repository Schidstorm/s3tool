#!/bin/bash

set -euo pipefail

cmd=$1
binaryName=debug-build/$cmd

mkdir -p debug-build

CGO_ENABLED=1 GO111MODULE=on go build -tags codes -gcflags="all=-N -l" -o $binaryName ./cmd/$cmd

env $(cat minio.env) ./$binaryName
