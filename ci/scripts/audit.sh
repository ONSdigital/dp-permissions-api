---
platform: linux

image_resource:
  type: docker-image
  source:
    repository: onsdigital/dp-concourse-tools-nancy
    tag: latest

inputs:
  - name: dp-permissions-api
    path: dp-permissions-api

run:
  path: dp-permissions-api/ci/scripts/audit.sh 