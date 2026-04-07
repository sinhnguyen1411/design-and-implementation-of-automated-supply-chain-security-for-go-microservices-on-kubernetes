"use strict";

const CASE_ORDER = [
  "VALID_ALLOW",
  "NEG_UNSIGNED_DENY",
  "NEG_MISSING_SBOM_DENY",
  "NEG_CVE_THRESHOLD_DENY",
];

const CASE_HINTS = {
  VALID_ALLOW: "Signed + attested + SBOM annotation + high_critical=0",
  NEG_UNSIGNED_DENY: "Expected deny due to missing signature evidence",
  NEG_MISSING_SBOM_DENY: "Expected deny due to missing SBOM digest annotation",
  NEG_CVE_THRESHOLD_DENY: "Expected deny due to high_critical greater than 0",
};

const SBOM_SOURCES = [
  "../../demo/sbom.spdx.json",
  "../../.tmp-sbom.json",
  "./demo-data/sbom.spdx.json",
];
const EVIDENCE_BASE_PATHS = [
  "../../demo/evidence/",
  "./demo-data/evidence/",
];

const dom = {
  runCombo: document.getElementById("run-combo"),
  runComboToggle: document.getElementById("run-combo-toggle"),
  runComboValue: document.getElementById("run-combo-value"),
  runComboPanel: document.getElementById("run-combo-panel"),
  runSearch: document.getElementById("run-search"),
  runOptions: document.getElementById("run-options"),
  refreshRuns: document.getElementById("refresh-runs"),
  statusBanner: document.getElementById("status-banner"),
  matrixGrid: document.getElementById("matrix-grid"),
  runMeta: document.getElementById("run-meta"),
  artifactContext: document.getElementById("artifact-context"),
  artifactList: document.getElementById("artifact-list"),
  summaryPreview: document.getElementById("summary-preview"),
  sbomDonut: document.getElementById("sbom-donut"),
  sbomLegend: document.getElementById("sbom-legend"),
  sbomTopList: document.getElementById("sbom-top-list"),
  sbomTotal: document.getElementById("sbom-total"),
};

const state = {
  comboOpen: false,
  activeEvidenceBasePath: "",
  availableRunIds: [],
  filteredRunIds: [],
  selectedRunId: "",
  latestRunData: {
    runId: "",
    basePath: "",
    casesByName: {},
  },
};

function setStatus(message, type) {
  dom.statusBanner.textContent = message;
  dom.statusBanner.dataset.type = type || "info";
}

function isValidRunId(value) {
  return /^[0-9]{8}-[0-9]{6}$/.test(value);
}

function sortRunsDesc(runIds) {
  return [...runIds].sort((a, b) => b.localeCompare(a));
}

function setComboDisplayValue(value) {
  dom.runComboValue.textContent = value || "Select run";
}

function setRunQueryParam(runId) {
  const url = new URL(window.location.href);
  if (runId) {
    url.searchParams.set("run", runId);
  } else {
    url.searchParams.delete("run");
  }
  window.history.replaceState(null, "", url.toString());
}

function setComboOpen(open) {
  if (dom.runComboToggle.disabled) {
    return;
  }

  state.comboOpen = open;
  dom.runComboPanel.hidden = !open;
  dom.runComboToggle.setAttribute("aria-expanded", String(open));

  if (open) {
    dom.runSearch.focus();
    dom.runSearch.select();
  }
}

function closeCombo() {
  setComboOpen(false);
}

function renderRunOptionList() {
  dom.runOptions.innerHTML = "";

  if (state.filteredRunIds.length === 0) {
    const li = document.createElement("li");
    li.className = "combo-empty";
    li.textContent = "No run id matched.";
    dom.runOptions.appendChild(li);
    return;
  }

  for (const runId of state.filteredRunIds) {
    const li = document.createElement("li");
    const button = document.createElement("button");
    button.type = "button";
    button.className = "combo-option";
    button.setAttribute("role", "option");
    button.setAttribute("aria-selected", String(runId === state.selectedRunId));
    button.dataset.runId = runId;
    button.textContent = runId;
    li.appendChild(button);
    dom.runOptions.appendChild(li);
  }
}

