#!/usr/bin/env bash
set -euo pipefail

# Delete the thesis Kind cluster used by bootstrap/demo scripts.
# Idempotent: succeeds even if the cluster does not exist.

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-devsecops}"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required" >&2
  exit 1
fi

if kind get clusters | grep -Fxq "$KIND_CLUSTER_NAME"; then
  echo "Deleting Kind cluster: $KIND_CLUSTER_NAME"
  kind delete cluster --name "$KIND_CLUSTER_NAME"
  echo "Deleted."
else
  echo "Kind cluster '$KIND_CLUSTER_NAME' does not exist. Nothing to do."
fi
