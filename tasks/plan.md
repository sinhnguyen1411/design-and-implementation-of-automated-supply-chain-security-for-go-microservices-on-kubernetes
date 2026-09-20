# Kế Hoạch Nghiên Cứu & Triển Khai Thực Nghiệm: Dependency Firewall (DevGuard) Cho 23 Go Microservices

## 1. Mục Tiêu & Cơ Sở Lý Luận
- **Vấn đề cốt lõi:** Các công cụ quét mã nguồn tĩnh (SCA - Software Composition Analysis) truyền thống (Trivy, Snyk, Grype) hoạt động theo mô hình *thụ động (passive / post-download)*. Chúng chỉ phân tích sau khi dependency đã được tải về máy hoặc runner CI. Với các cuộc tấn công chuỗi cung ứng hiện đại (Typosquatting, Malicious Release, Account Takeover có payload chạy ngay khi resolve/install), việc kiểm tra sau tải là quá muộn.
- **Giải pháp - Dependency Firewall (Active / Pre-ingestion Quarantine Proxy):** Đóng vai trò là Inline Network Gatekeeper đứng giữa máy phát triển/CI runner và Upstream Registries (`proxy.golang.org`, Docker Hub). Proxy chặn đứng yêu cầu (`HTTP 403 Forbidden`) TRƯỚC KHI bất kỳ byte mã độc nào chạm vào đĩa cứng.
- **Cơ chế sẵn có trong DevGuard (`CyberDev/core/controllers/dependencyfirewall/`):**
  1. Go Module Proxy Protocol (`@v/list`, `.info`, `.mod`, `.zip`).
  2. Policy Rule Engine (Gitignore-style pattern matching, allowlist/denylist với `!`).
  3. Real-time Malicious Package Check (truy vấn CSDL OSV/GitHub Advisories).
  4. Quarantine / Cooldown Window (`MinReleaseAge` theo giờ để chống Zero-Day maintainer takeover).
  5. Disk-backed LRU Caching với TTL và header `X-Cache: HIT / MISS`.
  6. Scoped Secret Authentication cô lập đa tổ chức/dự án.

---

## 2. Kế Hoạch 4 Giai Đoạn

### Giai Đoạn 1: Thiết Kế & Chuẩn Hóa Kịch Bản Kiểm Thử (Policy & Test Spec)
- Xây dựng ma trận kiểm thử (Test Matrix) cho 4 kịch bản: Rule Blocking, Malicious Package Blocking, Cooldown Quarantine, Pass-through Caching.

### Giai Đoạn 2: Kích Hoạt DevGuard Proxy & Kiểm Thử 3 Tầng Chặn (Local Harness & PoC)
- Chạy unit/integration test của DevGuard Dependency Firewall.
- Tạo test harness mô phỏng/chạy proxy DevGuard và kiểm tra phản hồi HTTP 403/200.
- Thu thập bằng chứng thực nghiệm bằng ảnh chụp desktop Windows 11 Windows Terminal thật 100%.

### Giai Đoạn 3: Tích Hợp Thực Nghiệm Trên Go Microservices & Đo Đạc Hiệu Năng
- Cấu hình biến môi trường `GOPROXY` trỏ qua DevGuard cho microservice (`user-service`).
- Thử nghiệm A/B tải package cấm vs package sạch.
- Đo lường độ trễ mạng và hiệu năng cache (`X-Cache: HIT` vs `MISS`).
- Chụp ảnh màn hình thật khi lệnh `go get` / `go mod download` bị firewall chặn đứng.

### Giai Đoạn 4: Đóng Gói Dockerfile CI/CD & Cập Nhật Báo Cáo Luận Văn
- Chuẩn hóa Dockerfile hardened với `GOPROXY` cho 23 microservices.
- Soạn thảo và cập nhật Mục 8 trong Báo cáo Khả thi (`CyberDev_Experimental_Feasibility_Report.md`, `.docx`, `.pdf`).
- Đồng bộ hóa và Git commit trên cả hai kho lưu trữ.
