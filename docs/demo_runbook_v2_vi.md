# Demo Runbook V2 (VI) - Thesis Supply Chain Security (10-15 phút)

Tài liệu này là kịch bản trình diễn kỹ thuật cho luận văn, theo chế độ **hybrid live + evidence**.
Mục tiêu là chứng minh nhất quán mô hình:
- `verify-before-release` trong CI,
- `enforce-before-run` tại Kubernetes admission.

## 1) Demo Interface (Decision-Complete)

### Input bắt buộc
- Kubernetes context: `docker-desktop`
- Namespace: `stock-trading`
- Evidence run-id: mặc định lấy **latest** trong `demo/evidence/`

### Output bắt buộc
- Matrix verdict có đủ 4 case chính và đều `PASS`:
  - `VALID_ALLOW`
  - `NEG_UNSIGNED_DENY`
  - `NEG_MISSING_SBOM_DENY`
  - `NEG_CVE_THRESHOLD_DENY`
- Regression check `VALID_ALLOW_RECHECK` là `PASS`
- Deny reason map khớp policy contract:
  - unsigned -> verify image/signature policy deny
  - missing SBOM annotation -> require SBOM policy deny
  - `high_critical != 0` -> CVE threshold policy deny

### Fallback contract (bắt buộc)
- Mỗi bước live quan trọng cho phép tối đa **60-90 giây**.
- Nếu quá thời gian hoặc gặp lỗi hạ tầng, chuyển ngay sang evidence file tương ứng trong `demo/evidence/<run-id>/...`.
- Không debug sâu trong lúc trình bày; ưu tiên giữ narrative và bằng chứng.

## 2) Chuẩn bị trước giờ demo (T-15 đến T-5)

