# Implementation Roadmap (Milestones and Issues)

This roadmap aligns repository execution work with the thesis specification in `docs/thesis_spec_en.md`.

## Tracking Links
- Milestones: https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestones
- Roadmap issues label: https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/issues?q=is%3Aissue%20is%3Aopen%20label%3Athesis-roadmap
- Thesis traceability table: `docs/thesis_spec_en.md` -> section "Thesis-to-Implementation Traceability"

## Thesis Mapping
| Thesis Section | Focus | Roadmap Coverage |
|---|---|---|
| I. Topic Overview | Problem framing and motivation | M1-M6 context and evidence trail |
| II. Objectives and Novelty | Verifiable and enforceable supply-chain pipeline | M2-M4 core controls + M6 packaging |
| III. Method and Implementation | Build, SBOM, scan, sign, attest, enforce | M1-M5 technical execution |
| IV. Evaluation Method | Valid/invalid artifact experiments | M5 evidence scenarios |
| V. Limitations | Practical boundaries and constraints | M6 documentation and reproducibility notes |

## Recommended Local Development Stack
The recommended local environment for this thesis is `Docker Desktop + kind + kubectl + Helm + Lens`.

This combination is appropriate for an implementation-focused supply chain security study because it balances reproducibility, operational realism, and demo efficiency on a single workstation. `Docker Desktop` provides the local container runtime baseline required to build and test Go microservice images. `kind` provisions a lightweight Kubernetes cluster using Docker-backed nodes, which is sufficient for validating admission enforcement, artifact verification behavior, and repeatable cluster bootstrap without requiring external cloud infrastructure.

`kubectl` should remain the primary operator interface for scripted validation, evidence capture, and policy troubleshooting because it exposes the exact control-plane and workload state needed for thesis evaluation. `Helm` is suitable for packaging and repeatedly deploying the microservice, policy components, and supporting security controls with versioned configuration, which improves repeatability across milestones and demo resets. `Lens` is recommended as a secondary observability interface for live inspection of workloads, logs, events, and Helm releases during demonstrations; it improves presentation clarity but does not replace the cluster runtime or CLI-based verification path.

This stack is therefore aligned with the thesis requirement to implement a verifiable and enforceable end-to-end pipeline, while keeping the environment local, reproducible, and practical for controlled experiments.

## Milestones
| Milestone | Due date (UTC) | Link | Objective |
|---|---|---|---|
| M1 - Baseline Hardening and Reproducible Build | 2026-04-07 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/1 | Harden runtime/build baseline and deployment defaults |
| M2 - CI Pipeline (SBOM and Vulnerability Scan) | 2026-04-21 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/2 | Automate SBOM, Go vulnerability checks, and vulnerability controls |
| M3 - Artifact Signing and Provenance Attestation | 2026-05-05 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/3 | Establish artifact integrity and provenance trust |
| M4 - Kubernetes Admission Enforcement | 2026-05-19 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/4 | Enforce supply-chain policies at cluster admission |
| M5 - End-to-End Demo and Evidence Collection | 2026-06-02 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/5 | Validate with positive/negative demo scenarios and evidence |
| M6 - Thesis Packaging and Reusability | 2026-06-16 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/6 | Finalize thesis artifacts and reusable playbook |

## Issue Breakdown
### M1 - Baseline Hardening and Reproducible Build
- #1 Harden Dockerfile and runtime defaults for least privilege
- #2 Align Kubernetes base manifests with production-safe defaults

### M2 - CI Pipeline (SBOM and Vulnerability Scan)
- #3 Generate SPDX SBOM in CI and publish artifacts
- #4 Add Grype scan with threshold-based security gate
- #13 Integrate govulncheck into fail-fast CI security gates
- #14 Document and enforce Go dependency integrity baseline

### M3 - Artifact Signing and Provenance Attestation
- #5 Integrate keyless Cosign signing for image digests
- #6 Generate and verify SLSA provenance attestation

### M4 - Kubernetes Admission Enforcement
- #7 Enforce signature and provenance checks via Kyverno
- #8 Enforce SBOM and CVE metadata requirements in admission

### M5 - End-to-End Demo and Evidence Collection
- #9 Automate Kind and Kyverno bootstrap for reproducible demo
- #10 Add adversarial deployment scenarios and evidence checklist

### M6 - Thesis Packaging and Reusability
- #11 Build thesis traceability matrix (objective to evidence)
- #12 Publish reusable onboarding guide for additional Go microservices

## Continuity Note
- Existing milestones and issues are preserved; alignment is incremental.
- New issues are additive and scoped to missing objectives from the thesis (Go vulnerability gate and dependency integrity controls).
