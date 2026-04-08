# Admission Matrix Summary

- Run ID: 20260406-154000
- Kubernetes context: docker-desktop
- Namespace: stock-trading
- Signed image digest: ttl.sh/stock-trading-matrix-signed-55c87360239d@sha256:735dce323a6793dfb7f8eb400cdb21ace9b4816195cf14bf1a4cb18540534278
- Unsigned image digest: ttl.sh/stock-trading-matrix-unsigned-efb75666c310@sha256:735dce323a6793dfb7f8eb400cdb21ace9b4816195cf14bf1a4cb18540534278
- SBOM digest: 6017F0A21EF550D6B19EB0814D814BB3FEB95AD56DAE0792C7C4A44455EAB3CA

| Case | Expected | Actual | Verdict | Reason |
|---|---|---|---|---|
| VALID_ALLOW | Allowed | Allowed | PASS | Deployment became Available. |
| NEG_UNSIGNED_DENY | Denied | Denied | PASS | Admission deny evidence detected at apply phase. |
| NEG_MISSING_SBOM_DENY | Denied | Denied | PASS | Admission deny evidence detected in events/ReplicaSet describe. |
| NEG_CVE_THRESHOLD_DENY | Denied | Denied | PASS | Admission deny evidence detected in events/ReplicaSet describe. |
