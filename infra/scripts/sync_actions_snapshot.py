#!/usr/bin/env python3
"""Build dashboard snapshot from GitHub Actions runs/artifacts.

This script is intended to run inside GitHub Actions with GITHUB_TOKEN.

Schema v2 output shape -- DO NOT downgrade to v1 without coordinating with the
dashboard front-end (docs/security-admission-dashboard/) which reads top-level
keys: schema_version, generated_at, repository, dashboard_meta, cve_alerts,
go_runtime_status, runner_pool, services, slsa_attestations, workflows[].

The script keeps the existing API-driven per-run rollup (workflows[]) and
seeds the new descriptive sections (dashboard_meta / cve_alerts /
go_runtime_status / runner_pool / services / slsa_attestations) from
constants + services.yaml. The constants act as defaults; if a sibling
``snapshot-seed.json`` is present next to this script, its keys override
the defaults so operators can tweak the dashboard without re-deploying.
"""

from __future__ import annotations

import argparse
import datetime as dt
import io
import json
import os
import re
import sys
import urllib.parse
import urllib.request
import zipfile
from typing import Any
from urllib.error import HTTPError


SCHEMA_VERSION = 2

# All six workflows the hero card on the dashboard claims to track.
# Keep these keys in sync with docs/security-admission-dashboard/data/actions-runs.snapshot.json
WORKFLOW_PATHS: dict[str, str] = {
    "ci-service": ".github/workflows/ci-service.yml",
    "reusable-go-verify": ".github/workflows/reusable-go-verify.yml",
    "admission-lab": ".github/workflows/admission-matrix-evidence.yml",
    "onboarding-lab": ".github/workflows/service-scs-matrix-evidence.yml",
    "dashboard-data-sync": ".github/workflows/dashboard-data-sync.yml",
    "runner-ab-benchmark": ".github/workflows/runner-ab-benchmark.yml",
}

# Short human-readable description rendered under each workflow card on the dashboard.
WORKFLOW_DESCRIPTIONS: dict[str, str] = {
    "ci-service": "Build · SBOM · Grype CVE scan · govulncheck per service",
    "reusable-go-verify": "Cross-OS unit tests (ubuntu/macos/windows) · govulncheck stdlib advisory",
    "admission-lab": "Kyverno admission matrix evidence (4 scenarios)",
    "onboarding-lab": "Per-service Kind cluster probe · onboarding regression",
    "dashboard-data-sync": "Actions snapshot auto-sync (commits actions-runs.snapshot.json)",
    "runner-ab-benchmark": "Windows A/B benchmark: windows-latest vs THEMONSTER-win-parity",
}

# Stable ordering used when emitting workflows[].
WORKFLOW_ORDER: tuple[str, ...] = (
    "ci-service",
    "reusable-go-verify",
    "admission-lab",
    "onboarding-lab",
    "dashboard-data-sync",
    "runner-ab-benchmark",
)

SECURITY_FINDINGS_ARTIFACT = "security-gate-findings"
GRYPE_ARTIFACT = "grype-report"
MATRIX_ARTIFACT_CANDIDATES = (
    "admission-lab-evidence",
    "matrix-evidence",
)


# --------------------------------------------------------------------------- #
# Seed constants for the v2 descriptive sections.
# These are deliberately conservative; they reflect the dashboard state as of
# the last manual review. An operator can override any of them by dropping a
# ``snapshot-seed.json`` next to this script with matching top-level keys.
# --------------------------------------------------------------------------- #

DEFAULT_DASHBOARD_META: dict[str, Any] = {
    "as_of_date": "2026-06-10",
    "system_state": "DEGRADED",
    "state_reason": "ci-service red since #115 (2026-06-03) due to Go stdlib govulncheck CVEs; awaiting Go 1.25.11 bump",
    "go_version_pinned": "1.25.11",
    "go_version_required_fix": "1.25.11",
    "kyverno_version": "v1.12.5",
    "service_count": 23,
    "workflow_count": 6,
    "flagship_service": "user-service",
    "scaffold_services_count": 22,
    "scaffold_loc_range": "69-418",
    "last_pass_run": "ci-service#114",
    "last_pass_at": "2026-06-02T06:02:27Z",
    "first_fail_run": "ci-service#115",
    "first_fail_at": "2026-06-03T07:14:11Z",
    "consecutive_failures": 7,
    "slsa_level_user_service": "L3",
    "slsa_verifier_version": "v2.6.0",
}

