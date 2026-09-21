# Kế Hoạch Tích Hợp Lâu Dài: Chuẩn Hóa Văn Phong Tiếng Việt Học Thuật & Tự Động Hóa Xuất Bản Đa Định Dạng

Tài liệu này xác lập quy chuẩn kỹ thuật, cấu trúc công cụ, và quy trình vận hành vĩnh viễn nhằm duy trì tính chuẩn xác học thuật, tính toàn vẹn số liệu thực nghiệm, và tự động hóa xuất bản đa định dạng (`.md`, `.docx`, `.pdf`) cho toàn bộ tài liệu kỹ thuật của dự án **DevGuard** trên cả hai kho lưu trữ:
1. **Repository Nền tảng:** `CyberDev`
2. **Repository Luận văn K8s:** `design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes`

---

## 1. Mục Tiêu & Nguyên Tắc Cốt Lõi

1. **Zero Blacklist Violations (Không Khoan Nhượng Khẩu Ngữ):**
   - Tuyệt đối không để lọt các từ ngữ dịch thô, tiếng lóng kỹ thuật hoặc khẩu ngữ đời thường (ví dụ: `"tường lửa thư viện"`, `"máy dev"`, `"lọt lộ"`, `"chặn oan"`, `"ép chạy xanh"`, `"bắt chẹt"`).
   - Giữ nguyên các thuật ngữ tiếng Anh chuẩn quốc tế đã định nghĩa trong [`docs/VIETNAMESE_TECHNICAL_STYLE_GUIDE.md`](./VIETNAMESE_TECHNICAL_STYLE_GUIDE.md) (*Dependency Firewall*, *Zero-Trust*, *Air-Gapped*, *CycloneDX SBOM*, *VEX*, *SLSA Provenance*, *Admission Controller*, *In-line Proxy*).

2. **Single Source of Truth & 100% Fact-Preserving:**
   - Tệp Markdown gốc ([`docs/CyberDev_Experimental_Feasibility_Report.md`](./CyberDev_Experimental_Feasibility_Report.md)) là nguồn chân lý duy nhất.
   - Các ấn phẩm định dạng DOCX và PDF, cùng các tài nguyên hình ảnh và telemetry JSON trên repository luận văn, đều được tự động đồng bộ phái sinh từ nguồn này.
   - Tuyệt đối không thay đổi bất kỳ số liệu thực nghiệm, thời gian benchmark, mã CVE hay cấu trúc bảng biểu kỹ thuật nào trong quá trình biên tập câu từ.

3. **Multi-Layer Enforcement (3 Tầng Kiểm Soát Tự Động):**
   - **Tầng 1 (Soạn thảo & Agent):** Hướng dẫn phong cách kỹ thuật + Agent Skills (`style-guide-vi`, `translationese-cleaner-vi`).
   - **Tầng 2 (Cục bộ):** Git Pre-commit Hook tự động ngăn chặn commit sai quy chuẩn.
   - **Tầng 3 (Hệ thống):** CI/CD GitHub Actions Quality Gate tự động kiểm tra trên mọi Pull Request và Push.

---

## 2. Kiến Trúc Bộ Công Cụ (Tooling Stack)

| Thành phần | Đường dẫn | Chức năng chính |
| :--- | :--- | :--- |
| **Quy chuẩn Văn phong** | `docs/VIETNAMESE_TECHNICAL_STYLE_GUIDE.md` | Bộ quy tắc thuật ngữ, cấu trúc câu, danh sách Blacklist/Whitelist và chuẩn typography tiếng Việt học thuật. |
| **Bộ Linter Tự Động** | `scripts/lint_vietnamese_report.py` | Kiểm tra tự động lỗi từ khóa Blacklist, quy chuẩn dấu câu, hỗ trợ chế độ `--fix` và `--github-annotation`. |
| **Bộ Xuất DOCX Chuyên Dụng** | `scripts/export_docx_report.py` | Tạo tài liệu Microsoft Word với bảng biểu chuẩn, tô màu callout, căn lề và định dạng code block học thuật. |
| **Bộ Xuất PDF Vector** | `scripts/export_feasibility_report.py` | Render HTML trung gian sang PDF vector phân trang sắc nét bằng Microsoft Edge Headless. |
| **Bộ Đồng Bộ Đa Repo** | `scripts/sync_to_thesis_repo.py` | Đồng bộ toàn bộ Markdown, DOCX, PDF, telemetry JSON, thư mục khóa và hình ảnh sang repo luận văn K8s. |
| **Bộ Kiểm Toán & Xác Thực** | `scripts/verify_multi_format_sync.py` | Kiểm tra tính nhất quán, đo kích thước tệp, đếm số mục, số bảng biểu và 10 tệp telemetry JSON. |
| **Pipeline Một Chạm** | `scripts/build_and_sync_all_reports.py` | Điều phối toàn bộ quy trình: *Lint $\rightarrow$ DOCX $\rightarrow$ PDF $\rightarrow$ Sync $\rightarrow$ Verify* chỉ với 1 lệnh duy nhất. |
| **Git Pre-commit Hook** | `scripts/install_git_hooks.py` | Tự động cài đặt hook vào `.git/hooks/pre-commit` để bảo vệ kho mã trước khi commit. |
| **GitHub Actions CI Gate** | `.github/workflows/docs-vietnamese-quality-gate.yml` | Quality gate tự động kích hoạt trên GitHub CI. |

