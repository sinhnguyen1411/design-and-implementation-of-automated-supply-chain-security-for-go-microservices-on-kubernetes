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
6. Scan SBOM with Grype and compute the count of fixable High/Critical findings (`grype-report.json`).
7. Sign image with Cosign (keyless OIDC), attach SBOM, and attest SLSA-style provenance.
8. Render Kustomize overlay annotations from scan/SBOM outputs for deployment.
9. Verify Cosign signature in registry and upload evidence artifacts.

Registry-backed publishing behavior:
- On `main`, if `GHCR_TOKEN` is configured, the workflow pushes to GHCR, signs the image, attaches SBOM, and attests provenance.
- Without `GHCR_TOKEN`, the workflow falls back to local verification mode: build, test, SBOM, and vulnerability scanning still run, but registry push/sign/attestation steps are skipped.
- On `pull_request`, the workflow also runs in local verification mode.

### Fail-fast Behavior
- Dependency integrity is a hard gate. The pipeline fails when:
  - downloaded module checksums do not verify,
  - or dependency graph requires non-readonly module updates.
- `go mod tidy -diff` is captured as an audit signal (warning) for follow-up cleanup and reproducibility tracking.
- `govulncheck` is a hard gate. If it reports actionable Go vulnerabilities, the job fails.
- Grype is a hard gate on fixable High/Critical findings. Findings marked `wont-fix`, `not-fixed`, or `unknown` remain in the report for auditability but do not block the pipeline by themselves.
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
- `clusterpolicy-cve-threshold.yaml`: requires `security.grype.io/high_critical: "0"` where the annotation represents the count of fixable High/Critical findings.
- `clusterpolicy-require-sbom.yaml`: requires `security.stock-trading.dev/sbom-digest`.

Apply policies:
```bash
kubectl apply -k deploy/policies/kyverno
```

## Local Bootstrap (Kind)
Bootstrap a local cluster with Kyverno and policies:
```bash
./scripts/devsecops_kind_bootstrap.sh
```

## Policy Contract (Admission-Time Requirements)
For workloads labeled `app.kubernetes.io/name=user-service`, admission is expected to enforce:
- Signature verification for the target image.
- SLSA-style provenance attestation presence.
- `security.grype.io/high_critical` annotation must be `"0"`.
- `security.stock-trading.dev/sbom-digest` annotation must exist and be non-empty.

Required Pod annotations for compliant deployments:
- `security.grype.io/high_critical: "<fixable High/Critical count from CI>"` (expected `0`)
- `security.stock-trading.dev/sbom-digest: "<sbom-sha256-or-oci-ref>"`

When available, apply generated CI overlay:
```bash
kubectl apply -k deploy/kubernetes/overlays/ci
```

## Automated Admission Matrix (Docker Desktop)
Run the thesis-aligned matrix on `docker-desktop`:
```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/admission_matrix_demo.ps1 -Context docker-desktop -Namespace stock-trading -ExportDir demo/evidence -ResetNamespace
```

Fixed matrix cases:
- `VALID_ALLOW`: signed + attested + required annotations, expect allow.
- `NEG_UNSIGNED_DENY`: unsigned image, expect deny.
- `NEG_MISSING_SBOM_DENY`: missing `security.stock-trading.dev/sbom-digest`, expect deny.
- `NEG_CVE_THRESHOLD_DENY`: `security.grype.io/high_critical != "0"`, expect deny.
- Regression re-check: rerun valid case after deny cases (`VALID_ALLOW_RECHECK`), expect allow.

Manual cryptographic verification (optional):
```bash
cosign verify --key <path-to-cosign.pub> <signed-image-digest>
cosign verify-attestation --type slsaprovenance --key <path-to-cosign.pub> <signed-image-digest>
```

## Evidence Collection Guidance
The matrix script exports a run directory with:
- `matrix-summary.md`: pass/fail table for all cases (`case -> expected -> actual -> verdict`).
- `matrix-index.json`: machine-readable evidence index for CI/manual parsing.
- `regression-valid-allow.json`: machine-readable result for post-deny valid-admission re-check.
- Per-case evidence files:
  - `kubectl-apply.txt`
  - `kubectl-wait.txt`
  - `events.txt`
  - `workloads.txt`
  - `describe-deployment.txt`
  - `describe-replicasets.txt`
  - `describe-pods.txt`
  - `kyverno-logs.txt`

Complementary CI evidence to collect:
- CI logs for dependency integrity, test, `govulncheck`, Grype, signing, and attestation steps.
- Uploaded artifacts listed above.
- Kubernetes deny/allow events and policy controller logs.

