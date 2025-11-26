#!/bin/bash
set -e

# Change to project root
cd "$(dirname "$0")/.."

# Configuration - update this to your Docker Hub username
IMAGE_NAME="shinde11/filebrowser-tunnel"

echo "Step 1: Building binaries..."
./scripts/build.sh

echo ""
echo "Step 2: Building Docker image..."
docker build -t ${IMAGE_NAME}:latest .

echo ""
echo "Step 3: Pushing Docker image..."
docker push ${IMAGE_NAME}:latest

echo ""
echo "=========================================="
echo "Success! Image pushed to ${IMAGE_NAME}:latest"
echo ""
echo "Server deployment:"
echo "  docker run -d -p 80:80 -e DOMAIN=https://your-domain.com ${IMAGE_NAME}:latest"
echo ""
echo "Client usage:"
echo "  curl -sL https://your-domain.com | sh"
echo "  curl -sL https://your-domain.com | sh -s -- /path/to/serve"
echo "=========================================="

