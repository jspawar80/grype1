#!/bin/bash

if [ -z "$DOCKER_USERNAME" ] || [ -z "$DOCKER_TOKEN" ]; then
    echo "DOCKER_USERNAME and DOCKER_TOKEN must be set"
    exit 1
fi

echo "$DOCKER_TOKEN" | docker login --username "$DOCKER_USERNAME" --password-stdin

# Print Docker CLI version
echo "Docker CLI version:"
docker version
echo

# Take the image to scan from command-line arguments
IMAGE_TO_SCAN="$2"

# Make sure an image is specified
if [ -z "$IMAGE_TO_SCAN" ]; then
    echo "Error: No image specified for scanning."
    exit 1
fi

# Print Docker image layer information
echo "Docker image layer information for $IMAGE_TO_SCAN:"
docker history $IMAGE_TO_SCAN
echo

if [ "$1" = "grype" ]; then
    SCANNER="grype"
    OUTPUT_EXT=".txt"
    echo "grype version:"
    grype version
    echo "Updating grype DB..."
    grype db delete
    grype db update
    echo "Grype DB status:"
    grype db status
    echo
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

FILENAME=$(basename "${IMAGE_TO_SCAN//:/_}")
FILENAME=${FILENAME//_/:}
export TZ='UTC'
TIMESTAMP=$(date "+%Y-%m-%d")

$SCANNER $IMAGE_TO_SCAN > "/output/${DOCKER_USERNAME}:${FILENAME}:${TIMESTAMP}:${1}${OUTPUT_EXT}"

# Print the output to stdout in addition to writing to the file
cat "/output/${DOCKER_USERNAME}:${FILENAME}:${TIMESTAMP}:${1}${OUTPUT_EXT}"
