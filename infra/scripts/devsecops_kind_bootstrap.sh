#!/usr/bin/env bash
set -euo pipefail

# Bootstrap a Kind cluster with Kyverno and apply supply-chain policies.
# Requires: kind and kubectl.

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-devsecops}"
KYVERNO_VERSION="${KYVERNO_VERSION:-v1.12.5}"
KYVERNO_ROLLOUT_TIMEOUT="${KYVERNO_ROLLOUT_TIMEOUT:-240s}"
RESET_CLUSTER="${RESET_CLUSTER:-false}"
KIND_CONTEXT="kind-${KIND_CLUSTER_NAME}"
KYVERNO_INSTALL_URL="https://raw.githubusercontent.com/kyverno/kyverno/${KYVERNO_VERSION}/config/release/install.yaml"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required" >&2
  exit 1
fi

if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl is required" >&2
  exit 1
fi

cluster_exists() {
  kind get clusters | grep -Fxq "$KIND_CLUSTER_NAME"
}

resolve_kyverno_deployment() {
  local names
  names="$(kubectl -n kyverno get deploy -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' 2>/dev/null || true)"

  if echo "$names" | grep -Fxq "kyverno-admission-controller"; then
    echo "deploy/kyverno-admission-controller"
    return
  fi

  if echo "$names" | grep -Fxq "kyverno"; then
    echo "deploy/kyverno"
    return
  fi

  echo "$names" | head -n 1 | sed '/^$/d; s#^#deploy/#'
}

wait_for_kyverno() {
  local deploy_ref
  for _ in $(seq 1 30); do
    deploy_ref="$(resolve_kyverno_deployment)"
    if [[ -n "${deploy_ref}" ]]; then
      if kubectl -n kyverno rollout status "${deploy_ref}" --timeout="${KYVERNO_ROLLOUT_TIMEOUT}"; then
        echo "Kyverno ready via ${deploy_ref}"
        return 0
      fi
    fi
    sleep 4
  done

  echo "Kyverno deployment did not become ready in time." >&2
  kubectl -n kyverno get deploy -o wide || true
  return 1
}

echo "[1/4] Preparing Kind cluster: $KIND_CLUSTER_NAME"
if [[ "${RESET_CLUSTER}" == "true" ]] && cluster_exists; then
  echo "RESET_CLUSTER=true -> deleting existing cluster ${KIND_CLUSTER_NAME}"
  kind delete cluster --name "$KIND_CLUSTER_NAME"
fi

if cluster_exists; then
  echo "Kind cluster ${KIND_CLUSTER_NAME} already exists. Reusing."
else
  echo "Creating Kind cluster ${KIND_CLUSTER_NAME}"
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
fi

if kubectl config get-contexts -o name | grep -Fxq "$KIND_CONTEXT"; then
  kubectl config use-context "$KIND_CONTEXT" >/dev/null
fi

echo "[2/4] Installing or reconciling Kyverno (${KYVERNO_VERSION})"
kubectl create namespace kyverno || true
kubectl apply -f "$KYVERNO_INSTALL_URL"
wait_for_kyverno

echo "[3/4] Applying supply-chain policies"
kubectl apply -k deploy/policies/kyverno

echo "[4/4] Verifying applied policies"
kubectl get clusterpolicies
echo "Done."
echo "Tip: use RESET_CLUSTER=true ./scripts/devsecops_kind_bootstrap.sh for clean rerun."
echo "Tip: use ./scripts/devsecops_kind_reset.sh to delete the Kind cluster."
