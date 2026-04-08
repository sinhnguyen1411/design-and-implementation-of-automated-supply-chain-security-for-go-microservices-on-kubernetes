# Admission Matrix Summary

- Run ID: 20260406-153527
- Kubernetes context: docker-desktop
- Namespace: stock-trading
- Signed image digest: ttl.sh/stock-trading-matrix-signed-0d18061e3624@sha256:296d1a9fd84592607802a6d5faef9ab5f7c1be73b87f2f14d725df6f06939e53
- Unsigned image digest: ttl.sh/stock-trading-matrix-unsigned-563d4df90f77@sha256:296d1a9fd84592607802a6d5faef9ab5f7c1be73b87f2f14d725df6f06939e53
- SBOM digest: AD7522744D9FEA68D305C1BCADF30F804616C3573613979BA73D2DCE368BAE76

| Case | Expected | Actual | Verdict | Reason |
|---|---|---|---|---|
| VALID_ALLOW | Allowed | Allowed | PASS | Deployment became Available. |
| NEG_UNSIGNED_DENY | Denied | Denied | PASS | Admission deny evidence detected at apply phase. |
| NEG_MISSING_SBOM_DENY | Denied | Denied | PASS | Admission deny evidence detected in events/ReplicaSet describe. |
| NEG_CVE_THRESHOLD_DENY | Denied | Denied | PASS | Admission deny evidence detected in events/ReplicaSet describe. |
