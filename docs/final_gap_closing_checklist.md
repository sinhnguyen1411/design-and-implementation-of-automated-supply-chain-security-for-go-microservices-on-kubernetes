# Final Gap-Closing Checklist (Issue-by-Issue)

Last updated: 2026-04-14

Repository: `sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes`

GitHub issue snapshot:
- Open: `#1` to `#12`
- Closed: `#13`, `#14`

This checklist is the execution contract to close the thesis package with auditable evidence.

## Global Exit Criteria
- [ ] All milestone issues `#1` to `#14` are closed on GitHub.
- [x] One canonical evidence source is defined and used consistently in docs/scripts.
- [ ] Traceability status in thesis docs matches real issue status and evidence links.
- [x] Final demo run has `VALID_ALLOW`, 3 deny cases, and `VALID_ALLOW_RECHECK` all `PASS`.

## Direct Mapping to Proposal I-V (As of 2026-04-14)
This is the primary final-package view for thesis closure.

Status legend (allowed values only): `Implemented | Partial | Missing`

### Section I - Topic Overview
| Proposal Item | Status | What Exists | What Is Missing | Evidence Links | Blocking Issue(s) |
|---|---|---|---|---|---|
| I.1 Background and problem statement | Implemented | Problem framing and attack motivation are documented and thesis-aligned. | None for this item. | [docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet (1)_readable.md](docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet%20(1)_readable.md#i-tong-quan-de-tai), [docs/thesis_spec_en.md](docs/thesis_spec_en.md#i-topic-overview) | - |
| I.2 Existing standards/frameworks/tooling | Implemented | SLSA, SBOM, Cosign, admission control positioning is documented in thesis and CI/admission flow docs. | None for this item. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#2-existing-standards-frameworks-and-tooling), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md#ci-workflow-githubworkflowssecure-supply-chainyml) | - |
| I.3 Current practical state | Implemented | Fragmented real-world state and end-to-end integration gap are explicitly stated. | None for this item. | [docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet (1)_readable.md](docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet%20(1)_readable.md#3-thuc-trang-truoc-thoi-diem-hien-tai), [docs/thesis_spec_en.md](docs/thesis_spec_en.md#3-current-practical-state) | - |

### Section II - Thesis Objectives
| Proposal Item | Status | What Exists | What Is Missing | Evidence Links | Blocking Issue(s) |
|---|---|---|---|---|---|
| II.1 General objective | Partial | End-to-end CI plus admission baseline exists with defined deny/allow matrix. | Final objective closure still depends on open implementation issues and synchronized evidence links. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#1-general-objective), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md#policy-contract-admission-time-requirements) | #11 |
| II.2.1 Dependency transparency/control | Partial | Go module integrity controls and SBOM generation are implemented in workflow and documented. | Need final thesis package proof-chain from real CI artifact to deploy annotation in one canonical place. | [.github/workflows/secure-supply-chain.yml](.github/workflows/secure-supply-chain.yml), [docs/go_dependency_integrity_baseline.md](docs/go_dependency_integrity_baseline.md), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#3-generate-spdx-sbom-in-ci-and-publish-artifacts-open) | #3 |
| II.2.2 Automated risk elimination (govulncheck + image scan) | Partial | govulncheck and Grype fail-fast gates are present. | Need explicit final-package fail-case and pass-case CI evidence references. | [.github/workflows/secure-supply-chain.yml](.github/workflows/secure-supply-chain.yml), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md#fail-fast-behavior), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#4-add-grype-scan-with-threshold-based-security-gate-open) | #4 |
| II.2.3 Signing + provenance attestation | Partial | Keyless Cosign sign plus SLSA-style attestation and verify steps are implemented in CI. | Need final, stable evidence chain linking one canonical CI run and acceptance narrative in thesis package. | [.github/workflows/secure-supply-chain.yml](.github/workflows/secure-supply-chain.yml), [infra/policies/kyverno/clusterpolicy-verify-images.yaml](infra/policies/kyverno/clusterpolicy-verify-images.yaml), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#5-integrate-keyless-cosign-signing-for-image-digests-open) | #5, #6 |
| II.2.4 Admission enforcement in Kubernetes | Partial | Kyverno policies enforce signature/attestation plus SBOM/CVE metadata checks. | Need finalized deterministic policy/evidence contract for thesis sign-off. | [infra/policies/kyverno/clusterpolicy-verify-images.yaml](infra/policies/kyverno/clusterpolicy-verify-images.yaml), [infra/policies/kyverno/clusterpolicy-require-sbom.yaml](infra/policies/kyverno/clusterpolicy-require-sbom.yaml), [infra/policies/kyverno/clusterpolicy-cve-threshold.yaml](infra/policies/kyverno/clusterpolicy-cve-threshold.yaml) | #7, #8 |
| II.2.5 Repeatable CI/CD pipeline | Partial | Pipeline stages are standardized and CI overlay rendering exists. | Deployment-consumer contract and bootstrap reproducibility still require closure. | [.github/workflows/secure-supply-chain.yml](.github/workflows/secure-supply-chain.yml), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md), [infra/scripts/devsecops_kind_bootstrap.sh](infra/scripts/devsecops_kind_bootstrap.sh) | #3, #8, #9 |
| II.2.6 Demonstration on real Go microservice | Partial | `user-service` demo matrix with allow/deny/regression is available from canonical run bundles under `demo/evidence`, with screenshot appendix pack completed. | Need final issue-state synchronization and GitHub closure trail for thesis package finalization. | [infra/scripts/admission_matrix_demo.ps1](infra/scripts/admission_matrix_demo.ps1), [demo/evidence/20260414-210227/matrix-summary.md](demo/evidence/20260414-210227/matrix-summary.md), [docs/lens_capture_checklist.md](docs/lens_capture_checklist.md) | #11 |
| II.3 Improvement over current state | Partial | Model direction is clearly from scan-and-alert to verify-and-enforce. | Need final synchronized traceability status and evidence references that match live issue state. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#2-implementation-objectives-and-novelty), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#11-build-thesis-traceability-matrix-objective-to-evidence-open) | #11 |

### Section III - Scope of Work and Method
| Proposal Item | Status | What Exists | What Is Missing | Evidence Links | Blocking Issue(s) |
|---|---|---|---|---|---|
| III.1 Requirement analysis and model definition | Implemented | End-to-end risk/control model is specified across Dev -> CI/CD -> Registry -> Admission. | None for this item. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#1-requirement-analysis-and-model-definition), [docs/scs_architecture_diagram.html](docs/scs_architecture_diagram.html) | - |
| III.2 Security architecture design | Implemented | Architecture sequence Build -> SBOM -> Scan -> Sign/Attest -> Push -> Verify/Enforce is documented. | None for this item. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#2-security-architecture-design), [README.md](README.md#architecture-overview) | - |
| III.3 Go sample microservice construction | Partial | `user-service` implementation, Docker packaging, and least-privilege runtime are present; Go toolchain baseline is aligned (`1.25.9`). | Need final hardening rationale and runtime-evidence packaging for thesis closure. | [services/user-service/Dockerfile](services/user-service/Dockerfile), [services/user-service/deploy/kubernetes/base/deployment.yaml](services/user-service/deploy/kubernetes/base/deployment.yaml), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#1-harden-dockerfile-and-runtime-defaults-for-least-privilege-open) | #1 |
| III.4 Secure CI/CD pipeline integration | Partial | CI stages and fail-fast controls are implemented end-to-end. | Final thesis packaging still needs explicit reproducible artifact mapping in docs for SBOM/scan gates. | [.github/workflows/secure-supply-chain.yml](.github/workflows/secure-supply-chain.yml), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md#ci-evidence-artifacts) | #3, #4 |
| III.5 Kubernetes enforcement setup | Partial | Kyverno bootstrap and policy set are implemented. | Base-manifest contract plus policy/enforcement alignment and reproducible bootstrap behavior are not fully closed. | [infra/scripts/devsecops_kind_bootstrap.sh](infra/scripts/devsecops_kind_bootstrap.sh), [infra/policies/kyverno/kustomization.yaml](infra/policies/kyverno/kustomization.yaml), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#2-align-kubernetes-base-manifests-with-production-safe-defaults-open) | #2, #7, #8, #9 |
| III.6 Validation testing | Partial | Admission matrix script exports structured per-case evidence and regression re-check; canonical real run and screenshot appendix are available. | Need clean Kind bootstrap evidence and final thesis cross-reference locking. | [infra/scripts/admission_matrix_demo.ps1](infra/scripts/admission_matrix_demo.ps1), [demo/evidence/20260414-210227/matrix-index.json](demo/evidence/20260414-210227/matrix-index.json), [docs/lens_screenshots/README.md](docs/lens_screenshots/README.md) | #9, #11 |
| III.7 Documentation standardization and synthesis | Partial | Roadmap, runbook-style docs, traceability section, and onboarding guide are available. | Traceability status/evidence links must be synchronized with open issue reality before final thesis package sign-off. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#thesis-to-implementation-traceability), [docs/implementation_roadmap.md](docs/implementation_roadmap.md), [docs/go_microservice_onboarding_guide.md](docs/go_microservice_onboarding_guide.md) | #11 |

### Section IV - Evaluation Method
| Proposal Item | Status | What Exists | What Is Missing | Evidence Links | Blocking Issue(s) |
|---|---|---|---|---|---|
| IV.1 Pipeline effectiveness (detect/block in build-scan-sign path) | Partial | Workflow enforces fail-fast logic and exports security artifacts. | Final package still needs explicit both-direction gate evidence references (fail and pass) in one canonical thesis section. | [.github/workflows/secure-supply-chain.yml](.github/workflows/secure-supply-chain.yml), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md#fail-fast-behavior), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#4-add-grype-scan-with-threshold-based-security-gate-open) | #4 |
| IV.2 Integrity and origin verification | Partial | CI has keyless signature verification and attestation verification; admission has keyless trust policy. | Need final evidence package linkage from a canonical CI run to thesis evaluation section. | [.github/workflows/secure-supply-chain.yml](.github/workflows/secure-supply-chain.yml), [infra/policies/kyverno/clusterpolicy-verify-images.yaml](infra/policies/kyverno/clusterpolicy-verify-images.yaml), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md#manual-cryptographic-verification-optional) | #5, #6 |
| IV.3 Admission enforcement effectiveness | Partial | Deny/allow evidence exists for unsigned, missing SBOM, and CVE threshold cases, plus recheck allow in canonical run bundle with screenshot index. | Need final trust-model/policy closure consistency and issue board synchronization. | [demo/evidence/20260414-210227/matrix-summary.md](demo/evidence/20260414-210227/matrix-summary.md), [docs/lens_capture_checklist.md](docs/lens_capture_checklist.md), [docs/lens_screenshots/README.md](docs/lens_screenshots/README.md) | #7, #11 |
| IV.4 Repeatability for other Go microservices | Partial | Reuse guide and second-service dry-run evidence now exist with captured deviations and remediation notes. | Need final issue closure trace on GitHub and final thesis cross-link lock. | [docs/go_microservice_onboarding_guide.md](docs/go_microservice_onboarding_guide.md), [demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md](demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md), [docs/final_gap_closing_checklist.md](docs/final_gap_closing_checklist.md#12-publish-reusable-onboarding-guide-for-additional-go-microservices-ready-to-close) | #12, #11 |

### Section V - Study Limitations
| Proposal Item | Status | What Exists | What Is Missing | Evidence Links | Blocking Issue(s) |
|---|---|---|---|---|---|
| V.1 Business-domain limitation | Implemented | Scope is constrained to a simple user-service domain. | None for this item. | [docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet (1)_readable.md](docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet%20(1)_readable.md#v-gioi-han-cua-de-tai), [docs/thesis_spec_en.md](docs/thesis_spec_en.md#v-study-limitations) | - |
| V.2 Security scope limitation | Implemented | Focus is supply-chain controls, not full application-layer security hardening. | None for this item. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#2-security-scope-limitation), [README.md](README.md#notes) | - |
| V.3 Environment limitation | Implemented | Local-cluster validation scope is explicitly documented (Kind/Minikube/Docker Desktop). | None for this item. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#3-environment-limitation), [docs/devsecops_ci_admission.md](docs/devsecops_ci_admission.md#local-bootstrap-kind) | - |
| V.4 Compliance-depth limitation | Implemented | Practical adoption of SLSA/SBOM/signing/provenance is documented without full enterprise-level claim. | None for this item. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#4-compliance-depth-limitation), [docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet (1)_readable.md](docs/Phu_luc_de_cuong_tot_nghiep_chi_tiet%20(1)_readable.md#v-gioi-han-cua-de-tai) | - |
| V.5 Performance limitation | Implemented | Evaluation emphasis is correctness/automation/enforcement rather than CI performance benchmarking. | None for this item. | [docs/thesis_spec_en.md](docs/thesis_spec_en.md#5-performance-limitation), [docs/demo_evidence.md](docs/demo_evidence.md#conclusion) | - |

### Issue <-> Proposal Crosswalk
| Issue | Proposal Item(s) |
|---|---|
| #1 | III.3 |
| #2 | III.5 |
| #3 | II.2.1, III.4 |
| #4 | II.2.2, IV.1 |
| #5 | II.2.3, IV.2 |
| #6 | II.2.3, IV.2 |
| #7 | II.2.4, III.5, IV.3 |
| #8 | II.2.4, III.5 |
| #9 | III.5, III.6 |
| #10 | III.6, IV.3 |
| #11 | III.7 |
| #12 | IV.4 |

## Issue-by-Issue Closure Plan

### #1 Harden Dockerfile and runtime defaults for least privilege (Open)
Current state:
- Hardened multi-stage build and distroless runtime are in place.
- Runtime security settings are in Kubernetes deployment (`runAsNonRoot`, `readOnlyRootFilesystem`, dropped caps).

Gap to close:
- Build/runtime hardening rationale is not centralized in a final thesis-ready checklist.
- Need cluster-backed runtime evidence snippet showing restricted pod security context in active deployment.

Checklist:
- [x] Align Go version across `Dockerfile` and `.github/workflows/secure-supply-chain.yml`.
- [ ] Add a short hardening rationale section in docs with residual risks and non-goals.
- [ ] Capture evidence: `docker build` success + deployment running with restricted security context.

Close when:
- [ ] Issue comment links to updated docs and command outputs proving least-privilege runtime.

### #2 Align Kubernetes base manifests with production-safe defaults (Open)
Current state:
- Base manifest includes strong security defaults and probes.
- Base deployment annotations are intentionally empty by design (`security.grype.io/high_critical`, `security.stock-trading.dev/sbom-digest`) so base-only deploy is non-compliant when admission policies are enabled.
- Cluster-backed evidence captured: `demo/evidence/20260414-204808-contract-base-overlay` (`02_apply_base.txt`, `02_kyverno_logs.txt`, `contract-summary.md`).

Gap to close:
- Need cluster-backed evidence showing deterministic behavior: base-only denied, CI overlay deployment allowed.

Checklist:
- [x] Replace raw `TBD` with documented contract strategy.
- [x] Option A selected: keep base annotations empty and require CI overlay for compliant deploy.
- [x] Update docs to explain compliant path (`overlay/ci`) vs baseline path.
- [x] Validate `kubectl apply -k deploy/kubernetes/base` behavior and expected policy result in a policy-enabled cluster run.

Close when:
- [x] Manifest behavior is deterministic and explicitly documented.

### #3 Generate SPDX SBOM in CI and publish artifacts (Open)
Current state:
- SBOM generation and artifact upload exist in workflow.

Gap to close:
- Need explicit proof chain from CI run to deployment metadata digest in final package.

Checklist:
- [ ] Add retrieval instructions for SBOM artifact from a real CI run.
- [ ] Add explicit command snippet to recompute and compare SBOM digest against deployment annotation.
- [ ] Store one reference CI run URL and artifact IDs in docs.

Close when:
- [ ] A reviewer can reproduce SBOM digest mapping from CI artifact to deployment annotation.

### #4 Add Grype scan with threshold-based security gate (Open)
Current state:
- Grype JSON report and security gate are implemented.
- `govulncheck` gate is integrated.

Gap to close:
- Need explicit fail-case evidence in final package, not only pass-case logic.

Checklist:
- [ ] Capture one CI run (or controlled local simulation) where threshold gate fails.
- [ ] Capture one run where gate passes and shows `high_critical=0`.
- [ ] Add both run references to docs (`pass` and `fail`) with concise explanation.

Close when:
- [ ] Gate behavior is evidenced in both directions (deny and allow).

### #5 Integrate keyless Cosign signing for image digests (Open)
Current state:
- CI performs keyless `cosign sign` for GHCR mode.

Gap to close:
- Admission policy trust model currently relies on static public key verification.
- Keyless signing and runtime verification policy are not aligned end-to-end.

Checklist:
- [ ] Decide one trust model and document it:
- [ ] Model A: full keyless (`certificate identity/issuer` policy).
- [ ] Model B: key-pair signing in CI with managed key distribution.
- [ ] Update Kyverno verify policy to match chosen model.
- [ ] Add verification commands and expected output in docs for chosen model.

Close when:
- [ ] Signing method in CI and admission verification method are cryptographically consistent.

### #6 Generate and verify SLSA provenance attestation (Open)
Current state:
- CI generates provenance predicate and attests image.

Gap to close:
- CI does not yet enforce `cosign verify-attestation` as a post-attest check.

Checklist:
- [ ] Add `cosign verify-attestation` step in CI for release path.
- [ ] Upload verification output as artifact/log evidence.
- [ ] Link verification output in thesis evidence docs.

Close when:
- [ ] Attestation is both generated and verified in the same CI trust path.

### #7 Enforce signature and provenance checks via Kyverno (Open)
Current state:
- Kyverno verify policy exists for `user-service` image path with key-based attestor.

Gap to close:
- Policy does not clearly enforce the same identity semantics as CI keyless flow.

Checklist:
- [ ] Refactor verify policy according to chosen trust model from issue `#5`.
- [ ] Validate deny scenarios:
- [ ] unsigned image,
- [ ] wrong signer/identity,
- [ ] missing or invalid provenance.
- [ ] Capture Kyverno event/log snippets for each deny reason.

Close when:
- [ ] Policy blocks all trust violations with deterministic deny messages.

### #8 Enforce SBOM and CVE metadata requirements in admission (Open)
Current state:
- `security.grype.io/high_critical` and `security.stock-trading.dev/sbom-digest` are enforced.
- Deployment contract is now documented: CI produces `services/<name>/deploy/kubernetes/overlays/ci/*` in artifact bundle and deployment consumers apply that overlay.
- Overlay-based allow path has refreshed evidence (`demo/evidence/20260414-204808-contract-base-overlay/04_apply_overlay_allow.txt`, `04_describe_deploy.txt`).

Gap to close:
- Need one refreshed cluster-backed run proving contract stability (deny/allow messages unchanged after contract cleanup).

Checklist:
- [x] Define canonical deployment contract:
- [x] where overlay is produced,
- [x] how it is consumed,
- [x] how values are audited.
- [x] Add one command flow from CI output to `kubectl apply -k ...` with expected annotations.
- [ ] Validate deny messages remain stable after contract cleanup.

Close when:
- [ ] Annotation contract is stable, documented, and reproducible.

### #9 Automate Kind and Kyverno bootstrap for reproducible demo (Open)
Current state:
- Bootstrap script is idempotent and supports dynamic Kyverno deployment-name readiness checks.
- Teardown/reset script exists for clean reruns (`infra/scripts/devsecops_kind_reset.sh`).
- Fallback evidence loop was executed on `docker-desktop` because `kind` CLI is unavailable on current host (`demo/evidence/20260414-204808-contract-base-overlay/contract-summary.md`).

Gap to close:
- Need one clean-environment **Kind** bootstrap + rerun evidence log captured in docs/evidence bundle (still pending on a host with `kind` installed).

Checklist:
- [x] Harmonize bootstrap readiness checks with currently installed Kyverno deployment names.
- [x] Add teardown/reset script path and document idempotent rerun behavior.
- [ ] Validate bootstrap on a clean environment and capture output log.

Close when:
- [ ] Clean bootstrap + rerun works without manual patching.

### #10 Add adversarial deployment scenarios and evidence checklist (Ready to close)
Current state:
- Matrix script and multiple evidence runs exist.
- Canonical evidence root is `demo/evidence` (real cluster-backed runs).
- Fresh canonical matrix run captured: `demo/evidence/20260414-210227` with all cases `PASS`.
- Dashboard supports Actions snapshot as primary source (`docs/security-admission-dashboard/data/actions-runs.snapshot.json`), with fallback bundled datasets under `docs/security-admission-dashboard/demo-data/evidence/`.

Gap to close:
- Need GitHub issue closure trail and final cross-doc linkage after appendix pack completion.

Checklist:
- [x] Declare one canonical evidence root path and update scripts/docs/dashboard accordingly.
- [x] Regenerate one canonical, real cluster-backed run-id for final thesis appendix.
- [x] Separate clearly:
- [x] real evidence set for thesis claims,
- [x] bundled synthetic demo set for offline UI.
- [x] Complete screenshot checklist and publish screenshot-evidence index (`docs/lens_capture_checklist.md`, `docs/lens_screenshots/README.md`).

Close when:
- [x] Final evidence bundle is canonical, consistent, and fully reproducible.

### #11 Build thesis traceability matrix (objective to evidence) (Open)
Current state:
- Traceability matrix exists in `docs/thesis_spec_en.md`.
- Traceability section now includes an explicit as-of marker and issue-state-aligned statuses (`Partial` where related issues remain open).
- Objective-level evidence register is now documented in `docs/traceability_evidence_register.md`.

Gap to close:
- Need final CI run URL(s) and artifact IDs inserted per objective row before final thesis submission.

Checklist:
- [x] Update traceability status to match real issue states.
- [ ] For each objective row, add direct evidence links:
- [ ] CI run URL,
- [x] artifact path,
- [x] admission evidence run-id.
- [x] Add a final "as-of date" marker in traceability section.

Close when:
- [ ] Traceability table is synchronized with issue tracker and evidence links.

### #12 Publish reusable onboarding guide for additional Go microservices (Ready to close)
Current state:
- Reuse guide is present and detailed.
- Second-service simulation evidence is captured in `demo/evidence/20260414-213541-onboarding-second-service/`.

Gap to close:
- Need final GitHub closure trail tying simulation evidence to thesis package references.

Checklist:
- [x] Run a dry-run onboarding simulation using a second sample service namespace/image naming convention.
- [x] Capture deviations and patch guide steps accordingly.
- [x] Add a short "known assumptions/limitations" block for adopters.

Close when:
- [x] A second-service simulation follows the guide with no undocumented blockers.

### #13 Integrate govulncheck into fail-fast CI security gates (Closed)
Current state:
- Closed and implemented.

Keep-closed checklist:
- [ ] Ensure future workflow edits keep govulncheck as hard gate.
- [ ] Keep report artifact upload and docs link intact.

### #14 Document and enforce Go dependency integrity baseline (Closed)
Current state:
- Closed and implemented.

Keep-closed checklist:
- [ ] Preserve dependency-integrity hard checks in CI.
- [ ] Keep `go mod tidy -diff` as audit signal unless policy is intentionally changed.

## Suggested Execution Order (To Minimize Rework)
1. `#5` + `#6` + `#7` (cryptographic trust model alignment)
2. `#8` + `#2` + `#9` (deployment contract and reproducible cluster flow)
3. `#11` (final thesis packaging linkage with CI run URLs/artifact IDs)
4. `#1` + `#3` + `#4` polishing (consistency, fail/pass evidence completeness)

## Final Sign-off Checklist
- [ ] One final CI run on `main` captured with full artifact set.
- [ ] One final admission matrix run-id captured from real cluster.
- [ ] Docs cross-links updated (`README`, `thesis_spec_en`, runbook, evidence docs).
- [ ] GitHub issues `#1` to `#12` closed with evidence links.
