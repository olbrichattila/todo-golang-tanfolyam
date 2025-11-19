#!/bin/bash

set -e

APP_NAME=todo
IMAGE_NAME=$APP_NAME-image
TAR_NAME=$APP_NAME.tar

docker build -t $IMAGE_NAME .

# Save image to tar (zip is optional)
docker save $IMAGE_NAME -o $TAR_NAME
# zip $TAR_NAME.zip $TAR_NAME  # Optional
echo "Image saved to $TAR_NAME"
