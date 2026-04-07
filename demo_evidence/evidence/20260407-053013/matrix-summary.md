# Admission Matrix Summary

- Run ID: 20260407-053013
- Kubernetes context: kind-devsecops
- Namespace: stock-trading
- Signed image digest: ttl.sh/stock-trading-matrix-signed-6ad79074707d@sha256:d720573cef839b9c12ee46e845f24241726b744ee584c67a002ece9866c81ba4
- Unsigned image digest: ttl.sh/stock-trading-matrix-unsigned-6b81f8271238@sha256:d720573cef839b9c12ee46e845f24241726b744ee584c67a002ece9866c81ba4
- SBOM digest: 11E20656760325BF89B0640AE156BEBEC93A1A317594E4C06A295D0EB0A799AF

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
