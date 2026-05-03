# Demo Evidence - Admission Matrix on Docker Desktop

## Scope
This evidence package validates thesis admission criteria for:
- invalid artifacts are denied,
- valid artifacts are admitted,
- valid admission still works after deny scenarios.

## Canonical Evidence Policy
- Canonical thesis evidence root: `demo/evidence/`
- Dashboard fallback-only dataset: `docs/security-admission-dashboard/demo-data/evidence/`
- Latest canonical matrix run (as of 2026-04-14): `demo/evidence/20260414-210227`

## Environment and Run Metadata
- Date: 2026-04-14
- Kubernetes context: `docker-desktop`
- Namespace: `stock-trading`
- Script: `infra/scripts/admission_matrix_demo.ps1`
- Command:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File infra/scripts/admission_matrix_demo.ps1 -Context docker-desktop -Namespace stock-trading -ExportDir demo/evidence -ResetNamespace
```

- Evidence directory: `demo/evidence/20260414-210227`
- Matrix summary file: `demo/evidence/20260414-210227/matrix-summary.md`
- Matrix JSON index: `demo/evidence/20260414-210227/matrix-index.json`
- Regression JSON result: `demo/evidence/20260414-210227/regression-valid-allow.json`
- Signed digest used:
  - `ttl.sh/stock-trading-matrix-signed-2369c9df4111@sha256:430529879d903a9514281956ea05030d56f82dcac8458751e57be0b410c2e8ef`
- Unsigned digest used:
  - `ttl.sh/stock-trading-matrix-unsigned-07dbc73a2438@sha256:430529879d903a9514281956ea05030d56f82dcac8458751e57be0b410c2e8ef`
- SBOM digest annotation used:
  - `6705A0D11524755F016E83CD19091668898658A01C356EE8661C7A3BEEAD2AAB`

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

Per-case directories under `demo/evidence/20260414-210227/`:
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