function filterRunOptions(keyword) {
  const term = String(keyword || "").trim().toLowerCase();
  if (!term) {
    state.filteredRunIds = [...state.availableRunIds];
  } else {
    state.filteredRunIds = state.availableRunIds.filter((runId) => runId.toLowerCase().includes(term));
  }
  renderRunOptionList();
}

function renderRunOptions(runIds) {
  state.availableRunIds = [...runIds];
  state.filteredRunIds = [...runIds];
  dom.runSearch.value = "";

  if (runIds.length === 0) {
    dom.runComboToggle.disabled = true;
    setComboDisplayValue("No runs found");
    renderRunOptionList();
    return;
  }

  dom.runComboToggle.disabled = false;
  if (!state.selectedRunId || !runIds.includes(state.selectedRunId)) {
    state.selectedRunId = "";
    setComboDisplayValue("Select run");
  } else {
    setComboDisplayValue(state.selectedRunId);
  }
  renderRunOptionList();
}

function parseRunIdsFromDirectoryListing(html) {
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, "text/html");
  const anchors = Array.from(doc.querySelectorAll("a[href]"));
  const runSet = new Set();

  for (const anchor of anchors) {
    const hrefRaw = String(anchor.getAttribute("href") || "");
    const href = decodeURIComponent(hrefRaw).replace(/^\.\//, "");
    const match = href.match(/^([0-9]{8}-[0-9]{6})\/?$/);
    if (match) {
      runSet.add(match[1]);
    }
  }

  return sortRunsDesc(Array.from(runSet));
}

async function fetchJson(path) {
  const response = await fetch(path, { cache: "no-store" });
  if (!response.ok) {
    throw new Error(`Failed to fetch ${path}: HTTP ${response.status}`);
  }
  return response.json();
}

async function fetchText(path) {
  const response = await fetch(path, { cache: "no-store" });
  if (!response.ok) {
    throw new Error(`Failed to fetch ${path}: HTTP ${response.status}`);
  }
  return response.text();
}

function byCaseName(list) {
  const result = {};
  if (!Array.isArray(list)) {
    return result;
  }
  for (const item of list) {
    if (item && typeof item.case === "string") {
      result[item.case] = item;
    }
  }
  return result;
}

function parseRunMetadata(summaryText) {
  const lines = summaryText.split(/\r?\n/);
  const pairs = {};
  for (const line of lines) {
    const match = line.match(/^- ([^:]+):\s*(.*)$/);
    if (match) {
      pairs[match[1].trim()] = match[2].trim();
    }
  }
  return {
    runId: pairs["Run ID"] || "-",
    context: pairs["Kubernetes context"] || "-",
    namespace: pairs["Namespace"] || "-",
    signedDigest: pairs["Signed image digest"] || "-",
    unsignedDigest: pairs["Unsigned image digest"] || "-",
    sbomDigest: pairs["SBOM digest"] || "-",
  };
}

function renderRunMeta(meta) {
  const entries = [
    ["Run ID", meta.runId],
    ["Context", meta.context],
    ["Namespace", meta.namespace],
    ["Signed Digest", meta.signedDigest],
    ["Unsigned Digest", meta.unsignedDigest],
    ["SBOM Digest", meta.sbomDigest],
  ];
  dom.runMeta.innerHTML = "";
  for (const [label, value] of entries) {
    const wrapper = document.createElement("div");
    const dt = document.createElement("dt");
    const dd = document.createElement("dd");
    dt.textContent = label;
    dd.textContent = value || "-";
    wrapper.appendChild(dt);
    wrapper.appendChild(dd);
    dom.runMeta.appendChild(wrapper);
  }
}

function verdictClass(verdict) {
  if (verdict === "PASS") return "pass";
  if (verdict === "FAIL") return "fail";
  return "neutral";
}

function signalClass(value) {
  const normalized = String(value || "").toLowerCase();
  if (normalized.includes("allow")) return "signal-allowed";
  if (normalized.includes("deny")) return "signal-denied";
  return "signal-neutral";
}

