Nội dung đề tài:

I. Tổng quan đề tài

1. Bối cảnh và vấn đề đặt ra  
Chuỗi cung ứng phần mềm đang trở thành mục tiêu tấn công trọng yếu, khi rủi ro không chỉ nằm ở mã nguồn mà còn ở dependency, quy trình build và artifact trước khi triển khai. Các sự cố như SolarWinds cho thấy chỉ một mắt xích trong pipeline bị xâm phạm cũng có thể gây hậu quả diện rộng. Cùng với đó, số vụ tấn công chuỗi cung ứng phần mềm tiếp tục gia tăng, khiến bảo mật chuỗi cung ứng trở thành yêu cầu cấp thiết trong phát triển phần mềm hiện đại [2][5]. Executive Order 14028 và các hướng dẫn liên quan cũng thúc đẩy mạnh việc minh bạch thành phần phần mềm thông qua SBOM [3].

2. Các chuẩn, khung và công cụ đã có trên thị trường  
Hiện nay đã có các nền tảng và công cụ hỗ trợ bảo mật chuỗi cung ứng phần mềm. SLSA cung cấp khung bảo đảm tính toàn vẹn của artifact, provenance và khả năng truy vết quá trình build [4]. SBOM hỗ trợ minh bạch thành phần và phụ thuộc của phần mềm, phục vụ quản lý rủi ro và lỗ hổng [3]. Bên cạnh đó, các kỹ thuật như quét CVE, sinh SBOM, ký image và xác minh artifact đã có thể tích hợp vào CI/CD pipeline [1][5]. Tuy nhiên, khó khăn lớn hiện nay không phải là thiếu công cụ, mà là thiếu cách tích hợp chúng thành một quy trình thống nhất và có khả năng kiểm chứng hiệu quả trong thực tế [1][2][4].

3. Thực trạng trước thời điểm hiện tại  
Mặc dù đã có chuẩn, chính sách định hướng và công cụ kỹ thuật tương đối đầy đủ, việc triển khai trong thực tế vẫn còn rời rạc. Nhiều tổ chức mới dừng ở mức quét lỗ hổng, tạo SBOM hoặc ký artifact như các bước riêng lẻ, chưa hình thành một pipeline end-to-end có khả năng kiểm chứng và cưỡng chế trước khi triển khai [1][2][3][4]. Đối với các hệ thống microservices trên Kubernetes, khoảng trống này càng rõ do áp lực phát hành nhanh và thiếu mô hình tích hợp thực dụng, dễ áp dụng. Vì vậy, vấn đề trọng tâm hiện nay là xây dựng một pipeline tự động có thể kết hợp build, scan, SBOM, ký và kiểm chứng để ngăn artifact không đạt yêu cầu đi vào môi trường chạy thực tế [1][2][4][5].

Kết luận tổng quan: đề tài hướng tới xây dựng một baseline DevSecOps cho Golang microservices trên Kubernetes có khả năng kiểm chứng và cưỡng chế, thay vì chỉ dừng ở mức quét và cảnh báo.

II. Mục tiêu của đề tài

1. Mục tiêu tổng quát  
Xây dựng một pipeline bảo mật chuỗi cung ứng phần mềm tự động cho microservice Golang triển khai trên Kubernetes, nhằm kiểm soát dependency, xác minh nguồn gốc và tính toàn vẹn của container image, đánh giá rủi ro bảo mật trước khi triển khai, và chủ động từ chối các artifact không đáp ứng yêu cầu. Mục tiêu của đề tài không chỉ là phát hiện nguy cơ, mà là bảo đảm artifact không an toàn không thể chạy trong cụm [1], [3], [4], [6]-[10].

2. Mục tiêu cụ thể  
1) Kiểm soát thành phần và dependency của dịch vụ Golang thông qua Go Modules, `go.sum` và checksum database, đồng thời tự động sinh SBOM cho mỗi bản build để phục vụ truy vết và kiểm kê thành phần phần mềm [3], [6], [7].  
2) Tự động phát hiện và loại bỏ rủi ro trước phát hành bằng cách quét lỗ hổng ở dependency/mã Go với `govulncheck`, kết hợp quét lỗ hổng container image; nếu vượt ngưỡng cho phép thì pipeline dừng và artifact không được phát hành [1], [5], [6].  
3) Ký số container image và tạo attestation/provenance cho từng bản build để mỗi artifact đều có thể xác minh được nguồn gốc, gắn với commit và quy trình build hợp lệ, qua đó hạn chế việc tạo image ngoài luồng kiểm soát [1], [4], [8].  
4) Cưỡng chế chính sách bảo mật tại Kubernetes Admission để cụm từ chối triển khai image chưa ký, thiếu attestation/SBOM hợp lệ hoặc không đáp ứng chính sách bảo mật đã xác lập [8]-[10].  
5) Chuẩn hóa toàn bộ quy trình thành pipeline CI/CD có thể lặp lại gồm các bước build, sinh SBOM, quét lỗ hổng, ký, đẩy image và kiểm chứng khi triển khai, từ đó có thể tái sử dụng cho các microservice Golang khác [1], [2], [4], [6], [8]-[10].  
6) Minh họa trên một microservice Golang thực tế bằng cách triển khai một user-service mẫu trên Kubernetes để kiểm chứng tính khả thi của mô hình theo luồng end-to-end từ commit đến deploy.

