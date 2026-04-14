# Lens Capture Checklist (Canonical Run)

As-of: `2026-04-14`  
Context: `docker-desktop`  
Namespace: `stock-trading`  
Canonical run-id: `20260414-210227` (`demo/evidence/20260414-210227`)

## Screenshot Targets (Completed)
| No. | Filename | Status | Screenshot Artifact | Evidence Source | Pass Signal |
|---|---|---|---|---|---|
| 01 | `01_valid_allow_deployment_available.png` | DONE | [docs/lens_screenshots/01_valid_allow_deployment_available.png](lens_screenshots/01_valid_allow_deployment_available.png) | [demo/evidence/20260414-210227/VALID_ALLOW/workloads.txt](../demo/evidence/20260414-210227/VALID_ALLOW/workloads.txt) | `user-service` Available `1/1` |
| 02 | `02_valid_allow_events_success.png` | DONE | [docs/lens_screenshots/02_valid_allow_events_success.png](lens_screenshots/02_valid_allow_events_success.png) | [demo/evidence/20260414-210227/VALID_ALLOW/events.txt](../demo/evidence/20260414-210227/VALID_ALLOW/events.txt) | `SuccessfulCreate`, `Started` |
| 03 | `03_unsigned_deny_apply_error.png` | DONE | [docs/lens_screenshots/03_unsigned_deny_apply_error.png](lens_screenshots/03_unsigned_deny_apply_error.png) | [demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt](../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kubectl-apply.txt) | contains `no signatures found` |
| 04 | `04_unsigned_deny_admission_report.png` | DONE | [docs/lens_screenshots/04_unsigned_deny_admission_report.png](lens_screenshots/04_unsigned_deny_admission_report.png) | [demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kyverno-logs.txt](../demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/kyverno-logs.txt) | deny mapped to verify-image policy |
| 05 | `05_missing_sbom_rs_failedcreate.png` | DONE | [docs/lens_screenshots/05_missing_sbom_rs_failedcreate.png](lens_screenshots/05_missing_sbom_rs_failedcreate.png) | [demo/evidence/20260414-210227/NEG_MISSING_SBOM_DENY/describe-replicasets.txt](../demo/evidence/20260414-210227/NEG_MISSING_SBOM_DENY/describe-replicasets.txt) | `FailedCreate` |
| 06 | `06_missing_sbom_events_reason.png` | DONE | [docs/lens_screenshots/06_missing_sbom_events_reason.png](lens_screenshots/06_missing_sbom_events_reason.png) | [demo/evidence/20260414-210227/NEG_MISSING_SBOM_DENY/events.txt](../demo/evidence/20260414-210227/NEG_MISSING_SBOM_DENY/events.txt) | contains SBOM deny reason |
| 07 | `07_missing_sbom_policy_report.png` | DONE | [docs/lens_screenshots/07_missing_sbom_policy_report.png](lens_screenshots/07_missing_sbom_policy_report.png) | [demo/evidence/20260414-210227/NEG_MISSING_SBOM_DENY/describe-replicasets.txt](../demo/evidence/20260414-210227/NEG_MISSING_SBOM_DENY/describe-replicasets.txt) | `require-sbom-digest` failed |
| 08 | `08_cve_deny_rs_failedcreate.png` | DONE | [docs/lens_screenshots/08_cve_deny_rs_failedcreate.png](lens_screenshots/08_cve_deny_rs_failedcreate.png) | [demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt](../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt) | `FailedCreate` |
| 09 | `09_cve_deny_events_reason.png` | DONE | [docs/lens_screenshots/09_cve_deny_events_reason.png](lens_screenshots/09_cve_deny_events_reason.png) | [demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/events.txt](../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/events.txt) | contains deny webhook reason |
| 10 | `10_cve_deny_policy_report.png` | DONE | [docs/lens_screenshots/10_cve_deny_policy_report.png](lens_screenshots/10_cve_deny_policy_report.png) | [demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt](../demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt) | `require-high-critical-zero` failed |
| 11 | `11_valid_recheck_available.png` | DONE | [docs/lens_screenshots/11_valid_recheck_available.png](lens_screenshots/11_valid_recheck_available.png) | [demo/evidence/20260414-210227/VALID_ALLOW_RECHECK/workloads.txt](../demo/evidence/20260414-210227/VALID_ALLOW_RECHECK/workloads.txt) | Available `1/1` |
| 12 | `12_matrix_summary_all_pass.png` | DONE | [docs/lens_screenshots/12_matrix_summary_all_pass.png](lens_screenshots/12_matrix_summary_all_pass.png) | [demo/evidence/20260414-210227/matrix-summary.md](../demo/evidence/20260414-210227/matrix-summary.md) | all scenarios `PASS` |

## Apply Commands (Canonical Run)
```powershell
kubectl config use-context docker-desktop
kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260414-210227/VALID_ALLOW/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260414-210227/NEG_UNSIGNED_DENY/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260414-210227/NEG_MISSING_SBOM_DENY/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260414-210227/NEG_CVE_THRESHOLD_DENY/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260414-210227/VALID_ALLOW_RECHECK/deployment.yaml
```

## Notes
- Screenshot pack lives under `docs/lens_screenshots/` with exact filenames above.
- Each screenshot is traceable to canonical raw evidence in `demo/evidence/20260414-210227/`.

