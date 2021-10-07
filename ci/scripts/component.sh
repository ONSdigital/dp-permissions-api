#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-permissions-api
  make test-component
popd
