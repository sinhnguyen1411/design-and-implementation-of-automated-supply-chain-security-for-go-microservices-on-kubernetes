#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DevGuard Dependency Firewall - 23 Go Microservices CI Gatekeeper
Scans and evaluates all 23 microservices' go.mod dependencies against DevGuard In-Line Firewall.
"""

import os
import sys
import json
import time
import re
import requests

if sys.stdout.encoding != 'utf-8':
    sys.stdout.reconfigure(encoding='utf-8', errors='replace')

SERVICES_DIR = "services"
PROXY_BASE_URL = os.environ.get("DEVGUARD_PROXY_URL", "http://localhost:8080/api/v1/dependency-proxy")
ASSET_SECRET = os.environ.get("DEVGUARD_ASSET_SECRET", "4a882f42-bdb4-49a0-b5fe-0509eb92faa0")
PROXY_ENDPOINT = f"{PROXY_BASE_URL}/{ASSET_SECRET}/go"

def parse_go_mod(file_path):
    """Parse direct dependencies from go.mod"""
    dependencies = []
    if not os.path.exists(file_path):
        return dependencies

    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Match require (...) blocks or single require lines
    in_require_block = False
    for line in content.splitlines():
        line = line.strip()
        if line.startswith("require ("):
            in_require_block = True
            continue
        elif in_require_block and line == ")":
            in_require_block = False
            continue

        if in_require_block or line.startswith("require "):
            if line.startswith("require "):
                line = line[len("require "):].strip()
            
            # Skip comments or indirect dependencies
            if line.startswith("//") or "// indirect" in line:
                continue
            
            parts = line.split()
            if len(parts) >= 2:
                mod_path = parts[0]
                version = parts[1]
                # Filter out internal/local modules
                if not mod_path.startswith("github.com/sinhnguyen1411/stock-trading-be"):
                    dependencies.append({"path": mod_path, "version": version})

    return dependencies

def check_dependency(mod_path, version):
    """Query DevGuard GOPROXY for a specific module version"""
    url = f"{PROXY_ENDPOINT}/{mod_path}/@v/{version}.info"
    start = time.perf_counter()
    try:
        r = requests.get(url, timeout=5)
        elapsed_ms = (time.perf_counter() - start) * 1000
        cache_status = r.headers.get("X-Cache", "MISS")
        is_blocked = (r.status_code == 403)
        trigger = (
            "X-Not-Allowed-Package" if "x-not-allowed-package" in [k.lower() for k in r.headers] else
            "X-Malicious-Package" if "x-malicious-package" in [k.lower() for k in r.headers] else
            "X-Too-New-Package" if "x-too-new-package" in [k.lower() for k in r.headers] else
            "CLEAN"
        )
        return {
            "module": mod_path,
            "version": version,
            "status_code": r.status_code,
            "latency_ms": round(elapsed_ms, 2),
            "cache_status": cache_status,
            "is_blocked": is_blocked,
            "trigger": trigger,
            "passed": not is_blocked
        }
    except Exception as e:
        elapsed_ms = (time.perf_counter() - start) * 1000
        return {
            "module": mod_path,
            "version": version,
            "status_code": 0,
            "latency_ms": round(elapsed_ms, 2),
            "cache_status": "ERROR",
            "is_blocked": False,
            "trigger": str(e),
            "passed": False
        }

def run_gatekeeper():
    print("=" * 90)
    print("  KIỂM TRA AN NINH THƯ VIỆN CỦA 23 GO MICROSERVICES")
    print("=" * 90)

    services = [d for d in os.listdir(SERVICES_DIR) if os.path.isdir(os.path.join(SERVICES_DIR, d))]
    services.sort()

    total_services = len(services)
    print(f"Bắt đầu quét {total_services} microservices trong thư mục '{SERVICES_DIR}/'...\n")

    matrix_report = []
    overall_total_deps = 0
    overall_passed_deps = 0
    overall_blocked_deps = 0
    service_pass_count = 0

    # Header bảng kết quả
    print(f"{'STT':<4} | {'Tên Microservice':<26} | {'Tổng':<5} | {'Hợp lệ':<7} | {'Bị chặn':<8} | {'Thời gian':<12} | {'Kết quả'}")
    print("-" * 90)

    for idx, svc_name in enumerate(services, 1):
        mod_file = os.path.join(SERVICES_DIR, svc_name, "go.mod")
        deps = parse_go_mod(mod_file)
        
        svc_results = []
        for dep in deps:
            res = check_dependency(dep["path"], dep["version"])
            svc_results.append(res)

        total_deps = len(svc_results)
        passed_deps = sum(1 for r in svc_results if r["passed"])
        blocked_deps = sum(1 for r in svc_results if r["is_blocked"])
        avg_lat = sum(r["latency_ms"] for r in svc_results) / total_deps if total_deps > 0 else 0.0

        overall_total_deps += total_deps
        overall_passed_deps += passed_deps
        overall_blocked_deps += blocked_deps

        status_str = "ĐẠT [100%]" if blocked_deps == 0 else f"LỖI [{blocked_deps} gói cấm]"
        if blocked_deps == 0:
            service_pass_count += 1

        print(f"{idx:<4} | {svc_name:<26} | {total_deps:<5} | {passed_deps:<7} | {blocked_deps:<8} | {avg_lat:>8.2f} ms | {status_str}")

        matrix_report.append({
            "service": svc_name,
            "go_mod_path": mod_file,
            "total_dependencies": total_deps,
            "passed": passed_deps,
            "blocked": blocked_deps,
            "avg_latency_ms": round(avg_lat, 2),
            "compliance_status": "COMPLIANT" if blocked_deps == 0 else "NON_COMPLIANT",
            "dependencies": svc_results
        })

    print("-" * 90)
    print(f"Tổng số microservice kiểm tra    : {total_services} / {total_services}")
    print(f"Số service đạt chuẩn (không lỗi) : {service_pass_count} / {total_services} ({(service_pass_count/total_services)*100:.1f}%)")
    print(f"Tổng số lượt kiểm tra thư viện   : {overall_total_deps}")
    print(f"Số gói an toàn được thông qua    : {overall_passed_deps}")
    print(f"Số gói vi phạm bị chặn           : {overall_blocked_deps}")
    print("=" * 90)

    # Lưu kết quả ra file JSON
    out_file = "docs/dependency_firewall_23_services_audit.json"
    with open(out_file, "w", encoding="utf-8") as f:
        json.dump({
            "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
            "total_services": total_services,
            "compliant_services": service_pass_count,
            "total_dependencies_audited": overall_total_deps,
            "total_passed": overall_passed_deps,
            "total_blocked": overall_blocked_deps,
            "services_matrix": matrix_report
        }, f, indent=2)
    print(f"[ĐÃ LƯU] File kết quả: {out_file}")

if __name__ == "__main__":
    run_gatekeeper()