function matrixMissingCard(caseName) {
  return `
    <article class="matrix-card" data-missing="true" data-verdict="MISSING">
      <div class="matrix-head">
        <div>
          <div class="case-id">${caseName}</div>
          <div class="matrix-detail">${CASE_HINTS[caseName]}</div>
        </div>
        <span class="pill neutral">MISSING</span>
      </div>
      <div class="matrix-detail">
        Case data is not present in matrix-index.json for this run.
      </div>
    </article>
  `;
}

function renderMatrix(basePath, casesByName) {
  let html = "";
  for (const caseName of CASE_ORDER) {
    const item = casesByName[caseName];
    if (!item) {
      html += matrixMissingCard(caseName);
      continue;
    }

    const verdict = String(item.verdict || "UNKNOWN");
    const actual = String(item.actual || "UNKNOWN");
    const expected = String(item.expected || "UNKNOWN");
    const reason = String(item.reason || "-");
    const includeSbom = String(item.include_sbom);
    const highCritical = String(item.high_critical ?? "-");
    const applyExit = String(item.apply_exit_code ?? "-");
    const waitExit = String(item.wait_exit_code ?? "-");
    const image = String(item.image || "-");
    const artifactCount = item.artifacts && typeof item.artifacts === "object"
      ? Object.keys(item.artifacts).length
      : 0;

    html += `
      <article class="matrix-card" data-case="${caseName}" data-verdict="${verdict}">
        <div class="matrix-head">
          <div>
            <div class="case-id">${caseName}</div>
            <div class="matrix-detail">${CASE_HINTS[caseName]}</div>
          </div>
          <span class="pill ${verdictClass(verdict)}">${verdict}</span>
        </div>
        <div class="matrix-detail">
          <strong>Expected:</strong> ${expected}<br />
          <strong>Actual:</strong> <span class="${signalClass(actual)}">${actual}</span><br />
          <strong>Reason:</strong> ${reason}
          <dl>
            <dt>include_sbom</dt><dd>${includeSbom}</dd>
            <dt>high_critical</dt><dd>${highCritical}</dd>
            <dt>apply_exit_code</dt><dd>${applyExit}</dd>
            <dt>wait_exit_code</dt><dd>${waitExit}</dd>
            <dt>image</dt><dd>${image}</dd>
          </dl>
        </div>
        <div class="matrix-actions">
          <span class="artifact-count">Artifacts: ${artifactCount}</span>
          <button type="button" data-action="show-artifacts" data-case="${caseName}">
            View Artifacts
          </button>
        </div>
      </article>
    `;
  }

  dom.matrixGrid.innerHTML = html;
  for (const button of dom.matrixGrid.querySelectorAll("button[data-action='show-artifacts']")) {
    button.addEventListener("click", () => {
      const caseName = button.getAttribute("data-case") || "";
      renderArtifacts(basePath, caseName);
    });
  }
}

function renderArtifacts(basePath, caseName) {
  const item = state.latestRunData.casesByName[caseName];
  dom.artifactList.innerHTML = "";

  if (!item || !item.artifacts || typeof item.artifacts !== "object") {
    dom.artifactContext.textContent = `No artifact map for case ${caseName}.`;
    return;
  }

  dom.artifactContext.textContent = `Showing artifacts for ${caseName}`;
  const keys = Object.keys(item.artifacts);
  if (keys.length === 0) {
    const li = document.createElement("li");
    li.textContent = "No artifact files listed.";
    dom.artifactList.appendChild(li);
    return;
  }

  for (const key of keys) {
    const relativePath = String(item.artifacts[key]);
    const absolutePath = `${basePath}/${relativePath}`;
    const li = document.createElement("li");
    const link = document.createElement("a");
    link.href = absolutePath;
    link.target = "_blank";
    link.rel = "noopener noreferrer";
    link.textContent = `${key}: ${relativePath}`;
    li.appendChild(link);
    dom.artifactList.appendChild(li);
  }
}

