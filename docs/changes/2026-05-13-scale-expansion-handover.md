# Scale Expansion Handover (2026-05-13)

This handover note captures the major changes already delivered for the high-scalable milestone so other developers can quickly sync.

## Delivered Scope
- Expanded service count from 4 to 10 by adding 6 trading-core services:
  - `market-data-service`
  - `pricing-service`
  - `execution-service`
  - `settlement-service`
  - `compliance-service`
  - `notification-service`
- Each new service includes:
  - `cmd/server/main.go`
  - domain logic under `internal/<domain>/...`
  - `GET /healthz` and one domain POST endpoint
  - unit tests for happy-path + guardrail
  - `go.mod`, `go.sum`, `Dockerfile`, and Kubernetes base manifests

## CI and Registry Changes
- `services.yaml` extended to 10 services with profile tags and Go version alignment.
- `ci-service` now supports nightly full-matrix execution via `schedule` while preserving changed-only behavior for push/PR.
- `onboarding-lab` matrix is now generated dynamically from `services.yaml` (nightly all-services coverage).

## Docs Updates
- README now includes:
  - 10-service landscape
  - CI scalability strategy (changed-only + nightly full)
  - scale evidence guidance
- Onboarding guide updated with CI execution modes and services registry ownership.

## Key Commits For Review
- `10dd63c` feat(scale): add 6 trading-core services and nightly full matrix
- `684a576` ci: use workflow_dispatch choice inputs for dropdown UI
- `03eed24` ci: add runner_target to avoid self-hosted wait hangs
- `c2df7b9` ci: add gate_mode to toggle security gate enforcement

## Validation Snapshot
- `go test ./...` passed across all service modules in the monorepo.
- Workflow YAMLs parse successfully after updates.
