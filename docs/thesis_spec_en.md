# Thesis Specification (English)

## I. Topic Overview

### 1. Background and Problem Statement
Software supply chain attacks have become a high-impact threat category because risk no longer resides only in application source code, but also in dependencies, build systems, and deployment artifacts. Incidents such as SolarWinds showed that compromise at a single pipeline stage can propagate broadly across downstream systems. At the same time, software supply chain attack frequency continues to increase, making supply chain security a practical requirement in modern software delivery [2][5].

Policy and standards pressure has also increased. Executive Order 14028 and related guidance emphasize software component transparency and Software Bill of Materials (SBOM) practices as a core baseline for risk visibility and response [3].

### 2. Existing Standards, Frameworks, and Tooling
The market already provides mature building blocks. SLSA provides an integrity-oriented framework for artifact trust, provenance, and build traceability [4]. SBOM standards provide component and dependency transparency for inventory and vulnerability management [3]. In addition, vulnerability scanning, SBOM generation, image signing, and artifact verification can all be integrated into CI/CD pipelines [1][5].

The central challenge is no longer tool availability; it is practical end-to-end integration. Many teams still lack a unified, testable, and enforceable workflow that connects these controls from commit to deployment in real operational environments [1][2][4].

### 3. Current Practical State
Despite standards and policy guidance, real-world adoption is often fragmented. Organizations frequently run isolated controls (for example, vulnerability scanning only, SBOM generation only, or signing only) rather than implementing an end-to-end pipeline that can both verify and enforce trust before deployment [1][2][3][4]. For Kubernetes-based microservice systems, this gap is more visible due to rapid release cycles and weak integration patterns.

Therefore, the core problem addressed in this thesis is to design and validate an automated pipeline that unifies build, scan, SBOM, signing, attestation, and admission enforcement so that non-compliant artifacts are prevented from running in the cluster [1][2][4][5].

## II. Thesis Objectives

### 1. General Objective
Design and implement an automated software supply chain security pipeline for Go microservices on Kubernetes to:
- control dependencies,
- verify container image origin and integrity,
- assess security risk before release,
- and proactively reject artifacts that do not meet policy requirements.

The objective is not only risk detection, but operational prevention: unsafe artifacts must not be admitted into the cluster runtime [1][3][4][6]-[10].

### 2. Implementation Objectives and Novelty
The thesis targets a verifiable, enforceable, and repeatable model through the following technical objectives:

1. **Dependency transparency and control for Go services**
   Use Go Modules, `go.sum`, and checksum database concepts as dependency integrity foundations, and generate SBOM for every build to support traceability and software inventory [3][6][7].

2. **Automated pre-release risk elimination**
   Detect and block risk in dependencies/source and images by combining `govulncheck` for Go vulnerabilities with container vulnerability scanning. If risk exceeds threshold, the pipeline fails and artifact publication is blocked [1][5][6].

3. **Artifact signing and provenance attestation**
   Sign container images and issue provenance/attestation per build so each artifact can be linked to a valid commit and build process, reducing the chance of out-of-band or untrusted image promotion [1][4][8].

4. **Admission-time policy enforcement on Kubernetes**
   Enforce admission controls so Kubernetes rejects unsigned images, missing or invalid attestation/SBOM metadata, or artifacts that violate established security policy [8][9][10].

5. **Repeatable CI/CD supply-chain pipeline**
   Standardize an end-to-end sequence:
   1) Code,
   2) Build,
   3) SBOM,
   4) Vulnerability scan,
   5) Sign and attest,
   6) Push image,
   7) Verify and enforce at deployment,
   with reuse potential for other Go microservices [1][2][4][6][8][9][10].

6. **Demonstration on a real Go microservice sample at scale**
   Validate feasibility using a Kubernetes-deployed user-service as the canonical demo target and demonstrate the full commit-to-deploy trust flow. Scalability is validated across 23 Go microservices with a shared CI matrix, confirming the pipeline is reusable beyond the primary demo service.

**Novelty claim**: the thesis integrates commonly isolated controls into a single continuous control chain from dependency management through deployment admission. This shifts security posture from "scan and alert" to "verify and enforce" with reproducible artifacts and operating procedure outputs [1][2][4][6][8][9][10].

### Objective Conclusion
The expected result is a DevSecOps baseline for Go microservices on Kubernetes that is verifiable and enforceable, not only advisory [6][8][9][10].

## III. Scope of Work and Method
The work is implementation-driven (design, integration, experiment, and evaluation), not purely theoretical [1][2][4][6][8][9][10].

### 1. Requirement Analysis and Model Definition
Identify primary supply chain risks:
- third-party dependency risk,
- vulnerable images,
- untrusted artifact origin,
- unauthorized deployment.

