# Speaker Notes V2 (1 trang) - Thesis Demo

Thời lượng mục tiêu: 10-15 phút.
Phong cách: technical thesis, ngắn gọn, tập trung cơ chế kiểm soát và bằng chứng.

## Opening (20-30 giây)
- "Mục tiêu demo là chứng minh chuỗi trust liên tục: verify trước khi release và enforce trước khi workload được chạy trên Kubernetes."
- "Mọi kết luận đều dựa trên evidence artifact, không chỉ quan sát live."

## Scene 1 - Architecture (1-2 phút)
- Key line: "Đây không phải scan rời rạc, mà là control chain end-to-end từ CI tới admission."
- Chỉ vào các điểm: dependency integrity -> govulncheck -> SBOM -> Grype gate -> Cosign sign/attest -> Kyverno verify.
- Chốt: "Nếu vi phạm ở bất kỳ điểm nào, artifact không được đi tiếp."

## Scene 2 - Pre-flight (1 phút)
- Key line: "Trước khi test case, cluster và policy đã ở trạng thái enforce-ready."
- Nêu ngắn 3 tín hiệu: context đúng, cluster reachable, 3 policy chính đã nạp.
- Tránh đi sâu debug/lý thuyết policy YAML.

## Scene 3 - VALID_ALLOW (2 phút)
- Key line: "Artifact hợp lệ được admission cho chạy."
- Nói khi chờ lệnh: "Case này có đủ signature, provenance, SBOM annotation, và high_critical=0."
- Chốt khi thấy `condition met`: "Đây là đường đi chuẩn được phép."

## Scene 4 - Negative Matrix (4-5 phút)
- Key line: "Admission deny có lý do policy cụ thể, không phải deny mơ hồ."
- `NEG_UNSIGNED_DENY`: "Không có signature -> deny với signal `no signatures found`."
- `NEG_MISSING_SBOM_DENY`: "Thiếu SBOM digest annotation -> deny bởi rule require-sbom-digest."
- `NEG_CVE_THRESHOLD_DENY`: "high_critical != 0 -> deny bởi rule require-high-critical-zero."
- Chốt: "3 deny reasons map 1-1 với policy contract."

## Scene 5 - Regression Recheck (1 phút)
- Key line: "Sau chuỗi deny, artifact hợp lệ vẫn được admit bình thường."
- Chốt khi thấy available: "Hệ thống vừa enforce chặt, vừa không phá đường deploy hợp lệ."

## Scene 6 - Evidence Board (2-3 phút)
- Key line: "Demo này có thể kiểm tra lại độc lập qua evidence bundle."
- Chỉ 3 phần: matrix summary, matrix index JSON, artifact per case.
- Chốt cuối: "Kết quả đạt: 1 allow, 3 deny đúng policy, 1 regression allow pass."

## Q&A anchors (nếu hội đồng hỏi)
- "Vì sao cần cả CI gate và admission?"
  - CI chặn rủi ro sớm; admission chặn bypass/manual deploy.
- "Điểm mới của thesis là gì?"
  - Tích hợp các control rời rạc thành chuỗi verify+enforce có evidence.
- "Có tái dùng cho service khác không?"
  - Có, cùng policy contract và cùng định dạng evidence bundle.

## Backup line khi live lỗi
- "Để giữ timeline, tôi chuyển sang evidence của cùng case trong run-id hiện tại; đây là output đã được script thu tự động và có thể mở từng artifact để đối chiếu."
