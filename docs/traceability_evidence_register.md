# Traceability Evidence Register

As-of: `2026-04-14`

Issue snapshot at this date:
- Open: `#1` to `#9`, `#11`
- Closed: `#10`, `#12`, `#13`, `#14`

Purpose:
- Provide a direct objective-to-evidence map for thesis packaging.
- Keep evidence links repo-local and auditable.
- Separate what is already available from what still depends on the final CI run.

## Canonical Evidence Sources
- Admission evidence root: `demo/evidence/20260414-210227/`
- Screenshot appendix pack: `docs/lens_screenshots/`
- Dashboard fallback-only dataset: `docs/security-admission-dashboard/demo-data/evidence/`
- Final CI run URL: `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294`
- Final CI artifact IDs: `dependency-integrity-report=6431038166`, `govulncheck-report=6431053554`, `sbom=6431094849`, `grype-report=6431095016`, `cosign-bundle=6431095189`

## Objective-to-Evidence Matrix
| Objective Item | Related Issue(s) | CI Run URL (Final) | CI Artifact / Path Evidence (Repo-Local) | Admission Run-ID Evidence | Status |
|---|---|---|---|---|---|
| Dependency transparency and control | #3, #14 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | Artifact `dependency-integrity-report` (ID `6431038166`), [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [go_dependency_integrity_baseline.md](go_dependency_integrity_baseline.md), [devsecops_ci_admission.md](devsecops_ci_admission.md#ci-evidence-artifacts) | [../demo/evidence/20260414-210227/matrix-index.json](../demo/evidence/20260414-210227/matrix-index.json) | Partial |
| Go vulnerability and image risk gating | #4, #13 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | Artifacts `govulncheck-report` (ID `6431053554`) and `grype-report` (ID `6431095016`), [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [devsecops_ci_admission.md](devsecops_ci_admission.md#fail-fast-behavior) | [../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt](../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt) | Partial |
| Image signing and provenance attestation | #5, #6, #7 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | Artifacts `sbom` (ID `6431094849`) and `cosign-bundle` (ID `6431095189`), [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [../deploy/policies/kyverno/clusterpolicy-verify-images.yaml](../deploy/policies/kyverno/clusterpolicy-verify-images.yaml), [devsecops_ci_admission.md](devsecops_ci_admission.md#policy-contract-admission-time-requirements) | [../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt](../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt) | Partial |
| Admission enforcement of trust controls | #7, #8, #10 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | [lens_capture_checklist.md](lens_capture_checklist.md), [lens_screenshots/README.md](lens_screenshots/README.md), [demo_evidence.md](demo_evidence.md) | [../demo/evidence/20260414-210227/matrix-summary.md](../demo/evidence/20260414-210227/matrix-summary.md), [../demo/evidence/20260414-210227/regression-valid-allow.json](../demo/evidence/20260414-210227/regression-valid-allow.json) | Partial |
| End-to-end reproducible pipeline | #3, #4, #9 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [devsecops_ci_admission.md](devsecops_ci_admission.md#deployment-contract-from-ci-output), [final_gap_closing_checklist.md](final_gap_closing_checklist.md#9-automate-kind-and-kyverno-bootstrap-for-reproducible-demo-open) | [../demo/evidence/20260414-210227/matrix-index.json](../demo/evidence/20260414-210227/matrix-index.json) | Partial |
| Reusability and thesis packaging | #11, #12 | `https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/24406893294` | [thesis_spec_en.md](thesis_spec_en.md#thesis-to-implementation-traceability), [final_gap_closing_checklist.md](final_gap_closing_checklist.md), [go_microservice_onboarding_guide.md](go_microservice_onboarding_guide.md) | [../demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md](../demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md), [lens_capture_checklist.md](lens_capture_checklist.md) | Partial |

## Remaining Fill-Ins Before Final Submission
- Close Issue `#11` after final reviewer check of this register and linked evidence bundle.
- Optionally mirror the same CI run URL + artifact IDs into `docs/thesis_spec_en.md` for one-click examiner navigation.
