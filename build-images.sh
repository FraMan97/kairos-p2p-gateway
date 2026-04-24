#!/bin/bash

echo "Building Gateway image..."
docker build -t kairos-gateway:local -f Dockerfile .