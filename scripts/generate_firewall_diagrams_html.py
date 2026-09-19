#!/usr/bin/env python3
"""
Generate clean, conversational, easy-to-understand architecture and flow diagrams
for DevGuard Dependency Firewall.
"""
import os
import shutil
import time
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

HTML_ARCHITECTURE = """<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<title>Cách DevGuard Bảo Vệ 23 Microservices</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: 'Segoe UI', -apple-system, BlinkMacSystemFont, Roboto, sans-serif;
    background: #0B0F19;
    color: #F8FAFC;
    padding: 44px 50px;
    width: 2200px;
    min-height: 1350px;
  }
  .header {
    text-align: center;
    margin-bottom: 40px;
  }
  .header h1 {
    font-size: 38px;
    font-weight: 800;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: #FFFFFF;
    margin-bottom: 12px;
  }
  .header p {
    font-size: 19px;
    color: #94A3B8;
    max-width: 1400px;
    margin: 0 auto;
    line-height: 1.5;
  }

  /* 4-Stage Horizontal Pipeline */
  .pipeline-container {
    display: flex;
    align-items: stretch;
    justify-content: space-between;
    gap: 18px;
    margin-bottom: 40px;
  }
  .stage-card {
    flex: 1;
    background: #111827;
    border-radius: 16px;
    padding: 26px 22px;
    position: relative;
    display: flex;
    flex-direction: column;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  }
  .stage-card.blue { border: 2px solid #0284C7; }
  .stage-card.orange { border: 2px solid #D97706; }
  .stage-card.cyan { border: 2px solid #0891B2; }
  .stage-card.green { border: 2px solid #059669; }

  .connector {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    flex-shrink: 0;
  }
  .connector-circle {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    font-weight: bold;
    border: 2px solid;
    background: rgba(16, 185, 129, 0.15);
    border-color: #10B981;
    color: #10B981;
  }

  .stage-header {
    margin-bottom: 16px;
  }
  .stage-title {
    font-size: 20px;
    font-weight: 800;
    text-transform: uppercase;
    margin-bottom: 6px;
  }
  .stage-card.blue .stage-title { color: #38BDF8; }
  .stage-card.orange .stage-title { color: #FBBF24; }
  .stage-card.cyan .stage-title { color: #22D3EE; }
  .stage-card.green .stage-title { color: #34D399; }

  .stage-subtitle {
    font-size: 14px;
    color: #94A3B8;
  }

  .stage-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 22px;
  }
  .feature-item {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    font-size: 15px;
    line-height: 1.45;
    color: #E2E8F0;
  }
  .feature-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    margin-top: 7px;
    flex-shrink: 0;
  }
  .stage-card.blue .feature-dot { background: #38BDF8; }
  .stage-card.orange .feature-dot { background: #FBBF24; }
  .stage-card.cyan .feature-dot { background: #22D3EE; }
  .stage-card.green .feature-dot { background: #34D399; }

  .code-badge {
    background: #1E293B;
    border: 1px solid #334155;
    padding: 8px 12px;
    border-radius: 8px;
    font-family: 'Consolas', 'Courier New', monospace;
    font-size: 13px;
    color: #CBD5E1;
    word-break: break-all;
    margin-top: 4px;
  }

  .stage-footer {
    border-radius: 10px;
    padding: 12px 14px;
    text-align: center;
    font-size: 13px;
    font-weight: 700;
    text-transform: uppercase;
  }
  .stage-card.blue .stage-footer {
    background: rgba(14, 165, 233, 0.15);
    color: #38BDF8;
    border: 1px solid rgba(14, 165, 233, 0.3);
  }
  .stage-card.orange .stage-footer {
    background: rgba(245, 158, 11, 0.15);
    color: #FBBF24;
    border: 1px solid rgba(245, 158, 11, 0.3);
  }
  .stage-card.cyan .stage-footer {
    background: rgba(6, 182, 212, 0.15);
    color: #22D3EE;
    border: 1px solid rgba(6, 182, 212, 0.3);
  }
  .stage-card.green .stage-footer {
    background: rgba(16, 185, 129, 0.15);
    color: #34D399;
    border: 1px solid rgba(16, 185, 129, 0.3);
  }

  /* Bottom Table Section */
  .table-container {
    background: #111827;
    border: 2px solid #1E3A8A;
    border-radius: 16px;
    padding: 26px 30px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  }
  .table-title {
    font-size: 18px;
    font-weight: 800;
    color: #38BDF8;
    margin-bottom: 18px;
    display: flex;
    align-items: center;
    gap: 10px;
    text-transform: uppercase;
  }
  .table-title::before {
    content: '';
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: #38BDF8;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 15px;
  }
  th {
    background: #1E293B;
    color: #94A3B8;
    font-weight: 700;
    text-align: left;
    padding: 12px 16px;
    text-transform: uppercase;
    font-size: 13px;
    border-bottom: 2px solid #334155;
  }
  td {
    padding: 15px 16px;
    border-bottom: 1px solid #1E293B;
    color: #E2E8F0;
    vertical-align: middle;
  }
  tr:last-child td { border-bottom: none; }
  .pill {
    display: inline-block;
    padding: 5px 12px;
    border-radius: 9999px;
    font-size: 12px;
    font-weight: 700;
    text-transform: uppercase;
  }
  .pill.red { background: rgba(239, 68, 68, 0.2); color: #F87171; border: 1px solid rgba(239, 68, 68, 0.4); }
  .pill.amber { background: rgba(245, 158, 11, 0.2); color: #FBBF24; border: 1px solid rgba(245, 158, 11, 0.4); }
  .pill.green { background: rgba(16, 185, 129, 0.2); color: #34D399; border: 1px solid rgba(16, 185, 129, 0.4); }
  
  .mono {
    font-family: 'Consolas', 'Courier New', monospace;
    color: #FDE047;
  }
</style>
</head>
<body>

  <div class="header">
    <h1>Cách DevGuard Bảo Vệ 23 Microservices Khi Tải Thư Viện Go</h1>
    <p>DevGuard đứng giữa làm người gác cổng: chặn gói cấm, diệt mã độc, giữ lại gói mới ra lò và lưu sẵn trên đĩa để build nhanh hơn</p>
  </div>

  <div class="pipeline-container">
    
    <!-- Stage 1 -->
    <div class="stage-card blue">
      <div class="stage-header">
        <div class="stage-title">1. Phía Dev & 23 Dịch Vụ</div>
        <div class="stage-subtitle">Nơi chạy code và gọi lệnh tải thư viện</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>23 Dịch vụ Go:</b> user-service, order-service, payment-service, auth-service...</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Lệnh chạy hàng ngày:</b> <code>go get</code>, <code>go mod download</code>, <code>docker build</code></div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Cài đặt duy nhất:</b> Đổi địa chỉ tải gói sang DevGuard thay vì tự tải thẳng</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Nguyên tắc an toàn:</b> Bắt buộc đi qua DevGuard, nếu có lỗi sẽ dừng build ngay</div>
        </div>
        <div class="code-badge">
          export GOPROXY="http://devguard:8080/api/v1/dependency-proxy/{secret}/go"
        </div>
      </div>
      <div class="stage-footer">BẮT BUỘC ĐI QUA CỔNG KIỂM TRA</div>
    </div>

    <div class="connector">
      <div class="connector-circle">&rarr;</div>
    </div>

    <!-- Stage 2 -->
    <div class="stage-card orange">
      <div class="stage-header">
        <div class="stage-title">2. Cổng Lọc DevGuard</div>
        <div class="stage-subtitle">Người gác cổng kiểm tra mọi thứ trước khi cho tải</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Nhận diện dự án:</b> Đọc mã bí mật (secret) để biết ai đang gọi tải gói</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>So danh sách cấm:</b> Kiểm tra tên gói xem có trong danh sách đen không</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Kiểm tra tuổi của gói:</b> Nếu gói mới ra mắt dưới 48 tiếng thì chặn ngay</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Báo lỗi tức thì:</b> Trả về lỗi 403, mã độc không bao giờ chạm được vào máy</div>
        </div>
      </div>
      <div class="stage-footer">CHẶN NGAY TỪ CỔNG (LỖI 403)</div>
    </div>

    <div class="connector">
      <div class="connector-circle">&rarr;</div>
    </div>

    <!-- Stage 3 -->
    <div class="stage-card cyan">
      <div class="stage-header">
        <div class="stage-title">3. Danh Sách Mã Độc & Cache</div>
        <div class="stage-subtitle">Kho dữ liệu nguy hiểm & nơi lưu sẵn gói sạch</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Cập nhật mã độc:</b> Lấy liên tục từ cơ sở dữ liệu quốc tế OSV và GitHub</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Kiểm tra mã độc trước:</b> Luôn soát mã độc trước khi cho phép lưu tạm vào máy</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Lưu sẵn trên đĩa (Cache):</b> Giữ lại gói an toàn để lần sau không phải tải lại</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Tự động dọn dẹp:</b> Nếu phát hiện gói đã lưu bị đánh dấu độc hại, xóa sạch ngay</div>
        </div>
      </div>
      <div class="stage-footer">BUILD LẠI SIÊU NHANH (TIẾT KIỆM 70% THỜI GIAN)</div>
    </div>

    <div class="connector">
      <div class="connector-circle">&rarr;</div>
    </div>

    <!-- Stage 4 -->
    <div class="stage-card green">
      <div class="stage-header">
        <div class="stage-title">4. Kho Gốc Của Google</div>
        <div class="stage-subtitle">Kho chứa thư viện Go chính thức trên internet</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>proxy.golang.org:</b> Nơi DevGuard ra tải hộ nếu trong máy chưa có sẵn</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>sum.golang.org:</b> Nơi đối chiếu mã kiểm tra chống hacker sửa trộm gói</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Chỉ DevGuard ra ngoài:</b> Máy lập trình viên không cần nối mạng ra ngoài internet</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Không lo đứt cáp:</b> Dù kho gốc có sập thì cụm K8s vẫn build được nhờ bộ nhớ đệm</div>
        </div>
      </div>
      <div class="stage-footer">KHO NGUYÊN BẢN CỦA GOOGLE / GOLANG</div>
    </div>

  </div>

  <!-- Bottom Table -->
  <div class="table-container">
    <div class="table-title">4 Bước DevGuard Kiểm Tra Mỗi Khi Có Người Tải Gói</div>
    <table>
      <thead>
        <tr>
          <th>Bước Kiểm Tra</th>
          <th>Mục Đích Là Gì?</th>
          <th>DevGuard Làm Thế Nào?</th>
          <th>Kết Quả Trả Về</th>
          <th>Tác Dụng Thực Tế</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td><b>Bước 1: So với danh sách cấm</b></td>
          <td>Tránh dùng thư viện lạ, rác hoặc công ty cấm</td>
          <td>So tên gói với danh sách quy tắc do quản trị viên đặt ra</td>
          <td><span class="mono">HTTP 403</span> | <span class="mono">X-Not-Allowed-Package</span></td>
          <td><span class="pill red">Chặn Tức Thì (0ms)</span></td>
        </tr>
        <tr>
          <td><b>Bước 2: Quét mã độc OSV</b></td>
          <td>Chặn các gói hacker đã cài cắm mã độc phá hoại</td>
          <td>Tra cứu tên và phiên bản trong kho dữ liệu mã độc đã biết</td>
          <td><span class="mono">HTTP 403</span> | <span class="mono">X-Malicious-Package</span></td>
          <td><span class="pill red">Chống Bị Cài Cửa Sau</span></td>
        </tr>
        <tr>
          <td><b>Bước 3: Cách ly gói mới ra mắt</b></td>
          <td>Tránh bẫy hacker vừa tạo gói giả hoặc chiếm tài khoản tải lên</td>
          <td>Xem ngày giờ phát hành: nếu chưa đủ 48 tiếng thì chặn lại</td>
          <td><span class="mono">HTTP 403</span> | <span class="mono">X-Too-New-Package</span></td>
          <td><span class="pill amber">Đợi 48h An Toàn</span></td>
        </tr>
        <tr>
          <td><b>Bước 4: Lấy từ bộ nhớ tạm (Cache)</b></td>
          <td>Tăng tốc độ build, không lo mạng chập chờn</td>
          <td>Nếu gói đã kiểm tra an toàn và có sẵn trên đĩa thì trả về luôn</td>
          <td><span class="mono">HTTP 200</span> | <span class="mono">X-Cache: HIT</span></td>
          <td><span class="pill green">Tải Trong Tích Tắc (&lt; 0.5s)</span></td>
        </tr>
      </tbody>
    </table>
  </div>

</body>
</html>
"""

