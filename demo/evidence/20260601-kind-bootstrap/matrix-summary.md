# Kind Bootstrap + Admission Matrix Evidence

- **Date:** 2026-06-01
- **Environment:** fresh local **Kind** cluster `devsecops` (node `kindest/node:v1.35.0`), context `kind-devsecops` — created clean via `RESET_CLUSTER=true infra/scripts/devsecops_kind_bootstrap.sh` (closes issue #9's "reproducible Kind bootstrap" requirement; previous evidence used docker-desktop).
- **Kyverno:** v1.12.5, all controllers rolled out, webhook responding.
- **Policies applied & READY (Enforce):** `verify-stock-trading-images`, `require-sbom-annotation`, `enforce-cve-threshold`.
- **Tooling:** kind v0.31.0, kubectl v1.34.1, Docker 29.3.1.

## Bootstrap

See [bootstrap.log](bootstrap.log) — clean-run from cluster creation through `kubectl get clusterpolicies` (all three `READY=True`).

## Admission matrix (see [admission-matrix.txt](admission-matrix.txt))

| Scenario | Image | Annotations | Expected | Actual |
|---|---|---|---|---|
| `A_VALID_ALLOW` | pause (non-stock-trading) | sbom-digest set, high_critical=`0` | Admit | ✅ **Admitted** (pod Running 1/1) |
| `B_NEG_MISSING_SBOM_DENY` | pause | high_critical=`0`, **no sbom-digest** | Deny | 🚫 **Denied** by `require-sbom-annotation` |
| `C_NEG_CVE_THRESHOLD_DENY` | pause | sbom-digest set, **high_critical=`2`** | Deny | 🚫 **Denied** by `enforce-cve-threshold` |
| `D_NEG_UNSIGNED_IMAGE_DENY` | `ghcr.io/sinhnguyen1411/stock-trading/order-service:dev` (unsigned) | both set | Deny | 🚫 **Denied** by `verify-stock-trading-images` (keyless signature + SLSA attestation both failed via GHCR/Rekor) |

All four outcomes match expectations. The deny cases require no GitHub-OIDC-signed image; the `verify-stock-trading-images` deny genuinely queried Sigstore Rekor and GHCR and rejected the unverified image.

> Note: a full `VALID_ALLOW` of a *signed* stock-trading image requires an image signed by the CI workflow's GitHub OIDC identity (produced only in CI, e.g. run 26732257799), which cannot be reproduced on a local-only cluster. The local matrix therefore demonstrates the enforce/deny behavior plus a compliant-annotation allow.
