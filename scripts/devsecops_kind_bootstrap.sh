#!/usr/bin/env bash
set -euo pipefail

# Bootstrap a Kind cluster with Kyverno and apply supply-chain policies.
# Requires: kind and kubectl.

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-devsecops}"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required" >&2
  exit 1
fi

if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl is required" >&2
  exit 1
fi

echo "[1/3] Creating Kind cluster: $KIND_CLUSTER_NAME"
kind create cluster --name "$KIND_CLUSTER_NAME" --config - <<'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 30080
        hostPort: 30080
      - containerPort: 30090
        hostPort: 30090
EOF

echo "[2/3] Installing Kyverno"
kubectl create namespace kyverno || true
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno/v1.12.5/config/release/install.yaml
kubectl -n kyverno rollout status deploy/kyverno --timeout=120s

echo "[3/3] Applying supply-chain policies"
kubectl apply -k deploy/policies/kyverno

echo "Done. Verify with: kubectl get clusterpolicies"
