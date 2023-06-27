#!/bin/bash

if [ -z "$DOCKER_USERNAME" ] || [ -z "$DOCKER_TOKEN" ]; then
    echo "DOCKER_USERNAME and DOCKER_TOKEN must be set"
    exit 1
fi

echo "$DOCKER_TOKEN" | docker login --username "$DOCKER_USERNAME" --password-stdin

if [ "$1" = "grype" ]; then
    SCANNER="grype"
    OUTPUT_EXT=".txt"
    echo "grype version:"
    grype version
elif [ "$1" = "trivy" ]; then
    SCANNER="trivy image --format json --timeout 2h"
    OUTPUT_EXT=".json"
    echo "trivy version:"
    trivy --version
elif [ "$1" = "blackduck" ]; then
    SCANNER="blackduck"
    OUTPUT_EXT=".txt"
    echo "blackduck version:"
    ./detect.sh --version
else
    echo "Invalid scanner tool"
    exit 1
fi

IMAGE_TO_SCAN="$2"
if [ -z "$IMAGE_TO_SCAN" ]; then
    echo "An image to scan must be provided"
    exit 1
fi

echo

#$SCANNER $IMAGE_TO_SCAN > "/output/$(basename "$IMAGE_TO_SCAN")_${1}${OUTPUT_EXT}"

#FILENAME=$(basename "${IMAGE_TO_SCAN//:/_}")
#FILENAME=${FILENAME//_/:}
#$SCANNER $IMAGE_TO_SCAN > "/output/${DOCKER_USERNAME}:${FILENAME}:${1}${OUTPUT_EXT}"

FILENAME=$(basename "${IMAGE_TO_SCAN//:/_}")
FILENAME=${FILENAME//_/:}
export TZ='UTC'
#TIMESTAMP=$(date "+%Y-%m-%d:%H:%M:%S")
TIMESTAMP=$(date "+%Y-%m-%d")

$SCANNER $IMAGE_TO_SCAN > "/output/${DOCKER_USERNAME}:${FILENAME}:${TIMESTAMP}:${1}${OUTPUT_EXT}"