Translate these into technical controls:
- SBOM generation,
- vulnerability scanning,
- image signing,
- provenance attestation,
- Kubernetes admission enforcement.

Model the end-to-end flow across four components:
1) Developer,
2) CI/CD,
3) Registry,
4) Kubernetes Admission [2][3][4][6][10].

### 2. Security Architecture Design
Design an integrated artifact lifecycle with the sequence:
1) Build,
2) SBOM,
3) Vulnerability Scan,
4) Sign and Attest,
5) Push,
6) Verify/Enforce at deployment.

Use open-source tooling to keep verification semantics consistent between CI/CD and Kubernetes admission [1][4][6][8][9][10].

### 3. Go Microservice Sample Construction
Use a practical Go user-service (registration, verification, login baseline) as the primary experiment target. Package with a multi-stage Dockerfile and least-privilege runtime hardening suitable for security-focused validation [1][6][7]. The monorepo has been scaled to 23 Go microservices — 10 core services (user-service, portfolio-service, order-service, risk-service, market-data-service, pricing-service, execution-service, settlement-service, compliance-service, notification-service) plus 13 extended services (apikey, kyc, watchlist, analytics, audit, fees, reporting, gateway, search, alert, data-feed, backtest, margin) — to validate repeatability of the supply-chain pipeline at scale; all services share the same Go `1.25.10` toolchain baseline and the complete list is maintained in `services.yaml`.

### 4. Secure CI/CD Pipeline Integration
Implement an automated pipeline executing:
1) Code,
2) Build,
3) SBOM,
4) Vulnerability Scan,
5) Sign and Attest,
6) Push.

Apply fail-fast behavior so policy violations stop the pipeline immediately and prevent non-compliant artifact publication [1][4][6][8].

### 5. Kubernetes Enforcement Setup
Provision a local validation cluster (Kind or Minikube) and configure admission controls at deployment time. Only signed images with valid attestation/metadata and policy compliance may run [8][9][10].

### 6. Validation Testing
Run two core scenario classes:
1) invalid artifact is denied,
2) valid artifact is admitted.

Evaluate with CI logs, artifact metadata, and admission events (`Denied` / `Allowed`) to verify correctness, automation depth, and risk-blocking capability [1][2][8][9][10].

### 7. Documentation Standardization and Synthesis
Consolidate architecture, pipeline, policy, and execution procedures into reproducible documentation, including explicit limitations and expansion paths [1][2][4].

### Method Conclusion
The method targets an implementable, verifiable, enforceable DevSecOps/supply-chain model across:
1) Developer,
2) CI/CD,
3) Registry,
4) Kubernetes Admission [1][4][8][9][10].

## IV. Evaluation Method
Evaluation is based on controlled experiments with valid and invalid artifact deployment scenarios, not only conceptual analysis.

Primary criteria:
1. **Pipeline effectiveness**: ability to detect and block risky artifacts during Build, Vulnerability Scan, and Sign/Attest stages [1][6][8].
2. **Integrity and origin verification**: ability to validate signatures, provenance/attestation, and image origin [4][8].
3. **Admission enforcement**: ability of Kubernetes admission to automatically reject violating workloads across registry/admission/deploy path [8][9][10].
4. **Repeatability**: ability to reuse the pipeline with minimal manual changes for other Go microservices [1][2][6].

Evaluation scope is intentionally bounded to feasibility and reproducibility of the model over full-scale enterprise compliance benchmarking [1][4][8][9][10].

## V. Study Limitations
1. **Business-domain limitation**
   The sample service is intentionally simple (user management, verification, login) and not a large enterprise workload.

2. **Security scope limitation**
   The focus is supply-chain controls across:
   1) Dependency,
   2) Build,
   3) Image,
   4) Admission.

   It does not deeply evaluate application-layer controls such as SQL injection defense, brute-force mitigation, rate limiting, or business logic abuse.

3. **Environment limitation**
   Validation is performed on local Kubernetes clusters (Kind/Minikube), not mandatory production cloud deployment.

4. **Compliance-depth limitation**
   The design adopts practical SLSA/SBOM/signing/provenance principles, but does not claim full enterprise-grade compliance maturity levels [2][4].

5. **Performance limitation**
   The primary goal is correctness, automation, and enforcement; CI/CD performance optimization and system load benchmarking are outside scope.

## VI. Results and Findings

This section summarizes what the implemented pipeline demonstrably achieves, evaluated against the implementation objectives in Section II.2 and verified on the latest green full cross-OS CI run [`26732257799`](https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/26732257799) (23 services × ubuntu + macOS, all jobs success).

