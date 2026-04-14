# Traceability Evidence Register

As-of: `2026-04-14`

Issue snapshot at this date:
- Open: `#1` to `#12`
- Closed: `#13`, `#14`

Purpose:
- Provide a direct objective-to-evidence map for thesis packaging.
- Keep evidence links repo-local and auditable.
- Separate what is already available from what still depends on the final CI run.

## Canonical Evidence Sources
- Admission evidence root: `demo/evidence/20260414-210227/`
- Screenshot appendix pack: `docs/lens_screenshots/`
- Dashboard fallback-only dataset: `docs/security-admission-dashboard/demo-data/evidence/`
- Final CI run URL + artifact IDs: pending under Issue `#11`

## Objective-to-Evidence Matrix
| Objective Item | Related Issue(s) | CI Run URL (Final) | CI Artifact / Path Evidence (Repo-Local) | Admission Run-ID Evidence | Status |
|---|---|---|---|---|---|
| Dependency transparency and control | #3, #14 | Pending (`#11`) | [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [go_dependency_integrity_baseline.md](go_dependency_integrity_baseline.md), [devsecops_ci_admission.md](devsecops_ci_admission.md#ci-evidence-artifacts) | [../demo/evidence/20260414-210227/matrix-index.json](../demo/evidence/20260414-210227/matrix-index.json) | Partial |
| Go vulnerability and image risk gating | #4, #13 | Pending (`#11`) | [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [devsecops_ci_admission.md](devsecops_ci_admission.md#fail-fast-behavior) | [../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt](../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt) | Partial |
| Image signing and provenance attestation | #5, #6, #7 | Pending (`#11`) | [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [../deploy/policies/kyverno/clusterpolicy-verify-images.yaml](../deploy/policies/kyverno/clusterpolicy-verify-images.yaml), [devsecops_ci_admission.md](devsecops_ci_admission.md#policy-contract-admission-time-requirements) | [../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt](../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt) | Partial |
| Admission enforcement of trust controls | #7, #8, #10 | Pending (`#11`) | [lens_capture_checklist.md](lens_capture_checklist.md), [lens_screenshots/README.md](lens_screenshots/README.md) | [../demo/evidence/20260414-210227/matrix-summary.md](../demo/evidence/20260414-210227/matrix-summary.md), [../demo/evidence/20260414-210227/regression-valid-allow.json](../demo/evidence/20260414-210227/regression-valid-allow.json) | Partial |
| End-to-end reproducible pipeline | #3, #4, #9 | Pending (`#11`) | [.github/workflows/secure-supply-chain.yml](../.github/workflows/secure-supply-chain.yml), [devsecops_ci_admission.md](devsecops_ci_admission.md#deployment-contract-from-ci-output), [final_gap_closing_checklist.md](final_gap_closing_checklist.md#9-automate-kind-and-kyverno-bootstrap-for-reproducible-demo-open) | [../demo/evidence/20260414-210227/matrix-index.json](../demo/evidence/20260414-210227/matrix-index.json) | Partial |
| Reusability and thesis packaging | #11, #12 | Pending (`#11`) | [thesis_spec_en.md](thesis_spec_en.md#thesis-to-implementation-traceability), [final_gap_closing_checklist.md](final_gap_closing_checklist.md), [go_microservice_onboarding_guide.md](go_microservice_onboarding_guide.md) | [../demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md](../demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md), [lens_capture_checklist.md](lens_capture_checklist.md) | Partial |

## Remaining Fill-Ins Before Final Submission
- Insert one final `main` CI run URL for each objective row (or one canonical run URL reused across rows where valid).
- Add artifact IDs/download references from that CI run:
  - `dependency-integrity-report`
  - `govulncheck-report`
  - `sbom`
  - `grype-report`
  - `cosign-bundle`
- Confirm links are mirrored in `docs/thesis_spec_en.md` and `docs/final_gap_closing_checklist.md`.