function categorizePackage(pkg) {
  const refs = Array.isArray(pkg.externalRefs) ? pkg.externalRefs : [];
  const purlRef = refs.find((ref) => ref && String(ref.referenceType || "").toLowerCase() === "purl");
  const purl = purlRef ? String(purlRef.referenceLocator || "").toLowerCase() : "";

  if (purl.startsWith("pkg:golang/")) return "golang";
  if (purl.startsWith("pkg:deb/")) return "deb";
  return "other";
}

function renderSbomSummary(packages) {
  const counts = { golang: 0, deb: 0, other: 0 };
  for (const pkg of packages) {
    counts[categorizePackage(pkg)] += 1;
  }

  const total = packages.length;
  dom.sbomTotal.textContent = `Total packages: ${total}`;

  if (total === 0) {
    dom.sbomLegend.innerHTML = "<p>No package entries found in SBOM.</p>";
    dom.sbomTopList.innerHTML = "<li>No dependencies available.</li>";
    return;
  }

  const golangPct = (counts.golang / total) * 100;
  const debPct = (counts.deb / total) * 100;
  dom.sbomDonut.style.background = `conic-gradient(
    var(--accent-cyan) 0 ${golangPct}%,
    var(--accent-blue) ${golangPct}% ${golangPct + debPct}%,
    var(--accent-violet) ${golangPct + debPct}% 100%
  )`;

  const rows = [
    { key: "golang", color: "var(--accent-cyan)", count: counts.golang },
    { key: "deb", color: "var(--accent-blue)", count: counts.deb },
    { key: "other", color: "var(--accent-violet)", count: counts.other },
  ];

  dom.sbomLegend.innerHTML = rows.map((row) => {
    const percent = total === 0 ? "0.0" : ((row.count / total) * 100).toFixed(1);
    return `
      <div class="legend-row">
        <span class="legend-label"><span class="legend-dot" style="background:${row.color}"></span>${row.key}</span>
        <span>${row.count} (${percent}%)</span>
      </div>
    `;
  }).join("");

  const topPackages = packages
    .map((pkg) => {
      const version = pkg.versionInfo ? `@${pkg.versionInfo}` : "";
      return {
        name: String(pkg.name || "unknown"),
        text: `${String(pkg.name || "unknown")}${version}`,
      };
    })
    .sort((a, b) => a.name.localeCompare(b.name))
    .slice(0, 10);

  dom.sbomTopList.innerHTML = topPackages.map((pkg) => `<li>${pkg.text}</li>`).join("");
}

async function loadSbom() {
  for (const source of SBOM_SOURCES) {
    try {
      const data = await fetchJson(source);
      const packages = Array.isArray(data.packages) ? data.packages : [];
      renderSbomSummary(packages);
      return source;
    } catch (error) {
      // Continue with fallback source.
    }
  }

  dom.sbomDonut.style.background = "conic-gradient(var(--warn) 0 100%)";
  dom.sbomLegend.innerHTML = "<p>SBOM not found. Checked ../../demo/sbom.spdx.json and ../../.tmp-sbom.json.</p>";
  dom.sbomTopList.innerHTML = "<li>SBOM file is missing.</li>";
  dom.sbomTotal.textContent = "Total packages: -";
  return "";
}

async function loadRun(runId) {
  if (!isValidRunId(runId)) {
    setStatus("Invalid run ID. Expected format: YYYYMMDD-HHMMSS", "error");
    return;
  }

  if (!state.activeEvidenceBasePath) {
    setStatus("No evidence source is active. Click Refresh Runs.", "error");
    return;
  }
  const basePath = `${state.activeEvidenceBasePath}${runId}`;
  setStatus(`Loading evidence from ${basePath} ...`, "info");

  try {
    const [matrixIndex, summaryText] = await Promise.all([
      fetchJson(`${basePath}/matrix-index.json`),
      fetchText(`${basePath}/matrix-summary.md`),
    ]);

    let regression = null;
    try {
      regression = await fetchJson(`${basePath}/regression-valid-allow.json`);
    } catch (error) {
      regression = null;
    }

    const casesByName = byCaseName(matrixIndex);
    state.latestRunData = {
      runId,
      basePath,
      casesByName,
    };

    renderRunMeta(parseRunMetadata(summaryText));
    renderMatrix(basePath, casesByName);
    dom.summaryPreview.textContent = summaryText.trim() || "Summary file is empty.";

    const firstSelectable = CASE_ORDER.find((caseName) => Boolean(casesByName[caseName]));
    if (firstSelectable) {
      renderArtifacts(basePath, firstSelectable);
    } else {
      dom.artifactContext.textContent = "No known matrix cases found for this run.";
      dom.artifactList.innerHTML = "<li>No artifacts available.</li>";
    }

    if (regression && regression.verdict) {
      setStatus(`Loaded run ${runId}. Regression check: ${regression.verdict}.`, "success");
    } else {
      setStatus(`Loaded run ${runId}. Regression file is missing or unreadable.`, "warning");
    }
  } catch (error) {
    dom.matrixGrid.innerHTML = "";
    dom.artifactList.innerHTML = "";
    dom.summaryPreview.textContent = "No summary loaded.";
    setStatus(`Unable to load evidence for run ${runId}. Check that files exist under ../../demo/evidence/${runId}.`, "error");
  }
}

