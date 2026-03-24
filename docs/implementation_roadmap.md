# Implementation Roadmap (Milestones and Issues)

This roadmap converts the thesis plan into executable GitHub milestones and issues.

## Tracking Links
- Milestones: https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestones
- Roadmap issues (label): https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/issues?q=is%3Aissue%20is%3Aopen%20label%3Athesis-roadmap

## Milestones
| Milestone | Due date (UTC) | Link | Objective |
|---|---|---|---|
| M1 - Baseline Hardening and Reproducible Build | 2026-04-07 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/1 | Harden runtime/build baseline and deployment defaults |
| M2 - CI Pipeline (SBOM and Vulnerability Scan) | 2026-04-21 | https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/milestone/2 | Automate SBOM and vulnerability controls |
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

