# Bảng Theo Dõi Nhiệm Vụ: Dependency Firewall (DevGuard)

- [x] **Giai Đoạn 1: Thiết Kế Kiến Trúc & Đặc Tả Chính Sách (Architecture & Policy Spec)** *(Hoàn thành: 2026-09-19)*
  - [x] Phân tích giao thức GOPROXY của Go toolchain (`@v/list`, `.info`, `.mod`, `.zip`) và cách DevGuard intercept.
  - [x] Đặc tả 4 kịch bản kiểm thử (Rule Blocking, Malicious Check, Cooldown Quarantine, Pass-through LRU Cache).
  - [x] Thiết lập bảng ma trận đánh giá A/B Testing giữa SCA Thụ Động vs. Dependency Firewall Chủ Động.

- [x] **Giai Đoạn 2: Kích Hoạt DevGuard Proxy & Kiểm Thử 3 Tầng Chặn (PoC & Local Harness)** *(Hoàn thành: 2026-09-20T17:15)*
  - [x] Kiểm tra và chạy unit tests của package `CyberDev/core/controllers/dependencyfirewall/` (27/27 tests PASS).
  - [x] Thiết lập harness chạy DevGuard Dependency Firewall proxy cục bộ trên cổng `:8080`.
  - [x] Kiểm thử kịch bản 1: Chặn package theo Rules Pattern (`pkg:go/github.com/sirupsen/logrus*`).
  - [x] Kiểm thử kịch bản 2: Chặn package chứa mã độc trong CSDL OSV (`boltdb-go/bolt`, `MAL-2025-2545`).
  - [x] Kiểm thử kịch bản 3: Chặn package trong thời gian cách ly Cooldown (`MinReleaseAge = 48h`).
  - [x] Chụp ảnh thực tế giao diện DevGuard Web UI và PowerShell terminal 100% (không dùng ảnh trống / synthetic).

- [x] **Giai Đoạn 3: Tích Hợp Thực Nghiệm Trên Go Microservices (Microservice Integration & Benchmark)** *(Hoàn thành: 2026-09-20T17:46)*
  - [x] Cấu hình `GOPROXY=http://localhost:8080/api/v1/dependency-proxy/{secret}/go,off` cho `user-service`.
  - [x] Thực hiện lệnh `go mod download -x` với dependency độc hại -> Bắt lỗi HTTP 403 Forbidden trực tiếp tại terminal.
  - [x] Đo lường thời gian tải và hit rate của cache disk LRU (`X-Cache: HIT` vs `MISS` — tăng tốc 70.5 lần).
  - [x] Chụp ảnh thực tế terminal màn hình PowerShell khi lệnh tải Go bị chặn bởi DevGuard Dependency Firewall.

- [x] **Giai Đoạn 4: Chuẩn Hóa Dockerfile CI/CD & Cập Nhật Báo Cáo Luận Văn** *(Hoàn thành: 2026-09-20T18:00)*
  - [x] Tích hợp `ARG GOPROXY` vào multi-stage Dockerfile của microservices.
  - [x] Soạn thảo Chương/Mục 8 trong Báo cáo Khả thi (`CyberDev_Experimental_Feasibility_Report.md`).
  - [x] Chèn 8 ảnh chụp thực tế vào thư viện bằng chứng `docs/evidence/` và cập nhật gallery minh chứng.
  - [x] Đồng bộ hóa sang kho `CyberDev` và `thesis-microservices`, thực hiện Git commit sẵn sàng push.
