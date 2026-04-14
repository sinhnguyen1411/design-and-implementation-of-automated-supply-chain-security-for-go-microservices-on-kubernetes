# Onboarding Second-Service Dry-Run Summary

Run ID: `20260414-213541-onboarding-second-service`
Date: `2026-04-14`
Context: `docker-desktop`
Namespace under test: `portfolio-trading`

## Cases
| Case | Manifest | Expected | Actual | Verdict |
|---|---|---|---|---|
| Case A - new service label | `case_portfolio_label_deploy_only.yaml` | Shows whether current policy scope enforces a non-`user-service` label | Deployment admitted; Pod creation proceeded; runtime failed later with image pull error (`ErrImagePull`) | PASS (scope gap reproduced) |
| Case B - `user-service` label control case | `case_user_service_label_deploy_only.yaml` | Should be denied for missing required security annotations | ReplicaSet create denied by Kyverno (`require-sbom-digest`, `require-high-critical-zero`) | PASS (control path validated) |

## Key Signals
- Case A evidence: `11_case1_workloads.txt`, `12_case1_events.txt`
- Case B evidence: `16_case2_describe_rs.txt`, `17_case2_events.txt`
- Cleanup evidence: `18_cleanup_namespace_delete.txt`

## Conclusion
Current onboarding guide is now backed by a second-service simulation evidence set. The main reusable gap is policy/script scope parameterization from `user-service`-specific defaults.
