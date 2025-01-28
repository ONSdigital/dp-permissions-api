#!/bin/bash -eux

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0
npm install -g @redocly/cli

pushd dp-permissions-api
  make lint
popd
