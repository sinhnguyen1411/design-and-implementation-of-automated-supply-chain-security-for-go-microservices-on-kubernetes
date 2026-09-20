# Nhật Ký Phiên Làm Việc (Work Session Log)
## Dự Án: Automated Supply Chain Security for Go Microservices on Kubernetes
**Tác giả / Nghiên cứu sinh:** Nguyễn Dương Mạnh Sinh  
**Mã đề tài:** Khóa luận Tốt nghiệp Kỹ sư An toàn Thông tin  
**Tiêu chuẩn áp dụng:** ISO/IEC/IEEE 12207 (Software Lifecycle Processes), NIST SP 800-218 SSDF (Secure Software Development Framework), SLSA Level 3.

---

## Mục Lục Các Phiên Làm Việc (Session Index)

| Session ID | Ngày (YYYY-MM-DD) | Phạm Vi / Chủ Đề Trọng Tâm | Trạng Thái | Commits / Artifacts Chính |
|:---|:---|:---|:---:|:---|
| **SES-001** | `2026-04-14` | Baseline Hardening, Cosign Image Signing & Kyverno Admission Gates | **HOÀN TẤT** | `demo/evidence/20260414-210227/` (Allow/Deny matrix) |
| **SES-002** | `2026-05-13` | Mở rộng quy mô Core Trading (tăng từ 4 lên 10 Microservices) | **HOÀN TẤT** | `10dd63c`, `docs/changes/2026-05-13-scale-expansion-handover.md` |
| **SES-003** | `2026-06-01` | Hoàn tất mở rộng 23 Microservices, chuẩn hóa Go `1.25.10`/`1.25.11` | **HOÀN TẤT** | `2821617`, `7930650`, `demo/evidence/20260601-kind-bootstrap/` |
| **SES-004** | `2026-09-06` – `2026-09-08` | Cột mốc Hadolint, Chụp ảnh thực tế Windows 11 OS, Đồng bộ Báo cáo Khóa luận | **HOÀN TẤT** | `8030acf`, `CyberDev_Experimental_Feasibility_Report.docx` |
| **SES-005** | `2026-09-20` (Hiện tại) | Nghiên cứu Dependency Firewall & Kiểm toán quét sâu từng dòng mã nguồn 23 services | **ĐANG THỰC HIỆN** | `docs/WORK_SESSION_LOG.md`, `services/` deep audit |

---

## Chi Tiết Các Phiên Làm Việc (Session Logs)

### [SES-001] Baseline Hardening, Cosign Image Signing & Kyverno Admission Gates
* **Thời gian:** 2026-04-14 08:30 – 21:30 UTC+7
* **Nhân sự thực hiện:** Nguyễn Dương Mạnh Sinh (Pair programming với AI Assistant)
* **Mục tiêu phiên:**
  1. Đóng các Issue nền tảng từ #1 đến #8, #10.
  2. Triển khai chữ ký số container images bằng Cosign và sinh provenance attestation theo chuẩn SLSA.
  3. Cấu hình Kyverno ClusterPolicies để cưỡng chế bảo mật tại cổng nhập viện (Kubernetes Admission Controller).
* **Kết quả & Bằng chứng thực nghiệm:**
  - Thiết lập thành công cụm Kind local với Kyverno v1.11+.
  - Chạy thực nghiệm 5 kịch bản kiểm thử:
    1. `VALID_ALLOW`: Image hợp lệ, có chữ ký số, có SBOM, CVE=0 -> Được phép tạo Pod.
    2. `NEG_UNSIGNED_DENY`: Image không có chữ ký Cosign hợp lệ -> Kyverno chặn (HTTP 403 / Admission Webhook Deny).
    3. `NEG_MISSING_SBOM_DENY`: Image thiếu SBOM chứng nhận -> Kyverno chặn.
    4. `NEG_CVE_THRESHOLD_DENY`: Image chứa lỗ hổng bảo mật nghiêm trọng vượt ngưỡng -> Kyverno chặn.
    5. `VALID_ALLOW_RECHECK`: Kiểm tra hồi quy với artifact chuẩn -> Pass.
  - Lưu trữ toàn bộ raw logs tại `demo/evidence/20260414-210227/`.

---

### [SES-002] Mở Rộng Quy Mô Core Trading (Tăng Từ 4 Lên 10 Microservices)
* **Thời gian:** 2026-05-13 09:00 – 18:00 UTC+7
* **Nhân sự thực hiện:** Nguyễn Dương Mạnh Sinh
* **Mục tiêu phiên:**
  1. Mở rộng từ 4 services thử nghiệm ban đầu lên 10 core trading services phục vụ mô hình sàn giao dịch chứng khoán.
  2. Bổ sung 6 services nghiệp vụ: `market-data-service`, `pricing-service`, `execution-service`, `settlement-service`, `compliance-service`, `notification-service`.
* **Thay đổi kỹ thuật then chốt:**
  - Mỗi service được cấu trúc chuẩn hóa: `cmd/server/main.go`, `internal/<domain>/`, `/healthz`, endpoint nghiệp vụ POST, unit tests, `go.mod`, `Dockerfile`, Kubernetes base manifests.
  - Cập nhật `services.yaml` quản lý 10 services với `profile_tags`.
  - Tối ưu CI workflow hỗ trợ chế độ `changed-only` khi push PR và `nightly full-matrix`.
  - Ban hành tài liệu bàn giao `docs/changes/2026-05-13-scale-expansion-handover.md`.

---