**1. Achieved capabilities (verify-and-enforce chain).**
- **Dependency transparency (Obj. II.2.1 — Implemented):** every build emits an SPDX SBOM (`anchore/sbom-action`) attested to the image via `cosign attest --type spdxjson`; module integrity is gated by `go mod verify` / readonly resolution / `go mod tidy -diff`.
- **Pre-release risk elimination (Obj. II.2.2 — Implemented):** `govulncheck` (fail-fast) plus Grype with a threshold gate on *fixable* High/Critical findings block publication in `enforced` mode, validated across all 23 services.
- **Signing & provenance (Obj. II.2.3 — Implemented):** keyless Cosign signing (Sigstore/Fulcio/Rekor, GitHub OIDC) on every image; SLSA **Build L3** provenance generated for the flagship user-service via the official `slsa-github-generator` and independently confirmed by `slsa-verifier`.
- **Admission enforcement (Obj. II.2.4 — Implemented):** Kyverno `verify-images`, `require-sbom`, and `cve-threshold` policies. The fixed admission matrix passes all cases — `VALID_ALLOW` admitted; `NEG_UNSIGNED_DENY`, `NEG_MISSING_SBOM_DENY`, `NEG_CVE_THRESHOLD_DENY` denied (evidence: `demo/evidence/20260414-210227/matrix-summary.md`).
- **Repeatable pipeline & scale (Obj. II.2.5/II.2.6 — Implemented):** a single shared CI matrix drives all 23 Go microservices from commit to signed/attested artifact on a common Go 1.25.10 baseline, demonstrating reuse well beyond the primary demo service.

**2. Reproducibility finding.** Cross-architecture validation (ubuntu amd64 + macOS arm64) surfaced an architecture-dependent floating-point divergence in analytics drawdown computation that produced non-deterministic results on arm64. Resolving it at the source (commit `7930650`) hardened the pipeline's reproducibility guarantee — a practical illustration that "reproducible build verification" must hold across runner architectures, not a single platform.

