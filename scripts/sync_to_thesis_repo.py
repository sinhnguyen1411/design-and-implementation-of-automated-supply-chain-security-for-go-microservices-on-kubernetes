import os
import sys
import shutil
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
THESIS_ROOT = Path(r"C:\Users\ADMIN\Documents\design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes")

CYBERDEV_DOCS = CYBERDEV_ROOT / "docs"
THESIS_DOCS = THESIS_ROOT / "docs"

def sync_artifacts():
    print("=" * 82)
    print(" ĐỒNG BỘ HÓA TÀI SẢN KỸ THUẬT SANG REPOSITORY LUẬN VĂN CHÍNH")
    print(f" Nguồn : {CYBERDEV_DOCS}")
    print(f" Đích  : {THESIS_DOCS}")
    print("=" * 82)

    THESIS_DOCS.mkdir(parents=True, exist_ok=True)
    (THESIS_DOCS / "images").mkdir(parents=True, exist_ok=True)

    # 1. Synchronize All Images in docs/images
    src_images = CYBERDEV_DOCS / "images"
    dst_images = THESIS_DOCS / "images"
    print("\n[1] ĐỒNG BỘ DANH MỤC HÌNH ẢNH (DOCS/IMAGES):")
    img_count = 0
    for img_file in src_images.glob("*.*"):
        dst_file = dst_images / img_file.name
        shutil.copy2(img_file, dst_file)
        img_count += 1
    print(f"    [+] Đã sao chép thành công {img_count} hình ảnh sang {dst_images}")

    # 2. Synchronize Telemetry JSON and Data Files
    print("\n[2] ĐỒNG BỘ DỮ LIỆU TELEMETRY JSON & METADATA AN NINH:")
    json_files = [
        "airgap_verification_telemetry.json",
        "sast_airgap_verification_telemetry.json",
        "secrets_airgap_verification_telemetry.json",
        "iac_airgap_verification_telemetry.json",
        "container_airgap_verification_telemetry.json",
        "dast_airgap_verification_telemetry.json",
        "supply_chain_airgap_verification_telemetry.json",
        "policy_gate_airgap_verification_telemetry.json",
        "dependency_firewall_phase3_comparative.json",
        "dependency_firewall_23_services_audit.json",
        "all_pillars_executive_summary.json",
        "slsa_provenance_user_service.json",
        "slsa_provenance_tampered.json",
        "slsa_provenance_user_service.sig",
        "user_service_release.sig",
        "cyclonedx_sbom_user_service.json",
        "airgap_user_service_sbom.json",
        "vex_rules.json"
    ]
    json_count = 0
    for jf in json_files:
        src = CYBERDEV_DOCS / jf
        if src.exists():
            dst = THESIS_DOCS / jf
            shutil.copy2(src, dst)
            print(f"    [+] {jf:<45} -> Đồng bộ thành công ({src.stat().st_size:,} bytes)")
            json_count += 1
        else:
            print(f"    [-] {jf:<45} -> Không tìm thấy ở nguồn")

    # 3. Synchronize Subdirectories (sarif_reports, cosign_keys)
    print("\n[3] ĐỒNG BỘ THƯ MỤC CHỨNG TỰ & BÁO CÁO SARIF:")
    for sub in ["sarif_reports", "cosign_keys"]:
        s_src = CYBERDEV_DOCS / sub
        s_dst = THESIS_DOCS / sub
        if s_src.exists():
            s_dst.mkdir(parents=True, exist_ok=True)
            for f in s_src.glob("*.*"):
                shutil.copy2(f, s_dst / f.name)
            print(f"    [+] Đã đồng bộ thư mục {sub} ({len(list(s_src.glob('*.*')))} tệp)")

    # 4. Synchronize Main Report Files (MD, DOCX, PDF)
    print("\n[4] ĐỒNG BỘ BÁO CÁO TOÀN DIỆN (MD, DOCX, PDF):")
    # DOCX (pick newest between primary and Fixed fallback)
    docx_fixed = CYBERDEV_DOCS / "CyberDev_Experimental_Feasibility_Report_Fixed.docx"
    docx_primary = CYBERDEV_DOCS / "CyberDev_Experimental_Feasibility_Report.docx"
    if docx_fixed.exists() and (not docx_primary.exists() or docx_fixed.stat().st_mtime > docx_primary.stat().st_mtime):
        docx_src = docx_fixed
    else:
        docx_src = docx_primary
    docx_dst = THESIS_DOCS / "CyberDev_Experimental_Feasibility_Report.docx"
    if docx_src.exists():
        shutil.copy2(docx_src, docx_dst)
        print(f"    [+] DOCX Report ({docx_src.name}) -> {docx_dst} ({docx_dst.stat().st_size / (1024*1024):.2f} MB)")

    # PDF
    pdf_src = CYBERDEV_DOCS / "CyberDev_Experimental_Feasibility_Report.pdf"
    pdf_dst = THESIS_DOCS / "CyberDev_Experimental_Feasibility_Report.pdf"
    if pdf_src.exists():
        shutil.copy2(pdf_src, pdf_dst)
        print(f"    [+] PDF Report  -> {pdf_dst} ({pdf_dst.stat().st_size / (1024*1024):.2f} MB)")

    # Markdown Report: Copy and update absolute path links to point locally to thesis repo
    md_src = CYBERDEV_DOCS / "CyberDev_Experimental_Feasibility_Report.md"
    md_dst = THESIS_DOCS / "CyberDev_Experimental_Feasibility_Report.md"
    if md_src.exists():
        content = md_src.read_text(encoding="utf-8")
        # Replace CyberDev path with thesis repo path for seamless viewing
        old_prefix = "file:///c:/Users/ADMIN/Documents/CyberDev/docs/images/"
        new_prefix = "file:///c:/Users/ADMIN/Documents/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/docs/images/"
        updated_content = content.replace(old_prefix, new_prefix)
        md_dst.write_text(updated_content, encoding="utf-8")
        print(f"    [+] MD Report   -> {md_dst} ({len(updated_content):,} characters)")

    # 4.1 Synchronize Vietnamese Style Guide, Integration Plan & Tooling
    tool_files = [
        ("docs", "VIETNAMESE_TECHNICAL_STYLE_GUIDE.md"),
        ("docs", "VIETNAMESE_TOOLING_LONG_TERM_INTEGRATION_PLAN.md"),
        ("scripts", "lint_vietnamese_report.py"),
        ("scripts", "build_and_sync_all_reports.py"),
        ("scripts", "install_git_hooks.py")
    ]
    print("\n[4.1] ĐỒNG BỘ CÔNG CỤ VĂN PHONG & QUY CHUẨN KỸ THUẬT TIẾNG VIỆT:")
    for folder, fname in tool_files:
        t_src = CYBERDEV_ROOT / folder / fname
        t_dst = THESIS_ROOT / folder / fname
        if t_src.exists():
            t_dst.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(t_src, t_dst)
            print(f"    [+] {fname:<48} -> Đồng bộ thành công")

    # 5. Verify Target Directory
    print("\n[5] KIỂM TRA TÍNH TOÀN VẸN TẠI THƯ MỤC ĐÍCH:")
    files_to_check = [
        md_dst,
        docx_dst,
        pdf_dst,
        THESIS_DOCS / "airgap_verification_telemetry.json",
        THESIS_DOCS / "policy_gate_airgap_verification_telemetry.json",
        THESIS_DOCS / "dependency_firewall_phase3_comparative.json"
    ]
    all_ok = True
    for f in files_to_check:
        if f.exists() and f.stat().st_size > 0:
            print(f"    [OK] {f.name:<45} | Size: {f.stat().st_size:,} bytes")
        else:
            print(f"    [FAIL] {f.name:<45} | Missing or 0 bytes")
            all_ok = False

    print("\n" + "=" * 82)
    if all_ok:
        print(" HOÀN TẤT ĐỒNG BỘ HÓA TOÀN DIỆN SANG REPOSITORY LUẬN VĂN THÀNH CÔNG RỰC RỠ!")
    else:
        print(" CÓ LỖI XẢY RA TRONG QUÁ TRÌNH ĐỒNG BỘ.")
    print("=" * 82)

if __name__ == "__main__":
    sync_artifacts()
