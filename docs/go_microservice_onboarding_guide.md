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
| Go module path | `github.com/sinhnguyen1411/stock-trading-be/services/user-service` | `<your-module-path>` |
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

## Second-Service Dry-Run Evidence (2026-04-14)
Dry-run and live validation evidence for a second-service simulation is captured at:
- `demo/evidence/20260414-213541-onboarding-second-service/`

Simulation parameters:
- Candidate service name: `portfolio-service`
- Candidate namespace: `portfolio-trading`
- Candidate label path: `app.kubernetes.io/name=portfolio-service`

Observed results:
- Case A (`portfolio-service` label): deployment passed admission and reached ReplicaSet/Pod creation; failure was image pull (`ErrImagePull`) rather than admission deny.
- Case B (`user-service` label on the same deployment): ReplicaSet create was denied by admission due to missing `security.stock-trading.dev/sbom-digest` and missing `security.grype.io/high_critical`.

Primary evidence files:
- `10_case1_apply.txt`, `11_case1_workloads.txt`, `12_case1_events.txt`
- `14_case2_apply.txt`, `16_case2_describe_rs.txt`, `17_case2_events.txt`
- `18_cleanup_namespace_delete.txt`

## Known Assumptions and Limitations for Reuse
- Kyverno image verification policy (`infra/policies/kyverno/clusterpolicy-verify-images.yaml`) now uses a wildcard pattern `ghcr.io/sinhnguyen1411/stock-trading/*` that covers all services under the registry path — no per-service policy change is needed for image signing verification.
- SBOM and CVE annotation policies still use label selectors that may need updating for new service names.
- `infra/scripts/admission_matrix_demo.ps1` is still hardcoded to `user-service` naming and stock-trading-specific metadata/paths.
- New services are registered in `services.yaml` at the repo root — the CI pipeline (`secure-supply-chain.yml`) automatically discovers and builds registered services via matrix strategy.
- Server-side `kubectl apply --dry-run=server` for namespaced resources still requires the target namespace to exist.

## Required Patches When Onboarding a New Service
1. Update policy match scope for the new service:
   - `infra/policies/kyverno/clusterpolicy-require-sbom.yaml`
   - `infra/policies/kyverno/clusterpolicy-cve-threshold.yaml`
   - `infra/policies/kyverno/clusterpolicy-verify-images.yaml`
2. Update deployment/manifests labels, namespace, image repo, and annotation flow consistently.
3. Parameterize or clone `infra/scripts/admission_matrix_demo.ps1` for the new service name/namespace and image path.
4. Re-run deny/allow matrix and append evidence in the same bundle format used by this repository.
