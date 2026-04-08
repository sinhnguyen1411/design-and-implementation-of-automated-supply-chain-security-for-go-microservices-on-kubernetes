# Admission Matrix Summary

- Run ID: 20260408-091500
- Kubernetes context: docker-desktop
- Namespace: stock-trading
- Signed image digest: ttl.sh/stock-trading-matrix-signed-a4fb2d70be3a@sha256:1d0d86bd74ef2aa3b1a20f902ed65b00e8a2eb3234f26a237fcd85fd1fbece55
- Unsigned image digest: ttl.sh/stock-trading-matrix-unsigned-a4fb2d70be3a@sha256:1d0d86bd74ef2aa3b1a20f902ed65b00e8a2eb3234f26a237fcd85fd1fbece55
- SBOM digest: 51F58E0422206693D853BC5FA90780226C4AED4A7A96E0B3DB2A6E059CC17E2F

| Case | Expected | Actual | Verdict | Reason |
|---|---|---|---|---|
| VALID_ALLOW | Allowed | Allowed | PASS | Deployment became Available. |
| NEG_UNSIGNED_DENY | Denied | Denied | PASS | Admission deny evidence detected at apply phase. |
| NEG_MISSING_SBOM_DENY | Denied | Denied | PASS | Admission deny evidence detected in events/ReplicaSet describe. |
| NEG_CVE_THRESHOLD_DENY | Denied | Denied | PASS | Admission deny evidence detected in events/ReplicaSet describe. |

## Regression Re-check

| Check | Expected | Actual | Verdict | Reason |
|---|---|---|---|---|
| VALID_ALLOW_RECHECK | Allowed | Allowed | PASS | Deployment became Available. |
