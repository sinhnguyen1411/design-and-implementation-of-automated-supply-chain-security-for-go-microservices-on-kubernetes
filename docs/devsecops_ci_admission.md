# DevSecOps CI + Admission Flow

This document describes the secure supply-chain pipeline (`dependency-integrity -> test -> govulncheck -> build -> SBOM -> scan -> sign -> attest -> push`) and Kyverno-based admission enforcement.

## CI Workflow (`.github/workflows/secure-supply-chain.yml`)
The workflow runs on pushes/PRs to `Thesis-SCS` and on manual dispatch.

Pipeline stages:
1. Run dependency-integrity checks for Go modules:
   - `go mod download`
   - `go mod verify`
   - `go list -deps -mod=readonly ./...`
   - `go mod tidy -diff` (audit-only signal)
2. Run Go unit tests (`go test ./...`).
3. Run Go vulnerability analysis with `govulncheck` (fail-fast).
4. Build and push image to GHCR.
5. Generate SBOM with Syft (`sbom.spdx.json`).
6. Scan SBOM with Grype and count High/Critical findings (`grype-report.json`).
7. Sign image with Cosign (keyless OIDC), attach SBOM, and attest SLSA-style provenance.
8. Render Kustomize overlay annotations from scan/SBOM outputs for deployment.
9. Verify Cosign signature in registry and upload evidence artifacts.

### Fail-fast Behavior
- Dependency integrity is a hard gate. The pipeline fails when:
  - downloaded module checksums do not verify,
  - or dependency graph requires non-readonly module updates.
- `go mod tidy -diff` is captured as an audit signal (warning) for follow-up cleanup and reproducibility tracking.
- `govulncheck` is a hard gate. If it reports actionable Go vulnerabilities, the job fails.
- Grype findings are captured, and enforcement can be switched on with `SECURITY_GATE=true`.
- Deployment-side admission policies still enforce runtime constraints even if deployment is attempted manually.

### CI Evidence Artifacts
- `dependency-integrity-report`: output of module integrity checks (`go mod verify`, readonly resolution, tidy diff).
- `govulncheck-report`: Go vulnerability scan output.
- `sbom`: `sbom.spdx.json`.
- `grype-report`: `grype-report.json`.
- `cosign-bundle`: provenance file, SBOM, scan report, and generated Kustomize overlay.

### Required Permissions/Secrets
- `GITHUB_TOKEN` with `packages:write` and `id-token:write`.
- For key-based signing scenarios, add private-key secrets as needed (`COSIGN_PRIVATE_KEY`, `COSIGN_PASSWORD`).

## Admission Policies (Kyverno)
Resources under `deploy/policies/kyverno/`:
- `cosign-public-key.yaml`: ConfigMap containing Cosign public key material.
- `clusterpolicy-verify-images.yaml`: verifies signed images and provenance.
- `clusterpolicy-cve-threshold.yaml`: requires `security.grype.io/high_critical: "0"`.
- `clusterpolicy-require-sbom.yaml`: requires `security.stock-trading.dev/sbom-digest`.

Apply policies:
```bash
kubectl apply -k deploy/policies/kyverno
```

Publish Cosign key:
```bash
kubectl -n kyverno create configmap cosign-public-key --from-file=cosign.pub=./cosign.pub
```

## Local Demo (Kind)
Bootstrap a local cluster with Kyverno and policies:
```bash
COSIGN_PUB_PATH=./cosign.pub ./scripts/devsecops_kind_bootstrap.sh
```

## Deployment Metadata Requirements
Required Pod annotations:
- `security.grype.io/high_critical: "<value from CI>"` (expected `0` for compliant builds)
- `security.stock-trading.dev/sbom-digest: "<sbom-sha256-or-oci-ref>"`

When available, apply generated CI overlay:
```bash
kubectl apply -k deploy/kubernetes/overlays/ci
```

## Verifying Enforcement
- Deploy unsigned or unannotated images: admission should deny.
- Deploy signed/scanned image with required annotations: admission should allow.
- Verify cryptographic evidence manually:
  ```bash
  cosign verify --keyless ghcr.io/sinhnguyen1411/stock-trading/user-service:dev
  cosign verify-attestation --type slsaprovenance --keyless ghcr.io/sinhnguyen1411/stock-trading/user-service:dev
  ```

## Evidence Collection Guidance
Collect:
- CI logs for dependency integrity, test, `govulncheck`, Grype, signing, and attestation steps.
- Uploaded artifacts listed above.
- Kubernetes denial/allow events (`kubectl events`, Kyverno policy reports/logs).