**3. Objective-closure status.** Of the eleven tracked GitHub objectives, six are closed with linked CI evidence (#3 SBOM, #4 Grype gate, #5 Cosign, #6 SLSA, #7 Kyverno verify, #8 SBOM/CVE admission). Three remain open as scoped engineering gaps rather than design gaps: per-service runtime hardening (#1 — 9/23 Dockerfiles distroless, 1/23 manifests carry a full `securityContext`), base-manifest alignment across all services (#2), and a clean Kind-bootstrap evidence capture on a host with the `kind` CLI (#9). The traceability objective (#11) stays open pending those closures.

## VII. Future Work and Recommendations

1. **Repo-wide runtime hardening.** Propagate the user-service hardening baseline (`runAsNonRoot`, `readOnlyRootFilesystem`, dropped capabilities, distroless base) to all 23 services' Dockerfiles and Kubernetes base manifests, closing #1 and #2.
2. **Reproducible cluster evidence.** Capture a clean `kind` bootstrap + admission-matrix run on a host with the `kind` CLI to replace the docker-desktop evidence environment, closing #9; alternatively, formally adopt docker-desktop as canonical and document the rationale.
3. **Negative-path automation.** Add an automated controlled-failure scenario (a deliberately vulnerable image) so the threshold gate's *block* behavior is evidenced alongside the pass-case (follow-up to #4).
4. **Uniform SLSA L3.** Extend the official SLSA Build L3 generator path from the flagship user-service to all services, rather than the lighter in-pipeline provenance attestation currently used for the non-flagship set.
5. **Policy-engine breadth.** Evaluate Sigstore Policy Controller as a complement/alternative to Kyverno for signature/attestation verification, and compare enforcement semantics.
6. **Beyond local validation.** Promote the model from a local single-node cluster to a multi-environment / managed-Kubernetes setting to study performance, key management (KMS-backed vs keyless), and multi-tenant trust boundaries.
7. **CI maintenance & guardrails.** Complete the Node-24 action migration (the evidence-pipeline `upload/download-artifact` majors deferred for stability) and add a lint/guardrail to catch cross-architecture numeric fragility before it reaches the nightly matrix.

## Thesis-to-Implementation Traceability
As-of `2026-06-01` (GitHub state snapshot: issues `#1` to `#9`, `#11` open; `#10`, `#12`, `#13`, `#14` closed). Scale: 23 services, Go `1.25.10` baseline, CI green on `main` (full cross-OS matrix run [26732257799](https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/26732257799), commit `7930650`).

Status legend: `Implemented | Partial | Missing`

Detailed objective-to-evidence register: `docs/traceability_evidence_register.md`.

| Objective Item | Control/Mechanism | Evidence Artifact or Log (Current) | Related Issue(s) | Status | Open Gap to Close |
|---|---|---|---|---|---|
| Dependency transparency and control | Go modules (`go.mod`, `go.sum`), checksum rationale, CI integrity checks (`go mod verify`, readonly graph, `go mod tidy -diff` audit), SBOM generation | `go.mod`, `go.sum`, `.github/workflows/ci-service.yml`, `docs/go_dependency_integrity_baseline.md`; CI run: [runs/26732257799](https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/26732257799); SBOM artifact: `...user-service_79306504f27c.spdx.json=7322461363` | #3, #14 | Implemented | SBOM generated and published per service; digest-to-annotation mapping enforced via Kyverno require-sbom policy. |
| Go vulnerability and image risk gating | `govulncheck` + Grype fail-fast threshold on fixable High/Critical findings | `.github/workflows/ci-service.yml`, `docs/devsecops_ci_admission.md` (fail-fast contract); pass-case: [runs/26732257799](https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/26732257799) (`user-service-grype-report=7322465750`, `user-service-security-gate-findings=7322465817`) | #4, #13 | Partial | Add explicit fail-case evidence (controlled simulation where gate blocks). |
| Image signing and provenance attestation | Cosign signing + provenance attestation and verification path | `.github/workflows/ci-service.yml`, `infra/policies/kyverno/clusterpolicy-verify-images.yaml`, `docs/devsecops_ci_admission.md`; SLSA L3 evidence: `user-service-slsa-l3-digest=7322475856`, `slsa-l3-verifier-evidence=7322602265` | #5, #6, #7 | Implemented | Keyless Cosign signing + SLSA L3 provenance generated and verified in-CI (provenance-l3 + verify-slsa-l3 jobs); Kyverno verify-images enforces at admission. |
| Admission enforcement of trust controls | Kyverno verifyImages + CVE/SBOM annotation policies + automated deny/allow matrix | `infra/policies/kyverno/*`, `infra/scripts/admission_matrix_demo.ps1`, `demo/evidence/20260414-210227/matrix-summary.md`, `docs/lens_capture_checklist.md` | #7, #8, #10 | Partial | Keep deny-message stability validation and finalize issue closure records on GitHub. |
| End-to-end reproducible pipeline | CI pipeline from test/govulncheck to push + deployment annotation overlay + matrix evidence export; 23-service scalability validated | `.github/workflows/ci-service.yml`, `services/user-service/deploy/kubernetes/overlays/ci`, `docs/devsecops_ci_admission.md`, `demo/evidence/20260414-210227/matrix-index.json`; CI: [runs/26732257799](https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/actions/runs/26732257799) | #3, #4, #9 | Implemented | Clean Kind bootstrap + admission matrix captured 2026-06-01 (`demo/evidence/20260601-kind-bootstrap/`). |
| Reusability and thesis packaging | Traceability matrix + reusable onboarding playbook for additional Go services; 23-service CI matrix demonstrates repeatability at scale | `docs/implementation_roadmap.md`, `docs/go_microservice_onboarding_guide.md`, `docs/final_gap_closing_checklist.md`, `demo/evidence/20260414-213541-onboarding-second-service/onboarding-summary.md` | #11, #12 | Partial | Close open issue trail with final evidence links after Kind bootstrap run. |

## References
[1] D. Patel, "Software supply chain security: Implementing SLSA compliance in CI/CD pipelines," International Journal for Research Trends and Innovation, vol. 10, no. 7, Jan. 2025.

[2] M. Tamanna, S. Hamer, M. Tran, S. Fahl, Y. Acar, and L. Williams, "Analyzing challenges in deployment of the SLSA framework for software supply chain security," Dec. 2024.

[3] National Institute of Standards and Technology (NIST), "Improving the Nation's Cybersecurity: NIST's Responsibilities under Executive Order 14028-Software Supply Chain Security Guidance," U.S. Department of Commerce, July 2022.

[4] The Linux Foundation, "Safeguarding artifact integrity across any software supply chain: What is SLSA?," Open Source Security Foundation, 2025.

[5] D. I. Jonathan, "Supply chain security in modern software: SBOMs, SLSA, and beyond," EM360Tech, Sept. 3, 2025.

[6] J. Qiu, "Vulnerability Management for Go," The Go Blog, Sep. 6, 2022.

[7] K. Hockman, "Module Mirror and Checksum Database Launched," The Go Blog, Aug. 29, 2019.

[8] Sigstore, "Cosign Quickstart," Sigstore Documentation, accessed Mar. 18, 2026.

[9] Sigstore, "Policy Controller Overview," Sigstore Documentation, accessed Mar. 18, 2026.

[10] Kubernetes, "Validating Admission Policy," Kubernetes Documentation, 2024.


