# Go Dependency Integrity Baseline

This document defines the dependency-integrity baseline enforced by the CI pipeline for this repository.

## Objectives
- Ensure Go module dependencies are tamper-evident and reproducible.
- Prevent implicit dependency graph mutations during CI runs.
- Keep dependency state auditable alongside SBOM and vulnerability evidence.

## Enforced Controls (Fail-Fast)
1. `go mod download`
   - Resolves module dependencies from configured module sources.
2. `go mod verify`
   - Verifies downloaded module content against cryptographic checksums.
3. `go list -deps -mod=readonly ./...`
   - Resolves dependency graph in readonly mode to prevent implicit `go.mod` / `go.sum` updates.

If any control above fails, the CI job fails.

## Audit Signal
- `go mod tidy -diff` runs as an audit signal.
- If drift is detected, CI records a warning in the dependency report artifact for follow-up remediation.

## Evidence Artifact
- Workflow artifact name: `dependency-integrity-report`
- Default report file: `dependency-integrity-report.txt`

The report is uploaded even when earlier integrity checks fail (`if: always()`), so thesis evidence can include both pass and fail cases.

## Relation to Thesis Controls
- Maps to objective: dependency transparency and control.
- Complements SBOM generation and vulnerability scanning by establishing trust in module inputs before build/sign/attestation.
