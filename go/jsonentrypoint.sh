#!/bin/bash

if [ -z "$DOCKER_USERNAME" ] || [ -z "$DOCKER_TOKEN" ]; then
    echo "DOCKER_USERNAME and DOCKER_TOKEN must be set"
    exit 1
fi

echo "$DOCKER_TOKEN" | docker login --username "$DOCKER_USERNAME" --password-stdin

# Take the image to scan from command-line arguments
IMAGE_TO_SCAN="$2"

# Make sure an image is specified
if [ -z "$IMAGE_TO_SCAN" ]; then
    echo "Error: No image specified for scanning."
    exit 1
fi

if [ "$1" = "grype" ]; then
    SCANNER="grype -o json"
    OUTPUT_EXT=".json"
    echo "Updating grype DB..."
    grype db delete
    grype db update
elif [ "$1" = "trivy" ]; then
    SCANNER="trivy image --format json --timeout 2h"
    OUTPUT_EXT=".json"
elif [ "$1" = "blackduck" ]; then
    SCANNER="blackduck"
    OUTPUT_EXT=".txt"
else
    echo "Invalid scanner tool"
    exit 1
fi

FILENAME=$(basename "${IMAGE_TO_SCAN//:/_}")
FILENAME=${FILENAME//_/:}
export TZ='UTC'
TIMESTAMP=$(date "+%Y-%m-%d")
OUTPUT_FILE="/output/${DOCKER_USERNAME}:${FILENAME}:${TIMESTAMP}:${1}${OUTPUT_EXT}"

# Print Docker CLI version
echo "Docker CLI version:" > $OUTPUT_FILE
docker version >> $OUTPUT_FILE
echo >> $OUTPUT_FILE

# Print Docker image layer information
echo "Docker image layer information for $IMAGE_TO_SCAN:" >> $OUTPUT_FILE
docker history $IMAGE_TO_SCAN >> $OUTPUT_FILE
echo >> $OUTPUT_FILE

# Print grype version and DB status
if [ "$1" = "grype" ]; then
    echo "grype version:" >> $OUTPUT_FILE
    grype version >> $OUTPUT_FILE
    echo "Grype DB status:" >> $OUTPUT_FILE
    grype db status >> $OUTPUT_FILE
    echo >> $OUTPUT_FILE
fi

# Perform scan and append to file
$SCANNER $IMAGE_TO_SCAN >> $OUTPUT_FILE

# Print the output to stdout in addition to writing to the file
cat $OUTPUT_FILE