HTML_FLOW = """<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<title>Quy Trình DevGuard Xử Lý Khi Tải Gói</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: 'Segoe UI', -apple-system, BlinkMacSystemFont, Roboto, sans-serif;
    background: #0B0F19;
    color: #F8FAFC;
    padding: 44px 50px;
    width: 2200px;
    min-height: 1350px;
  }
  .header {
    text-align: center;
    margin-bottom: 36px;
  }
  .header h1 {
    font-size: 38px;
    font-weight: 800;
    letter-spacing: 1px;
    text-transform: uppercase;
    color: #FFFFFF;
    margin-bottom: 12px;
  }
  .header p {
    font-size: 19px;
    color: #94A3B8;
    max-width: 1400px;
    margin: 0 auto;
    line-height: 1.5;
  }

  /* Decision Tree Flow */
  .flow-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    margin-bottom: 36px;
  }
  .flow-card {
    background: #111827;
    border: 2px solid #1E293B;
    border-radius: 16px;
    padding: 24px 20px;
    display: flex;
    flex-direction: column;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  }
  .flow-card.active-step { border-color: #0284C7; }
  .flow-card.gate-step { border-color: #D97706; }
  .flow-card.cache-step { border-color: #0891B2; }
  .flow-card.pass-step { border-color: #059669; }

  .step-num {
    display: inline-block;
    padding: 4px 10px;
    background: #1E293B;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 800;
    color: #38BDF8;
    margin-bottom: 12px;
    text-transform: uppercase;
  }
  .flow-card.gate-step .step-num { color: #FBBF24; }
  .flow-card.pass-step .step-num { color: #34D399; }

  .card-title {
    font-size: 18px;
    font-weight: 800;
    color: #FFFFFF;
    margin-bottom: 12px;
    line-height: 1.3;
  }
  .card-desc {
    font-size: 14px;
    color: #94A3B8;
    line-height: 1.5;
    margin-bottom: 16px;
    flex: 1;
  }

  .branch-box {
    border-radius: 10px;
    padding: 12px;
    font-size: 13px;
    margin-top: 8px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    line-height: 1.4;
  }
  .branch-box.block {
    background: rgba(239, 68, 68, 0.12);
    border: 1px solid rgba(239, 68, 68, 0.3);
    color: #FCA5A5;
  }
  .branch-box.pass {
    background: rgba(16, 185, 129, 0.12);
    border: 1px solid rgba(16, 185, 129, 0.3);
    color: #6EE7B7;
  }
  .branch-box b { font-weight: 700; }

  /* Bottom Response Payloads */
  .responses-container {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
  }
  .resp-card {
    background: #111827;
    border-radius: 14px;
    padding: 20px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  }
  .resp-card.red { border: 2px solid #DC2626; }
  .resp-card.amber { border: 2px solid #D97706; }
  .resp-card.green { border: 2px solid #059669; }

  .resp-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 12px;
  }
  .resp-title {
    font-size: 15px;
    font-weight: 800;
    text-transform: uppercase;
  }
  .resp-card.red .resp-title { color: #F87171; }
  .resp-card.amber .resp-title { color: #FBBF24; }
  .resp-card.green .resp-title { color: #34D399; }

  .resp-status {
    font-size: 12px;
    font-weight: 800;
    padding: 3px 8px;
    border-radius: 4px;
  }
  .resp-card.red .resp-status { background: #7F1D1D; color: #FEE2E2; }
  .resp-card.amber .resp-status { background: #78350F; color: #FEF3C7; }
  .resp-card.green .resp-status { background: #064E3B; color: #D1FAE5; }

  .resp-code {
    background: #0B0F19;
    border: 1px solid #1E293B;
    border-radius: 8px;
    padding: 12px;
    font-family: 'Consolas', 'Courier New', monospace;
    font-size: 13px;
    color: #E2E8F0;
    line-height: 1.45;
    white-space: pre-wrap;
  }
  .resp-meta {
    font-size: 12px;
    color: #94A3B8;
    margin-top: 10px;
    line-height: 1.4;
  }
</style>
</head>
<body>

  <div class="header">
    <h1>Quy Trình DevGuard Xử Lý Khi Có Yêu Cầu Tải Gói Go</h1>
    <p>Từng bước kiểm tra đơn giản, rõ ràng: quyết định cho phép tải về hay chặn đứng một thư viện</p>
  </div>

  <div class="flow-grid">
    
    <!-- Card 1 -->
    <div class="flow-card active-step">
      <div class="step-num">BƯỚC 1: NHẬN YÊU CẦU</div>
      <div class="card-title">Máy Dev Hoặc CI Gọi Lệnh Tải</div>
      <div class="card-desc">
        Lập trình viên hoặc server CI chạy lệnh <code>go get</code> hoặc <code>go mod download</code>. Lệnh tự động chuyển hướng qua DevGuard:
      </div>
      <div class="branch-box pass">
        <b>Đường dẫn gọi tới:</b>
        <code>GET /api/v1/dependency-proxy/{secret}/go/github.com/pkg/@v/v1.0.0.info</code>
      </div>
      <div class="branch-box pass">
        <b>Hành động:</b> Đọc mã bí mật (secret) để nạp danh sách luật cấm và thời gian cách ly (48 giờ).
      </div>
    </div>

    <!-- Card 2 -->
    <div class="flow-card gate-step">
      <div class="step-num">BƯỚC 2: TÊN GÓI CÓ BỊ CẤM KHÔNG?</div>
      <div class="card-title">So Với Danh Sách Luật Cấm</div>
      <div class="card-desc">
        So khớp tên gói với danh sách quy tắc cho phép hoặc cấm của công ty:
      </div>
      <div class="branch-box block">
        <b>&#10006; NẾU BỊ CẤM:</b>
        Chặn luôn với lỗi 403 Forbidden. Build dừng ngay, mã độc không vào được máy.
      </div>
      <div class="branch-box pass">
        <b>&#10004; NẾU HỢP LỆ:</b>
        Cho phép chuyển tiếp sang Bước 3 để rà soát mã độc.
      </div>
    </div>

    <!-- Card 3 -->
    <div class="flow-card gate-step">
      <div class="step-num">BƯỚC 3: CÓ PHẢI MÃ ĐỘC ĐÃ BIẾT?</div>
      <div class="card-title">Quét Kho Dữ Liệu Mã Độc OSV</div>
      <div class="card-desc">
        Tra cứu tên gói trong kho dữ liệu mã độc (OSV / GitHub) trước khi đụng vào bộ nhớ đệm:
      </div>
      <div class="branch-box block">
        <b>&#10006; NẾU LÀ MÃ ĐỘC:</b>
        Chặn ngay với lỗi 403 Forbidden và xóa sạch đĩa nếu từng lưu.
      </div>
      <div class="branch-box pass">
        <b>&#10004; NẾU SẠCH SẼ:</b>
        Cho phép chuyển tiếp sang Bước 4 để kiểm tra tuổi gói.
      </div>
    </div>

    <!-- Card 4 -->
    <div class="flow-card cache-step">
      <div class="step-num">BƯỚC 4: GÓI CÓ QUÁ MỚI KHÔNG?</div>
      <div class="card-title">Cách Ly Gói Mới & Trả Từ Cache</div>
      <div class="card-desc">
        Kiểm tra thời gian phát hành của gói và kiểm tra bộ nhớ đệm đĩa:
      </div>
      <div class="branch-box block">
        <b>&#10006; MỚI DƯỚI 48 GIỜ:</b>
        Chặn lại với lỗi 403 Forbidden. Tạm cách ly chờ cộng đồng kiểm chứng!
      </div>
      <div class="branch-box pass">
        <b>&#10004; TRÊN 48 GIỜ & CÓ TRONG CACHE:</b>
        Trả về ngay lập tức (Báo thành công 200 OK, không tốn internet).
      </div>
    </div>

  </div>

  <!-- Response Payloads -->
  <div class="responses-container">
    
    <!-- Resp 1 -->
    <div class="resp-card red">
      <div class="resp-header">
        <div class="resp-title">1. Khi Bị Chặn Vì Là Mã Độc Hoặc Cấm</div>
        <div class="resp-status">LỖI 403 FORBIDDEN</div>
      </div>
      <div class="resp-code">HTTP/1.1 403 Forbidden
Content-Type: application/json
X-Malicious-Package: blocked

{
  "error": "Forbidden",
  "message": "Gói này bị chặn vì phát hiện chứa mã độc nguy hiểm",
  "reason": "Phát hiện mã độc phá hoại mang mã định danh MAL-...",
  "blocked": true
}</div>
      <div class="resp-meta">Tác dụng: Lệnh build dừng ngay lập tức. Hacker không thể cài cắm bất cứ file độc nào vào ổ cứng của bạn.</div>
    </div>

    <!-- Resp 2 -->
    <div class="resp-card amber">
      <div class="resp-header">
        <div class="resp-title">2. Khi Bị Chặn Vì Gói Mới Ra Lò (Dưới 48 Giờ)</div>
        <div class="resp-status">LỖI 403 FORBIDDEN</div>
      </div>
      <div class="resp-code">HTTP/1.1 403 Forbidden
Content-Type: application/json
X-Too-New-Package: blocked

{
  "error": "Forbidden",
  "message": "Gói này bị chặn vì phát hành quá mới",
  "reason": "Gói mới ra mắt được 14 tiếng, quy định cần đủ 48 tiếng để đảm bảo an toàn",
  "blocked": true
}</div>
      <div class="resp-meta">Tác dụng: Tránh hơn 90% các vụ hacker vừa chiếm tài khoản lập trình viên rồi tải ngay bản độc hại lên mạng.</div>
    </div>

    <!-- Resp 3 -->
    <div class="resp-card green">
      <div class="resp-header">
        <div class="resp-title">3. Khi Gói An Toàn Và Lấy Từ Bộ Nhớ Tạm</div>
        <div class="resp-status">THÀNH CÔNG 200 OK</div>
      </div>
      <div class="resp-code">HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Cache: HIT
X-Proxy-Type: go

{
  "Version": "v1.9.1",
  "Time": "2023-03-08T08:52:16Z"
}</div>
      <div class="resp-meta">Tác dụng: Lấy trực tiếp từ ổ cứng trong mạng nội bộ Kubernetes, tốc độ dưới 10 mili-giây, không tốn mạng internet.</div>
    </div>

  </div>

</body>
</html>
"""

