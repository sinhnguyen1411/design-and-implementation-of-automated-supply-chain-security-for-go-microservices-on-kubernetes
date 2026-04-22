# Bundled Dashboard Demo Data

This folder contains static fallback datasets for the security admission dashboard.

## Source of Truth and Fallback Contract
- Primary source (thesis live view): `docs/security-admission-dashboard/data/actions-runs.snapshot.json`.
- Snapshot producer: `.github/workflows/dashboard-data-sync.yml`.
- Snapshot inputs: workflow runs from:
- `secure-supply-chain` (build-time CVE/gate/sign data).
- `admission-matrix-evidence` (runtime admission matrix evidence).
- `service-scs-matrix-evidence` (matrix-style runtime evidence across multiple services).
- This `demo-data/evidence` folder is fallback-only when snapshot is missing/unavailable (offline preview mode).
- `demo/evidence` remains useful for local script outputs, but dashboard default is Actions snapshot first.

## Included Run IDs
- `20260406-154444`
- `20260407-203000`

## Data Shape
Each run includes:
- `matrix-summary.md`
- `matrix-index.json`
- `regression-valid-allow.json`
- `security-gate-findings.sample.json` (fallback CVE list for dashboard panel when live files are unavailable)
- Per-case artifact files under:
  - `VALID_ALLOW/`
  - `NEG_UNSIGNED_DENY/`
  - `NEG_MISSING_SBOM_DENY/`
  - `NEG_CVE_THRESHOLD_DENY/`
  - `VALID_ALLOW_RECHECK/`

The JSON structure mirrors the output of `scripts/admission_matrix_demo.ps1`.
