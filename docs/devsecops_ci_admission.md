# DevSecOps CI + Admission Flow

This document describes the secure supply-chain pipeline (`dependency-integrity -> test -> govulncheck -> build -> SBOM -> scan -> sign -> attest -> push`) and Kyverno-based admission enforcement.

## CI Workflow (`.github/workflows/secure-supply-chain.yml`)
The workflow runs on pushes/PRs to `main` and on manual dispatch.

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
6. Scan SBOM with Grype and fail on High/Critical findings (`grype-report.json`).
7. Sign image with Cosign (keyless OIDC), attach SBOM, and attest SLSA-style provenance.
8. Render Kustomize overlay annotations from scan/SBOM outputs for deployment.
9. Verify Cosign signature in registry and upload evidence artifacts.

### Fail-fast Behavior
- Dependency integrity is a hard gate. The pipeline fails when:
  - downloaded module checksums do not verify,
  - or dependency graph requires non-readonly module updates.
- `go mod tidy -diff` is captured as an audit signal (warning) for follow-up cleanup and reproducibility tracking.
- `govulncheck` is a hard gate. If it reports actionable Go vulnerabilities, the job fails.
- Grype is a hard gate. If High/Critical findings are detected, the pipeline fails and artifact publication stops.
- Deployment-side admission policies still enforce runtime constraints even if deployment is attempted manually.

### CI Evidence Artifacts
- `dependency-integrity-report`: output of module integrity checks (`go mod verify`, readonly resolution, tidy diff).
- `govulncheck-report`: Go vulnerability scan output.
- `sbom`: `sbom.spdx.json`.
- `grype-report`: `grype-report.json`.
- `cosign-bundle`: provenance file, SBOM, scan report, and generated Kustomize overlay.

### Required Permissions/Secrets
- `GITHUB_TOKEN` with `packages:write` and `id-token:write`.
- Optional fallback secrets for GHCR pushes when repository package permissions are restricted:
  - `GHCR_USERNAME`
  - `GHCR_TOKEN`

## Admission Policies (Kyverno)
Resources under `deploy/policies/kyverno/`:
- `clusterpolicy-verify-images.yaml`: verifies signed images and provenance.
- `clusterpolicy-cve-threshold.yaml`: requires `security.grype.io/high_critical: "0"`.
- `clusterpolicy-require-sbom.yaml`: requires `security.stock-trading.dev/sbom-digest`.

Apply policies:
```bash
kubectl apply -k deploy/policies/kyverno
```

## Local Demo (Kind)
Bootstrap a local cluster with Kyverno and policies:
```bash
./scripts/devsecops_kind_bootstrap.sh
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

