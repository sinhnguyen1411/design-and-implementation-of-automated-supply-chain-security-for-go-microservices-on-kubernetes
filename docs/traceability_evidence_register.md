# Traceability Evidence Register

As-of: `2026-05-19`

Issue snapshot at this date:
- Open: `#1` to `#9`, `#11`
- Closed: `#10`, `#12`, `#13`, `#14`

Scale context: Repository now contains **10 Go microservices** (user-service, portfolio-service, order-service, risk-service, market-data-service, pricing-service, execution-service, settlement-service, compliance-service, notification-service). All services aligned to Go `1.25.10` baseline as of commit `2821617`.

Purpose:
- Provide a direct objective-to-evidence map for thesis packaging.
- Keep evidence links repo-local and auditable.
- Separate what is already available from what still depends on the final CI run.

## Canonical Evidence Sources
- Dashboard Actions snapshot (source of truth): `docs/security-admission-dashboard/data/actions-runs.snapshot.json`
- Snapshot sync workflow: `.github/workflows/dashboard-data-sync.yml`
- Runtime matrix evidence workflow: `.github/workflows/admission-matrix-evidence.yml`
- Multi-service SCS matrix workflow: `.github/workflows/service-scs-matrix-evidence.yml`
- Admission evidence root: `demo/evidence/20260414-210227/`
- Screenshot appendix pack: `docs/lens_screenshots/`
- Dashboard fallback-only dataset: `docs/security-admission-dashboard/demo-data/evidence/`
- Baseline CI run URL (2026-04-14): `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294`
- Baseline CI artifact IDs (aggregated, single-service era): `dependency-integrity-report=6431038166`, `govulncheck-report=6431053554`, `sbom=6431094849`, `grype-report=6431095016`, `cosign-bundle=6431095189`
- Latest green CI run URL (push, 2026-05-19): `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/25811075803`
- Latest green CI run URL (workflow_dispatch, 2026-05-19): `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/25811079788`
- Latest CI artifact IDs — **user-service canonical** (run 25811075803, per-service structure):
  - `user-service-sbom=6975304852`
  - `user-service-grype-report=6975305135`
  - `user-service-security-gate-findings=6975305453`
  - `user-service-supply-chain-artifacts=6975305801` (cosign bundle, SLSA attestation, Kustomize overlay)
  - `user-service-ubuntu-latest-verify=6975283688`
  - `user-service-macos-latest-verify=6975281971`
- Note: CI v2 (10-service era) exports per-service artifacts; `{service}-supply-chain-artifacts` replaces the old `cosign-bundle` bundle.

## Objective-to-Evidence Matrix
| Objective Item | Related Issue(s) | CI Run URL (Final) | CI Artifact / Path Evidence (Repo-Local) | Admission Run-ID Evidence | Status |
|---|---|---|---|---|---|
| Dependency transparency and control | #3, #14 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | Artifact `dependency-integrity-report` (ID `6431038166`), [.github/workflows/ci-service.yml](../.github/workflows/ci-service.yml), [go_dependency_integrity_baseline.md](go_dependency_integrity_baseline.md), [devsecops_ci_admission.md](devsecops_ci_admission.md#ci-evidence-artifacts) | [../demo/evidence/20260414-210227/matrix-index.json](../demo/evidence/20260414-210227/matrix-index.json) | Partial |
| Go vulnerability and image risk gating | #4, #13 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | Artifacts `govulncheck-report` (ID `6431053554`) and `grype-report` (ID `6431095016`), [.github/workflows/ci-service.yml](../.github/workflows/ci-service.yml), [devsecops_ci_admission.md](devsecops_ci_admission.md#fail-fast-behavior) | [../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt](../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt) | Partial |
| Image signing and provenance attestation | #5, #6, #7 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | Artifacts `sbom` (ID `6431094849`) and `cosign-bundle` (ID `6431095189`), [.github/workflows/ci-service.yml](../.github/workflows/ci-service.yml), [../infra/policies/kyverno/clusterpolicy-verify-images.yaml](../infra/policies/kyverno/clusterpolicy-verify-images.yaml), [devsecops_ci_admission.md](devsecops_ci_admission.md#policy-contract-admission-time-requirements) | [../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt](../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt) | Partial |
| Admission enforcement of trust controls | #7, #8, #10 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | [lens_capture_checklist.md](lens_capture_checklist.md), [lens_screenshots/README.md](lens_screenshots/README.md), [demo_evidence.md](demo_evidence.md) | [../demo/evidence/20260414-210227/matrix-summary.md](../demo/evidence/20260414-210227/matrix-summary.md), [../demo/evidence/20260414-210227/regression-valid-allow.json](../demo/evidence/20260414-210227/regression-valid-allow.json) | Partial |
| End-to-end reproducible pipeline | #3, #4, #9 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | [.github/workflows/ci-service.yml](../.github/workflows/ci-service.yml), [devsecops_ci_admission.md](devsecops_ci_admission.md#deployment-contract-from-ci-output), [final_gap_closing_checklist.md](final_gap_closing_checklist.md#9-automate-kind-and-kyverno-bootstrap-for-reproducible-demo-open) | [../demo/evidence/20260414-210227/matrix-index.json](../demo/evidence/20260414-210227/matrix-index.json) | Partial |
| Reusability and thesis packaging | #11, #12 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | [thesis_spec_en.md](thesis_spec_en.md#thesis-to-implementation-traceability), [final_gap_closing_checklist.md](final_gap_closing_checklist.md), [go_microservice_onboarding_guide.md](go_microservice_onboarding_guide.md) | [../demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md](../demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md), [lens_capture_checklist.md](lens_capture_checklist.md) | Partial |

## Remaining Fill-Ins Before Final Submission
- Close Issue `#11` after final reviewer check of this register and linked evidence bundle.
- Mirror the latest CI run URLs (runs/25811075803, runs/25811079788) into `docs/thesis_spec_en.md` traceability table for examiner navigation.
- Note: `windows-parity-smoke` is intentionally skipped when `runner_target=gh-hosted` — this is by design and does not affect supply-chain gate status.
- Capture one clean Kind bootstrap + admission matrix run-id from a host with `kind` CLI installed (issue #9).

