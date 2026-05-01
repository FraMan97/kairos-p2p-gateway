#!/bin/bash

helm upgrade --install kairos-gateway ./helm \
  -f ./helm/values-dev.yaml \
  -f ./helm/values-prod.yaml \
  --namespace production --create-namespace

kubectl rollout restart deployment/kairos-gateway -n production

echo "Deploy Gateway completed!"