async function chooseRun(runId, options) {
  const opts = Object.assign({
    load: true,
    close: true,
    updateQuery: true,
  }, options || {});

  if (!runId) {
    return;
  }

  state.selectedRunId = runId;
  setComboDisplayValue(runId);
  renderRunOptionList();

  if (opts.updateQuery) {
    setRunQueryParam(runId);
  }

  if (opts.close) {
    closeCombo();
  }

  if (opts.load) {
    await loadRun(runId);
  }
}

async function discoverAndLoadRuns(preferredRunId) {
  setStatus("Scanning evidence run directories...", "info");

  for (const evidenceBasePath of EVIDENCE_BASE_PATHS) {
    try {
      const html = await fetchText(evidenceBasePath);
      const runIds = parseRunIdsFromDirectoryListing(html);
      if (runIds.length === 0) {
        continue;
      }

      state.activeEvidenceBasePath = evidenceBasePath;
      renderRunOptions(runIds);

      let selectedRun = runIds[0];
      if (preferredRunId && isValidRunId(preferredRunId) && runIds.includes(preferredRunId)) {
        selectedRun = preferredRunId;
      } else if (preferredRunId && isValidRunId(preferredRunId) && !runIds.includes(preferredRunId)) {
        setStatus(`Run ${preferredRunId} not found. Loading latest run ${selectedRun}.`, "warning");
      }

      await chooseRun(selectedRun, { load: true, close: true, updateQuery: true });
      return;
    } catch (error) {
      // Try next source.
    }
  }

  state.activeEvidenceBasePath = "";
  renderRunOptions([]);
  setStatus("No evidence runs found. Provide demo/evidence runs or use bundled demo-data.", "error");
}

function initComboEvents() {
  dom.runComboToggle.addEventListener("click", () => {
    setComboOpen(!state.comboOpen);
  });

  dom.runSearch.addEventListener("input", () => {
    filterRunOptions(dom.runSearch.value);
  });

  dom.runSearch.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      closeCombo();
      dom.runComboToggle.focus();
    }
  });

  dom.runOptions.addEventListener("click", (event) => {
    const target = event.target;
    if (!(target instanceof HTMLElement)) {
      return;
    }
    const option = target.closest(".combo-option");
    if (!option) {
      return;
    }
    const runId = String(option.dataset.runId || "");
    if (!runId) {
      return;
    }
    void chooseRun(runId, { load: true, close: true, updateQuery: true });
  });

  document.addEventListener("click", (event) => {
    const target = event.target;
    if (!(target instanceof Node)) {
      return;
    }
    if (!dom.runCombo.contains(target)) {
      closeCombo();
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && state.comboOpen) {
      closeCombo();
    }
  });
}

function init() {
  initComboEvents();

  dom.refreshRuns.addEventListener("click", () => {
    void discoverAndLoadRuns(state.selectedRunId);
  });

  const params = new URLSearchParams(window.location.search);
  const runFromQuery = String(params.get("run") || "").trim();

  void discoverAndLoadRuns(runFromQuery);
  void loadSbom();
}

init();
