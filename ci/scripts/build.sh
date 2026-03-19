#!/bin/bash -eux

pushd dp-permissions-api
  make build-go
  cp build/dp-permissions-api Dockerfile.concourse ../build
popd
