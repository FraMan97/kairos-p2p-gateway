#!/bin/bash

# 0. Kill any hanging processes occupying port 8085 and 8081
echo "Cleaning up previous port-forwards..."
fuser -k 3000/tcp || true

# 1. Build the images
./build-images.sh

# 2. Load the images into the local Kind cluster
echo "Loading images into Kind..."
kind load docker-image kairos-gateway:local --name kairos-vault

# 3. Install the P2P infrastructure with Helm
echo "Installing Helm Chart..."
helm upgrade --install kairos-gateway ./helm -f ./helm/values-dev.yaml \
  --set gateway.image=kairos-gateway:local \
  --set gateway.pullPolicy=Never 

# 4. Wait for Kubernetes to finish the job
echo "Waiting for all pods to be ready (this may take a few seconds)..."
kubectl rollout status deployment/kairos-gateway --timeout=120s

# 5. Configure Port Forwards (in background)
echo "Configuring Node API Port-Forward on localhost:8086..."
kubectl port-forward svc/kairos-gateway 8086:3000 > /dev/null 2>&1 &

echo "Deploy completed! The Gateway is running."
echo "- Gateway API is on http://localhost:8086"