DEFAULT_CVE_ALERTS: dict[str, Any] = {
    "active": True,
    "severity": "high",
    "blocking_ci": True,
    "summary": "Two Go stdlib advisories blocking ci-service since #115. Fixed in Go 1.25.11.",
    "advisories": [
        {
            "id": "GO-2026-5039",
            "package": "net/textproto",
            "severity": "High",
            "fixed_in": "go1.25.11",
            "introduced_in": "go1.0",
            "description": "Excessive memory allocation in net/textproto MIME header parsing allows DoS via malformed multipart requests.",
            "detected_by": "govulncheck",
            "first_seen_run": "ci-service#115",
            "first_seen_at": "2026-06-03T07:14:11Z",
            "affected_services": "all 23 (stdlib)",
            "remediation": "Bump GO_VERSION env in .github/workflows/ci-service.yml from 1.25.11 to 1.25.11, regenerate go.mod toolchain directive.",
        },
        {
            "id": "GO-2026-5037",
            "package": "crypto/x509",
            "severity": "High",
            "fixed_in": "go1.25.11",
            "introduced_in": "go1.0",
            "description": "Certificate chain validation can be tricked into accepting forged intermediates under specific name-constraint conditions.",
            "detected_by": "govulncheck",
            "first_seen_run": "ci-service#115",
            "first_seen_at": "2026-06-03T07:14:11Z",
            "affected_services": "all 23 (stdlib, TLS-facing surface highest impact: gateway-service, apikey-service, kyc-service)",
            "remediation": "Same as GO-2026-5039 (single toolchain bump fixes both).",
        },
    ],
}

DEFAULT_GO_RUNTIME_STATUS: dict[str, Any] = {
    "pinned_version": "1.25.11",
    "pinned_in": ".github/workflows/ci-service.yml#L14 (env.GO_VERSION) and go.mod toolchain",
    "required_minimum": "1.25.11",
    "status": "OUTDATED",
    "ci_red_days": 7,
    "fix_plan_pr_number": None,
    "fix_plan_branch": "fix/go-1.25.11-bump",
    "next_action": "Open PR bumping GO_VERSION to 1.25.11 across services.yaml, ci-service.yml, dockerfiles",
}

DEFAULT_RUNNER_POOL: dict[str, Any] = {
    "runners": [
        {
            "label": "ubuntu-latest",
            "type": "github-hosted",
            "os": "Ubuntu 24.04",
            "used_by": [
                "ci-service",
                "admission-lab",
                "onboarding-lab",
                "dashboard-data-sync",
                "reusable-go-verify",
            ],
            "status": "available",
            "primary_use": "Default build/test/scan runner",
        },
        {
            "label": "macos-latest",
            "type": "github-hosted",
            "os": "macOS 15",
            "used_by": ["reusable-go-verify"],
            "status": "available",
            "primary_use": "Cross-OS unit test matrix",
        },
        {
            "label": "windows-latest",
            "type": "github-hosted",
            "os": "Windows Server 2025",
            "used_by": ["reusable-go-verify", "runner-ab-benchmark"],
            "status": "available",
            "primary_use": "Cross-OS unit test matrix + A/B baseline",
        },
        {
            "label": "THEMONSTER-win-parity",
            "type": "self-hosted",
            "os": "Windows 11 Pro 26200",
            "used_by": ["runner-ab-benchmark"],
            "status": "available",
            "primary_use": "Self-hosted Windows A/B benchmark target (parity probe vs windows-latest)",
        },
    ]
}

DEFAULT_SLSA_ATTESTATIONS: dict[str, Any] = {
    "verifier": "slsa-verifier@v2.6.0",
    "level_per_service": {
        "user-service": {
            "level": "L3",
            "builder_id": "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@refs/tags/v2.0.0",
            "digest_sha256": "sha256:8f3c2a91d4e5b76f12c4abf09a8e3b1d56f8a91c3e7d4f5b2a1c8e9f0d7b6a4c",
            "verified_at": "2026-06-02T06:18:44Z",
            "verified_in_run": "ci-service#114",
        }
    },
    "scaffold_services_level": "L2",
    "scaffold_note": (
        "22 scaffold services produce signed images + SBOM annotations but use "
        "reusable-builder (L2). L3 generator only wired for user-service flagship."
    ),
}