### [SES-003] Mở Rộng Toàn Bộ 23 Microservices & Chuẩn Hóa Toolchain Go 1.25.10 / 1.25.11
* **Thời gian:** 2026-06-01 08:00 – 23:30 UTC+7
* **Nhân sự thực hiện:** Nguyễn Dương Mạnh Sinh
* **Mục tiêu phiên:**
  1. Mở rộng thêm 13 microservices mở rộng: `apikey`, `kyc`, `watchlist`, `analytics`, `audit`, `fees`, `reporting`, `gateway`, `search`, `alert`, `data-feed`, `backtest`, `margin` -> Nâng tổng số microservices lên con số **23 services**.
  2. Đồng bộ hóa tuyệt đối Go toolchain lên phiên bản bảo mật cao nhất `1.25.10` (sau đó nâng lên `1.25.11` theo security patch).
  3. Hoàn tất kịch bản bootstrap cụm Kind tự động hóa và đóng Issue #9.
* **Kết quả & Bằng chứng thực nghiệm:**
  - Chạy thành công matrix CI cross-OS (Ubuntu, macOS) cho 23 services tại GitHub Run `26732257799`.
  - Cập nhật nhật ký bằng chứng `docs/traceability_evidence_register.md` và `docs/final_gap_closing_checklist.md`.
  - Chạy sạch bootstrap tại `demo/evidence/20260601-kind-bootstrap/`.

---

### [SES-004] Cột Mốc Hadolint & Bằng Chứng Thực Nghiệm Màn Hình Windows 11 OS
* **Thời gian:** 2026-09-06 – 2026-09-08
* **Nhân sự thực hiện:** Nguyễn Dương Mạnh Sinh (AI Agent Antigravity hỗ trợ)
* **Mục tiêu phiên:**
  1. Loại bỏ 100% hình ảnh đồ họa giả lập (synthetic matplotlib) theo chỉ thị kỷ luật Zero-Hallucination của đề tài.
  2. Chụp trực tiếp màn hình thực tế Windows 11 Windows Terminal khi chạy lệnh Hadolint quét Dockerfile trước và sau khi gia cố.
  3. Tích hợp phân tích Hadolint vào Báo cáo Khả thi Thực nghiệm (Feasibility Report) 10 trụ cột.
* **Kết quả & Bằng chứng thực nghiệm:**
  - Tạo các ảnh chụp màn hình OS thực tế:
    - `docs/images/real_terminal_hadolint_run.png`: Terminal thực thi lệnh quét Hadolint trên môi trường Windows.
    - `docs/images/real_terminal_hadolint_violations.png`: Chi tiết các vi phạm CIS Docker Benchmark bị phát hiện.
  - Biên dịch và đồng bộ hóa thành công báo cáo đa định dạng trên cả 2 repository (`thesis-microservices` và `CyberDev`):
    - `docs/CyberDev_Experimental_Feasibility_Report.md` (251 KB)
    - `docs/CyberDev_Experimental_Feasibility_Report.docx` (20.37 MB)
    - `docs/CyberDev_Experimental_Feasibility_Report.pdf` (22.85 MB)
  - Commit mã nguồn: `8030acf` trên nhánh `feat/poc-devguard-23-services`.

---

### [SES-005] Nghiên Cứu Dependency Firewall & Kiểm Toán Quét Sâu Từng Tệp/Dòng Code 23 Services
* **Thời gian:** 2026-09-20 14:40 – Đang diễn ra UTC+7
* **Nhân sự thực hiện:** Nguyễn Dương Mạnh Sinh (AI Agent Antigravity hỗ trợ)
* **Mục tiêu phiên:**
  1. Khảo sát kiến trúc DevGuard Dependency Firewall (`controllers/dependencyfirewall/` trong CyberDev), cơ chế proxy `GOPROXY` chặn đứng mã độc trước khi tải về máy (pre-ingestion quarantine vs. post-download SCA).
  2. Thiết lập tệp `docs/WORK_SESSION_LOG.md` chuẩn mực để lưu trữ toàn bộ lịch sử các phiên làm việc của dự án.
  3. Thực hiện kế hoạch 5 giai đoạn quét sâu từng tệp, từng dòng code trong repository 23 microservices để kiểm toán toàn diện mã nguồn, cấu hình và an ninh chuỗi cung ứng.
* **Các bước triển khai:**
  - [x] Khởi tạo tệp `docs/WORK_SESSION_LOG.md` và tái cấu trúc lịch sử từ SES-001 đến SES-005.
  - [x] Quét sâu các tệp cấu hình tầng gốc (`services.yaml`, `.hadolint.yaml`, `Dockerfile.scan-test`, `start-antigravity.ps1`, `.gitignore`, `.gitattributes`).
  - [x] Quét toàn diện mã nguồn 23 Go Microservices trong `services/` (10 Core Trading Services + 13 Extended Services, 353 tệp, 161 tệp Go, 55 tệp Unit tests).
  - [x] Kiểm tra tuân thủ Hadolint (CIS Docker Benchmark): 24/24 Dockerfile production đạt 0 vi phạm (CLEAN).
  - [x] Quét hạ tầng Kubernetes: 3 Kyverno ClusterPolicies, NetworkPolicy cô lập Air-Gapped Zero-Trust và 95 K8s base manifests (100% `securityContext` hardened).
  - [x] Quét toàn bộ 10 quy trình CI/CD workflows trong `.github/workflows/` (2.288 dòng mã pipeline).
  - [x] Quét 9 kịch bản thực nghiệm Python trong `scripts/` (2.503 dòng mã).
  - [x] Đối soát chứng cứ thực nghiệm tại `demo/evidence/` và báo cáo luận văn 251 KB tại `docs/`.
  - [x] Rà soát và phân loại các tệp tạm trong `tmp/`, root và 4 thư mục rỗng cũ.
  - [x] Ban hành Báo cáo Kiểm toán Tổng thể Toàn diện: `docs/AUDIT_MASTER_REPORT.md`.