### 2.1. Chạy matrix để tạo/gia hạn evidence (khuyến nghị)

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/admission_matrix_demo.ps1 -Context docker-desktop -Namespace stock-trading -ExportDir demo/evidence -ResetNamespace
```

### 2.2. Chọn run-id mặc định (latest)

```bash
LATEST_RUN_ID=$(ls -1 demo/evidence 2>/dev/null | grep -E '^[0-9]{8}-[0-9]{6}$' | sort -r | head -n 1)
echo "$LATEST_RUN_ID"
RUN_ID="${LATEST_RUN_ID:-20260407-203000}"
echo "$RUN_ID"
```

Nếu không có run mới, dùng baseline đã ghi trong docs:
- `20260406-154444`
- Hoặc dùng bundled dataset: `20260407-203000` trong `docs/security-admission-dashboard/demo-data/evidence/`

Lưu ý: bundled dataset dùng để trình bày dashboard/evidence; các bước deploy live vẫn nên dùng run trong `demo/evidence/`.

### 2.3. Mở sẵn tab phục vụ trình bày
- `docs/scs_architecture_diagram.html`
- `docs/security-admission-dashboard/index.html?run=<LATEST_RUN_ID>`
- `demo/evidence/<LATEST_RUN_ID>/matrix-summary.md`
- `docs/demo_speaker_notes_v2_vi.md` (1 trang đọc khi trình bày)

Lưu ý dashboard là static page, cần chạy static server từ repo root nếu mở trực tiếp bằng file path không đọc được dữ liệu:

```bash
python3 -m http.server 8080
# mở: http://localhost:8080/docs/security-admission-dashboard/index.html?run=<LATEST_RUN_ID>
```

Dashboard sẽ ưu tiên đọc dữ liệu từ `demo/evidence/`. Nếu không có, dashboard tự fallback sang bundled dataset tại `docs/security-admission-dashboard/demo-data/evidence/`.

## 3) Timeline Runbook 10-15 phút (6 scene cố định)

### Scene 1 - Context + Architecture (1-2 phút)

Thông điệp kỹ thuật chính: chuỗi kiểm soát là liên tục từ CI đến Admission, không phải scan rời rạc.

Thao tác:
- Mở `docs/scs_architecture_diagram.html`
- Chỉ vào luồng: `dependency integrity -> govulncheck -> SBOM -> Grype gate -> sign -> attest -> admission verify`

Lời thoại mẫu:
- "Luận văn này tập trung vào chuỗi trust liên tục: CI xác minh trước khi phát hành, Kubernetes enforce trước khi workload được chạy."

Expected signal:
- Hội đồng nhìn thấy rõ fail-fast path ở CI và deny path ở admission trong cùng một kiến trúc.

Fallback:
- Nếu không mở được HTML, nói theo cấu trúc tương đương trong `docs/devsecops_ci_admission.md` phần pipeline stages + admission policies.

### Scene 2 - Pre-flight (1 phút)

Thông điệp kỹ thuật chính: trước khi demo case, môi trường đã sẵn sàng và policy đã nạp.

Command:

```bash
kubectl config use-context docker-desktop
kubectl -n stock-trading get ns stock-trading
kubectl get nodes
kubectl get clusterpolicy verify-user-service-images require-sbom-annotation enforce-cve-threshold
kubectl -n kyverno get deploy
```

Expected signal:
- Context đúng `docker-desktop`
- Namespace tồn tại
- Node `Ready`
- Có các policy liên quan verify image / sbom / cve threshold

Fallback (60-90s):
- Mở `docs/demo_evidence.md` phần **Pre-check Result** để xác nhận pre-check đã pass trong run evidence.

### Scene 3 - VALID_ALLOW (2 phút)

Thông điệp kỹ thuật chính: artifact hợp lệ (signed + attested + annotation đúng) được admission cho chạy.

Command:

```bash
kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f "demo/evidence/${RUN_ID}/VALID_ALLOW/deployment.yaml"
kubectl -n stock-trading wait --for=condition=Available deployment/user-service --timeout=90s
kubectl -n stock-trading get deploy user-service
```

Expected signal:
- `deployment.apps/user-service condition met`
- Deployment `Available` (1/1)

Fallback (60-90s):
- Mở `demo/evidence/<RUN_ID>/VALID_ALLOW/kubectl-wait.txt` và chỉ ra dòng `condition met`.

### Scene 4 - 3 Negative Cases (4-5 phút)

Thông điệp kỹ thuật chính: admission deny theo policy cụ thể, có lý do rõ và tái kiểm chứng được.

### 4.1 NEG_UNSIGNED_DENY

```bash
kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f "demo/evidence/${RUN_ID}/NEG_UNSIGNED_DENY/deployment.yaml"
```

Expected deny reason:
- Chuỗi lỗi chứa `no signatures found`

Fallback evidence:
- `demo/evidence/<RUN_ID>/NEG_UNSIGNED_DENY/kubectl-apply.txt`

### 4.2 NEG_MISSING_SBOM_DENY

```bash
kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f "demo/evidence/${RUN_ID}/NEG_MISSING_SBOM_DENY/deployment.yaml"
```

Expected deny reason:
- Vi phạm `security.stock-trading.dev/sbom-digest`
- Có `require-sbom-digest` / `missing SBOM reference annotation`

Fallback evidence:
- `demo/evidence/<RUN_ID>/NEG_MISSING_SBOM_DENY/describe-replicasets.txt`

### 4.3 NEG_CVE_THRESHOLD_DENY

```bash
kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f "demo/evidence/${RUN_ID}/NEG_CVE_THRESHOLD_DENY/deployment.yaml"
```

Expected deny reason:
- Vi phạm `security.grype.io/high_critical must be '0'`
- Có `require-high-critical-zero`

Fallback evidence:
- `demo/evidence/<RUN_ID>/NEG_CVE_THRESHOLD_DENY/describe-replicasets.txt`

Deny reason map cần nói rõ khi trình bày:

| Case | Policy semantics | Signal cần thấy |
|---|---|---|
| `NEG_UNSIGNED_DENY` | verify image/signature | `no signatures found` |
| `NEG_MISSING_SBOM_DENY` | require SBOM annotation | `require-sbom-digest` |
| `NEG_CVE_THRESHOLD_DENY` | enforce CVE threshold | `require-high-critical-zero` |

### Scene 5 - VALID_ALLOW_RECHECK (1 phút)

Thông điệp kỹ thuật chính: sau các deny case, workload hợp lệ vẫn được admit bình thường (không bị side effect).

Command:

```bash
kubectl -n stock-trading delete deploy user-service --ignore-not-found=true
kubectl -n stock-trading apply -f "demo/evidence/${RUN_ID}/VALID_ALLOW_RECHECK/deployment.yaml"
kubectl -n stock-trading wait --for=condition=Available deployment/user-service --timeout=90s
```

Expected signal:
- `condition met` và workload chạy lại thành công.

Fallback:
- `demo/evidence/<RUN_ID>/VALID_ALLOW_RECHECK/kubectl-wait.txt`

### Scene 6 - Evidence Board (2-3 phút)

Thông điệp kỹ thuật chính: kết quả không chỉ live “tại chỗ” mà còn được đóng gói evidence có thể kiểm tra lại.

Thao tác:
- Mở dashboard: `docs/security-admission-dashboard/index.html?run=<RUN_ID>`
- Mở `matrix-summary.md`
- Chỉ ra 4 case chính `PASS` + regression `PASS`

Expected signal:
- Matrix hiển thị đúng thứ tự case và verdict.
- Có thể mở artifact theo từng case (`kubectl-apply`, `describe-*`, `kyverno-logs`).

Fallback:
- Nếu dashboard không load, mở trực tiếp:
  - `demo/evidence/<RUN_ID>/matrix-summary.md`
  - `demo/evidence/<RUN_ID>/matrix-index.json`
  - `demo/evidence/<RUN_ID>/regression-valid-allow.json`

## 4) Dry-run Checklist (Pass/Fail)

| Check | Tiêu chí Pass | Kết quả |
|---|---|---|
| Pre-flight | cluster reachable + policy loaded | `[] PASS` / `[] FAIL` |
| `VALID_ALLOW` | admitted + Available | `[] PASS` / `[] FAIL` |
| `NEG_UNSIGNED_DENY` | deny do signature | `[] PASS` / `[] FAIL` |
| `NEG_MISSING_SBOM_DENY` | deny do missing SBOM annotation | `[] PASS` / `[] FAIL` |
| `NEG_CVE_THRESHOLD_DENY` | deny do high_critical != 0 | `[] PASS` / `[] FAIL` |
| `VALID_ALLOW_RECHECK` | admitted sau deny sequence | `[] PASS` / `[] FAIL` |
| Evidence board | matrix + summary + regression file đọc được | `[] PASS` / `[] FAIL` |

## 5) Acceptance Criteria của runbook

- Người vận hành mới có thể chạy demo theo tài liệu mà không cần quyết định bổ sung.
- Mỗi scene có đúng 1 thông điệp kỹ thuật chính để tránh quá tải.
- Mọi bước live quan trọng đều có fallback rõ ràng trong 60-90 giây.
- Narrative bám sát traceability của thesis và policy contract hiện hành.

## 6) Liên kết chuẩn hoá với tài liệu hiện có

- Pipeline + admission contract: `docs/devsecops_ci_admission.md`
- Evidence baseline và deny reason: `docs/demo_evidence.md`
- Lens capture flow: `docs/lens_capture_checklist.md`
- Thesis objective và traceability: `docs/thesis_spec_en.md`
- Speaker notes 1 trang: `docs/demo_speaker_notes_v2_vi.md`
- Bundled dashboard data: `docs/security-admission-dashboard/demo-data/`
