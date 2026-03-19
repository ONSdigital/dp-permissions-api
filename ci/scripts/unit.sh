#!/bin/bash -eux

pushd dp-permissions-api
  make test-go
popd
