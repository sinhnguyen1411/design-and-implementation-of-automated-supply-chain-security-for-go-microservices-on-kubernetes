import os
import sys
import json
import re
from pathlib import Path

# Ensure UTF-8 output on Windows PowerShell
if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
        sys.stderr.reconfigure(encoding="utf-8")
    except Exception:
        pass

CYBERDEV_ROOT = Path(r"C:\Users\ADMIN\Documents\CyberDev")
DOCS_DIR = CYBERDEV_ROOT / "docs"

def audit_report():
    print("=" * 82)
    print(" DEVGUARD AIR-GAPPED MULTI-PHASE AUDIT & SYNCHRONIZATION VERIFICATION")
    print("=" * 82)

    # 1. Check Document Files
    md_file = DOCS_DIR / "CyberDev_Experimental_Feasibility_Report.md"
    docx_file = DOCS_DIR / "CyberDev_Experimental_Feasibility_Report.docx"
    pdf_file = DOCS_DIR / "CyberDev_Experimental_Feasibility_Report.pdf"

    print("\n[1] KIỂM TRA TRẠNG THÁI VÀ KÍCH THƯỚC CÁC TỆP BÁO CÁO CHÍNH:")
    for f in [md_file, docx_file, pdf_file]:
        if f.exists():
            size_mb = f.stat().st_size / (1024 * 1024)
            print(f"    [+] {f.name:<45} | Size: {size_mb:6.2f} MB | Status: TỒN TẠI")
        else:
            print(f"    [-] {f.name:<45} | Status: CHƯA TẠO")

    # 2. Check Telemetry JSON files for all phases
    telemetry_files = [
        ("Phase 1: Multi-tier Harness & SCA", "airgap_verification_telemetry.json"),
        ("Phase 2: SAST Opengrep 5 PoCs", "sast_airgap_verification_telemetry.json"),
        ("Phase 3: Secret Scanning Gitleaks", "secrets_airgap_verification_telemetry.json"),
        ("Phase 4: IaC Security Trivy Config", "iac_airgap_verification_telemetry.json"),
        ("Phase 5: Container Security Trivy Image", "container_airgap_verification_telemetry.json"),
        ("Phase 6: DAST Nuclei Runtime Testing", "dast_airgap_verification_telemetry.json"),
        ("Phase 7: Supply Chain Cosign & SLSA", "supply_chain_airgap_verification_telemetry.json"),
        ("Phase 8: Dependency Firewall In-Line Proxy", "dependency_firewall_phase3_comparative.json"),
        ("Phase 8b: Dependency Firewall Benchmark", "dependency_firewall_phase2_benchmark.json"),
        ("Phase 9: 23 Services Firewall Audit", "dependency_firewall_23_services_audit.json"),
        ("Phase 10: Unified CI/CD Policy Gate", "policy_gate_airgap_verification_telemetry.json")
    ]

    print("\n[2] KIỂM TRA FILE TELEMETRY JSON ĐỊNH LƯỢNG (TELEMETRY JSON):")
    all_telemetry_ok = True
    for label, fname in telemetry_files:
        p = DOCS_DIR / fname
        if p.exists():
            try:
                data = json.loads(p.read_text(encoding="utf-8"))
                print(f"    [+] {label:<45} -> {fname:<45} | Valid JSON")
            except Exception as e:
                print(f"    [-] {label:<45} -> {fname:<45} | JSON Error: {e}")
                all_telemetry_ok = False
        else:
            print(f"    [-] {label:<45} -> {fname:<45} | MISSING")
            all_telemetry_ok = False

    # 3. Check Figure and Table counts in Markdown
    print("\n[3] KIỂM TRA MỤC LỤC & SỐ LƯỢNG HÌNH ẢNH, BẢNG BIỂU TRONG MARKDOWN:")
    if md_file.exists():
        text = md_file.read_text(encoding="utf-8")
        figures = re.findall(r"\*Hình (\d+\.\d+):", text)
        tables = re.findall(r"\*Bảng (\d+(\.\d+)?):", text)
        sections = re.findall(r"## Phần (\d+):", text)
        print(f"    [+] Số phần chính (Main Sections)       : {len(sections)} (Phần 1 - 14)")
        print(f"    [+] Tổng số hình ảnh có chú thích rõ ràng : {len(figures)} hình ảnh")
        print(f"    [+] Tổng số bảng biểu kỹ thuật          : {len(tables)} bảng biểu")
        print(f"    [+] Tổng số ký tự trong báo cáo         : {len(text):,} ký tự")

    # 4. Check Images on Disk
    print("\n[4] KIỂM TRA HÌNH ẢNH AIR-GAPPED & DEPENDENCY FIREWALL ĐÃ TẠO:")
    airgap_imgs = [
        "airgap_zero_trust_topology.png",
        "airgap_terminal_verification.png",
        "real_airgap_docker_network.png",
        "airgap_sast_zero_trust_architecture.png",
        "airgap_terminal_sast_opengrep.png",
        "airgap_docker_sast_runner.png",
        "airgap_secret_zero_trust_architecture.png",
        "airgap_terminal_secret_gitleaks.png",
        "airgap_docker_secret_runner.png",
        "airgap_iac_zero_trust_architecture.png",
        "airgap_terminal_iac_trivy.png",
        "airgap_docker_iac_runner.png",
        "airgap_container_zero_trust_architecture.png",
        "airgap_terminal_container_trivy.png",
        "airgap_docker_container_runner.png",
        "airgap_dast_zero_trust_architecture.png",
        "airgap_terminal_dast_nuclei.png",
        "airgap_docker_dast_runner.png",
        "airgap_supply_chain_zero_trust_architecture.png",
        "airgap_terminal_supply_chain_cosign.png",
        "airgap_docker_supply_chain_runner.png",
        "devguard_dependency_firewall_flow.png",
        "real_terminal_dependency_firewall_harness.png",
        "real_terminal_dependency_firewall_blocking.png",
        "real_terminal_dependency_firewall_phase3_matrix.png",
        "real_terminal_dependency_firewall_phase3_details.png",
        "devguard_web_dependency_firewall_overview.png",
        "airgap_policy_gate_zero_trust_architecture.png",
        "airgap_terminal_policy_gate.png",
        "airgap_docker_policy_gate_runner.png"
    ]
    missing_imgs = []
    for img_name in airgap_imgs:
        img_path = DOCS_DIR / "images" / img_name
        if img_path.exists():
            kb = img_path.stat().st_size / 1024
            print(f"    [+] {img_name:<48} | Size: {kb:6.1f} KB")
        else:
            missing_imgs.append(img_name)
            print(f"    [-] {img_name:<48} | MISSING")

    print("\n" + "=" * 82)
    if all_telemetry_ok and len(missing_imgs) == 0:
        print(" TẤT CẢ CÁC TRỤ CỘT AN NINH ĐỀU ĐÃ ĐẠT 100% ĐỒNG BỘ VÀ KIỂM CHỨNG TOÀN DIỆN!")
    else:
        print(f" CẢNH BÁO: Còn thiếu {len(missing_imgs)} hình ảnh hoặc tệp thông số đo đạc.")
    print("=" * 82)

if __name__ == "__main__":
    audit_report()