---

## 3. Quy Trình Vận Hành Thường Nhật (Developer Workflow)

### 3.1. Khi Chỉnh Sửa Tài Liệu Báo Cáo
1. Chỉnh sửa nội dung trực tiếp tại `docs/CyberDev_Experimental_Feasibility_Report.md`.
2. Kiểm tra nhanh bằng bộ linter:
   ```bash
   python scripts/lint_vietnamese_report.py docs/CyberDev_Experimental_Feasibility_Report.md
   ```
3. Nếu có cảnh báo typography hoặc từ khóa đơn giản, có thể dùng chế độ tự động sửa:
   ```bash
   python scripts/lint_vietnamese_report.py docs/CyberDev_Experimental_Feasibility_Report.md --fix
   ```

### 3.2. Khi Cần Xuất Bản & Đồng Bộ Toàn Bộ
Chỉ cần thực thi một lệnh duy nhất:
```bash
python scripts/build_and_sync_all_reports.py
```
Quy trình sẽ tự động:
- [x] Kiểm tra độ sạch của văn phong (báo lỗi nếu có vi phạm).
- [x] Xuất bản tệp `.docx` mới nhất (khoảng ~21.9 MB).
- [x] Xuất bản tệp `.pdf` mới nhất (khoảng ~24.2 MB).
- [x] Đồng bộ sang repository luận văn K8s.
- [x] Xác thực đối soát chéo và báo cáo trạng thái hoàn thành.

### 3.3. Cài Đặt Git Hook Bảo Vệ Cục Bộ
Chạy script cài đặt một lần trên mỗi máy phát triển:
```bash
python scripts/install_git_hooks.py
```
Sau khi cài đặt, mỗi khi thực hiện `git commit` có chứa thay đổi trong thư mục `docs/`, hệ thống sẽ tự động quét. Nếu phát hiện vi phạm Blacklist, quá trình commit sẽ bị dừng lại và thông báo vị trí dòng vi phạm để sửa.

---

## 4. Lộ Trình Nâng Cấp Kỹ Thuật Dài Hạn

1. **Giai đoạn 1 (Hiện tại - Đã Hoàn Thành):**
   - Xây dựng bộ quy chuẩn phong cách và từ điển Blacklist chuyên ngành an ninh chuỗi cung ứng.
   - Hoàn thành linter cơ bản phát hiện từ cấm và typography.
   - Tự động hóa xuất DOCX và Edge Headless PDF.
   - Đồng bộ hóa 100% hai repository.

2. **Giai đoạn 2 (Tích hợp CI/CD & Hooks - Đang Triển Khai):**
   - Cài đặt Pre-commit hook tự động.
   - Thiết lập GitHub Actions Workflow làm Quality Gate bắt buộc trước khi merge PR.
   - Hỗ trợ chế độ `--fix` an toàn trong linter.

3. **Giai đoạn 3 (Mở Rộng NLP & Chấm Điểm Học Thuật):**
   - Tích hợp thêm thư viện kiểm tra hình vị tiếng Việt (`pyvi` hoặc từ điển Hunspell tiếng Việt) để phát hiện lỗi gõ telex sót.
   - Xây dựng chỉ số đo lường độ mạch lạc học thuật (Academic Prose Density Score).
