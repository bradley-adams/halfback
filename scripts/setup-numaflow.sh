#!/usr/bin/env bash
set -euo pipefail

echo "Creating numaflow-system namespace (if not present)..."
kubectl create namespace numaflow-system --dry-run=client -o yaml | kubectl apply -f -

echo "Installing Numaflow..."
kubectl apply -n numaflow-system -f https://raw.githubusercontent.com/numaproj/numaflow/stable/config/install.yaml

echo "Waiting for Numaflow pods to be ready..."
kubectl wait --for=condition=Ready pods --all -n numaflow-system --timeout=120s

echo "Done. Current pods:"
kubectl get pods -n numaflow-system