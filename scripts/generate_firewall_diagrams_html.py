#!/usr/bin/env python3
"""
Generate publication-quality, high-resolution architecture and flow diagrams
for DevGuard Dependency Firewall matching the exact aesthetic of airgap_zero_trust_topology.png.
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
<title>DevGuard Dependency Firewall Architecture</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: 'Segoe UI', -apple-system, BlinkMacSystemFont, Roboto, sans-serif;
    background: #0B0F19;
    color: #F8FAFC;
    padding: 48px 56px;
    width: 2200px;
    min-height: 1400px;
  }
  .header {
    text-align: center;
    margin-bottom: 48px;
  }
  .header h1 {
    font-size: 38px;
    font-weight: 800;
    letter-spacing: 1.5px;
    text-transform: uppercase;
    color: #FFFFFF;
    margin-bottom: 12px;
  }
  .header p {
    font-size: 19px;
    color: #94A3B8;
    max-width: 1400px;
    margin: 0 auto;
  }

  /* 4-Stage Horizontal Pipeline */
  .pipeline-container {
    display: flex;
    align-items: stretch;
    justify-content: space-between;
    gap: 20px;
    margin-bottom: 48px;
  }
  .stage-card {
    flex: 1;
    background: #111827;
    border-radius: 16px;
    padding: 28px 24px;
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
    width: 44px;
    flex-shrink: 0;
  }
  .connector-circle {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    font-weight: bold;
    border: 2px solid;
  }
  .connector-circle.arrow {
    background: rgba(16, 185, 129, 0.15);
    border-color: #10B981;
    color: #10B981;
  }
  .connector-circle.deny {
    background: rgba(239, 68, 68, 0.15);
    border-color: #EF4444;
    color: #EF4444;
  }

  .stage-header {
    margin-bottom: 18px;
  }
  .stage-title {
    font-size: 20px;
    font-weight: 800;
    letter-spacing: 0.5px;
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
    gap: 14px;
    margin-bottom: 24px;
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
    padding: 12px 16px;
    text-align: center;
    font-size: 14px;
    font-weight: 700;
    letter-spacing: 0.5px;
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
    padding: 28px 32px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  }
  .table-title {
    font-size: 19px;
    font-weight: 800;
    color: #38BDF8;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    gap: 12px;
    text-transform: uppercase;
    letter-spacing: 0.8px;
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
    padding: 14px 18px;
    text-transform: uppercase;
    font-size: 13px;
    letter-spacing: 0.8px;
    border-bottom: 2px solid #334155;
  }
  td {
    padding: 16px 18px;
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
    letter-spacing: 0.5px;
  }
  .pill.red { background: rgba(239, 68, 68, 0.2); color: #F87171; border: 1px solid rgba(239, 68, 68, 0.4); }
  .pill.amber { background: rgba(245, 158, 11, 0.2); color: #FBBF24; border: 1px solid rgba(245, 158, 11, 0.4); }
  .pill.green { background: rgba(16, 185, 129, 0.2); color: #34D399; border: 1px solid rgba(16, 185, 129, 0.4); }
  .pill.blue { background: rgba(14, 165, 233, 0.2); color: #38BDF8; border: 1px solid rgba(14, 165, 233, 0.4); }
  
  .mono {
    font-family: 'Consolas', 'Courier New', monospace;
    color: #FDE047;
  }
</style>
</head>
<body>

  <div class="header">
    <h1>Kiến Trúc Tường Lửa Phụ Thuộc DevGuard (Dependency Firewall Architecture)</h1>
    <p>Mô hình In-line Pre-ingestion Proxy kiểm soát toàn diện chuỗi cung ứng mã nguồn cho 23 Go Microservices trên Kubernetes</p>
  </div>

  <div class="pipeline-container">
    
    <!-- Stage 1 -->
    <div class="stage-card blue">
      <div class="stage-header">
        <div class="stage-title">1. Môi Trường Build & Clients</div>
        <div class="stage-subtitle">Các điểm phát sinh nhu cầu tải phụ thuộc Go</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>23 Go Microservices:</b> user-service, order-service, payment-service, auth-service, ...</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Developer Workstation:</b> Lệnh <code>go get</code> / <code>go mod download</code></div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>CI/CD Runner:</b> GitHub Actions / GitLab CI Runner nội bộ</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Kubernetes Builder:</b> Docker BuildKit / Kaniko multi-stage</div>
        </div>
        <div class="code-badge">
          ENV GOPROXY="http://devguard:8080/api/v1/dependency-proxy/{secret}/go"
        </div>
      </div>
      <div class="stage-footer">FAIL-CLOSED ENFORCEMENT (CHẶN 100% DIRECT)</div>
    </div>

    <div class="connector">
      <div class="connector-circle arrow">&rarr;</div>
    </div>

    <!-- Stage 2 -->
    <div class="stage-card orange">
      <div class="stage-header">
        <div class="stage-title">2. DevGuard Firewall Gateway</div>
        <div class="stage-subtitle">Cổng kiểm soát an ninh in-line độc quyền (Core)</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Reverse Proxy Router:</b> Phân tích giao thức GOPROXY (<code>/@v/list</code>, <code>.info</code>, <code>.mod</code>, <code>.zip</code>)</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Secret Scope Resolver:</b> Xác thực token phân cấp Asset / Project / Organization</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Rule Pattern Matcher:</b> Whitelist / Blacklist PURL (hỗ trợ wildcard <code>*</code> và phủ định <code>!</code>)</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Cooldown Quarantine Engine:</b> Ngăn chặn gói mới phát hành &lt; 48 giờ (Phòng vệ Zero-Day)</div>
        </div>
      </div>
      <div class="stage-footer">HTTP 403 FORBIDDEN (CHẶN PRE-INGESTION)</div>
    </div>

    <div class="connector">
      <div class="connector-circle arrow">&rarr;</div>
    </div>

    <!-- Stage 3 -->
    <div class="stage-card cyan">
      <div class="stage-header">
        <div class="stage-title">3. Dữ Liệu & Bộ Đệm Cục Bộ</div>
        <div class="stage-subtitle">Kho tình báo mã độc & đệm đĩa tăng tốc K8s</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>PostgreSQL VulnDB:</b> Bảng <code>malicious_packages</code> & components đồng bộ OSV / OpenSSF</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Malicious Checker:</b> Tra cứu mã độc <code>MAL-*</code> TRƯỚC KHI truy xuất bộ nhớ đệm cache</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Disk-backed LRU Cache:</b> Giới hạn 1GB, thời gian sống 7 ngày cho file bất biến (immutable)</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Cache Poisoning Defense:</b> Xóa đệm đĩa ngay lập tức khi phát hiện mã độc mới</div>
        </div>
      </div>
      <div class="stage-footer">X-CACHE: HIT (TĂNG TỐC BUILD 75%)</div>
    </div>

    <div class="connector">
      <div class="connector-circle arrow">&rarr;</div>
    </div>

    <!-- Stage 4 -->
    <div class="stage-card green">
      <div class="stage-header">
        <div class="stage-title">4. Upstream Go Ecosystem</div>
        <div class="stage-subtitle">Hệ sinh thái Golang chính thức toàn cầu</div>
      </div>
      <div class="stage-content">
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Official Go Proxy:</b> <code>https://proxy.golang.org</code> cung cấp nguồn module sạch</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Checksum Database:</b> <code>https://sum.golang.org</code> bảo đảm tính toàn vẹn cryptographic</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Egress Security Transport:</b> DevGuard lọc và cô lập kết nối ra bên ngoài qua proxy</div>
        </div>
        <div class="feature-item">
          <div class="feature-dot"></div>
          <div><b>Left-Pad Resiliency:</b> Ổn định build ngay cả khi registry công cộng bị sập hoặc xóa gói</div>
        </div>
      </div>
      <div class="stage-footer">HTTP 200 OK (CHỨNG THỰC NGUYÊN BẢN)</div>
    </div>

  </div>

  <!-- Bottom Table -->
  <div class="table-container">
    <div class="table-title">Ma Trận 4 Chốt Chặn An Ninh Thực Nghiệm & Cơ Chế Phòng Thủ (Security Gates)</div>
    <table>
      <thead>
        <tr>
          <th>Chốt Chặn An Ninh</th>
          <th>Mục Tiêu Bảo Vệ</th>
          <th>Cơ Chế Kỹ Thuật (DevGuard Core)</th>
          <th>Mã Trả Về & Header HTTP</th>
          <th>Trạng Thái Bảo Vệ</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td><b>Chốt Chặn 1: Rule Engine</b></td>
          <td>Ngăn chặn các package không được phê duyệt hoặc vi phạm chính sách tổ chức</td>
          <td>Hàm <span class="mono">CheckNotAllowedPackage</span> so khớp PURL với danh sách quy tắc (Whitelist / Blacklist)</td>
          <td><span class="mono">HTTP 403</span> | <span class="mono">X-Not-Allowed-Package: blocked</span></td>
          <td><span class="pill red">Chặn Tức Thì (0ms đĩa)</span></td>
        </tr>
        <tr>
          <td><b>Chốt Chặn 2: Malicious DB</b></td>
          <td>Chặn đứng các gói độc hại đã bị gắn mã định danh khai thác OSV / OpenSSF</td>
          <td>Hàm <span class="mono">checkMaliciousPackage</span> truy vấn bảng <span class="mono">malicious_packages</span> trước khi đọc Cache</td>
          <td><span class="mono">HTTP 403</span> | <span class="mono">X-Malicious-Package: blocked</span></td>
          <td><span class="pill red">Chống Cache Poisoning</span></td>
        </tr>
        <tr>
          <td><b>Chốt Chặn 3: Cooldown Quarantine</b></td>
          <td>Cách ly các gói mới phát hành để phòng ngừa tấn công Account Hijacking & Zero-Day</td>
          <td>Bóc tách <span class="mono">Time</span> trong metadata <span class="mono">.info</span>; tính <span class="mono">time.Since(Time) &lt; MinReleaseAge</span></td>
          <td><span class="mono">HTTP 403</span> | <span class="mono">X-Too-New-Package: blocked</span></td>
          <td><span class="pill amber">Cách Ly An Toàn (48h)</span></td>
        </tr>
        <tr>
          <td><b>Chốt Chặn 4: Disk LRU Cache</b></td>
          <td>Bảo đảm tính sẵn sàng cao, miễn nhiễm sự cố sập registry và tăng tốc độ build K8s</td>
          <td>Lưu trữ các file bất biến (<span class="mono">.info</span>, <span class="mono">.mod</span>, <span class="mono">.zip</span>) trong đĩa cục bộ với TTL 7 ngày</td>
          <td><span class="mono">HTTP 200</span> | <span class="mono">X-Cache: HIT</span></td>
          <td><span class="pill green">Tăng Tốc 75% (0.8s)</span></td>
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
<title>DevGuard Dependency Firewall Request Flow</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: 'Segoe UI', -apple-system, BlinkMacSystemFont, Roboto, sans-serif;
    background: #0B0F19;
    color: #F8FAFC;
    padding: 48px 56px;
    width: 2200px;
    min-height: 1400px;
  }
  .header {
    text-align: center;
    margin-bottom: 40px;
  }
  .header h1 {
    font-size: 38px;
    font-weight: 800;
    letter-spacing: 1.5px;
    text-transform: uppercase;
    color: #FFFFFF;
    margin-bottom: 12px;
  }
  .header p {
    font-size: 19px;
    color: #94A3B8;
    max-width: 1400px;
    margin: 0 auto;
  }

  /* Decision Tree Flow */
  .flow-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    margin-bottom: 40px;
  }
  .flow-card {
    background: #111827;
    border: 2px solid #1E293B;
    border-radius: 16px;
    padding: 24px 20px;
    display: flex;
    flex-direction: column;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
    position: relative;
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
    gap: 24px;
  }
  .resp-card {
    background: #111827;
    border-radius: 14px;
    padding: 22px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  }
  .resp-card.red { border: 2px solid #DC2626; }
  .resp-card.amber { border: 2px solid #D97706; }
  .resp-card.green { border: 2px solid #059669; }

  .resp-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 14px;
  }
  .resp-title {
    font-size: 16px;
    font-weight: 800;
    text-transform: uppercase;
  }
  .resp-card.red .resp-title { color: #F87171; }
  .resp-card.amber .resp-title { color: #FBBF24; }
  .resp-card.green .resp-title { color: #34D399; }

  .resp-status {
    font-size: 13px;
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
    padding: 14px;
    font-family: 'Consolas', 'Courier New', monospace;
    font-size: 13px;
    color: #E2E8F0;
    line-height: 1.5;
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
    <h1>Luồng Điều Phối & Cây Quyết Định Cách Ly Phụ Thuộc (Decision Tree & Request Flow)</h1>
    <p>Quy trình tuần tự xử lý yêu cầu tải module Go và cơ chế kích hoạt 3 chốt chặn an ninh tại cổng mạng</p>
  </div>

  <div class="flow-grid">
    
    <!-- Card 1 -->
    <div class="flow-card active-step">
      <div class="step-num">BƯỚC 1: TIẾP NHẬN YÊU CẦU</div>
      <div class="card-title">Go Client Khởi Tạo Yêu Cầu</div>
      <div class="card-desc">
        Microservice thực thi <code>go mod download</code> hoặc <code>go get</code>. Yêu cầu tải metadata bản phát hành được chuyển hướng đến cổng DevGuard:
      </div>
      <div class="branch-box pass">
        <b>Endpoint Yêu Cầu:</b>
        <code>GET /api/v1/dependency-proxy/{secret}/go/github.com/pkg/@v/v1.0.0.info</code>
      </div>
      <div class="branch-box pass">
        <b>Hành động:</b> Giải mã Secret Scope & nạp cấu hình (Rules, MinReleaseAge = 48h).
      </div>
    </div>

    <!-- Card 2 -->
    <div class="flow-card gate-step">
      <div class="step-num">BƯỚC 2: CHỐT CHẶN CHÍNH SÁCH</div>
      <div class="card-title">Kiểm Tra Quy Tắc (Rule Engine)</div>
      <div class="card-desc">
        So khớp PURL (<code>pkg:go/...</code>) theo thứ tự định nghĩa giống .gitignore (quy tắc cuối cùng quyết định, hỗ trợ <code>!</code> cho whitelist):
      </div>
      <div class="branch-box block">
        <b>&#10006; VI PHẠM BLACKLIST:</b>
        Trả về ngay <b>HTTP 403 Forbidden</b>. Build dừng ngay lập tức, mã độc không chạm đĩa!
      </div>
      <div class="branch-box pass">
        <b>&#10004; HỢP LỆ THEO RULE:</b>
        Cho phép đi tiếp đến Chốt chặn Kiểm tra Mã độc (Malicious DB).
      </div>
    </div>

    <!-- Card 3 -->
    <div class="flow-card gate-step">
      <div class="step-num">BƯỚC 3: ĐỐI SOÁT MÃ ĐỘC</div>
      <div class="card-title">Truy Vấn OSV Malicious Feed</div>
      <div class="card-desc">
        Truy vấn bảng <code>malicious_packages</code> TRƯỚC KHI đọc Cache để triệt tiêu hoàn toàn nguy cơ Cache Poisoning:
      </div>
      <div class="branch-box block">
        <b>&#10006; KHỚP MÃ ĐỘC MAL-*:</b>
        Xóa ngay khỏi cache đĩa nếu có, trả về <b>HTTP 403 (X-Malicious-Package)</b>.
      </div>
      <div class="branch-box pass">
        <b>&#10004; SẠCH SẼ (CLEAN):</b>
        Cho phép tiến hành kiểm tra bộ đệm đĩa cục bộ (Disk Cache).
      </div>
    </div>

    <!-- Card 4 -->
    <div class="flow-card cache-step">
      <div class="step-num">BƯỚC 4: BỘ ĐỆM & QUARANTINE</div>
      <div class="card-title">Đánh Giá Cooldown & Caching</div>
      <div class="card-desc">
        Kiểm tra trạng thái đĩa cache hoặc tải từ upstream <code>proxy.golang.org</code> kèm đánh giá thời gian phát hành:
      </div>
      <div class="branch-box block">
        <b>&#10006; QUÁ MỚI (&lt; 48 GIỜ):</b>
        Trả về <b>HTTP 403 (X-Too-New-Package)</b>. Kích hoạt cơ chế cách ly Zero-Day!
      </div>
      <div class="branch-box pass">
        <b>&#10004; AN TOÀN (&ge; 48 GIỜ):</b>
        Lưu Cache đĩa & trả về <b>HTTP 200 OK (X-Cache: HIT/MISS)</b>.
      </div>
    </div>

  </div>

  <!-- Response Payloads -->
  <div class="responses-container">
    
    <!-- Resp 1 -->
    <div class="resp-card red">
      <div class="resp-header">
        <div class="resp-title">1. Phản Ứng Khi Chặn Rule / Mã Độc</div>
        <div class="resp-status">HTTP 403 FORBIDDEN</div>
      </div>
      <div class="resp-code">HTTP/1.1 403 Forbidden
Content-Type: application/json
X-Malicious-Package: blocked

{
  "error": "Forbidden",
  "message": "This package has been blocked by the malicious package firewall",
  "reason": "Package github.com/fake-org/malicious-package is flagged as malicious (ID: MAL-FAKE-TEST-GO-GITHUB-COM-FAKE-ORG-MALICIOUS-PACKAGE)",
  "blocked": true
}</div>
      <div class="resp-meta">Hiệu ứng: Chặn đứng tức thì tại bước <code>go mod download</code>. Runner thoát với exit code 1, bảo vệ máy chủ an toàn tuyệt đối.</div>
    </div>

    <!-- Resp 2 -->
    <div class="resp-card amber">
      <div class="resp-header">
        <div class="resp-title">2. Phản Ứng Cách Ly Cooldown (Zero-Day)</div>
        <div class="resp-status">HTTP 403 FORBIDDEN</div>
      </div>
      <div class="resp-code">HTTP/1.1 403 Forbidden
Content-Type: application/json
X-Too-New-Package: blocked

{
  "error": "Forbidden",
  "message": "This package has been blocked because it was released too recently",
  "reason": "Package github.com/gin-gonic/gin was released 14h20m ago, which is less than the required minimum of 48 hours",
  "blocked": true
}</div>
      <div class="resp-meta">Hiệu ứng: Ngăn chặn 90%+ cuộc tấn công Account Hijacking xảy ra trong 48 giờ đầu khi cộng đồng chưa kịp cảnh báo CVE.</div>
    </div>

    <!-- Resp 3 -->
    <div class="resp-card green">
      <div class="resp-header">
        <div class="resp-title">3. Phản Ứng Khi Gói An Toàn & Cache HIT</div>
        <div class="resp-status">HTTP 200 OK</div>
      </div>
      <div class="resp-code">HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Cache: HIT
X-Proxy-Type: go

{
  "Version": "v1.9.1",
  "Time": "2023-03-08T08:52:16Z"
}</div>
      <div class="resp-meta">Hiệu ứng: Phục vụ trực tiếp từ bộ đệm đĩa nội bộ trong cụm K8s, thời gian đáp ứng dưới 10ms, giảm tải 100% băng thông internet ngoài.</div>
    </div>

  </div>

</body>
</html>
"""

def render_html_with_selenium(html_content: str, output_png_path: str, width=2200, height=1400):
    temp_html = output_png_path.replace(".png", "_temp.html")
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
        time.sleep(1) # wait for fonts & rendering
        # Get actual height
        required_height = driver.execute_script("return Math.max(document.body.scrollHeight, document.body.offsetHeight, document.documentElement.clientHeight);")
        required_width = driver.execute_script("return Math.max(document.body.scrollWidth, document.body.offsetWidth, document.documentElement.clientWidth);")
        driver.set_window_size(required_width, required_height)
        time.sleep(0.5)
        driver.save_screenshot(output_png_path)
        print(f"Successfully rendered high-quality image: {output_png_path} ({required_width}x{required_height})")
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

    print("=== Rendering Architecture Diagram ===")
    render_html_with_selenium(HTML_ARCHITECTURE, arch_img, width=2200, height=1350)

    print("=== Rendering Flow & Decision Diagram ===")
    render_html_with_selenium(HTML_FLOW, flow_img, width=2200, height=1350)

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
