# Admission Matrix Summary

- Run ID: 20260414-210227
- Kubernetes context: docker-desktop
- Namespace: stock-trading
- Signed image digest: ttl.sh/stock-trading-matrix-signed-2369c9df4111@sha256:430529879d903a9514281956ea05030d56f82dcac8458751e57be0b410c2e8ef
- Unsigned image digest: ttl.sh/stock-trading-matrix-unsigned-07dbc73a2438@sha256:430529879d903a9514281956ea05030d56f82dcac8458751e57be0b410c2e8ef
- SBOM digest: 6705A0D11524755F016E83CD19091668898658A01C356EE8661C7A3BEEAD2AAB

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
