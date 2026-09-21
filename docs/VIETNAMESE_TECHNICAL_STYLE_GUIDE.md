# QUY CHUẨN VĂN PHONG KỸ THUẬT & THUẬT NGỮ TIẾNG VIỆT
## Đề tài: Thiết kế và Hiện thực Hệ thống Bảo mật Chuỗi Cung ứng Tự động cho Go Microservices trên Kubernetes (CyberDev)

Tài liệu này đóng vai trò là kim chỉ nam (*Style Guide & Glossary*) bắt buộc áp dụng trong toàn bộ quá trình biên tập, rà soát và hiệu đính câu chữ của Báo cáo Nghiên cứu và Luận văn Tốt nghiệp.

---

## 1. NGUYÊN TẮC BẢO TOÀN DỮ KIỆN KỸ THUẬT (FACT-PRESERVING)
1. **Tuyệt đối không thay đổi số liệu thực nghiệm:**
   - Thời gian phản hồi, độ trễ (`ms`), tỷ lệ phần trăm (`%`), tỷ lệ tăng tốc (ví dụ: `44.6 lần`).
   - Số lượng vi dịch vụ (`23 Go microservices`), số lượng gói phụ thuộc (`37 gói`), số lượng quy tắc, số lượng CVE.
   - Các mốc thời gian, kích thước tệp tin (`MB`, `KB`, `bytes`), mã băm SHA-256 (`h1:`, `sha256:...`).
2. **Bảo toàn nguyên vẹn mã lệnh và định danh hệ thống:**
   - Không dịch tên cờ lệnh, biến môi trường: `GOPROXY`, `GOSUMDB`, `CGO_ENABLED`, `COSIGN_EXPERIMENTAL`.
   - Giữ nguyên cú pháp terminal, câu lệnh CLI: `go mod download`, `cosign verify`, `gitleaks detect`, `trivy image`.
   - Giữ nguyên đường dẫn tệp tin, tên tệp: `go.mod`, `go.sum`, `Dockerfile`, `kyverno.yaml`, `telemetry.json`.
   - Giữ nguyên định dạng chú thích: `Hình X.Y: ...` và `Bảng X.Y: ...`.

---

## 2. BẢNG QUY CHUẨN THUẬT NGỮ QUỐC TẾ (BẮT BUỘC GIỮ TIẾNG ANH)
Các thuật ngữ kỹ thuật chuyên sâu dưới đây **bắt buộc giữ nguyên tiếng Anh** để đảm bảo tính chính xác học thuật, tránh việc dịch máy thô cứng gây hiểu sai bản chất:

| Thuật ngữ Tiếng Anh | Ngữ cảnh sử dụng | Giải thích / Chú giải lần đầu |
| :--- | :--- | :--- |
| **Dependency Firewall** | Tường lửa kiểm soát gói phụ thuộc | Lớp bảo vệ kiểm soát và chặn lọc các gói phụ thuộc độc hại |
| **In-line Proxy** | Cơ chế proxy trung gian trực tiếp | Proxy kiểm soát luồng dữ liệu theo thời gian thực (in-line) |
| **Air-Gapped** | Môi trường mạng cô lập | Mạng cách ly hoàn toàn với Internet công cộng |
| **Zero-Trust** | Kiến trúc không tin cậy mặc định | Nguyên tắc xác thực và kiểm soát liên tục mọi thực thể |
| **Cooldown Period / Quarantine** | Thời gian cách ly gói mới phát hành | Khoảng thời gian trì hoãn để kiểm định an toàn gói mới |
| **Shift-Left** | Chiến lược dịch trái an ninh | Tích hợp kiểm thử bảo mật ngay từ giai đoạn đầu phát triển |
| **SBOM** | Software Bill of Materials | Danh mục thành phần phần mềm nguồn mở |
| **VEX** | Vulnerability Exploitability eXchange | Báo cáo khả năng khai thác lỗ hổng bảo mật |
| **Attestation** | Chứng thư xác thực / Chứng chỉ số | Dữ liệu siêu hình học chứng thực các bước kiểm tra trong CI/CD |
| **Provenance** | Nguồn gốc xuất xứ phần mềm | Thông tin chứng minh nguồn gốc bản build và mã nguồn |
| **Policy Gate** | Cổng kiểm soát chính sách | Chốt chặn kiểm duyệt trước khi nạp ứng dụng vào Kubernetes |
| **Admission Controller** | Bộ kiểm soát nạp tài nguyên Kubernetes | Thành phần chặn và thẩm định request gửi tới Kubernetes API |
| **Mutating / Validating Webhook** | Webhook biến đổi / Thẩm định K8s | Webhook can thiệp và kiểm tra tính hợp lệ của Pod/Deployment |
| **Multi-stage Build** | Đóng gói Docker đa tầng | Kỹ thuật tối ưu hóa dung lượng và loại bỏ công cụ thừa |
| **Microservices** | Kiến trúc vi dịch vụ | Hệ thống phân tán bao gồm nhiều dịch vụ độc lập |
| **Cache Poisoning** | Đầu độc bộ nhớ đệm | Hành vi tấn công làm sai lệch hoặc chèn mã độc vào cache |
| **Repojacking / Typosquatting** | Chiếm đoạt repo / Đặt tên nhầm lẫn | Kỹ thuật tấn công chuỗi cung ứng mạo danh thư viện |

