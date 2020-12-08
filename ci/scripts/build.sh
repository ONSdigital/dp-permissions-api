#!/bin/bash -eux

pushd dp-permissions-api
  make build
  cp build/dp-permissions-api Dockerfile.concourse ../build
popd
