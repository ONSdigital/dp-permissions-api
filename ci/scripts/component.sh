#!/bin/bash -eux

set +x
git config --global url."https://${GITHUB_USERNAME}:${GITHUB_TOKEN}@github.com/${GITHUB_REPO_PATH}".insteadOf "https://github.com/${GITHUB_REPO_PATH}"
set -x

cwd=$(pwd)

pushd $cwd/dp-permissions-api
  make test-component
popd
