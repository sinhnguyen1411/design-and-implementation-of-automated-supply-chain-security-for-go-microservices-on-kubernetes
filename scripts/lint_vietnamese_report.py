#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DevGuard Vietnamese Academic & Technical Prose Linter
Tool: lint_vietnamese_report.py
Mục đích: Tự động phát hiện lỗi chính tả, typography (khoảng trắng trước dấu câu),
          và các từ khóa dịch máy thô cứng (Blacklist) trong tài liệu báo cáo Markdown.
"""

import re
import sys
import argparse
from pathlib import Path

# Đảm bảo in tiếng Việt trên console Windows không bị lỗi UnicodeEncodeError
if sys.platform == "win32":
    sys.stdout.reconfigure(encoding="utf-8")
    sys.stderr.reconfigure(encoding="utf-8")

# Bảng từ khóa cấm / dịch máy thô cứng (Blacklist) -> Gợi ý thay thế
BLACKLIST_RULES = [
    (r"(?i)\bTường lửa Thư viện\b", "Dependency Firewall"),
    (r"(?i)\bmáy dev\b", "máy trạm phát triển (Developer Workstation)"),
    (r"(?i)\bmáy trạm dev\b", "máy trạm phát triển"),
    (r"(?i)\blọt vào máy trạm\b", "vượt qua chốt chặn an ninh / lọt vào mã nguồn"),
    (r"(?i)\blọt vào máy dev\b", "lọt vào mã nguồn"),
    (r"(?i)\bchạy mất mạng\b", "hoạt động trong môi trường mạng cô lập (Air-Gapped)"),
    (r"(?i)\bkhi mất mạng\b", "khi ngắt kết nối Internet / trong môi trường mạng cô lập"),
    (r"(?i)\bMất mạng hoàn toàn\b", "Ngắt kết nối Internet hoàn toàn (Air-Gapped)"),
    (r"(?i)\bgói mới ra\b", "gói thư viện mới phát hành"),
    (r"(?i)\bchặn cách ly \(mới ra\)\b", "cách ly theo chính sách Cooldown Period"),
    (r"(?i)\blàm độc cache\b", "đầu độc bộ nhớ đệm (Cache Poisoning)"),
    (r"(?i)\bủy quyền nội tuyến\b", "In-line Proxy"),
    (r"(?i)\btiến hành việc\b", "thực hiện / triển khai"),
    (r"(?i)\bcó khả năng của việc\b", "có khả năng"),
    (r"(?i)\bđược thực thi bởi\b", "hệ thống thực thi"),
]

# Quy tắc lỗi Typography (khoảng trắng trước dấu câu ngoài code block)
TYPO_RULES = [
    (r"(\w)\s+([,;:?!])", r"\1\2", "Khoảng trắng thừa trước dấu câu"),
    (r"\(\s+([^\s])", r"(\1", "Khoảng trắng thừa sau dấu mở ngoặc đơn"),
    (r"([^\s])\s+\)", r"\1)", "Khoảng trắng thừa trước dấu đóng ngoặc đơn"),
]

def lint_file(file_path: Path, auto_fix: bool = False):
    if not file_path.exists():
        print(f"[ERROR] Không tìm thấy tệp tin: {file_path}")
        return 1

    content = file_path.read_text(encoding="utf-8")
    lines = content.splitlines()
    
    in_code_block = False
    violations = []
    fixed_lines = []

    print(f"\n{'='*80}")
    print(f" DEVGUARD VIETNAMESE PROSE & TERMINOLOGY LINTER")
    print(f" Tệp kiểm tra : {file_path}")
    print(f" Chế độ       : {'TỰ ĐỘNG SỬA (--fix)' if auto_fix else 'CHỈ KIỂM TRA (Audit)'}")
    print(f"{'='*80}\n")

    for idx, line in enumerate(lines, start=1):
        # Theo dõi khối code ```
        if line.strip().startswith("```"):
            in_code_block = not in_code_block
            fixed_lines.append(line)
            continue

        if in_code_block:
            fixed_lines.append(line)
            continue

        current_line = line

        # 1. Kiểm tra Blacklist Thuật ngữ
        for pattern, replacement in BLACKLIST_RULES:
            matches = list(re.finditer(pattern, current_line))
            for m in matches:
                violations.append({
                    "line": idx,
                    "type": "BLACKLIST",
                    "matched": m.group(0),
                    "suggestion": replacement,
                    "context": current_line.strip()
                })
                if auto_fix:
                    current_line = re.sub(pattern, replacement, current_line)

        # 2. Kiểm tra Typography (loại trừ inline code `...`)
        # Tách tạm inline code để không can thiệp vào mã lệnh
        parts = current_line.split("`")
        for i in range(0, len(parts), 2):  # Chỉ xét phần văn bản thuần (chỉ số chẵn)
            text_part = parts[i]
            for pattern, rep, desc in TYPO_RULES:
                matches = list(re.finditer(pattern, text_part))
                for m in matches:
                    violations.append({
                        "line": idx,
                        "type": "TYPOGRAPHY",
                        "matched": m.group(0),
                        "suggestion": desc,
                        "context": current_line.strip()
                    })
                if auto_fix:
                    parts[i] = re.sub(pattern, rep, text_part)

        if auto_fix:
            current_line = "`".join(parts)

        fixed_lines.append(current_line)

    # Hiển thị kết quả vi phạm
    blacklist_count = sum(1 for v in violations if v["type"] == "BLACKLIST")
    typo_count = sum(1 for v in violations if v["type"] == "TYPOGRAPHY")

    if violations:
        print(f"[!] PHÁT HIỆN {len(violations)} ĐIỂM CẦN CẢI THIỆN ({blacklist_count} lỗi thuật ngữ Blacklist, {typo_count} lỗi typography):\n")
        # Nhóm theo dòng hiển thị tối đa 25 lỗi tiêu biểu
        for v in violations[:35]:
            tag = "[BLACKLIST]" if v["type"] == "BLACKLIST" else "[TYPOGRAPHY]"
            print(f"  Line {v['line']:>4} | {tag:<12} : '{v['matched']}' ➔ Gợi ý: {v['suggestion']}")
            print(f"        Ngữ cảnh: \"{v['context'][:90]}...\"\n")

        if len(violations) > 35:
            print(f"  ... và còn {len(violations) - 35} vị trí khác tương tự.\n")
    else:
        print("[+] TUYỆT VỜI! Không phát hiện lỗi thuật ngữ Blacklist hoặc lỗi typography nào!\n")

    if auto_fix and violations:
        file_path.write_text("\n".join(fixed_lines) + "\n", encoding="utf-8")
        print(f"[+] ĐÃ TỰ ĐỘNG SỬA VÀ LƯU LẠI VÀO: {file_path}")

    print(f"{'='*80}")
    print(f" TỔNG KẾT: Blacklist: {blacklist_count} | Typography: {typo_count} | Tổng: {len(violations)}")
    print(f"{'='*80}\n")

    return 0 if len(violations) == 0 else 1

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Linter văn phong tiếng Việt cho báo cáo CyberDev")
    parser.add_argument("file", nargs="?", default="docs/CyberDev_Experimental_Feasibility_Report.md", help="Đường dẫn file markdown")
    parser.add_argument("--fix", action="store_true", help="Tự động sửa các lỗi typography và thay thế thuật ngữ Blacklist chuẩn")
    args = parser.parse_args()

    sys.exit(lint_file(Path(args.file), auto_fix=args.fix))
