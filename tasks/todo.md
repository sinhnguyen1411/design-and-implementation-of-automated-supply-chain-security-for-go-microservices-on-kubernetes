# Bảng Theo Dõi Nhiệm Vụ: Dependency Firewall (DevGuard)

- [ ] **Giai Đoạn 1: Thiết Kế Kiến Trúc & Đặc Tả Chính Sách (Architecture & Policy Spec)**
  - [ ] Phân tích giao thức GOPROXY của Go toolchain (`@v/list`, `.info`, `.mod`, `.zip`) và cách DevGuard intercept.
  - [ ] Đặc tả 4 kịch bản kiểm thử (Rule Blocking, Malicious Check, Cooldown Quarantine, Pass-through LRU Cache).
  - [ ] Thiết lập bảng ma trận đánh giá A/B Testing giữa SCA Thụ Động vs. Dependency Firewall Chủ Động.

- [ ] **Giai Đoạn 2: Kích Hoạt DevGuard Proxy & Kiểm Thử 3 Tầng Chặn (PoC & Local Harness)**
  - [ ] Kiểm tra và chạy unit tests của package `CyberDev/core/controllers/dependencyfirewall/`.
  - [ ] Thiết lập harness chạy DevGuard Dependency Firewall proxy cục bộ với cổng riêng.
  - [ ] Kiểm thử kịch bản 1: Chặn package theo Rules Pattern (`pkg:golang/github.com/vulnerable-corp/*`).
  - [ ] Kiểm thử kịch bản 2: Chặn package chứa mã độc trong CSDL OSV (`malicious_packages`).
  - [ ] Kiểm thử kịch bản 3: Chặn package trong thời gian cách ly Cooldown (`MinReleaseAge = 24h`).
  - [ ] Chụp ảnh thực tế Windows 11 Windows Terminal 100% không dùng AI/synthetic.

- [ ] **Giai Đoạn 3: Tích Hợp Thực Nghiệm Trên Go Microservices (Microservice Integration & Benchmark)**
  - [ ] Cấu hình `GOPROXY=http://localhost:8080/api/v1/dependency-proxy/{secret}/go,direct` cho `user-service`.
  - [ ] Thực hiện lệnh `go get` / `go mod download` với dependency độc hại -> Bắt lỗi HTTP 403 Forbidden tại terminal.
  - [ ] Đo lường thời gian tải và hit rate của cache disk LRU (`X-Cache: HIT` vs `MISS`).
  - [ ] Chụp ảnh thực tế terminal màn hình Windows khi lệnh build Go bị chặn bởi DevGuard Dependency Firewall.

- [ ] **Giai Đoạn 4: Chuẩn Hóa Dockerfile CI/CD & Cập Nhật Báo Cáo Luận Văn**
  - [ ] Tích hợp `ARG GOPROXY` vào multi-stage Dockerfile của microservices.
  - [ ] Soạn thảo Chương/Mục 8 trong Báo cáo Khả thi (`CyberDev_Experimental_Feasibility_Report.md`).
  - [ ] Chèn ảnh chụp thực tế vào báo cáo, biên dịch ra file Word (`.docx`) và PDF (`.pdf`).
  - [ ] Đồng bộ hóa sang kho `CyberDev` và `thesis-microservices`, thực hiện Git commit & push.