3. Mục tiêu cải tiến so với thực trạng hiện tại  
Đề tài hướng tới tích hợp các cơ chế bảo mật đang được sử dụng rời rạc thành một chuỗi kiểm soát thống nhất từ dependency, build đến deploy. Điểm cải tiến chính là chuyển từ mô hình chỉ quét và cảnh báo sang mô hình có khả năng kiểm chứng và cưỡng chế, trong đó SBOM, quét lỗ hổng, ký image và attestation trở thành điều kiện bắt buộc trước khi artifact được chạy trên Kubernetes. Kết quả đầu ra không chỉ là mô hình lý thuyết mà còn là pipeline, cấu hình và quy trình có thể tái lập trong thực tế [1], [2], [4], [6], [8]-[10].

Kết luận mục tiêu: đề tài hướng tới xây dựng một baseline DevSecOps cho Golang microservices trên Kubernetes có khả năng kiểm chứng và cưỡng chế, thay vì chỉ dừng ở mức quét và cảnh báo [6], [8]-[10].

III. Nội dung và phương pháp thực hiện

Đề tài được triển khai theo hướng thiết kế -> tích hợp -> thực nghiệm -> đánh giá, nhằm xây dựng và kiểm chứng một pipeline bảo mật chuỗi cung ứng phần mềm cho microservice Golang trên Kubernetes, thay vì chỉ phân tích lý thuyết [1], [2], [4], [6], [8], [9], [10].

1. Phân tích yêu cầu và mô hình hóa  
Xác định các rủi ro chính trong chuỗi cung ứng phần mềm: dependency bên thứ ba, lỗ hổng image, artifact không rõ nguồn gốc và triển khai trái phép. Từ đó, chuyển hóa thành các yêu cầu kỹ thuật gồm: SBOM, quét lỗ hổng, ký image, attestation/provenance và policy enforcement trên Kubernetes.  
Luồng tổng thể được mô hình hóa theo chuỗi: Dev -> CI/CD -> Registry -> Kubernetes Admission [2], [3], [4], [6], [10].

2. Thiết kế kiến trúc giải pháp  
Thiết kế pipeline bảo mật theo hướng tích hợp các bước kiểm soát vào toàn bộ vòng đời artifact:  
Build -> SBOM -> Scan -> Sign/Attest -> Push -> Verify/Enforce.  
Đồng thời lựa chọn và kết nối các công cụ mã nguồn mở để bảo đảm artifact được kiểm tra theo cùng một chuẩn từ CI/CD đến Kubernetes Admission [1], [4], [6], [8], [9], [10].

3. Xây dựng microservice Golang mẫu  
Sử dụng một microservice Golang làm đối tượng thực nghiệm, ví dụ user-service với các chức năng cơ bản như đăng ký, xác thực và đăng nhập. Dịch vụ được đóng gói bằng Docker multi-stage và hạn chế đặc quyền khi chạy container để phù hợp với yêu cầu kiểm thử bảo mật [1], [6], [7].

4. Tích hợp pipeline CI/CD bảo mật  
Xây dựng pipeline tự động thực hiện chuỗi:  
Code -> Build -> SBOM -> Vulnerability Scan -> Sign/Attest -> Push.  
Pipeline được cấu hình fail-fast, nghĩa là dừng ngay khi phát hiện vi phạm bảo mật vượt ngưỡng, nhằm ngăn artifact không đạt chuẩn được phát hành [1], [4], [6], [8].

5. Thiết lập Kubernetes và cơ chế cưỡng chế  
Triển khai cụm thử nghiệm bằng Kind/Minikube, sau đó cấu hình admission control để kiểm tra artifact tại thời điểm deploy. Cơ chế này được mô hình hóa theo quan hệ: CI/CD <-> Registry <-> Kubernetes Admission.  
Chỉ những image đã ký, có attestation hợp lệ và đáp ứng chính sách bảo mật mới được phép chạy trong cụm [8], [9], [10].

6. Kiểm thử  
Thực hiện hai nhóm kịch bản:  
- Artifact không hợp lệ -> bị từ chối triển khai  
- Artifact hợp lệ -> được chấp nhận triển khai  

Kết quả được đánh giá thông qua log CI/CD, metadata artifact và sự kiện Admission Denied/Allowed để kiểm chứng tính đúng đắn, mức tự động hóa và khả năng ngăn chặn artifact rủi ro [1], [2], [8], [9], [10].

