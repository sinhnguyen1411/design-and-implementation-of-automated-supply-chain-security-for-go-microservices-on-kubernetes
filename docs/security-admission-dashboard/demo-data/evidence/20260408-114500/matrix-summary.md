# Admission Matrix Summary

- Run ID: 20260408-114500
- Kubernetes context: docker-desktop
- Namespace: stock-trading
- Signed image digest: ttl.sh/stock-trading-matrix-signed-c853c0f8a3b5@sha256:8d87c68bf36634dbcd4d161ac23ac4ec22f81f326f45b820e6564e127493db8f
- Unsigned image digest: ttl.sh/stock-trading-matrix-unsigned-71b8f514cbd1@sha256:8d87c68bf36634dbcd4d161ac23ac4ec22f81f326f45b820e6564e127493db8f
- SBOM digest: E3E11664C6A175D35675B517F9C281E5D1F0FF937F572575042017AC30D6EE08

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