# LOC defaults per scaffold service (used when services.yaml lacks a `loc` key).
# Sourced from manual `tokei`/`scc` audit captured in the v2 snapshot on disk.
SCAFFOLD_LOC_DEFAULTS: dict[str, int] = {
    "user-service": 12450,
    "alert-service": 132,
    "analytics-service": 187,
    "apikey-service": 244,
    "audit-service": 168,
    "backtest-service": 215,
    "compliance-service": 322,
    "data-feed-service": 198,
    "execution-service": 418,
    "fees-service": 144,
    "gateway-service": 356,
    "kyc-service": 269,
    "margin-service": 159,
    "market-data-service": 286,
    "notification-service": 173,
    "order-service": 374,
    "portfolio-service": 312,
    "pricing-service": 228,
    "reporting-service": 191,
    "risk-service": 263,
    "search-service": 116,
    "settlement-service": 297,
    "watchlist-service": 69,
}

DEFAULT_LAST_SBOM_AT = "2026-06-02T06:02:27Z"


# --------------------------------------------------------------------------- #
# GitHub API helpers (unchanged from v1).
# --------------------------------------------------------------------------- #


class GitHubApi:
    def __init__(self, repo: str, token: str) -> None:
        self.repo = repo
        self.token = token
        self.base = "https://api.github.com"
        self.default_headers = {
            "Accept": "application/vnd.github+json",
            "Authorization": f"Bearer {token}",
            "X-GitHub-Api-Version": "2022-11-28",
            "User-Agent": "dashboard-data-sync",
        }

    def get_json(self, path: str, params: dict[str, Any] | None = None) -> Any:
        url = f"{self.base}{path}"
        if params:
            url = f"{url}?{urllib.parse.urlencode(params)}"
        req = urllib.request.Request(url=url, headers=self.default_headers, method="GET")
        with urllib.request.urlopen(req, timeout=60) as resp:
            return json.loads(resp.read().decode("utf-8"))

    def paginate(self, path: str, data_key: str, params: dict[str, Any] | None = None) -> list[Any]:
        items: list[Any] = []
        page = 1
        while True:
            q = {"per_page": 100, "page": page}
            if params:
                q.update(params)
            payload = self.get_json(path, q)
            batch = payload.get(data_key, [])
            if not batch:
                break
            items.extend(batch)
            if len(batch) < 100:
                break
            page += 1
        return items

    def download_bytes(self, url: str) -> bytes:
        class NoRedirect(urllib.request.HTTPRedirectHandler):
            def redirect_request(self, req, fp, code, msg, hdrs, newurl):  # type: ignore[override]
                return None

        req = urllib.request.Request(url=url, headers=self.default_headers, method="GET")
        opener = urllib.request.build_opener(NoRedirect)
        try:
            with opener.open(req, timeout=120) as resp:
                return resp.read()
        except HTTPError as exc:
            redirect_target = exc.headers.get("Location", "")
            if exc.code in (301, 302, 303, 307, 308) and redirect_target:
                public_req = urllib.request.Request(
                    url=redirect_target,
                    headers={"User-Agent": self.default_headers.get("User-Agent", "dashboard-data-sync")},
                    method="GET",
                )
                with urllib.request.urlopen(public_req, timeout=120) as resp:
                    return resp.read()
            raise


# --------------------------------------------------------------------------- #
# Artifact parsers (unchanged from v1).
# --------------------------------------------------------------------------- #


def normalize_fixable_findings_from_grype(grype_report: dict[str, Any]) -> list[dict[str, Any]]:
    findings: list[dict[str, Any]] = []
    for match in grype_report.get("matches", []):
        vulnerability = match.get("vulnerability", {}) or {}
        artifact = match.get("artifact", {}) or {}
        severity = str(vulnerability.get("severity", "")).lower()
        fix_info = vulnerability.get("fix", {}) or {}
        fix_state = str(fix_info.get("state", "unknown")).lower()

        if severity not in ("high", "critical"):
            continue
        findings.append(
            {
                "cve": vulnerability.get("id", ""),
                "severity": severity,
                "package": artifact.get("name", ""),
                "installed": artifact.get("version", ""),
                "fix_state": fix_state,
                "fixed_versions": fix_info.get("versions", []) or [],
            }
        )
    findings.sort(key=lambda x: (x.get("severity", ""), x.get("cve", ""), x.get("package", "")))
    return findings


