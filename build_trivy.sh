#!/bin/bash
DOCKER_USERNAME='jspawar80'
DOCKER_PASSWORD='dckr_pat_ecD6s3cZVofGBBE7jSAR_ykbxL4'

DOCKER_TAG='interlynk_scanner_trivy'
if [ -z "${DOCKER_USERNAME}" ] || [ -z "${DOCKER_PASSWORD}" ]; then
  echo "Docker credentials are not set."
  exit 1
fi

if [ -z "${DOCKER_TAG}" ]; then
  echo "Docker tag is not provided."
  exit 1
fi

LATEST_TRIVY=$(curl --silent "https://api.github.com/repos/aquasecurity/trivy/releases/latest" | jq -r '.tag_name')
echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

docker build -t ${DOCKER_USERNAME}/${DOCKER_TAG}:latest -f Dockerfiletrivy .
docker tag ${DOCKER_USERNAME}/${DOCKER_TAG}:latest ${DOCKER_USERNAME}/${DOCKER_TAG}:${LATEST_TRIVY}
docker push ${DOCKER_USERNAME}/${DOCKER_TAG}:latest
docker push ${DOCKER_USERNAME}/${DOCKER_TAG}:${LATEST_TRIVY}
