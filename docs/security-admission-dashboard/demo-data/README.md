# Bundled Dashboard Demo Data

This folder contains static evidence datasets for the security admission dashboard.

## Purpose
- Provide deterministic demo runs when `demo/evidence` is not available.
- Keep the dashboard usable in repository-only viewing scenarios.
- Keep a clear contract: `demo/evidence` is canonical for thesis claims, while this folder is fallback-only for offline UI/demo continuity.

## Included Run IDs
- `20260406-154444`
- `20260407-203000`

## Data Shape
Each run includes:
- `matrix-summary.md`
- `matrix-index.json`
- `regression-valid-allow.json`
- Per-case artifact files under:
  - `VALID_ALLOW/`
  - `NEG_UNSIGNED_DENY/`
  - `NEG_MISSING_SBOM_DENY/`
  - `NEG_CVE_THRESHOLD_DENY/`
  - `VALID_ALLOW_RECHECK/`

The JSON structure mirrors the output of `scripts/admission_matrix_demo.ps1`.
