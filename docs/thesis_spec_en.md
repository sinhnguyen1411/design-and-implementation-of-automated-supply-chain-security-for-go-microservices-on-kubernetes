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

6. **Demonstration on a real Go microservice sample**
   Validate feasibility using a Kubernetes-deployed user-service and demonstrate the full commit-to-deploy trust flow.

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
Use a practical Go user-service (registration, verification, login baseline) as the experiment target. Package with a multi-stage Dockerfile and least-privilege runtime hardening suitable for security-focused validation [1][6][7].

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

## Thesis-to-Implementation Traceability
| Objective Item | Control/Mechanism | Evidence Artifact or Log | Related Issue(s) | Status |
|---|---|---|---|---|
| Dependency transparency and control | Go modules (`go.mod`, `go.sum`), checksum rationale, CI integrity checks (`go mod verify`, readonly graph, `go mod tidy -diff` audit), SBOM generation | `go.mod`, `go.sum`, `dependency-integrity-report.txt`, `sbom.spdx.json`, `docs/devsecops_ci_admission.md`, `docs/go_dependency_integrity_baseline.md` | #3, #12, #14 | Implemented |
| Go vulnerability and image risk gating | `govulncheck` + Grype fail-fast threshold on fixable High/Critical findings | CI logs, `govulncheck-report.txt`, `grype-report.json`, workflow outputs | #4, #13 | Implemented |
| Image signing and provenance attestation | Cosign keyless sign + SLSA-style attestation | `.github/workflows/secure-supply-chain.yml`, `provenance.json`, `docs/demo_evidence.md` | #5, #6 | Implemented |
| Admission enforcement of trust controls | Kyverno verifyImages + CVE/SBOM annotation policies + automated deny/allow matrix | `deploy/policies/kyverno/*`, `scripts/admission_matrix_demo.ps1`, `docs/demo_evidence.md` | #7, #8, #10 | Implemented |
| End-to-end reproducible pipeline | CI pipeline from test/govulncheck to push + deployment annotation overlay + matrix evidence export | `.github/workflows/secure-supply-chain.yml`, `deploy/kubernetes/overlays/ci` artifact, `scripts/admission_matrix_demo.ps1`, `docs/devsecops_ci_admission.md` | #3, #4, #9, #13 | Implemented |
| Reusability and thesis packaging | Traceability matrix + reusable onboarding playbook for additional Go services | `docs/implementation_roadmap.md`, `docs/go_microservice_onboarding_guide.md`, thesis docs | #11, #12 | Implemented |

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

