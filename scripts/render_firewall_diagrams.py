#!/usr/bin/env python3
"""
Render DevGuard Dependency Firewall architecture and flow diagrams into high-resolution PNG images.
"""
import os
import requests
import shutil

ARCH_DIAGRAM_MERMAID = """flowchart TB
    subgraph Clients ["1. KHÔNG GIAN PHÁT TRIỂN VÀ BUILD (CLIENTS)"]
        Dev["Developer Workstation<br/>(go get / go build)"]
        CI["CI/CD Pipeline Runner<br/>(GitHub Actions / GitLab CI)"]
        K8sBuild["Kubernetes Image Builder<br/>(Kaniko / Docker BuildKit)"]
        Services["23 Go Microservices<br/>(user-service, order-service, ...)"]
        
        Dev --> |"GOPROXY=http://devguard.../go"| Gateway
        CI --> |"GOPROXY=http://devguard.../go"| Gateway
        K8sBuild --> |"GOPROXY=http://devguard.../go"| Gateway
        Services -.-> |"go.mod dependencies"| K8sBuild
    end

    subgraph DevGuard ["2. DEVGUARD DEPENDENCY FIREWALL GATEWAY (CORE)"]
        Gateway["Reverse Proxy and Router<br/>(/api/v1/dependency-proxy/:secret/go/*)"]
        
        subgraph AuthPolicy ["Phân giải Danh tính và Cấu hình"]
            SecretResolver["Secret Scope Resolver<br/>(Asset / Project / Org Scope)"]
            ConfigLoader["Config Loader<br/>(Rules, MinReleaseAge Cooldown)"]
        end

        subgraph InspectionEngine ["Bộ ba Động cơ Kiểm soát (Inspection Engine)"]
            RuleEngine["Rule Pattern Matcher<br/>(matchPattern with * and ! overrides)"]
            MalChecker["Malicious Package Checker<br/>(vulndb.MaliciousPackageChecker)"]
            CooldownChecker["Cooldown / Quarantine Engine<br/>(time.Since releaseTime &lt; MinReleaseAge)"]
        end

        subgraph LocalCache ["Bộ đệm Đĩa Cục bộ (Local Disk Cache)"]
            LRUCache["Disk-backed LRU Cache<br/>(goCacheTTL: 7 days immutable)"]
        end

        Gateway --> SecretResolver
        SecretResolver --> ConfigLoader
        ConfigLoader --> RuleEngine
        RuleEngine --> MalChecker
        MalChecker --> LRUCache
        LRUCache --> CooldownChecker
    end

    subgraph DataLayer ["3. CƠ SỞ DỮ LIỆU VÀ TÌNH BÁO MÃ ĐỘC"]
        PostgresDB[("PostgreSQL<br/>(malicious_packages and components)")]
        OSVFeed["OSV / OpenSSF / GitHub Advisory<br/>(MAL-* Feed Stream)"]
        
        OSVFeed --> |"Đồng bộ định kỳ"| PostgresDB
        PostgresDB <--> |"Truy vấn PURL theo SemVer"| MalChecker
    end

    subgraph UpstreamRegistry ["4. UPSTREAM GOLANG ECOSYSTEM"]
        GoProxy["Upstream Go Proxy<br/>(https://proxy.golang.org)"]
        ChecksumDB["Go Checksum Database<br/>(https://sum.golang.org)"]
        
        CooldownChecker --> |"Fetch khi Cache Miss"| GoProxy
        GoProxy -.-> |"Xác thực toàn vẹn checksum"| ChecksumDB
    end

    classDef clientStyle fill:#e1f5fe,stroke:#0288d1,stroke-width:2px;
    classDef firewallStyle fill:#fff3e0,stroke:#f57c00,stroke-width:2px;
    classDef dataStyle fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px;
    classDef upstreamStyle fill:#e8f5e9,stroke:#388e3c,stroke-width:2px;

    class Dev,CI,K8sBuild,Services clientStyle;
    class Gateway,SecretResolver,ConfigLoader,RuleEngine,MalChecker,CooldownChecker,LRUCache firewallStyle;
    class PostgresDB,OSVFeed dataStyle;
    class GoProxy,ChecksumDB upstreamStyle;
"""