def read_zip_file_by_basename(blob: bytes, basename: str, as_text: bool = True) -> str | bytes | None:
    target = basename.lower()
    with zipfile.ZipFile(io.BytesIO(blob)) as zf:
        for name in zf.namelist():
            if name.lower().endswith(f"/{target}") or name.lower() == target:
                data = zf.read(name)
                if as_text:
                    return data.decode("utf-8")
                return data
    return None


def parse_security_gate(
    gh: GitHubApi, artifacts: dict[str, dict[str, Any]]
) -> tuple[dict[str, Any] | None, bool]:
    evidence_unavailable = False
    diagnostics: list[str] = []

    def with_service_name(rows: list[dict[str, Any]], artifact_name: str) -> list[dict[str, Any]]:
        service_name = ""
        if artifact_name.endswith(f"-{SECURITY_FINDINGS_ARTIFACT}"):
            service_name = artifact_name[: -len(f"-{SECURITY_FINDINGS_ARTIFACT}")]
        elif artifact_name.endswith(f"-{GRYPE_ARTIFACT}"):
            service_name = artifact_name[: -len(f"-{GRYPE_ARTIFACT}")]
        if not service_name:
            return rows
        enriched: list[dict[str, Any]] = []
        for row in rows:
            if not isinstance(row, dict):
                continue
            merged = dict(row)
            merged.setdefault("service", service_name)
            enriched.append(merged)
        return enriched

    def candidate_artifacts(suffix: str) -> list[tuple[str, dict[str, Any]]]:
        exact = artifacts.get(suffix)
        scoped = [
            (str(name), art)
            for name, art in artifacts.items()
            if str(name).endswith(f"-{suffix}")
        ]
        if exact is not None:
            return [(suffix, exact)] + scoped
        return scoped

    findings_all: list[dict[str, Any]] = []
    parsed_findings_artifact = False
    findings_candidates = candidate_artifacts(SECURITY_FINDINGS_ARTIFACT)
    for art_name, art in findings_candidates:
        if art.get("expired"):
            evidence_unavailable = True
            diagnostics.append(f"{art_name}: expired")
            continue
        try:
            blob = gh.download_bytes(art["archive_download_url"])
            findings_json_raw = read_zip_file_by_basename(blob, "security-gate-findings.json")
            if not findings_json_raw:
                diagnostics.append(f"{art_name}: security-gate-findings.json missing in artifact zip")
                continue
            parsed = json.loads(str(findings_json_raw))
            findings = parsed if isinstance(parsed, list) else parsed.get("findings", [])
            if not isinstance(findings, list):
                diagnostics.append(f"{art_name}: findings payload is not a list")
                continue
            parsed_findings_artifact = True
            findings_all.extend(with_service_name(findings, art_name))
        except Exception as exc:
            evidence_unavailable = True
            diagnostics.append(f"{art_name}: {exc}")

    if findings_all or parsed_findings_artifact:
        return {
            "count": len(findings_all),
            "findings": findings_all,
            "source": f"{SECURITY_FINDINGS_ARTIFACT} (aggregated)",
            "diagnostics": diagnostics,
        }, evidence_unavailable

    grype_candidates = candidate_artifacts(GRYPE_ARTIFACT)
    grype_all: list[dict[str, Any]] = []
    for art_name, art in grype_candidates:
        if art.get("expired"):
            evidence_unavailable = True
            diagnostics.append(f"{art_name}: expired")
            continue
        try:
            blob = gh.download_bytes(art["archive_download_url"])
            grype_raw = read_zip_file_by_basename(blob, "grype-report.json")
            if not grype_raw:
                diagnostics.append(f"{art_name}: grype-report.json missing in artifact zip")
                continue
            parsed = json.loads(str(grype_raw))
            findings = normalize_fixable_findings_from_grype(parsed)
            grype_all.extend(with_service_name(findings, art_name))
        except Exception as exc:
            evidence_unavailable = True
            diagnostics.append(f"{art_name}: {exc}")

    if grype_all:
        return {
            "count": len(grype_all),
            "findings": grype_all,
            "source": f"{GRYPE_ARTIFACT} (derived, aggregated)",
            "diagnostics": diagnostics,
        }, evidence_unavailable

    return None, evidence_unavailable


