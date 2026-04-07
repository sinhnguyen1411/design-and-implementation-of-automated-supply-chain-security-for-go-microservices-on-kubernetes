# Lens Capture Checklist (Ordered)

Use Kubernetes context `docker-desktop` and namespace `stock-trading`.

## Screenshot Targets
| No. | Filename | Where to capture in Lens / Terminal | Pass signal |
|---|---|---|---|
| 01 | `01_valid_allow_deployment_available.png` | `Workloads > Deployments` after applying `VALID_ALLOW` | `user-service` Available `1/1` |
| 02 | `02_valid_allow_events_success.png` | `Events` after `VALID_ALLOW` | `SuccessfulCreate`, `Started` |
| 03 | `03_unsigned_deny_apply_error.png` | Terminal output after applying `NEG_UNSIGNED_DENY` | contains `no signatures found` |
| 04 | `04_unsigned_deny_admission_report.png` | `Custom Resources > kyverno.io > Admission Report` | deny mapped to verify image policy |
| 05 | `05_missing_sbom_rs_failedcreate.png` | `Workloads > Replica Sets` after `NEG_MISSING_SBOM_DENY` | `FailedCreate` |
| 06 | `06_missing_sbom_events_reason.png` | `Events` after `NEG_MISSING_SBOM_DENY` | `require-sbom-digest` reason |
| 07 | `07_missing_sbom_policy_report.png` | `Custom Resources > wgpolicyk8s.io > Policy Report` | missing SBOM rule failed |
| 08 | `08_cve_deny_rs_failedcreate.png` | `Workloads > Replica Sets` after `NEG_CVE_THRESHOLD_DENY` | `FailedCreate` |
| 09 | `09_cve_deny_events_reason.png` | `Events` after `NEG_CVE_THRESHOLD_DENY` | `require-high-critical-zero` reason |
| 10 | `10_cve_deny_policy_report.png` | `Custom Resources > wgpolicyk8s.io > Policy Report` | CVE threshold rule failed |
| 11 | `11_valid_recheck_available.png` | `Workloads > Deployments` after `VALID_ALLOW_RECHECK` | Available `1/1` |
| 12 | `12_matrix_summary_all_pass.png` | Open `demo/evidence/<run-id>/matrix-summary.md` | all cases `PASS` |

## Apply Commands (run in repo root)
```powershell
kubectl config use-context docker-desktop
kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260406-154444/VALID_ALLOW/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260406-154444/NEG_UNSIGNED_DENY/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260406-154444/NEG_MISSING_SBOM_DENY/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260406-154444/NEG_CVE_THRESHOLD_DENY/deployment.yaml

kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f demo/evidence/20260406-154444/VALID_ALLOW_RECHECK/deployment.yaml
```

## Save Location
Put screenshots under `docs/lens_screenshots/` using exact filenames above.
