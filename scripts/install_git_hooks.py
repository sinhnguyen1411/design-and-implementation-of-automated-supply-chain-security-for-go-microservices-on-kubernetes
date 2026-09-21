#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DevGuard Pre-commit Hook Installer
Tool: install_git_hooks.py
Mục đích: Cài đặt Git hook pre-commit tự động trên cả 2 repository (CyberDev và Luận văn K8s)
          để kiểm tra văn phong học thuật tiếng Việt trước mỗi lần commit tài liệu.
"""

import os
import sys
from pathlib import Path

# Đảm bảo in tiếng Việt trên console Windows
if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
        sys.stderr.reconfigure(encoding="utf-8")
    except Exception:
        pass

PRE_COMMIT_SCRIPT_CONTENT = """#!/bin/sh
# DevGuard Git Pre-commit Hook for Vietnamese Technical Documentation

echo "================================================================================"
echo " [DevGuard Quality Gate] Kiểm tra Văn phong Kỹ thuật & Thuật ngữ Tiếng Việt..."
echo "================================================================================"

# Lấy danh sách các tệp Markdown đang được stage trong thư mục docs/
STAGED_DOCS=$(git diff --cached --name-only --diff-filter=ACM | grep -E '^docs/.*\\.md$')

if [ -z "$STAGED_DOCS" ]; then
    echo " [+] Không có tài liệu markdown nào trong docs/ thay đổi. Bỏ qua kiểm tra."
    exit 0
fi

FAILED=0
for file in $STAGED_DOCS; do
    if [ -f "$file" ]; then
        echo " -> Đang kiểm tra: $file"
        python scripts/lint_vietnamese_report.py "$file"
        if [ $? -ne 0 ]; then
            FAILED=1
        fi
    fi
done

if [ $FAILED -ne 0 ]; then
    echo ""
    echo " [LỖI NGHIÊM TRỌNG] Commit bị từ chối do vi phạm quy chuẩn văn phong tiếng Việt!"
    echo " Vui lòng xem dòng vi phạm chi tiết ở trên và sửa lại, hoặc chạy lệnh sau để tự động sửa:"
    echo "   python scripts/lint_vietnamese_report.py docs/CyberDev_Experimental_Feasibility_Report.md --fix"
    echo "================================================================================"
    exit 1
fi

echo " [+] TẤT CẢ TÀI LIỆU TIẾNG VIỆT ĐỀU ĐẠT CHUẨN HỌC THUẬT DEVGUARD (0 LỖI)!"
echo "================================================================================"
exit 0
"""

REPOS = [
    Path(r"C:\Users\ADMIN\Documents\CyberDev"),
    Path(r"C:\Users\ADMIN\Documents\design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes")
]

def install_hooks():
    print("=" * 80)
    print(" CÀI ĐẶT DEVGUARD GIT PRE-COMMIT HOOKS CHO CẢ HAI REPOSITORY")
    print("=" * 80)

    for repo in REPOS:
        if not repo.exists():
            print(f"[-] Không tìm thấy thư mục repository: {repo}")
            continue

        git_dir = repo / ".git"
        if not git_dir.exists():
            print(f"[-] Thư mục không phải git repo: {repo}")
            continue

        hooks_dir = git_dir / "hooks"
        hooks_dir.mkdir(parents=True, exist_ok=True)
        hook_path = hooks_dir / "pre-commit"

        # Ghi tệp hook với Unix line-endings (\n) để Git bash thực thi chuẩn xác
        hook_path.write_bytes(PRE_COMMIT_SCRIPT_CONTENT.encode("utf-8"))
        
        # Tạo thư mục .githooks version-controlled
        githooks_dir = repo / ".githooks"
        githooks_dir.mkdir(parents=True, exist_ok=True)
        (githooks_dir / "pre-commit").write_bytes(PRE_COMMIT_SCRIPT_CONTENT.encode("utf-8"))

        print(f"[+] Đã cài đặt thành công pre-commit hook tại:")
        print(f"    - {hook_path}")
        print(f"    - {githooks_dir / 'pre-commit'}")

    print("\n" + "=" * 80)
    print(" HOÀN TẤT CÀI ĐẶT HOOKS! TỪ NAY MỌI COMMIT SẼ ĐƯỢC TỰ ĐỘNG BẢO VỆ.")
    print("=" * 80)

if __name__ == "__main__":
    install_hooks()
