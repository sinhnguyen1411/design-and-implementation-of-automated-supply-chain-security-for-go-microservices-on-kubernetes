# Demo Evidence - Admission Matrix on Docker Desktop

## Scope
This evidence package validates thesis admission criteria for:
- invalid artifacts are denied,
- valid artifacts are admitted,
- valid admission still works after deny scenarios.

## Environment and Run Metadata
- Date: 2026-04-06
- Kubernetes context: `docker-desktop`
- Namespace: `stock-trading`
- Script: `scripts/admission_matrix_demo.ps1`
- Command:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/admission_matrix_demo.ps1 -Context docker-desktop -Namespace stock-trading -ExportDir demo/evidence -ResetNamespace
```

- Evidence directory: `demo/evidence/20260406-154444`
- Matrix summary file: `demo/evidence/20260406-154444/matrix-summary.md`
- Matrix JSON index: `demo/evidence/20260406-154444/matrix-index.json`
- Regression JSON result: `demo/evidence/20260406-154444/regression-valid-allow.json`
- Signed digest used:
  - `ttl.sh/stock-trading-matrix-signed-c853c0f8a3b5@sha256:8d87c68bf36634dbcd4d161ac23ac4ec22f81f326f45b820e6564e127493db8f`
- Unsigned digest used:
  - `ttl.sh/stock-trading-matrix-unsigned-71b8f514cbd1@sha256:8d87c68bf36634dbcd4d161ac23ac4ec22f81f326f45b820e6564e127493db8f`
- SBOM digest annotation used:
  - `E3E11664C6A175D35675B517F9C281E5D1F0FF937F572575042017AC30D6EE08`

## Pre-check Result
- `go test ./...` passed before matrix execution (captured in script output).
- Cluster reachability (`kubectl get nodes`) and Kyverno policy application passed.

## Matrix Verdict (case -> expected -> actual -> verdict)
| Case | Expected | Actual | Verdict |
|---|---|---|---|
| `VALID_ALLOW` | Allowed | Allowed | PASS |
| `NEG_UNSIGNED_DENY` | Denied | Denied | PASS |
| `NEG_MISSING_SBOM_DENY` | Denied | Denied | PASS |
| `NEG_CVE_THRESHOLD_DENY` | Denied | Denied | PASS |

Regression check (post-deny):

| Check | Expected | Actual | Verdict |
|---|---|---|---|
| `VALID_ALLOW_RECHECK` | Allowed | Allowed | PASS |

## Deny/Allow Reasons (raw evidence)
1. `VALID_ALLOW` admitted
   - Evidence: `VALID_ALLOW/kubectl-wait.txt`
   - Signal: `deployment.apps/user-service condition met`

2. `NEG_UNSIGNED_DENY` denied by signature verification
   - Evidence: `NEG_UNSIGNED_DENY/kubectl-apply.txt`
   - Signal excerpt:

```text
verify-local-matrix-images:
  autogen-verify-local-matrix-signature-and-attestation: 'failed to verify image ...
    .attestors[0].entries[0].keys: no signatures found'
```

3. `NEG_MISSING_SBOM_DENY` denied by required SBOM annotation
   - Evidence: `NEG_MISSING_SBOM_DENY/describe-replicasets.txt`
   - Signal excerpt:

```text
require-sbom-annotation:
  require-sbom-digest: 'validation error: Pod rejected: missing SBOM reference annotation
    (security.stock-trading.dev/sbom-digest) ...'
```

4. `NEG_CVE_THRESHOLD_DENY` denied by CVE threshold rule
   - Evidence: `NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt`
   - Signal excerpt:

```text
enforce-cve-threshold:
  require-high-critical-zero: 'validation error: Pod rejected: fixable High/Critical
    CVEs not cleared (security.grype.io/high_critical must be ''0'') ...'
```

5. `VALID_ALLOW_RECHECK` admitted after all deny cases
   - Evidence: `VALID_ALLOW_RECHECK/kubectl-wait.txt`
   - Signal: `deployment.apps/user-service condition met`

## Artifact Completeness Checklist
Each case contains all required artifact groups:
- `kubectl apply`: `kubectl-apply.txt`
- `kubectl events`: `events.txt`
- `describe workload`: `describe-deployment.txt`, `describe-replicasets.txt`, `describe-pods.txt`
- policy controller logs: `kyverno-logs.txt`

Per-case directories under `demo/evidence/20260406-154444/`:
- `VALID_ALLOW/`
- `NEG_UNSIGNED_DENY/`
- `NEG_MISSING_SBOM_DENY/`
- `NEG_CVE_THRESHOLD_DENY/`
- `VALID_ALLOW_RECHECK/`

## Conclusion
The run satisfies the minimum appendix-aligned admission validation target on `docker-desktop`:
- 1 valid workload admitted,
- 3 negative scenarios denied with explicit policy reasons,
- valid workload still admitted in post-deny regression re-check.