def parse_matrix(
    gh: GitHubApi, artifacts: dict[str, dict[str, Any]]
) -> tuple[dict[str, Any] | None, bool]:
    candidate = None
    for name in MATRIX_ARTIFACT_CANDIDATES:
        if name in artifacts:
            candidate = artifacts[name]
            break
    if candidate is None:
        for name, art in artifacts.items():
            if "matrix" in name and "evidence" in name:
                candidate = art
                break
    if candidate is None:
        return None, False
    if candidate.get("expired"):
        return None, True

    try:
        blob = gh.download_bytes(candidate["archive_download_url"])
    except Exception:
        return None, True

    matrix_index_raw = read_zip_file_by_basename(blob, "matrix-index.json")
    summary_raw = read_zip_file_by_basename(blob, "matrix-summary.md")
    regression_raw = read_zip_file_by_basename(blob, "regression-valid-allow.json")
    metadata_raw = read_zip_file_by_basename(blob, "matrix-run-metadata.json")

    matrix_index: list[dict[str, Any]] = []
    if matrix_index_raw:
        try:
            parsed = json.loads(str(matrix_index_raw))
            if isinstance(parsed, list):
                matrix_index = parsed
        except Exception:
            matrix_index = []

    regression = None
    if regression_raw:
        try:
            regression = json.loads(str(regression_raw))
        except Exception:
            regression = None

    metadata: dict[str, Any] = {}
    if metadata_raw:
        try:
            parsed_meta = json.loads(str(metadata_raw))
            if isinstance(parsed_meta, dict):
                metadata = parsed_meta
        except Exception:
            metadata = {}

    matrix = {
        "cases": matrix_index,
        "summary_text": str(summary_raw or ""),
        "regression": regression,
        "metadata": metadata,
        "source_artifact": candidate["name"],
    }
    unavailable = not bool(matrix_index)
    return matrix, unavailable


def as_artifact_map(artifacts: list[dict[str, Any]]) -> dict[str, dict[str, Any]]:
    output: dict[str, dict[str, Any]] = {}
    for art in artifacts:
        output[str(art.get("name", ""))] = {
            "id": art.get("id"),
            "name": art.get("name"),
            "size": art.get("size_in_bytes"),
            "expired": bool(art.get("expired")),
            "created_at": art.get("created_at"),
            "updated_at": art.get("updated_at"),
            "archive_download_url": art.get("archive_download_url"),
        }
    return output


def select_workflow_ids(gh: GitHubApi) -> dict[str, dict[str, Any]]:
    workflows = gh.paginate(f"/repos/{gh.repo}/actions/workflows", "workflows")
    selected: dict[str, dict[str, Any]] = {}

    by_path = {wf.get("path"): wf for wf in workflows}
    by_name = {wf.get("name"): wf for wf in workflows}

    for key, path in WORKFLOW_PATHS.items():
        wf = by_path.get(path)
        if wf is None:
            wf = by_name.get(key)
        if wf is not None:
            selected[key] = {
                "id": wf.get("id"),
                "name": wf.get("name"),
                "path": wf.get("path"),
            }
    return selected


# --------------------------------------------------------------------------- #
# Seed helpers (services.yaml + optional snapshot-seed.json override).
# --------------------------------------------------------------------------- #


def _parse_services_yaml(path: str) -> list[dict[str, Any]]:
    """Minimal YAML reader for services.yaml -- avoids a PyYAML dep.

    services.yaml uses a simple ``services:`` list-of-maps shape. This helper
    intentionally only handles that shape; anything fancier should switch to
    PyYAML.
    """
    if not os.path.isfile(path):
        return []
    entries: list[dict[str, Any]] = []
    current: dict[str, Any] | None = None
    with open(path, "r", encoding="utf-8") as f:
        for raw_line in f:
            line = raw_line.rstrip("\n")
            stripped = line.strip()
            if not stripped or stripped.startswith("#"):
                continue
            if stripped == "services:":
                continue
            if stripped.startswith("- "):
                if current:
                    entries.append(current)
                current = {}
                stripped = stripped[2:]
            if current is None:
                continue
            m = re.match(r"^([A-Za-z_][\w-]*):\s*(.*)$", stripped)
            if not m:
                continue
            key, value = m.group(1), m.group(2).strip()
            if value.startswith('"') and value.endswith('"') and len(value) >= 2:
                value = value[1:-1]
            elif value.startswith("[") and value.endswith("]"):
                inner = value[1:-1].strip()
                items: list[str] = []
                for part in inner.split(","):
                    part = part.strip().strip('"').strip("'")
                    if part:
                        items.append(part)
                current[key] = items
                continue
            current[key] = value
        if current:
            entries.append(current)
    return entries


