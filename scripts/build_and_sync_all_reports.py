#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DevGuard All-in-One Documentation & Feasibility Report Pipeline
Tool: build_and_sync_all_reports.py
Mục đích: Pipeline một chạm tự động hóa toàn bộ chu trình:
  [1] Kiểm định văn phong tiếng Việt (Linter)
  [2] Xuất bản tài liệu Word (.docx) chuyên dụng
  [3] Render tài liệu PDF vector qua Edge Headless (.pdf)
  [4] Đồng bộ toàn bộ tài sản sang repository Luận văn K8s
  [5] Kiểm toán đối soát chéo và xác thực toàn vẹn dữ liệu
"""

import sys
import time
import subprocess
from pathlib import Path

# Đảm bảo in tiếng Việt trên console Windows
if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
        sys.stderr.reconfigure(encoding="utf-8")
    except Exception:
        pass

SCRIPTS_DIR = Path(r"C:\Users\ADMIN\Documents\CyberDev\scripts")
DOCS_DIR = Path(r"C:\Users\ADMIN\Documents\CyberDev\docs")

def run_step(step_num: int, title: str, cmd: list, cwd: Path = SCRIPTS_DIR.parent):
    print("\n" + "=" * 80)
    print(f" BƯỚC [{step_num}/5]: {title.upper()}")
    print(f" Lệnh thực thi: {' '.join(cmd)}")
    print("=" * 80)
    
    start_time = time.time()
    res = subprocess.run(cmd, cwd=cwd)
    duration = time.time() - start_time
    
    if res.returncode != 0:
        print(f"\n[-] THẤT BẠI TẠI BƯỚC [{step_num}]: {title} (Thoát mã: {res.returncode})")
        print("    Pipeline đã dừng lại để tránh phát sinh dữ liệu không đồng nhất.")
        return False, duration
    
    print(f"\n[+] HOÀN THÀNH BƯỚC [{step_num}]: {title} trong {duration:.2f} giây.")
    return True, duration

def main():
    total_start = time.time()
    print("*" * 80)
    print(" DEVGUARD AUTOMATED REPORT BUILD & MULTI-REPO SYNC PIPELINE")
    print(f" Khởi động tại: {time.strftime('%Y-%m-%d %H:%M:%S')}")
    print("*" * 80)

    steps = [
        (1, "Kiểm định Văn phong Kỹ thuật & Thuật ngữ Tiếng Việt", [sys.executable, str(SCRIPTS_DIR / "lint_vietnamese_report.py"), str(DOCS_DIR / "CyberDev_Experimental_Feasibility_Report.md")]),
        (2, "Xuất bản Tài liệu Báo cáo Microsoft Word (.docx)", [sys.executable, str(SCRIPTS_DIR / "export_docx_report.py")]),
        (3, "Render Tài liệu PDF Vector Hoàn Chỉnh (.pdf) qua Edge Headless", [sys.executable, str(SCRIPTS_DIR / "export_feasibility_report.py")]),
        (4, "Đồng bộ Toàn diện Sang Repository Luận văn K8s", [sys.executable, str(SCRIPTS_DIR / "sync_to_thesis_repo.py")]),
        (5, "Kiểm toán & Xác thực Tính Toàn Vẹn Đa Định Dạng", [sys.executable, str(SCRIPTS_DIR / "verify_multi_format_sync.py")]),
    ]

    timings = {}
    for num, title, cmd in steps:
        ok, dur = run_step(num, title, cmd)
        timings[title] = dur
        if not ok:
            sys.exit(1)

    total_dur = time.time() - total_start
    print("\n" + "*" * 80)
    print(" BẢNG TỔNG KẾT THỜI GIAN THỰC THI PIPELINE:")
    for title, dur in timings.items():
        print(f"  - {title:<60}: {dur:6.2f}s")
    print(f"\n  >> TỔNG THỜI GIAN PIPELINE: {total_dur:.2f} GIÂY")
    print(" TẤT CẢ 5 BƯỚC ĐÃ THÀNH CÔNG RỰC RỠ VÀ HOÀN TOÀN ĐỒNG BỘ!")
    print("*" * 80 + "\n")

if __name__ == "__main__":
    main()
