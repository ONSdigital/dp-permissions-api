#!/bin/bash -eux

pushd dp-permissions-api
  make lint-go
popd