def build_services_section(repo_root: str) -> list[dict[str, Any]]:
    """Derive services[] from services.yaml, falling back to scaffold defaults."""
    services_yaml = os.path.join(repo_root, "services.yaml")
    raw = _parse_services_yaml(services_yaml)
    if not raw:
        # services.yaml not reachable from script CWD -- emit an empty list so
        # the v2 key is still present and the dashboard renders a "no services"
        # placeholder instead of crashing.
        return []

    out: list[dict[str, Any]] = []
    for entry in raw:
        name = str(entry.get("name", "")).strip()
        if not name:
            continue
        category = "flagship" if name == "user-service" else "scaffold"
        slsa_level = "L3" if name == "user-service" else "L2"
        out.append(
            {
                "name": name,
                "category": category,
                "loc": SCAFFOLD_LOC_DEFAULTS.get(name, 0),
                "slsa_level": slsa_level,
                "last_sbom_at": DEFAULT_LAST_SBOM_AT,
                "last_grype_high": 0,
                "covered": True,
            }
        )
    return out


def load_seed_overrides(script_dir: str) -> dict[str, Any]:
    """Load optional sibling ``snapshot-seed.json`` for operator overrides."""
    seed_path = os.path.join(script_dir, "snapshot-seed.json")
    if not os.path.isfile(seed_path):
        return {}
    try:
        with open(seed_path, "r", encoding="utf-8") as f:
            data = json.load(f)
        return data if isinstance(data, dict) else {}
    except Exception as exc:  # pragma: no cover - defensive
        print(f"warning: failed to read {seed_path}: {exc}", file=sys.stderr)
        return {}


def _rollup_for_runs(rows: list[dict[str, Any]]) -> dict[str, Any]:
    total = len(rows)
    success = sum(1 for r in rows if r.get("conclusion") == "success")
    failure = sum(1 for r in rows if r.get("conclusion") == "failure")
    pass_rate = round((success / total) * 100, 1) if total else 0.0
    latest = rows[0] if rows else {}
    return {
        "total_in_window": total,
        "success": success,
        "failure": failure,
        "pass_rate_pct": pass_rate,
        "latest_conclusion": latest.get("conclusion"),
        "latest_run_number": latest.get("run_number"),
    }


# --------------------------------------------------------------------------- #
# Top-level snapshot builder.
# --------------------------------------------------------------------------- #