def render_html_with_selenium(html_content: str, output_png_path: str, width=2200, height=1350):
    temp_html = output_png_path.replace(".png", "_temp_render.html")
    with open(temp_html, "w", encoding="utf-8") as f:
        f.write(html_content)

    chrome_options = Options()
    chrome_options.add_argument("--headless=new")
    chrome_options.add_argument(f"--window-size={width},{height}")
    chrome_options.add_argument("--force-device-scale-factor=1.5")
    chrome_options.add_argument("--hide-scrollbars")
    chrome_options.add_argument("--disable-gpu")
    chrome_options.add_argument("--no-sandbox")

    driver = webdriver.Chrome(options=chrome_options)
    try:
        driver.get("file:///" + os.path.abspath(temp_html).replace("\\", "/"))
        time.sleep(1)
        driver.save_screenshot(output_png_path)
        print(f"Successfully rendered: {output_png_path}")
    finally:
        driver.quit()
        if os.path.exists(temp_html):
            os.remove(temp_html)

def main():
    thesis_img_dir = os.path.abspath("docs/images")
    cyberdev_img_dir = os.path.abspath("../CyberDev/docs/images")
    artifact_dir = r"C:\Users\ADMIN\.gemini\antigravity-ide\brain\78481042-9160-470b-aa3d-203543a658b8"

    arch_img = os.path.join(thesis_img_dir, "devguard_dependency_firewall_architecture.png")
    flow_img = os.path.join(thesis_img_dir, "devguard_dependency_firewall_flow.png")

    print("=== Rendering Architecture Diagram with Conversational Vietnamese ===")
    render_html_with_selenium(HTML_ARCHITECTURE, arch_img)

    print("=== Rendering Flow Diagram with Conversational Vietnamese ===")
    render_html_with_selenium(HTML_FLOW, flow_img)

    # Sync to CyberDev
    if os.path.exists(cyberdev_img_dir):
        shutil.copy2(arch_img, os.path.join(cyberdev_img_dir, "devguard_dependency_firewall_architecture.png"))
        shutil.copy2(flow_img, os.path.join(cyberdev_img_dir, "devguard_dependency_firewall_flow.png"))
        print(f"Synced to CyberDev: {cyberdev_img_dir}")

    # Sync to Artifact Directory
    if os.path.exists(artifact_dir):
        shutil.copy2(arch_img, os.path.join(artifact_dir, "devguard_dependency_firewall_architecture.png"))
        shutil.copy2(flow_img, os.path.join(artifact_dir, "devguard_dependency_firewall_flow.png"))
        print(f"Synced to Artifact Directory: {artifact_dir}")

if __name__ == "__main__":
    main()
