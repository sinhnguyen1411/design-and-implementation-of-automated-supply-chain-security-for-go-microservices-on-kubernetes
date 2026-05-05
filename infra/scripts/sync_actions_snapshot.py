#!/usr/bin/env python3
"""Build dashboard snapshot from GitHub Actions runs/artifacts.

This script is intended to run inside GitHub Actions with GITHUB_TOKEN.
"""

from __future__ import annotations

import argparse
import datetime as dt
import io
import json
import os
import sys
import urllib.parse
import urllib.request
import zipfile
from typing import Any


WORKFLOW_PATHS = {
    "secure-supply-chain": ".github/workflows/secure-supply-chain.yml",
    "admission-matrix-evidence": ".github/workflows/admission-matrix-evidence.yml",
    "service-scs-matrix-evidence": ".github/workflows/service-scs-matrix-evidence.yml",
}

SECURITY_FINDINGS_ARTIFACT = "security-gate-findings"
GRYPE_ARTIFACT = "grype-report"
MATRIX_ARTIFACT_CANDIDATES = (
    "admission-matrix-evidence",
    "matrix-evidence",
)


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
        req = urllib.request.Request(url=url, headers=self.default_headers, method="GET")
        with urllib.request.urlopen(req, timeout=120) as resp:
            return resp.read()


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
        if fix_state in ("wont-fix", "not-fixed", "unknown"):
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

    findings_art = artifacts.get(SECURITY_FINDINGS_ARTIFACT)
    if findings_art:
        if findings_art.get("expired"):
            return None, True
        try:
            blob = gh.download_bytes(findings_art["archive_download_url"])
            findings_json_raw = read_zip_file_by_basename(blob, "security-gate-findings.json")
            if findings_json_raw:
                parsed = json.loads(str(findings_json_raw))
                findings = parsed if isinstance(parsed, list) else parsed.get("findings", [])
                if not isinstance(findings, list):
                    findings = []
                return {
                    "count": len(findings),
                    "findings": findings,
                    "source": SECURITY_FINDINGS_ARTIFACT,
                }, False
        except Exception:
            evidence_unavailable = True

    grype_art = artifacts.get(GRYPE_ARTIFACT)
    if grype_art:
        if grype_art.get("expired"):
            return None, True
        try:
            blob = gh.download_bytes(grype_art["archive_download_url"])
            grype_raw = read_zip_file_by_basename(blob, "grype-report.json")
            if grype_raw:
                parsed = json.loads(str(grype_raw))
                findings = normalize_fixable_findings_from_grype(parsed)
                return {
                    "count": len(findings),
                    "findings": findings,
                    "source": f"{GRYPE_ARTIFACT} (derived)",
                }, False
        except Exception:
            evidence_unavailable = True

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


def build_snapshot(repo: str, token: str, top_n: int) -> dict[str, Any]:
    gh = GitHubApi(repo=repo, token=token)
    workflow_meta = select_workflow_ids(gh)

    workflows_payload: list[dict[str, Any]] = []
    flat_runs: list[dict[str, Any]] = []

    for workflow_key in ("secure-supply-chain", "admission-matrix-evidence", "service-scs-matrix-evidence"):
        meta = workflow_meta.get(workflow_key)
        if not meta:
            workflows_payload.append(
                {
                    "workflow_key": workflow_key,
                    "workflow_name": workflow_key,
                    "workflow_id": None,
                    "workflow_path": WORKFLOW_PATHS[workflow_key],
                    "runs": [],
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
            artifacts_resp = gh.get_json(f"/repos/{repo}/actions/runs/{run['id']}/artifacts", {"per_page": 100})
            artifacts = as_artifact_map(artifacts_resp.get("artifacts", []))

            security_gate = None
            matrix = None
            evidence_unavailable = False

            if workflow_key == "secure-supply-chain":
                security_gate, unavailable = parse_security_gate(gh, artifacts)
                evidence_unavailable = evidence_unavailable or unavailable
            elif workflow_key in ("admission-matrix-evidence", "service-scs-matrix-evidence"):
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
                "runs": run_rows,
            }
        )

    flat_runs.sort(key=lambda x: str(x.get("created_at", "")), reverse=True)
    snapshot = {
        "schema_version": 1,
        "generated_at": dt.datetime.now(dt.timezone.utc).isoformat(),
        "repository": repo,
        "top_n_per_workflow": top_n,
        "workflows": workflows_payload,
        "runs": flat_runs,
    }
    return snapshot


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--repo", required=True, help="owner/repo")
    parser.add_argument("--output", required=True, help="snapshot JSON output path")
    parser.add_argument("--top-n", type=int, default=100)
    args = parser.parse_args()

    token = os.environ.get("GITHUB_TOKEN", "")
    if not token:
        print("GITHUB_TOKEN is required", file=sys.stderr)
        return 2

    snapshot = build_snapshot(args.repo, token, args.top_n)
    os.makedirs(os.path.dirname(args.output), exist_ok=True)
    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(snapshot, f, ensure_ascii=False, indent=2)
        f.write("\n")
    print(f"Wrote snapshot: {args.output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