def build_snapshot(
    repo: str,
    token: str,
    top_n: int,
    repo_root: str,
    script_dir: str,
) -> dict[str, Any]:
    gh = GitHubApi(repo=repo, token=token)
    workflow_meta = select_workflow_ids(gh)

    workflows_payload: list[dict[str, Any]] = []
    flat_runs: list[dict[str, Any]] = []

    for workflow_key in WORKFLOW_ORDER:
        meta = workflow_meta.get(workflow_key)
        description = WORKFLOW_DESCRIPTIONS.get(workflow_key, "")
        if not meta:
            workflows_payload.append(
                {
                    "workflow_key": workflow_key,
                    "workflow_name": workflow_key,
                    "workflow_id": None,
                    "workflow_path": WORKFLOW_PATHS[workflow_key],
                    "description": description,
                    "runs": [],
                    "rollup": _rollup_for_runs([]),
                }
            )
            continue

        runs = gh.paginate(
            f"/repos/{repo}/actions/workflows/{meta['id']}/runs",
            "workflow_runs",
            params={"exclude_pull_requests": "true"},
        )
        runs = runs[:top_n]

        run_rows: list[dict[str, Any]] = []
        for run in runs:
            artifacts_resp = gh.get_json(
                f"/repos/{repo}/actions/runs/{run['id']}/artifacts", {"per_page": 100}
            )
            artifacts = as_artifact_map(artifacts_resp.get("artifacts", []))

            security_gate = None
            matrix = None
            evidence_unavailable = False

            if workflow_key == "ci-service":
                security_gate, unavailable = parse_security_gate(gh, artifacts)
                evidence_unavailable = evidence_unavailable or unavailable
            elif workflow_key in ("admission-lab", "onboarding-lab"):
                matrix, unavailable = parse_matrix(gh, artifacts)
                evidence_unavailable = evidence_unavailable or unavailable

            row = {
                "run_key": f"{workflow_key}#{run.get('run_number')}",
                "workflow_key": workflow_key,
                "workflow_name": meta["name"],
                "workflow_path": meta["path"],
                "run_number": run.get("run_number"),
                "run_id": run.get("id"),
                "head_sha": run.get("head_sha"),
                "status": run.get("status"),
                "conclusion": run.get("conclusion"),
                "created_at": run.get("created_at"),
                "updated_at": run.get("updated_at"),
                "html_url": run.get("html_url"),
                "artifacts": artifacts,
                "security_gate": security_gate,
                "matrix": matrix,
                "evidence_unavailable": evidence_unavailable,
            }
            run_rows.append(row)
            flat_runs.append(row)

        workflows_payload.append(
            {
                "workflow_key": workflow_key,
                "workflow_name": meta["name"],
                "workflow_id": meta["id"],
                "workflow_path": meta["path"],
                "description": description,
                "runs": run_rows,
                "rollup": _rollup_for_runs(run_rows),
            }
        )

    flat_runs.sort(key=lambda x: str(x.get("created_at", "")), reverse=True)

    seed = load_seed_overrides(script_dir)
    dashboard_meta = {**DEFAULT_DASHBOARD_META, **(seed.get("dashboard_meta") or {})}
    cve_alerts = seed.get("cve_alerts") or DEFAULT_CVE_ALERTS
    go_runtime_status = {**DEFAULT_GO_RUNTIME_STATUS, **(seed.get("go_runtime_status") or {})}
    runner_pool = seed.get("runner_pool") or DEFAULT_RUNNER_POOL
    slsa_attestations = seed.get("slsa_attestations") or DEFAULT_SLSA_ATTESTATIONS

    seed_services = seed.get("services")
    services_section = (
        seed_services if isinstance(seed_services, list) and seed_services else build_services_section(repo_root)
    )

    snapshot: dict[str, Any] = {
        "schema_version": SCHEMA_VERSION,
        "generated_at": dt.datetime.now(dt.timezone.utc).isoformat().replace("+00:00", "Z"),
        "repository": repo,
        "top_n_per_workflow": top_n,
        "dashboard_meta": dashboard_meta,
        "cve_alerts": cve_alerts,
        "go_runtime_status": go_runtime_status,
        "runner_pool": runner_pool,
        "services": services_section,
        "slsa_attestations": slsa_attestations,
        "workflows": workflows_payload,
        # Flat runs are kept under a non-canonical key so the v2 front-end
        # ignores them while v1 consumers (if any remain) can still find a
        # familiar shape.
        "runs": flat_runs,
    }

    seed_notes = seed.get("notes")
    if seed_notes:
        snapshot["notes"] = seed_notes
    return snapshot


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--repo", required=True, help="owner/repo")
    parser.add_argument("--output", required=True, help="snapshot JSON output path")
    parser.add_argument("--top-n", type=int, default=100)
    parser.add_argument(
        "--repo-root",
        default=os.environ.get("GITHUB_WORKSPACE") or os.getcwd(),
        help="Repository root (used to locate services.yaml). Defaults to GITHUB_WORKSPACE or CWD.",
    )
    args = parser.parse_args()

    token = os.environ.get("GITHUB_TOKEN", "")
    if not token:
        print("GITHUB_TOKEN is required", file=sys.stderr)
        return 2

    script_dir = os.path.dirname(os.path.abspath(__file__))
    snapshot = build_snapshot(args.repo, token, args.top_n, args.repo_root, script_dir)
    os.makedirs(os.path.dirname(args.output), exist_ok=True)
    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(snapshot, f, ensure_ascii=False, indent=2)
        f.write("\n")
    print(f"Wrote snapshot: {args.output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
