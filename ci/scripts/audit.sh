#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-permissions-api
  make audit
popd