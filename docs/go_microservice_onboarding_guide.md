# Go Microservice Onboarding and Reuse Guide

This guide describes how to reuse the current supply-chain security baseline for another Go microservice with minimal redesign.

## Reuse Objective
- Keep the same control chain: dependency integrity, vulnerability gates, SBOM, signing, attestation, and admission enforcement.
- Change only service-specific parameters and deployment metadata.
- Produce evidence in the same format used by this repository.

## Parameters You Must Replace
For a new service, define and version these values first:

| Parameter | Example in current repo | New service value |
|---|---|---|
| Go module path | `github.com/sinhnguyen1411/stock-trading-be` | `<your-module-path>` |
| Image repository | `ghcr.io/sinhnguyen1411/stock-trading/user-service` | `<registry>/<org>/<service>` |
| Kubernetes app label | `app.kubernetes.io/name=user-service` | `<service-app-label>` |
| Namespace | `stock-trading` | `<service-namespace>` |
| Runtime config secret keys | `auth-access-token-secret`, `auth-refresh-token-secret` | `<service-secret-keys>` |
| SBOM annotation key | `security.stock-trading.dev/sbom-digest` | Keep key or publish approved replacement |
| CVE annotation key | `security.grype.io/high_critical` | Keep key or publish approved replacement |

## Pipeline and Secrets Setup
Set the following before enabling secure releases:
- GitHub Actions permissions: `packages:write`, `id-token:write`.
- Registry credentials (if not using default token path): `GHCR_USERNAME`, `GHCR_TOKEN`.
- Cosign keyless flow in CI for signing and attestation.
- CI artifact retention suitable for thesis/demo evidence review.

## Admission Policy Contract
Your deployment must satisfy this contract:
- Image is signed by the expected key identity.
- Image has SLSA provenance attestation.
- Pod annotation `security.grype.io/high_critical` exists and equals `"0"`.
- Pod annotation `security.stock-trading.dev/sbom-digest` exists and is non-empty.

If you rename annotation keys, update all three locations together:
- CI overlay rendering logic.
- Kyverno policies.
- Demo/evidence scripts and docs.

## Rollout Order for a New Service
1. Port Dockerfile and deployment baseline to least-privilege runtime.
2. Enable dependency-integrity + tests + `govulncheck` in CI.
3. Add SBOM and Grype gate.
4. Enable signing + attestation on release path.
5. Apply admission policies and validate against real deployment manifests.
6. Run the admission matrix and export evidence bundle.

## Porting in One Day Checklist
Use this sequence to onboard a new Go microservice quickly:

| Time box | Activity | Output |
|---|---|---|
| Hour 1 | Replace module/image/namespace/labels and secret keys | Updated configs and manifests |
| Hour 2 | Wire CI stages (integrity, test, govulncheck, SBOM, scan) | Passing CI pre-release gates |
| Hour 3 | Wire signing + attestation and verify manually | Verifiable signed digest |
| Hour 4 | Apply Kyverno contract policies in local cluster | Enforcement-ready admission |
| Hour 5 | Run admission matrix (`VALID_ALLOW` + 3 deny cases) | `matrix-summary.md` + `matrix-index.json` |
| Hour 6 | Capture docs and traceability updates | Review-ready evidence package |

## Definition of Done
The onboarding is complete when all statements are true:
- CI fails on dependency integrity drift or actionable Go vulnerability findings.
- CI fails on fixable High/Critical image findings above threshold.
- Release image is signed and has SLSA provenance attestation.
- Admission blocks unsigned image deployments.
- Admission blocks missing SBOM annotation deployments.
- Admission blocks non-zero CVE threshold annotation deployments.
- Valid signed and annotated deployment is admitted.
- Evidence bundle includes summary table, JSON index, and per-case logs.
- Traceability document marks reusability as implemented with links to evidence.
