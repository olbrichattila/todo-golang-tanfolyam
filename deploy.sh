#!/bin/bash
set -e

APP_NAME=todo
IMAGE_NAME=$APP_NAME-image
CONTAINER_NAME=$APP_NAME-container
TAR_NAME=$APP_NAME.tar
VOLUME_NAME=${APP_NAME}_data

echo "[*] Stopping old container (if exists)..."
docker stop $CONTAINER_NAME || true
docker rm $CONTAINER_NAME || true

echo "[*] Removing old image (if exists)..."
docker rmi $IMAGE_NAME || true

echo "[*] Loading new image..."
docker load -i $TAR_NAME

echo "[*] Starting new container..."
docker run -d \
  --name $CONTAINER_NAME \
  -p 8080:80 \
  $IMAGE_NAME

echo "[✔] Deployed successfully."
