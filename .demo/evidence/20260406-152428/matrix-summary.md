# Admission Matrix Summary

- Run ID: 20260406-152428
- Kubernetes context: docker-desktop
- Namespace: stock-trading
- Signed image digest: ttl.sh/stock-trading-matrix-signed-e478a4efa09b@sha256:f9a7107856903729e19dc08f4b6b0a998c0a547732d3ea66114154509d53dcca
- Unsigned image digest: ttl.sh/stock-trading-matrix-unsigned-4ff9b705ba1b@sha256:f9a7107856903729e19dc08f4b6b0a998c0a547732d3ea66114154509d53dcca
- SBOM digest: 91F0A7223E1D887BC34429A8C37607F5A17C049BB254A47BA38334FE7E5A3048

| Case | Expected | Actual | Verdict | Reason |
|---|---|---|---|---|
| VALID_ALLOW | Allowed | Allowed | PASS | Deployment became Available. |
| NEG_UNSIGNED_DENY | Denied | UnknownOrAllowed | FAIL | Expected Denied but no strong deny signal was found. |
| NEG_MISSING_SBOM_DENY | Denied | UnknownOrAllowed | FAIL | Expected Denied but no strong deny signal was found. |
| NEG_CVE_THRESHOLD_DENY | Denied | UnknownOrAllowed | FAIL | Expected Denied but no strong deny signal was found. |