SEQ_DIAGRAM_MERMAID = """sequenceDiagram
    autonumber
    actor Client as Go Client / CI Runner (user-service)
    participant Proxy as DevGuard Router (ProxyGo)
    participant Auth as Scope and Config Resolver
    participant Rules as Rule Engine (CheckNotAllowed)
    participant MalDB as Malicious DB (OSV/OpenSSF)
    participant Cache as Disk LRU Cache
    participant Upstream as Upstream (proxy.golang.org)

    Note over Client, Upstream: Bắt đầu tải phụ thuộc (go mod download / go get)
    Client->>Proxy: GET /api/v1/dependency-proxy/{secret}/go/github.com/pkg/@v/v1.0.0.info
    
    Proxy->>Auth: Giải mã UUID secret và lấy cấu hình (Rules, MinReleaseAge)
    Auth-->>Proxy: Trả về Config (minReleaseAge: 48h, rules: [...])
    
    Proxy->>Rules: Kiểm tra package identifier theo rules
    alt Bị chặn bởi Rule (Blacklist)
        Rules-->>Proxy: Match blocked pattern
        Proxy-->>Client: HTTP 403 Forbidden (X-Not-Allowed-Package: blocked)
        Note over Client: Build dừng ngay lập tức. Mã độc không chạm đĩa!
    else Hợp lệ theo Rule
        Rules-->>Proxy: Pass
    end

    Proxy->>MalDB: Tra cứu PURL trong bảng malicious_packages
    alt Phát hiện mã độc (MAL-*)
        MalDB-->>Proxy: Match Malicious Affected Component
        Proxy->>Cache: Xóa khỏi cache đĩa nếu có (chống Cache Poisoning)
        Proxy-->>Client: HTTP 403 Forbidden (X-Malicious-Package: blocked)
        Note over Client: Cảnh báo mã độc chuỗi cung ứng được kích hoạt!
    else Không có mã độc
        MalDB-->>Proxy: Clean
    end

    Proxy->>Cache: Kiểm tra file trong Cache đĩa (Freshness duoi 7 ngay)
    alt Cache HIT
        Cache-->>Proxy: Trả về Data + Cached ReleaseTime
        opt Nếu có cấu hình MinReleaseAge
            Proxy->>Proxy: So sánh: now - ReleaseTime có nhỏ hơn MinReleaseAge?
            alt Quá mới (duoi 48h)
                Proxy-->>Client: HTTP 403 Forbidden (X-Too-New-Package: blocked)
            end
        end
        Proxy-->>Client: HTTP 200 OK (X-Cache: HIT, Content-Type: text/plain)
    else Cache MISS
        Proxy->>Upstream: GET https://proxy.golang.org/github.com/pkg/@v/v1.0.0.info
        Upstream-->>Proxy: HTTP 200 OK (Data chứa ReleaseTime)
        
        Proxy->>Proxy: ExtractGoVersionAndReleaseTime(data)
        alt time.Since(ReleaseTime) nhỏ hơn MinReleaseAge (Chưa đủ 48h)
            Proxy-->>Client: HTTP 403 Forbidden (X-Too-New-Package: blocked)
            Note over Client: Gói bị cách ly (Quarantine Cooldown) phòng ngừa Zero-Day!
        else Đủ thời gian an toàn (tren 48h)
            Proxy->>Cache: Lưu vào Disk LRU Cache (kèm ReleaseTime)
            Proxy-->>Client: HTTP 200 OK (X-Cache: MISS, Dữ liệu module)
        end
    end
"""

def render_mermaid_to_png(mermaid_text: str, output_path: str):
    url = "https://kroki.io/mermaid/png"
    print(f"Rendering to {output_path} via {url}...")
    headers = {"Content-Type": "text/plain; charset=utf-8"}
    resp = requests.post(url, data=mermaid_text.encode("utf-8"), headers=headers, timeout=30)
    if resp.status_code != 200:
        raise RuntimeError(f"Kroki returned HTTP {resp.status_code}: {resp.text}")
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, "wb") as f:
        f.write(resp.content)
    print(f"Successfully saved: {output_path} ({len(resp.content)} bytes)")

def main():
    thesis_img_dir = os.path.abspath("docs/images")
    cyberdev_img_dir = os.path.abspath("../CyberDev/docs/images")

    arch_img_thesis = os.path.join(thesis_img_dir, "devguard_dependency_firewall_architecture.png")
    seq_img_thesis = os.path.join(thesis_img_dir, "devguard_dependency_firewall_flow.png")

    render_mermaid_to_png(ARCH_DIAGRAM_MERMAID, arch_img_thesis)
    render_mermaid_to_png(SEQ_DIAGRAM_MERMAID, seq_img_thesis)

    # Sync to CyberDev if directory exists
    if os.path.exists(os.path.dirname(cyberdev_img_dir)):
        os.makedirs(cyberdev_img_dir, exist_ok=True)
        shutil.copy2(arch_img_thesis, os.path.join(cyberdev_img_dir, "devguard_dependency_firewall_architecture.png"))
        shutil.copy2(seq_img_thesis, os.path.join(cyberdev_img_dir, "devguard_dependency_firewall_flow.png"))
        print(f"Synchronized diagrams to CyberDev: {cyberdev_img_dir}")

if __name__ == "__main__":
    main()