7. Chuẩn hóa tài liệu và tổng kết  
Tổng hợp kiến trúc, pipeline, policy Kubernetes và quy trình triển khai thành tài liệu có thể tái lập; đồng thời nêu giới hạn và hướng mở rộng của mô hình [1], [2], [4].

Kết luận: phương pháp thực hiện của đề tài hướng tới một mô hình DevSecOps/Supply Chain Security có thể chạy được, kiểm chứng được và cưỡng chế được, trong đó các bước bảo mật được tích hợp thống nhất từ Dev -> CI/CD -> Registry -> Kubernetes Admission [1], [4], [8], [9], [10].

IV. Phương pháp đánh giá

Hệ thống được đánh giá bằng thực nghiệm có kiểm soát thông qua các kịch bản triển khai artifact hợp lệ và không hợp lệ, thay vì chỉ phân tích lý thuyết. Các tiêu chí đánh giá chính gồm:

- Đánh giá pipeline: khả năng phát hiện và chặn artifact rủi ro ngay trong giai đoạn Build -> Scan -> Sign/Attest [1], [6], [8].  
- Đánh giá tính toàn vẹn: khả năng xác minh chữ ký, provenance/attestation và nguồn gốc của container image [4], [8].  
- Đánh giá enforcement: khả năng Kubernetes tự động từ chối các Pod vi phạm tại bước Registry -> Kubernetes Admission -> Deploy [8], [9], [10].  
- Đánh giá tính lặp lại: mức độ dễ tái sử dụng pipeline cho các microservice Golang khác với ít chỉnh sửa thủ công [1], [2], [6].

Kết luận: phương pháp đánh giá của đề tài tập trung vào khả năng phát hiện, xác minh và cưỡng chế trong toàn bộ chuỗi Build -> Registry -> Kubernetes Admission -> Deploy, còn phạm vi nghiên cứu được giới hạn ở mức đủ để kiểm chứng tính khả thi và khả năng tái lập của mô hình [1], [4], [8]-[10].

V. Giới hạn của đề tài

Đề tài tập trung vào việc chứng minh tính khả thi của một pipeline bảo mật chuỗi cung ứng phần mềm cho microservice Golang trên Kubernetes, vì vậy có các giới hạn sau:

- Giới hạn về nghiệp vụ: microservice thử nghiệm chỉ ở mức đơn giản, như quản lý người dùng, xác thực và đăng nhập, không hướng tới hệ thống nghiệp vụ lớn hoặc phức tạp.  
- Giới hạn về phạm vi bảo mật: đề tài tập trung vào Supply Chain Security theo chuỗi Dependency -> Build -> Image -> Admission, không đi sâu vào các bài toán bảo mật ứng dụng như SQL injection, brute force, rate limiting hay logic nghiệp vụ.  
- Giới hạn về môi trường triển khai: mô hình được kiểm thử trên cụm Kubernetes cục bộ như Kind/Minikube, không bắt buộc triển khai trên hạ tầng cloud thương mại.  
- Giới hạn về mức độ tuân thủ chuẩn: đề tài vận dụng tinh thần của SLSA, SBOM, signing và provenance ở mức thực dụng, không nhằm đạt đầy đủ mọi mức tuân thủ ở quy mô doanh nghiệp lớn [2], [4].  
- Giới hạn về hiệu năng: mục tiêu chính là kiểm chứng tính đúng đắn, tính tự động hóa và khả năng cưỡng chế của pipeline, không phải tối ưu hiệu năng CI/CD hay benchmark tải hệ thống.

Tài liệu tham khảo:

[1] D. Patel, "Software supply chain security: Implementing SLSA compliance in CI/CD pipelines," International Journal for Research Trends and Innovation, vol. 10, no. 7, Jan. 2025.  
[2] M. Tamanna, S. Hamer, M. Tran, S. Fahl, Y. Acar, and L. Williams, "Analyzing challenges in deployment of the SLSA framework for software supply chain security," Dec. 2024.  
[3] National Institute of Standards and Technology (NIST), "Improving the Nation's Cybersecurity: NIST's Responsibilities under Executive Order 14028-Software Supply Chain Security Guidance," U.S. Department of Commerce, July 2022.  
[4] The Linux Foundation, "Safeguarding artifact integrity across any software supply chain: What is SLSA?," Open Source Security Foundation, 2025.  
[5] D. I. Jonathan, "Supply chain security in modern software: SBOMs, SLSA, and beyond," EM360Tech, Sept. 3, 2025.  
[6] J. Qiu, "Vulnerability Management for Go," The Go Blog, Sep. 6, 2022.  
[7] K. Hockman, "Module Mirror and Checksum Database Launched," The Go Blog, Aug. 29, 2019.  
[8] Sigstore, "Cosign Quickstart," Sigstore Documentation, accessed Mar. 18, 2026.  
[9] Sigstore, "Policy Controller Overview," Sigstore Documentation, accessed Mar. 18, 2026.  
[10] Kubernetes, "Validating Admission Policy," Kubernetes Documentation, 2024.  