---

## 3. DANH MỤC THUẬT NGỮ TIẾNG VIỆT CHUẨN HÓA
Khi diễn đạt bằng tiếng Việt, bắt buộc dùng các từ ngữ trang trọng, chuẩn xác:

* **Vulnerability** ➔ *Lỗ hổng bảo mật*
* **Malicious Package / Malware** ➔ *Gói thư viện độc hại / Mã độc*
* **Integrity** ➔ *Tính toàn vẹn (dữ liệu / mã nguồn)*
* **Cryptographic Hash** ➔ *Mã băm mật mã (SHA-256)*
* **Local Cache** ➔ *Bộ nhớ đệm nội bộ / Cục bộ*
* **Developer Workstation** ➔ *Máy trạm phát triển*
* **Upstream Registry** ➔ *Kho lưu trữ thượng nguồn (công cộng)*
* **Compliance Posture** ➔ *Trạng thái tuân thủ an ninh*
* **False Positive** ➔ *Cảnh báo giả*
* **Supply Chain Security** ➔ *An ninh chuỗi cung ứng phần mềm*

---

## 4. DANH SÁCH TỪ NGỮ CẤM (BLACKLIST - TRANSLATIONESE)
Tuyệt đối loại bỏ các cụm từ dịch thô vụng về sau đây:

| ❌ TỪ CẤM (BỊ ĐỘNG / DỊCH MÁY THÔ CỨNG) |  THAY THẾ BẰNG |
| :--- | :--- |
| *Tường lửa Thư viện* | **Dependency Firewall** |
| *máy dev / máy trạm dev* | **máy trạm phát triển (Developer Workstation)** |
| *lọt vào máy trạm / lọt vào máy dev* | **vượt qua chốt chặn an ninh / lọt vào mã nguồn** |
| *chạy mất mạng / khi mất mạng* | **trong môi trường mạng cô lập (Air-Gapped)** |
| *được thực thi bởi / được tạo ra bởi* | (Chuyển thành câu chủ động) **hệ thống thực thi / công cụ tạo ra** |
| *có khả năng của việc kiểm tra* | **có khả năng kiểm tra** |
| *tiến hành việc phân tích* | **phân tích** |
| *ủy quyền nội tuyến* | **In-line Proxy** |
| *gói mới ra* | **gói thư viện mới phát hành** |
| *chặn cách ly (mới ra)* | **cách ly theo chính sách Cooldown Period** |
| *làm độc cache* | **đầu độc bộ nhớ đệm (Cache Poisoning)** |

---

## 5. QUY CHUẨN DẤU CÂU VÀ TRÌNH BÀY (TYPOGRAPHY)
1. **Không có khoảng trắng trước dấu câu:**
   - Sai: `gói thư viện ,` | `kết quả :` | `thực nghiệm ;` | `thành công .`
   - Đúng: `gói thư viện,` | `kết quả:` | `thực nghiệm;` | `thành công.`
2. **Khoảng trắng sau dấu câu:**
   - Sau dấu chấm, phẩy, hai chấm, chấm phẩy bắt buộc có 1 dấu cách đơn.
3. **Quy chuẩn dấu ngoặc đơn và ngoặc kép:**
   - Không có khoảng trắng sát mép trong: `( nội dung )` ➔ `(nội dung)`, `" từ khóa "` ➔ `"từ khóa"`.
4. **Viết hoa tiêu đề và nhãn bảng:**
   - Tiêu đề mục: Viết hoa chữ cái đầu mỗi từ chính hoặc viết hoa đầu câu nhất quán.
   - Nhãn minh họa: Luôn in nghiêng chú thích: `*Hình X.Y: Tên hình ảnh.*`, `*Bảng X.Y: Tên bảng số liệu.*`.
