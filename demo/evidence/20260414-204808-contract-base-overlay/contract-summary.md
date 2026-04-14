# Contract Evidence Summary (Base Deny -> Overlay Allow)

- Run ID: 20260414-204808
- Context: docker-desktop (fallback because kind CLI is unavailable)
- Evidence directory: demo/evidence/20260414-204808-contract-base-overlay
- Base apply exit code: 1
- Overlay apply exit code: 0
- Overlay wait exit code: 1
- Base denied verdict: True
- Overlay allowed verdict (admission): True
- Overlay image digest: ttl.sh/stock-trading-contract-b170c22c5cec@sha256:b63294a0359495f94f7db92649a372b9063a1bcfe4b4ed43ece3649be0dca793
- Overlay SBOM digest annotation: 8F26EFB2ED9689C4A1ED7C49C89F640B0C6175AEF1A0FD96512E15C627457021

## Notes
- `OVERLAY_WAIT_EXIT=1` is runtime readiness (CrashLoopBackOff), not admission denial.
- Admission allow signal is from successful `kubectl apply -k overlay-allow` plus Kyverno verify pass entries in `04_kyverno_logs.txt`.

## Key Evidence Files
- `02_apply_base.txt`
- `02_kyverno_logs.txt`
- `04_apply_overlay_allow.txt`
- `04_events.txt`
- `04_describe_deploy.txt`
- `04_kyverno_logs.txt`
- `overlay-allow/kustomization.yaml`
- `overlay-allow/patch-allow.yaml`
