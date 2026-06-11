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

const ALL_SERVICES_FALLBACK = [
  "user-service", "portfolio-service", "order-service", "risk-service",
  "market-data-service", "pricing-service", "execution-service", "settlement-service",
  "compliance-service", "notification-service", "apikey-service", "kyc-service",
  "watchlist-service", "analytics-service", "audit-service", "fees-service",
  "reporting-service", "gateway-service", "search-service", "alert-service",
  "data-feed-service", "backtest-service", "margin-service",
];

const WORKFLOW_META = {
  "ci-service": "Build · SBOM · Grype CVE scan per service",
  "reusable-go-verify": "Cross-OS unit tests · govulncheck",
  "admission-lab": "Kyverno admission matrix evidence",
  "onboarding-lab": "Per-service Kind cluster probe",
  "dashboard-data-sync": "Actions snapshot auto-sync",
  "runner-ab-benchmark": "Windows runner A/B benchmark",
};

const CASE_HINTS_VI = {
  VALID_ALLOW: "Da ky + co attestation + co annotation SBOM + high_critical=0",
  NEG_UNSIGNED_DENY: "Du kien bi tu choi do thieu bang chung chu ky",
  NEG_MISSING_SBOM_DENY: "Du kien bi tu choi do thieu annotation SBOM digest",
  NEG_CVE_THRESHOLD_DENY: "Du kien bi tu choi do high_critical lon hon 0",
};

const I18N = {
  en: {
    pageTitle: "Supply Chain Security Dashboard — 23 Microservices",
    eyebrow: "Supply Chain Security Dashboard",
    heroTitle: "DevSecOps Evidence Observatory — 23 Microservices on Kubernetes",
    heroCopy: "Automated supply chain security evidence: SBOM, CVE scanning, Kyverno admission control, SLSA provenance. Powered by 6 GitHub Actions workflows across 23 Go microservices.",
    pipelineStatus: "CI Pipeline Status",
    pipelineStatusNote: "Last run per workflow from Actions snapshot",
    serviceCoverage: "Service Coverage",
    serviceCoverageNote: "Supply chain artifact presence per service from latest ci-service run",
    refreshRuns: "Refresh Runs",
    runHelp: "Click to open command palette, search run-id, then select to load automatically.",
    runSearchPlaceholder: "Search run id...",
    overviewAvailableRuns: "Available Runs",
    overviewCiHealth: "CI Health",
    overviewMatrixVerdict: "Expectation Match",
    overviewAdmissionOutcome: "Admission Outcome",
    overviewEvidenceFootprint: "Evidence Footprint",
    overviewLastGreen: "Last Green Build",
    runMetadata: "Run Metadata",
    quickRawFiles: "Quick Raw Files",
    quickRawFilesNote: "Open scanner and SBOM outputs used for presentation evidence.",
    linkGrype: "Open `../../.tmp-grype.json`",
    linkGateFindings: "Open `../../security-gate-findings.json`",
    linkTmpSbom: "Open `../../.tmp-sbom.json`",
    linkDemoSbom: "Open `../../demo/sbom.spdx.json`",
    linkBundledSbom: "Open `./demo-data/sbom.spdx.json` (bundled)",
    linkActionsSnapshot: "Open `./data/actions-runs.snapshot.json` (Actions snapshot)",
    linkBundledNotes: "Open bundled evidence notes",
    linkAlertCve: "Jump to Active CVE Alerts",
    linkServicesYaml: "Open `services.yaml`",
    linkCiServiceYml: "Open `.github/workflows/ci-service.yml`",
    linkSlsaVerifier: "SLSA verifier docs",
    quickRawFilesSourceHead: "Source-of-truth files",
    activeCveAlerts: "Active CVE Alerts",
    activeCveAlertsNote: "Blocking advisories pinned at top.",
    goRuntimeStatus: "Go Runtime Status",
    goRuntimeStatusNote: "Pinned vs required toolchain · remediation plan.",
    ciTimeline: "7-Day CI Timeline",
    ciTimelineNote: "Last 10 ci-service runs · click a dot to open run on GitHub.",
    slsaTracker: "SLSA Attestation Tracker",
    slsaTrackerNote: "L3 user-service flagship · L2 scaffold fleet · slsa-verifier reference.",
    runnerPool: "Runner Pool",
    runnerPoolNote: "GitHub-hosted + self-hosted runners used by the 6 workflows.",
    runnerAb: "Runner A/B Benchmark",
    runnerAbNote: "windows-latest vs THEMONSTER-win-parity (self-hosted).",
    serviceInventory: "Service Inventory (23 microservices)",
    serviceInventoryNote: "Sortable detail: category · LOC · SLSA level · last SBOM · last grype high count.",
    invName: "Service",
    invCategory: "Category",
    invLoc: "LOC",
    invSlsa: "SLSA",
    invSbom: "Last SBOM",
    invGrype: "Last Grype High",
    invCovered: "Covered",
    sbomSubtitle: "SBOM from user-service flagship build #114 (last green) · stdlib not included in 3rd-party dep view.",
    admissionUnaffectedNote: "admission-lab is GREEN (last 3 runs all 4/4) — unaffected by the stdlib CVE because it uses pre-built signed test images.",
    ciRedStreak: "Red streak: {count} runs · since {firstFail}",
    redStreakBadge: "{count}-run red streak",
    relTimeNow: "just now",
    relTimeMinsAgo: "{n}m ago",
    relTimeHoursAgo: "{n}h ago",
    relTimeDaysAgo: "{n}d ago",
    securityGateFindings: "Security Gate Findings",
    cveLoading: "Loading CVE findings...",
    cveId: "CVE",
    cveSeverity: "Severity",
    cvePackage: "Package",
    cveInstalled: "Installed",
    cveFixed: "Fixed Versions",
    cveNoFindings: "No fixable high/critical CVEs found in loaded source.",
    cveSourceSummary: "Source: {source}. Fixable high/critical: {count}.",
    cveSourceUnavailable: "No CVE findings source found. Expected one of: ../../security-gate-findings.json, ../../.tmp-grype.json, or bundled sample.",
    cveSnapshotV2Note: "Table summarised in Active CVE Alerts above — per-run findings not embedded in v2 snapshot, see ci-summary artifact.",
    cveSnapshotV2JumpLink: "Jump to Active CVE Alerts",
    admissionMatrix: "Admission Matrix",
    admissionMatrixNote: "Fixed scenario order: VALID_ALLOW, NEG_UNSIGNED_DENY, NEG_MISSING_SBOM_DENY, NEG_CVE_THRESHOLD_DENY",
    sbomDependencyView: "SBOM Dependency View",
    sbomDependencyNote: "Grouped by package type: golang, deb, other",
    topDependencies: "Top Dependencies",
    artifactExplorer: "Artifact Explorer",
    matrixSummary: "Matrix Summary (Markdown)",
    selectRun: "Select run",
    noRunsFound: "No runs found",
    noRunMatched: "No run id matched.",
    runLabel: "Evidence Run ID",
    latestRuns: "Latest IDs: {runs}",
    noRunsDetected: "No evidence runs detected.",
    noRunLoaded: "No run loaded.",
    noPolicyDecision: "No policy decision loaded.",
    waitingArtifacts: "Waiting for artifact inventory.",
    passRate: "{rate}% PASS",
    caseSummary: "{passCount}/{totalCases} scenarios matched expected admission behavior.",
    admissionSummary: "{allowCount} allow / {denyCount} deny",
    trackedFiles: "Tracked files across {totalCases} matrix scenarios.",
    policyDerived: "Derived from actual admission outcomes in the selected run.",
    verdictMatched: "MATCHED EXPECTATION",
    verdictMismatch: "MISMATCH",
    metaLabels: ["Run ID", "Context", "Namespace", "Signed Digest", "Unsigned Digest", "SBOM Digest"],
    expected: "Expected",
    actual: "Actual",
    reason: "Reason",
    artifacts: "Artifacts",
    viewArtifacts: "View Artifacts",
    missing: "MISSING",
    missingCaseText: "Case data is not present in matrix-index.json for this run.",
    noArtifactMap: "No artifact map for case {caseName}.",
    showingArtifacts: "Showing artifacts for {caseName}",
    noArtifactFiles: "No artifact files listed.",
    chooseCase: "Choose a case from the matrix to inspect evidence files.",
    noKnownCases: "No known matrix cases found for this run.",
    noArtifactsAvailable: "No artifacts available.",
    summaryNotLoaded: "No summary loaded.",
    summaryEmpty: "Summary file is empty.",
    totalPackages: "Total packages: {total}",
    noPackagesFound: "No package entries found in SBOM.",
    noDependencies: "No dependencies available.",
    sbomNotFound: "SBOM not found. Checked ../../demo/sbom.spdx.json and ../../.tmp-sbom.json.",
    sbomMissing: "SBOM file is missing.",
    statusScanningRuns: "Scanning evidence run directories...",
    statusInvalidRun: "Invalid run ID. Expected format: YYYYMMDD-HHMMSS or gha:workflow#run_number",
    statusNoSource: "No evidence source is active. Click Refresh Runs.",
    statusLoadingEvidence: "Loading evidence from {basePath} ...",
    statusRunLoaded: "Loaded run {runId}. Regression check: {result}.",
    statusRegressionMissing: "Loaded run {runId}. Regression file is missing or unreadable.",
    statusLoadRunFailed: "Unable to load evidence for run {runId}. Check that files exist under ../../demo/evidence/{runId}.",
    statusRunFallback: "Run {preferredRunId} not found. Loading latest run {selectedRun}.",
    statusSnapshotRunMissing: "Snapshot run not found: {runId}",
    statusSnapshotLoading: "Loading GitHub Actions run {runId} ...",
    statusSnapshotLoaded: "Loaded {runId}. Regression: {result}.",
    statusSnapshotPartial: "Loaded {runId} with partial or expired artifacts.",
    statusSnapshotRegressionMissing: "Loaded {runId}. Regression evidence is missing.",
    statusSnapshotFallback: "Run {preferredRunId} not found in Actions snapshot. Loading {selectedRun}.",
    statusNoRuns: "No evidence runs found. Provide demo/evidence runs or use bundled demo-data."
  },
  vi: {
    pageTitle: "Bảng điều khiển bảo mật chuỗi cung ứng — 23 microservice",
    eyebrow: "Bảng điều khiển bảo mật chuỗi cung ứng",
    heroTitle: "Đài quan sát bằng chứng DevSecOps — 23 microservice trên Kubernetes",
    heroCopy: "Bằng chứng bảo mật chuỗi cung ứng tự động: SBOM, quét CVE, kiểm soát admission bằng Kyverno, chứng thực SLSA. Vận hành bởi 6 workflow GitHub Actions trên 23 microservice Go.",
    pipelineStatus: "Trạng thái pipeline CI",
    pipelineStatusNote: "Lần chạy gần nhất của mỗi workflow lấy từ snapshot Actions",
    serviceCoverage: "Mức độ bao phủ dịch vụ",
    serviceCoverageNote: "Sự hiện diện của artifact chuỗi cung ứng cho mỗi dịch vụ từ lần chạy ci-service mới nhất",
    refreshRuns: "Làm mới danh sách run",
    runHelp: "Nhấn để mở bảng lệnh, tìm theo run-id, sau đó chọn để tải tự động.",
    runSearchPlaceholder: "Tìm theo run id...",
    overviewAvailableRuns: "Số run khả dụng",
    overviewCiHealth: "Tình trạng CI",
    overviewMatrixVerdict: "Mức độ khớp kỳ vọng",
    overviewAdmissionOutcome: "Kết quả admission",
    overviewEvidenceFootprint: "Quy mô bằng chứng",
    overviewLastGreen: "Bản build xanh gần nhất",
    runMetadata: "Siêu dữ liệu của run",
    quickRawFiles: "Tệp thô truy cập nhanh",
    quickRawFilesNote: "Mở các tệp đầu ra của scanner và SBOM được dùng làm bằng chứng trình bày.",
    linkGrype: "Mở `../../.tmp-grype.json`",
    linkGateFindings: "Mở `../../security-gate-findings.json`",
    linkTmpSbom: "Mở `../../.tmp-sbom.json`",
    linkDemoSbom: "Mở `../../demo/sbom.spdx.json`",
    linkBundledSbom: "Mở `./demo-data/sbom.spdx.json` (đóng gói kèm)",
    linkActionsSnapshot: "Mở `./data/actions-runs.snapshot.json` (snapshot Actions)",
    linkBundledNotes: "Mở ghi chú bằng chứng đóng gói kèm",
    linkAlertCve: "Chuyển đến mục Cảnh báo CVE đang hoạt động",
    linkServicesYaml: "Mở `services.yaml`",
    linkCiServiceYml: "Mở `.github/workflows/ci-service.yml`",
    linkSlsaVerifier: "Tài liệu SLSA verifier",
    quickRawFilesSourceHead: "Tệp nguồn gốc",
    activeCveAlerts: "Cảnh báo CVE đang hoạt động",
    activeCveAlertsNote: "Các advisory đang chặn được ghim ở đầu.",
    goRuntimeStatus: "Trạng thái runtime Go",
    goRuntimeStatusNote: "Phiên bản toolchain đã ghim so với yêu cầu · kế hoạch khắc phục.",
    ciTimeline: "Dòng thời gian CI 7 ngày",
    ciTimelineNote: "10 lần chạy ci-service gần nhất · nhấn vào chấm để mở run trên GitHub.",
    slsaTracker: "Theo dõi chứng thực SLSA",
    slsaTrackerNote: "L3 cho flagship user-service · L2 cho cụm scaffold · tham chiếu slsa-verifier.",
    runnerPool: "Bể runner",
    runnerPoolNote: "Các runner GitHub-hosted + self-hosted được dùng bởi 6 workflow.",
    runnerAb: "Benchmark A/B runner",
    runnerAbNote: "windows-latest so với THEMONSTER-win-parity (self-hosted).",
    serviceInventory: "Danh mục dịch vụ (23 microservice)",
    serviceInventoryNote: "Chi tiết có thể sắp xếp: danh mục · LOC · mức SLSA · SBOM gần nhất · số CVE high của lần grype gần nhất.",
    invName: "Dịch vụ",
    invCategory: "Danh mục",
    invLoc: "LOC",
    invSlsa: "SLSA",
    invSbom: "SBOM gần nhất",
    invGrype: "Grype high gần nhất",
    invCovered: "Đã bao phủ",
    sbomSubtitle: "SBOM lấy từ bản build flagship user-service #114 (lần xanh gần nhất) · không bao gồm stdlib trong góc nhìn phụ thuộc bên thứ ba.",
    admissionUnaffectedNote: "admission-lab đang GREEN (3 lần chạy gần nhất đều 4/4) — không bị ảnh hưởng bởi CVE của stdlib vì sử dụng image kiểm thử đã được ký sẵn.",
    ciRedStreak: "Chuỗi đỏ: {count} lần · kể từ {firstFail}",
    redStreakBadge: "Chuỗi đỏ {count} lần",
    relTimeNow: "vừa xong",
    relTimeMinsAgo: "{n} phút trước",
    relTimeHoursAgo: "{n} giờ trước",
    relTimeDaysAgo: "{n} ngày trước",
    securityGateFindings: "Kết quả Security Gate",
    cveLoading: "Đang tải danh sách CVE...",
    cveId: "CVE",
    cveSeverity: "Mức nghiêm trọng",
    cvePackage: "Gói",
    cveInstalled: "Phiên bản đang cài",
    cveFixed: "Phiên bản đã vá",
    cveNoFindings: "Không phát hiện CVE high/critical có bản vá trong nguồn đã tải.",
    cveSourceSummary: "Nguồn: {source}. Số CVE high/critical có bản vá: {count}.",
    cveSourceUnavailable: "Không tìm thấy nguồn dữ liệu CVE. Cần một trong các tệp: ../../security-gate-findings.json, ../../.tmp-grype.json, hoặc mẫu đóng gói kèm.",
    cveSnapshotV2Note: "Bảng được tổng hợp trong mục Cảnh báo CVE đang hoạt động phía trên — chi tiết theo từng run không được nhúng trong snapshot v2, vui lòng xem artifact ci-summary.",
    cveSnapshotV2JumpLink: "Chuyển đến mục Cảnh báo CVE đang hoạt động",
    admissionMatrix: "Ma trận admission",
    admissionMatrixNote: "Thứ tự kịch bản cố định: VALID_ALLOW, NEG_UNSIGNED_DENY, NEG_MISSING_SBOM_DENY, NEG_CVE_THRESHOLD_DENY",
    sbomDependencyView: "Góc nhìn phụ thuộc SBOM",
    sbomDependencyNote: "Nhóm theo loại gói: golang, deb, khác",
    topDependencies: "Phụ thuộc hàng đầu",
    artifactExplorer: "Trình duyệt artifact",
    matrixSummary: "Tóm tắt ma trận (Markdown)",
    selectRun: "Chọn run",
    noRunsFound: "Không tìm thấy run nào",
    noRunMatched: "Không có run id nào khớp.",
    runLabel: "Run ID của bằng chứng",
    latestRuns: "ID gần nhất: {runs}",
    noRunsDetected: "Chưa phát hiện run bằng chứng nào.",
    noRunLoaded: "Chưa tải run nào.",
    noPolicyDecision: "Chưa tải quyết định chính sách nào.",
    waitingArtifacts: "Đang chờ kiểm kê artifact.",
    passRate: "{rate}% PASS",
    caseSummary: "{passCount}/{totalCases} kịch bản khớp với hành vi admission kỳ vọng.",
    admissionSummary: "{allowCount} cho phép / {denyCount} từ chối",
    trackedFiles: "Đã theo dõi tệp trên {totalCases} kịch bản của ma trận.",
    policyDerived: "Suy ra từ kết quả admission thực tế trong run đã chọn.",
    verdictMatched: "KHỚP KỲ VỌNG",
    verdictMismatch: "KHÔNG KHỚP",
    metaLabels: ["Run ID", "Ngữ cảnh", "Namespace", "Digest đã ký", "Digest chưa ký", "Digest SBOM"],
    expected: "Kỳ vọng",
    actual: "Thực tế",
    reason: "Lý do",
    artifacts: "Artifact",
    viewArtifacts: "Xem artifact",
    missing: "THIẾU",
    missingCaseText: "Dữ liệu kịch bản không có trong matrix-index.json của run này.",
    noArtifactMap: "Không có bản đồ artifact cho kịch bản {caseName}.",
    showingArtifacts: "Đang hiển thị artifact của {caseName}",
    noArtifactFiles: "Không có tệp artifact nào được liệt kê.",
    chooseCase: "Chọn một kịch bản từ ma trận để xem các tệp bằng chứng.",
    noKnownCases: "Không tìm thấy kịch bản ma trận đã biết cho run này.",
    noArtifactsAvailable: "Không có artifact nào khả dụng.",
    summaryNotLoaded: "Chưa tải bản tóm tắt nào.",
    summaryEmpty: "Tệp tóm tắt trống.",
    totalPackages: "Tổng số gói: {total}",
    noPackagesFound: "Không tìm thấy mục gói nào trong SBOM.",
    noDependencies: "Không có phụ thuộc nào khả dụng.",
    sbomNotFound: "Không tìm thấy SBOM. Đã kiểm tra ../../demo/sbom.spdx.json và ../../.tmp-sbom.json.",
    sbomMissing: "Thiếu tệp SBOM.",
    statusScanningRuns: "Đang quét các thư mục run bằng chứng...",
    statusInvalidRun: "Run ID không hợp lệ. Định dạng yêu cầu: YYYYMMDD-HHMMSS hoặc gha:workflow#run_number",
    statusNoSource: "Chưa có nguồn bằng chứng nào đang hoạt động. Nhấn Làm mới danh sách run.",
    statusLoadingEvidence: "Đang tải bằng chứng từ {basePath} ...",
    statusRunLoaded: "Đã tải run {runId}. Kiểm tra hồi quy: {result}.",
    statusRegressionMissing: "Đã tải run {runId}. Tệp kiểm tra hồi quy bị thiếu hoặc không đọc được.",
    statusLoadRunFailed: "Không thể tải bằng chứng cho run {runId}. Vui lòng kiểm tra các tệp tại ../../demo/evidence/{runId}.",
    statusRunFallback: "Không tìm thấy run {preferredRunId}. Đang tải run gần nhất {selectedRun}.",
    statusSnapshotRunMissing: "Không tìm thấy run trong snapshot: {runId}",
    statusSnapshotLoading: "Đang tải run GitHub Actions {runId} ...",
    statusSnapshotLoaded: "Đã tải {runId}. Hồi quy: {result}.",
    statusSnapshotPartial: "Đã tải {runId} với artifact một phần hoặc đã hết hạn.",
    statusSnapshotRegressionMissing: "Đã tải {runId}. Thiếu bằng chứng hồi quy.",
    statusSnapshotFallback: "Không tìm thấy run {preferredRunId} trong snapshot Actions. Đang tải {selectedRun}.",
    statusNoRuns: "Không tìm thấy run bằng chứng nào. Vui lòng cung cấp các run demo/evidence hoặc dùng demo-data đóng gói kèm."
  }
};


const SBOM_SOURCES = [
  "../../demo/sbom.spdx.json",
  "../../.tmp-sbom.json",
  "./demo-data/sbom.spdx.json",
];

const GATE_FINDINGS_SOURCES = [
  "../../security-gate-findings.json",
  "../../.tmp-grype.json",
  "./demo-data/security-gate-findings.sample.json",
];

const ACTIONS_SNAPSHOT_PATH = "./data/actions-runs.snapshot.json";

const EVIDENCE_BASE_PATHS = [
  "../../demo/evidence/",
  "./demo-data/evidence/",
];

const dom = {
  pipelineGrid: document.getElementById("pipeline-grid"),
  serviceGrid: document.getElementById("service-grid"),
  coverageCount: document.getElementById("coverage-count"),
  runCombo: document.getElementById("run-combo"),
  runComboToggle: document.getElementById("run-combo-toggle"),
  runComboValue: document.getElementById("run-combo-value"),
  runComboPanel: document.getElementById("run-combo-panel"),
  runSearch: document.getElementById("run-search"),
  runOptions: document.getElementById("run-options"),
  refreshRuns: document.getElementById("refresh-runs"),
  langEn: document.getElementById("lang-en"),
  langVi: document.getElementById("lang-vi"),
  statusBanner: document.getElementById("status-banner"),
  overviewRunCount: document.getElementById("overview-run-count"),
  overviewRunList: document.getElementById("overview-run-list"),
  overviewPassRate: document.getElementById("overview-pass-rate"),
  overviewCaseSummary: document.getElementById("overview-case-summary"),
  overviewAdmissionSummary: document.getElementById("overview-admission-summary"),
  overviewPolicyNote: document.getElementById("overview-policy-note"),
  overviewArtifactCount: document.getElementById("overview-artifact-count"),
  overviewArtifactNote: document.getElementById("overview-artifact-note"),
  matrixGrid: document.getElementById("matrix-grid"),
  runMeta: document.getElementById("run-meta"),
  artifactContext: document.getElementById("artifact-context"),
  artifactList: document.getElementById("artifact-list"),
  summaryPreview: document.getElementById("summary-preview"),
  sbomDonut: document.getElementById("sbom-donut"),
  sbomLegend: document.getElementById("sbom-legend"),
  sbomTopList: document.getElementById("sbom-top-list"),
  sbomTotal: document.getElementById("sbom-total"),
  cveSummary: document.getElementById("cve-summary"),
  cveTableBody: document.getElementById("cve-table-body"),
  cveSnapshotNote: document.getElementById("cve-snapshot-note"),
  cveTableWrap: document.getElementById("cve-table-wrap"),
  statsStrip: document.getElementById("stats-strip"),
  alertRibbon: document.getElementById("alert-ribbon"),
  statusSystem: document.getElementById("status-system"),
  alertCveSection: document.getElementById("alert-cve"),
  cveAlertsGrid: document.getElementById("cve-alerts-grid"),
  cveAlertsSummary: document.getElementById("cve-alerts-summary"),
  goRuntimeSection: document.getElementById("go-runtime"),
  goRuntimeBody: document.getElementById("go-runtime-body"),
  ciTimelineSection: document.getElementById("ci-timeline"),
  ciTimelineStrip: document.getElementById("ci-timeline-strip"),
  slsaTrackerSection: document.getElementById("slsa-tracker"),
  slsaTrackerBody: document.getElementById("slsa-tracker-body"),
  runnerPoolSection: document.getElementById("runner-pool"),
  runnerPoolGrid: document.getElementById("runner-pool-grid"),
  runnerAbSection: document.getElementById("runner-ab"),
  runnerAbBody: document.getElementById("runner-ab-body"),
  serviceInventorySection: document.getElementById("service-inventory"),
  serviceInventoryBody: document.getElementById("service-inventory-body"),
  admissionUnaffectedNote: document.getElementById("admission-unaffected-note"),
};

const state = {
  comboOpen: false,
  dataMode: "legacy",
  activeEvidenceBasePath: "",
  availableRunIds: [],
  filteredRunIds: [],
  runMetaById: {},
  snapshotRunById: {},
  snapshotData: null,
  selectedRunId: "",
  selectedCaseName: "",
  language: "en",
  lastLoadStatus: "",
  lastRegressionVerdict: "",
  latestRunData: {
    mode: "legacy",
    runId: "",
    basePath: "",
    runHtmlUrl: "",
    casesByName: {},
    summaryText: "",
  },
};

function escapeHtml(value) {
  return String(value == null ? "" : value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function t(key, vars = {}) {
  const table = I18N[state.language] || I18N.en;
  let text = table[key] || I18N.en[key] || key;
  for (const [name, value] of Object.entries(vars)) {
    text = text.replaceAll(`{${name}}`, String(value));
  }
  return text;
}

function getCaseHint(caseName) {
  if (state.language === "vi") {
    return CASE_HINTS_VI[caseName] || CASE_HINTS[caseName] || caseName;
  }
  return CASE_HINTS[caseName] || caseName;
}

function normalizeAdmissionDecision(value) {
  const raw = String(value || "").trim().toLowerCase();
  if (!raw) return "unknown";
  if (raw.includes("allow")) return "allow";
  if (raw.includes("deny") || raw.includes("denied")) return "deny";
  return "unknown";
}

function localizeAdmissionValue(value) {
  const normalized = normalizeAdmissionDecision(value);
  if (state.language !== "vi") {
    return value;
  }
  if (normalized === "allow") return "Cho phep";
  if (normalized === "deny") return "Tu choi";
  return value;
}

function localizeReason(reason) {
  const text = String(reason || "");
  if (state.language !== "vi") {
    return text;
  }

  if (text === "Deployment became Available.") {
    return "Deployment da dat trang thai Available.";
  }
  if (text === "Admission deny evidence detected at apply phase.") {
    return "Phat hien bang chung bi tu choi o buoc apply.";
  }
  if (text === "Admission deny evidence detected in events/ReplicaSet describe.") {
    return "Phat hien bang chung bi tu choi trong events/describe ReplicaSet.";
  }

  return text;
}

function setStatus(message, type) {
  dom.statusBanner.textContent = message;
  dom.statusBanner.dataset.type = type || "info";
}

function isLegacyRunId(value) {
  return /^[0-9]{8}-[0-9]{6}$/.test(value);
}

function isSnapshotRunId(value) {
  return /^gha:[A-Za-z0-9._-]+#\d+$/.test(String(value || ""));
}

function isValidRunId(value) {
  return isLegacyRunId(value) || isSnapshotRunId(value);
}

function sortRunsDesc(runIds) {
  return [...runIds].sort((a, b) => b.localeCompare(a));
}

function formatRunTimestamp(value) {
  const raw = String(value || "").trim();
  if (!raw) {
    return "-";
  }
  const parsed = new Date(raw);
  if (Number.isNaN(parsed.getTime())) {
    return raw;
  }
  return parsed.toISOString().replace("T", " ").replace(".000Z", "Z");
}

function setComboDisplayValue(value) {
  dom.runComboValue.textContent = value || t("selectRun");
}

function updateUrlParams() {
  const url = new URL(window.location.href);
  if (state.selectedRunId) {
    url.searchParams.set("run", state.selectedRunId);
  } else {
    url.searchParams.delete("run");
  }
  url.searchParams.set("lang", state.language);
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

function applyStaticTranslations() {
  document.title = t("pageTitle");
  document.documentElement.lang = state.language;

  for (const node of document.querySelectorAll("[data-i18n]")) {
    const key = node.getAttribute("data-i18n");
    if (key) {
      node.textContent = t(key);
    }
  }

  for (const node of document.querySelectorAll("[data-i18n-placeholder]")) {
    const key = node.getAttribute("data-i18n-placeholder");
    if (key) {
      node.setAttribute("placeholder", t(key));
    }
  }

  dom.runOptions.setAttribute("aria-label", t("runLabel"));
  dom.langEn.classList.toggle("is-active", state.language === "en");
  dom.langVi.classList.toggle("is-active", state.language === "vi");
}

function renderRunOptionList() {
  dom.runOptions.innerHTML = "";
  dom.runOptions.setAttribute("aria-label", t("runLabel"));

  if (state.filteredRunIds.length === 0) {
    const li = document.createElement("li");
    li.className = "combo-empty";
    li.textContent = t("noRunMatched");
    dom.runOptions.appendChild(li);
    return;
  }

  let currentGroup = "";
  for (const runId of state.filteredRunIds) {
    const meta = state.runMetaById[runId] || {};
    const workflowLabel = String(meta.workflowLabel || "");
    if (workflowLabel && workflowLabel !== currentGroup) {
      const groupLi = document.createElement("li");
      groupLi.className = "combo-group";
      groupLi.textContent = workflowLabel;
      dom.runOptions.appendChild(groupLi);
      currentGroup = workflowLabel;
    }

    const li = document.createElement("li");
    const button = document.createElement("button");
    button.type = "button";
    button.className = "combo-option";
    button.setAttribute("role", "option");
    button.setAttribute("aria-selected", String(runId === state.selectedRunId));
    button.dataset.runId = runId;
    button.textContent = String(meta.label || runId);
    if (meta.hint) {
      button.title = String(meta.hint);
    }
    li.appendChild(button);
    dom.runOptions.appendChild(li);
  }
}

function filterRunOptions(keyword) {
  const term = String(keyword || "").trim().toLowerCase();
  if (!term) {
    state.filteredRunIds = [...state.availableRunIds];
  } else {
    state.filteredRunIds = state.availableRunIds.filter((runId) => {
      const meta = state.runMetaById[runId] || {};
      const haystack = `${runId} ${meta.label || ""} ${meta.search || ""}`.toLowerCase();
      return haystack.includes(term);
    });
  }
  renderRunOptionList();
}

function renderRunOptions(runIds, runMetaById) {
  state.availableRunIds = [...runIds];
  state.filteredRunIds = [...runIds];
  state.runMetaById = runMetaById || {};
  dom.runSearch.value = "";
  dom.overviewRunCount.textContent = String(runIds.length);
  dom.overviewRunList.textContent = runIds.length > 0
    ? t("latestRuns", {
      runs: runIds
        .slice(0, 4)
        .map((runId) => String((state.runMetaById[runId] || {}).label || runId))
        .join(", "),
    })
    : t("noRunsDetected");

  if (runIds.length === 0) {
    dom.runComboToggle.disabled = true;
    setComboDisplayValue(t("noRunsFound"));
    renderRunOptionList();
    return;
  }

  dom.runComboToggle.disabled = false;
  if (!state.selectedRunId || !runIds.includes(state.selectedRunId)) {
    state.selectedRunId = "";
    setComboDisplayValue(t("selectRun"));
  } else {
    setComboDisplayValue(String((state.runMetaById[state.selectedRunId] || {}).label || state.selectedRunId));
  }
  renderRunOptionList();
}

function relativeTime(iso) {
  const raw = String(iso || "").trim();
  if (!raw) return "-";
  const then = Date.parse(raw);
  if (Number.isNaN(then)) return raw;
  const diffMs = Date.now() - then;
  const mins = Math.floor(diffMs / 60000);
  if (mins < 1) return t("relTimeNow");
  if (mins < 60) return t("relTimeMinsAgo", { n: mins });
  const hours = Math.floor(mins / 60);
  if (hours < 48) return t("relTimeHoursAgo", { n: hours });
  const days = Math.floor(hours / 24);
  return t("relTimeDaysAgo", { n: days });
}

function buildSparklineSvg(runs) {
  const items = (Array.isArray(runs) ? runs : []).slice(0, 10).reverse();
  if (items.length === 0) return "";
  const w = items.length * 12;
  let dots = "";
  items.forEach((run, i) => {
    const cx = 6 + i * 12;
    const concl = String(run.conclusion || run.status || "").toLowerCase();
    const fill = concl === "success" ? "var(--ok, #22c55e)" : concl === "failure" ? "var(--deny, #ef4444)" : "#64748b";
    const title = `#${run.run_number || "-"} ${concl}`;
    dots += `<circle cx="${cx}" cy="7" r="4" fill="${fill}"><title>${title}</title></circle>`;
  });
  return `<svg viewBox="0 0 ${w} 14" width="${w}" height="14" xmlns="http://www.w3.org/2000/svg">${dots}</svg>`;
}

function computeRedStreak(runs) {
  let count = 0;
  let firstFail = null;
  for (const run of runs) {
    if (String(run.conclusion || "").toLowerCase() === "failure") {
      count++;
      firstFail = run;
    } else {
      break;
    }
  }
  return { count, firstFail };
}

function renderPipelineStatus(snapshot) {
  if (!dom.pipelineGrid) return;
  const workflows = Array.isArray(snapshot.workflows) ? snapshot.workflows : [];
  if (workflows.length === 0) {
    dom.pipelineGrid.innerHTML = "";
    return;
  }
  dom.pipelineGrid.innerHTML = workflows.map((wf) => {
    const key = String(wf.workflow_key || "");
    const runs = Array.isArray(wf.runs) ? wf.runs : [];
    const last = runs[0] || null;
    const conclusion = last ? String(last.conclusion || last.status || "unknown") : "no data";
    const runNum = last ? `#${last.run_number}` : "-";
    const rel = last ? relativeTime(last.created_at) : "-";
    const url = last ? String(last.html_url || "") : "";
    const statusAttr = conclusion === "success" ? "success" : conclusion === "failure" ? "failure" : "neutral";
    const pillCls = conclusion === "success" ? "pass" : conclusion === "failure" ? "fail" : "neutral";
    const isRed = conclusion === "failure";
    const nameEl = url
      ? `<a href="${url}" target="_blank" rel="noopener noreferrer" class="pipeline-name">${key}</a>`
      : `<span class="pipeline-name">${key}</span>`;
    const desc = String(wf.description || WORKFLOW_META[key] || "");
    const rollup = wf.rollup || {};
    const passRate = (typeof rollup.pass_rate_pct === "number") ? `${rollup.pass_rate_pct.toFixed(1)}% pass` : "-";
    const spark = buildSparklineSvg(runs);
    let streakLine = "";
    if (key === "ci-service" && isRed) {
      const { count, firstFail } = computeRedStreak(runs);
      if (count > 0 && firstFail) {
        const firstFailDate = formatRunTimestamp(firstFail.created_at).slice(0, 10);
        const firstFailKey = `#${firstFail.run_number} (${firstFailDate})`;
        streakLine = `<div class="pipeline-streak">${t("ciRedStreak", { count, firstFail: firstFailKey })}</div>`;
      }
    }
    return `<div class="pipeline-card${isRed ? " is-red" : ""}" data-status="${statusAttr}">${nameEl}<div class="pipeline-desc">${desc}</div><div class="pipeline-sparkline">${spark}</div><div class="pipeline-meta"><span>${runNum} · ${rel}</span><span class="pill ${pillCls}">${conclusion}</span></div><div class="pipeline-rollup">${passRate}</div>${streakLine}</div>`;
  }).join("");
}

function renderServiceCoverage(snapshot) {
  if (!dom.serviceGrid) return;
  const services = Array.isArray(snapshot.services) ? snapshot.services : null;

  if (services && services.length > 0) {
    let okCount = 0;
    // Render flagship first then others
    const ordered = [...services].sort((a, b) => {
      const af = a.category === "flagship" ? 0 : 1;
      const bf = b.category === "flagship" ? 0 : 1;
      if (af !== bf) return af - bf;
      return (b.loc || 0) - (a.loc || 0);
    });
    dom.serviceGrid.innerHTML = ordered.map((svc) => {
      const has = svc.covered !== false;
      if (has) okCount++;
      const ci = has ? "ok" : "missing";
      const flagshipCls = svc.category === "flagship" ? " svc-chip-flagship" : "";
      const slsa = String(svc.slsa_level || "");
      const slsaCls = slsa === "L3" ? "slsa-pill-l3" : slsa === "L2" ? "slsa-pill-l2" : "slsa-pill-l1";
      const loc = Number(svc.loc || 0);
      const locTxt = loc >= 1000 ? `${(loc / 1000).toFixed(1)}k` : `${loc}`;
      const grypeDot = `<span class="grype-dot${(svc.last_grype_high === 0) ? " ok" : " bad"}" title="last grype high"></span>`;
      const flagshipLabel = svc.category === "flagship" ? `<span class="svc-flagship-label">flagship · SLSA ${slsa}</span>` : "";
      return `<span class="svc-chip${flagshipCls}" data-ci="${ci}"><span class="svc-dot"></span>${svc.name}<span class="loc-pill">${locTxt}</span><span class="slsa-pill ${slsaCls}">${slsa}</span>${grypeDot}${flagshipLabel}</span>`;
    }).join("");
    if (dom.coverageCount) dom.coverageCount.textContent = `(${okCount}/${services.length})`;
    return;
  }

  // Fallback (legacy v1): derive from ci-service artifacts
  const workflows = Array.isArray(snapshot.workflows) ? snapshot.workflows : [];
  const ciWf = workflows.find((wf) => wf.workflow_key === "ci-service");
  const runs = ciWf && Array.isArray(ciWf.runs) ? ciWf.runs : [];
  const seen = new Set();
  for (const run of runs.slice(0, 5)) {
    for (const name of Object.keys(run.artifacts || {})) seen.add(name);
  }
  let okCount = 0;
  dom.serviceGrid.innerHTML = ALL_SERVICES_FALLBACK.map((svc) => {
    const has = seen.has(`${svc}-grype-report`) || seen.has(`${svc}-security-gate-findings`) || seen.has(`${svc}-supply-chain-artifacts`);
    if (has) okCount++;
    const ci = has ? "ok" : "missing";
    return `<span class="svc-chip" data-ci="${ci}"><span class="svc-dot"></span>${svc}</span>`;
  }).join("");
  if (dom.coverageCount) dom.coverageCount.textContent = `(${okCount}/${ALL_SERVICES_FALLBACK.length})`;
}

function renderStatsStrip(meta, workflows, slsa) {
  if (!dom.statsStrip) return;
  meta = meta || {};
  workflows = Array.isArray(workflows) ? workflows : [];
  slsa = slsa || {};

  const serviceCount = Number(meta.service_count || 0);
  const goPinned = String(meta.go_version_pinned || "-");
  const goReq = String(meta.go_version_required_fix || "");
  const kyverno = String(meta.kyverno_version || "-");
  const wfCount = workflows.length || Number(meta.workflow_count || 0);
  const redWf = workflows.filter((wf) => {
    const r = (wf.rollup || {});
    return String(r.latest_conclusion || "").toLowerCase() === "failure";
  }).length;

  const goNeedsBump = goReq && goReq !== goPinned;
  const slsaLevel = String(meta.slsa_level_user_service ?? slsa?.level_per_service?.["user-service"]?.level ?? "L2");
  const chips = [
    `<span class="stat-chip">${serviceCount || 23} services</span>`,
    `<span class="stat-chip${goNeedsBump ? " is-red" : ""}">Go ${goPinned}${goNeedsBump ? ` → ${goReq}` : ""}</span>`,
    `<span class="stat-chip">Kyverno ${kyverno}</span>`,
    `<span class="stat-chip">Cosign keyless</span>`,
    `<span class="stat-chip is-gold">SLSA ${slsaLevel} user-service · L2 fleet</span>`,
    `<span class="stat-chip">Grype + govulncheck</span>`,
    `<span class="stat-chip${redWf > 0 ? " is-red" : ""}">${wfCount || 6} workflows${redWf > 0 ? ` · ${redWf} red` : ""}</span>`,
  ];
  dom.statsStrip.innerHTML = chips.join("");
}

function renderAlertRibbon(meta) {
  if (!dom.alertRibbon) return;
  meta = meta || {};
  const state = String(meta.system_state || "").toUpperCase();
  if (state === "DEGRADED") {
    const days = Number(meta.consecutive_failures || 0);
    const goReq = String(meta.go_version_required_fix || "");
    const txt = `SYSTEM DEGRADED — CI red ${days} ${days === 1 ? "day" : "days"} · Go stdlib CVE blocking${goReq ? ` · pending Go ${goReq} bump` : ""}`;
    dom.alertRibbon.textContent = txt;
    dom.alertRibbon.classList.remove("hidden");
    dom.alertRibbon.dataset.severity = "high";
  } else {
    dom.alertRibbon.classList.add("hidden");
    dom.alertRibbon.textContent = "";
  }
}

function renderSystemStatusLine(meta, workflows) {
  if (!dom.statusSystem) return;
  meta = meta || {};
  workflows = Array.isArray(workflows) ? workflows : [];
  const state = String(meta.system_state || "UNKNOWN").toUpperCase();
  const parts = [`System: ${state}`];
  for (const wf of workflows) {
    const r = wf.rollup || {};
    const concl = String(r.latest_conclusion || "").toLowerCase();
    if (!concl) continue;
    const tag = concl === "success" ? "GREEN" : concl === "failure" ? "RED" : concl.toUpperCase();
    parts.push(`${wf.workflow_key} ${tag}`);
  }
  dom.statusSystem.textContent = parts.join(" · ");
  dom.statusSystem.dataset.state = state.toLowerCase();
}

function renderCveAlerts(cveAlerts, snapshot) {
  if (!dom.alertCveSection || !dom.cveAlertsGrid) return;
  const advisories = (cveAlerts && Array.isArray(cveAlerts.advisories)) ? cveAlerts.advisories : [];
  if (!cveAlerts || !cveAlerts.active || advisories.length === 0) {
    dom.alertCveSection.classList.add("hidden");
    return;
  }
  dom.alertCveSection.classList.remove("hidden");
  if (dom.cveAlertsSummary) {
    dom.cveAlertsSummary.textContent = String(cveAlerts.summary || t("activeCveAlertsNote"));
  }

  // Build a map of run_key -> html_url so first_seen_run can link out
  const runUrlByKey = {};
  for (const wf of (snapshot && snapshot.workflows) || []) {
    for (const run of (wf.runs || [])) {
      const key = String(run.run_key || `${wf.workflow_key}#${run.run_number}`);
      runUrlByKey[key] = String(run.html_url || "");
    }
  }

  const now = Date.now();
  dom.cveAlertsGrid.innerHTML = advisories.map((adv) => {
    const firstAt = Date.parse(String(adv.first_seen_at || "")) || now;
    const daysSince = Math.max(0, Math.floor((now - firstAt) / 86400000));
    const runKey = String(adv.first_seen_run || "");
    const runUrl = runUrlByKey[runKey] || "";
    const runEl = runUrl
      ? `<a href="${escapeHtml(runUrl)}" target="_blank" rel="noopener noreferrer">${escapeHtml(runKey)}</a>`
      : `<span>${escapeHtml(runKey)}</span>`;
    const fixedTrim = (adv.fixed_in || "").replace(/^go/, "");
    const prSnippet = `# .github/workflows/ci-service.yml\nenv:\n  GO_VERSION: "${escapeHtml(fixedTrim)}"`;
    return `
      <article class="cve-alert-card" data-severity="${escapeHtml(String(adv.severity || "").toLowerCase())}">
        <header class="cve-alert-head">
          <div>
            <div class="cve-alert-id">${escapeHtml(adv.id)}</div>
            <div class="cve-alert-pkg">${escapeHtml(adv.package)}</div>
          </div>
          <span class="pill fail">${escapeHtml(adv.severity)}</span>
        </header>
        <p class="cve-alert-desc">${escapeHtml(adv.description)}</p>
        <dl class="cve-alert-meta">
          <dt>Fixed in</dt><dd>${escapeHtml(adv.fixed_in)}</dd>
          <dt>First seen</dt><dd>${runEl} · ${daysSince}d ago</dd>
          <dt>Affected</dt><dd>${escapeHtml(adv.affected_services)}</dd>
          <dt>Detected by</dt><dd>${escapeHtml(adv.detected_by)}</dd>
        </dl>
        <p class="cve-alert-fix"><strong>Fix:</strong> ${escapeHtml(adv.remediation)}</p>
        <pre class="cve-alert-snippet">${prSnippet}</pre>
      </article>
    `;
  }).join("");
}

function renderGoRuntimeStatus(goStatus, meta) {
  if (!dom.goRuntimeSection || !dom.goRuntimeBody) return;
  if (!goStatus) {
    dom.goRuntimeSection.classList.add("hidden");
    return;
  }
  dom.goRuntimeSection.classList.remove("hidden");
  const pinned = String(goStatus.pinned_version || "-");
  const required = String(goStatus.required_minimum || "-");
  const status = String(goStatus.status || "").toUpperCase();
  const days = Number(goStatus.ci_red_days || (meta && meta.consecutive_failures) || 0);
  const branch = String(goStatus.fix_plan_branch || "");
  const pinnedIn = String(goStatus.pinned_in || "");
  const next = String(goStatus.next_action || "");
  dom.goRuntimeBody.innerHTML = `
    <div class="go-runtime-row">
      <div class="go-runtime-pill go-runtime-pinned">
        <span class="go-runtime-label">Pinned</span>
        <span class="go-runtime-value">${pinned}</span>
      </div>
      <span class="go-runtime-arrow">→</span>
      <div class="go-runtime-pill go-runtime-required">
        <span class="go-runtime-label">Required</span>
        <span class="go-runtime-value">${required}</span>
      </div>
      <div class="go-runtime-pill go-runtime-status${status === "OUTDATED" ? " is-red" : ""}">
        <span class="go-runtime-label">Status</span>
        <span class="go-runtime-value">${status} · ${days}d</span>
      </div>
    </div>
    <dl class="go-runtime-detail">
      <dt>Pinned in</dt><dd><code>${pinnedIn}</code></dd>
      <dt>Fix branch</dt><dd><code>${branch || "-"}</code></dd>
      <dt>Next action</dt><dd>${next || "-"}</dd>
    </dl>
  `;
}

function renderCiTimeline(workflows) {
  if (!dom.ciTimelineSection || !dom.ciTimelineStrip) return;
  workflows = Array.isArray(workflows) ? workflows : [];
  const ci = workflows.find((wf) => wf.workflow_key === "ci-service");
  if (!ci) {
    dom.ciTimelineSection.classList.add("hidden");
    return;
  }
  const runs = Array.isArray(ci.runs) ? ci.runs.slice(0, 10) : [];
  if (runs.length === 0) {
    dom.ciTimelineSection.classList.add("hidden");
    return;
  }
  dom.ciTimelineSection.classList.remove("hidden");
  const ordered = [...runs].reverse();
  dom.ciTimelineStrip.innerHTML = ordered.map((run) => {
    const concl = String(run.conclusion || "").toLowerCase();
    const cls = concl === "success" ? "ok" : concl === "failure" ? "fail" : "neutral";
    const date = formatRunTimestamp(run.created_at).slice(0, 10);
    const url = String(run.html_url || "#");
    return `<a class="ci-tl-dot ci-tl-${cls}" href="${url}" target="_blank" rel="noopener noreferrer" title="#${run.run_number} ${concl} · ${date}">
      <span class="ci-tl-bullet"></span>
      <span class="ci-tl-label">#${run.run_number}</span>
      <span class="ci-tl-date">${date.slice(5)}</span>
    </a>`;
  }).join("");
}

function renderSlsaTracker(slsa, meta) {
  if (!dom.slsaTrackerSection || !dom.slsaTrackerBody) return;
  if (!slsa || !slsa.level_per_service) {
    dom.slsaTrackerSection.classList.add("hidden");
    return;
  }
  const flagshipKey = String((meta && meta.flagship_service) || "user-service");
  const flagship = slsa.level_per_service[flagshipKey];
  if (!flagship) {
    dom.slsaTrackerSection.classList.add("hidden");
    return;
  }
  dom.slsaTrackerSection.classList.remove("hidden");
  const verifier = String(slsa.verifier || "-");
  const level = String(flagship.level || "-");
  const builder = String(flagship.builder_id || "-");
  const builderShort = builder.length > 96 ? builder.slice(0, 93) + "..." : builder;
  const digest = String(flagship.digest_sha256 || "-");
  const verifiedAt = formatRunTimestamp(flagship.verified_at);
  const verifiedRun = String(flagship.verified_in_run || "-");
  const note = String(slsa.scaffold_note || "");
  dom.slsaTrackerBody.innerHTML = `
    <div class="slsa-card">
      <header class="slsa-head">
        <span class="slsa-flagship">${flagshipKey}</span>
        <span class="pill pass">SLSA ${level}</span>
      </header>
      <dl class="slsa-detail">
        <dt>Verifier</dt><dd><code>${verifier}</code></dd>
        <dt>Builder</dt><dd><code title="${builder}">${builderShort}</code></dd>
        <dt>Digest</dt><dd><code>${digest}</code></dd>
        <dt>Verified at</dt><dd>${verifiedAt}</dd>
        <dt>Verified in</dt><dd>${verifiedRun}</dd>
      </dl>
    </div>
    <p class="slsa-footnote">${note}</p>
  `;
}

function renderRunnerPool(pool) {
  if (!dom.runnerPoolSection || !dom.runnerPoolGrid) return;
  const runners = (pool && Array.isArray(pool.runners)) ? pool.runners : [];
  if (runners.length === 0) {
    dom.runnerPoolSection.classList.add("hidden");
    return;
  }
  dom.runnerPoolSection.classList.remove("hidden");
  dom.runnerPoolGrid.innerHTML = runners.map((r) => {
    const type = String(r.type || "github-hosted");
    const typeCls = type === "self-hosted" ? "is-self" : "is-hosted";
    const used = Array.isArray(r.used_by) ? r.used_by : [];
    const usedHtml = used.map((wf) => `<span class="runner-used-pill" data-wf="${wf}">${wf}</span>`).join("");
    return `
      <article class="runner-card ${typeCls}">
        <header class="runner-head">
          <span class="runner-label">${r.label}</span>
          <span class="pill ${type === "self-hosted" ? "neutral" : "pass"}">${type}</span>
        </header>
        <p class="runner-os">${r.os}</p>
        <p class="runner-use">${r.primary_use || ""}</p>
        <div class="runner-used-by">${usedHtml}</div>
      </article>
    `;
  }).join("");
}

function renderRunnerAb(workflows) {
  if (!dom.runnerAbSection || !dom.runnerAbBody) return;
  workflows = Array.isArray(workflows) ? workflows : [];
  const wf = workflows.find((w) => w.workflow_key === "runner-ab-benchmark");
  const runs = wf && Array.isArray(wf.runs) ? wf.runs : [];
  const latest = runs[0];
  if (!latest || !latest.benchmark_summary) {
    dom.runnerAbSection.classList.add("hidden");
    return;
  }
  dom.runnerAbSection.classList.remove("hidden");
  const b = latest.benchmark_summary || {};
  const url = String(latest.html_url || "#");
  dom.runnerAbBody.innerHTML = `
    <div class="ab-row">
      <div class="ab-card">
        <span class="ab-label">windows-latest</span>
        <span class="ab-value">${b.windows_latest_min}s min</span>
      </div>
      <span class="ab-vs">vs</span>
      <div class="ab-card ab-winner">
        <span class="ab-label">THEMONSTER-win-parity</span>
        <span class="ab-value">${b.themonster_min}s min</span>
      </div>
      <div class="ab-card ab-speedup">
        <span class="ab-label">Speedup</span>
        <span class="ab-value">${b.speedup_pct}%</span>
      </div>
    </div>
    <p class="ab-footnote">Source: <a href="${url}" target="_blank" rel="noopener noreferrer">runner-ab-benchmark#${latest.run_number}</a></p>
  `;
}

function renderServiceInventory(services) {
  if (!dom.serviceInventorySection || !dom.serviceInventoryBody) return;
  services = Array.isArray(services) ? services : [];
  if (services.length === 0) {
    dom.serviceInventorySection.classList.add("hidden");
    return;
  }
  dom.serviceInventorySection.classList.remove("hidden");
  const ordered = [...services].sort((a, b) => {
    const af = a.category === "flagship" ? 0 : 1;
    const bf = b.category === "flagship" ? 0 : 1;
    if (af !== bf) return af - bf;
    return (b.loc || 0) - (a.loc || 0);
  });
  dom.serviceInventoryBody.innerHTML = ordered.map((s) => {
    const slsa = String(s.slsa_level || "L1");
    const slsaCls = slsa === "L3" ? "slsa-pill-l3" : slsa === "L2" ? "slsa-pill-l2" : "slsa-pill-l1";
    const sbom = formatRunTimestamp(s.last_sbom_at).slice(0, 10);
    const grypeOk = Number(s.last_grype_high || 0) === 0;
    const grypeText = grypeOk ? "0 (ok)" : `${s.last_grype_high}`;
    const covered = s.covered === false ? "no" : "yes";
    return `<tr>
      <td>${s.name}${s.category === "flagship" ? " <span class=\"flagship-tag\">flagship</span>" : ""}</td>
      <td>${s.category}</td>
      <td>${s.loc}</td>
      <td><span class="slsa-pill ${slsaCls}">${slsa}</span></td>
      <td>${sbom}</td>
      <td class="${grypeOk ? "grype-ok" : "grype-bad"}">${grypeText}</td>
      <td>${covered}</td>
    </tr>`;
  }).join("");
}

function renderAdmissionUnaffectedNote(meta) {
  if (!dom.admissionUnaffectedNote) return;
  const state = String((meta && meta.system_state) || "").toUpperCase();
  if (state === "DEGRADED") {
    dom.admissionUnaffectedNote.classList.remove("hidden");
  } else {
    dom.admissionUnaffectedNote.classList.add("hidden");
  }
}

function renderDashboardMetaPanels(snapshot) {
  const meta = snapshot && snapshot.dashboard_meta;
  const workflows = (snapshot && Array.isArray(snapshot.workflows)) ? snapshot.workflows : [];
  const slsa = snapshot && snapshot.slsa_attestations;
  renderStatsStrip(meta, workflows, slsa);
  renderAlertRibbon(meta);
  renderSystemStatusLine(meta, workflows);
  renderCveAlerts(snapshot && snapshot.cve_alerts, snapshot);
  renderGoRuntimeStatus(snapshot && snapshot.go_runtime_status, meta);
  renderCiTimeline(workflows);
  renderSlsaTracker(slsa, meta);
  renderRunnerPool(snapshot && snapshot.runner_pool);
  renderRunnerAb(workflows);
  renderServiceInventory(snapshot && snapshot.services);
  renderAdmissionUnaffectedNote(meta);
  // Update hero copy with as_of_date
  if (meta && meta.as_of_date) {
    const heroNode = document.getElementById("hero-copy");
    if (heroNode) {
      heroNode.textContent = `Automated supply chain security evidence: SBOM, CVE scanning, Kyverno admission control, SLSA provenance. Powered by 6 GitHub Actions workflows across 23 Go microservices · Snapshot of ${meta.as_of_date}.`;
    }
  }
}

function renderOverview(casesByName) {
  const caseItems = CASE_ORDER.map((caseName) => casesByName[caseName]).filter(Boolean);
  const totalCases = caseItems.length;
  const passCount = caseItems.filter((item) => String(item.verdict || "").toUpperCase() === "PASS").length;
  const allowCount = caseItems.filter((item) => normalizeAdmissionDecision(item.actual) === "allow").length;
  const denyCount = caseItems.filter((item) => normalizeAdmissionDecision(item.actual) === "deny").length;

  // Card 1 (repurposed): CI Health from snapshot
  const snapshot = state.snapshotData || {};
  const workflows = Array.isArray(snapshot.workflows) ? snapshot.workflows : [];
  const ci = workflows.find((wf) => wf.workflow_key === "ci-service");
  if (ci) {
    const rollup = ci.rollup || {};
    const latestConcl = String(rollup.latest_conclusion || (ci.runs && ci.runs[0] && ci.runs[0].conclusion) || "-");
    const passRatePct = (typeof rollup.pass_rate_pct === "number") ? rollup.pass_rate_pct.toFixed(0) : "-";
    const { count: redStreak } = computeRedStreak(ci.runs || []);
    const lbl = document.querySelector("[data-i18n='overviewAvailableRuns'], [data-i18n='overviewCiHealth']");
    if (lbl) {
      lbl.setAttribute("data-i18n", "overviewCiHealth");
      lbl.textContent = t("overviewCiHealth");
    }
    dom.overviewRunCount.textContent = latestConcl.toUpperCase();
    dom.overviewRunList.textContent = redStreak > 0
      ? `${t("redStreakBadge", { count: redStreak })} · ${passRatePct}% pass last ${rollup.total_in_window || 10}`
      : `${passRatePct}% pass last ${rollup.total_in_window || 10}`;
  } else {
    dom.overviewRunCount.textContent = "-";
    dom.overviewRunList.textContent = t("noRunsDetected");
  }

  if (totalCases === 0) {
    dom.overviewPassRate.textContent = "-";
    dom.overviewCaseSummary.textContent = t("noRunLoaded");
    dom.overviewAdmissionSummary.textContent = "-";
    dom.overviewPolicyNote.textContent = t("noPolicyDecision");
  } else {
    const passRate = ((passCount / totalCases) * 100).toFixed(0);
    dom.overviewPassRate.textContent = t("passRate", { rate: passRate });
    dom.overviewCaseSummary.textContent = t("caseSummary", { passCount, totalCases });
    dom.overviewAdmissionSummary.textContent = t("admissionSummary", { allowCount, denyCount });
    dom.overviewPolicyNote.textContent = t("policyDerived");
  }

  // Card 4 (repurposed): Last Green Build from dashboard_meta
  const meta = snapshot.dashboard_meta || {};
  if (meta.last_pass_run) {
    const lbl = document.querySelector("[data-i18n='overviewEvidenceFootprint'], [data-i18n='overviewLastGreen']");
    if (lbl) {
      lbl.setAttribute("data-i18n", "overviewLastGreen");
      lbl.textContent = t("overviewLastGreen");
    }
    dom.overviewArtifactCount.textContent = `#${String(meta.last_pass_run).split("#")[1] || meta.last_pass_run}`;
    const rel = relativeTime(meta.last_pass_at);
    const goV = meta.go_version_pinned ? ` · Go ${meta.go_version_pinned}` : "";
    dom.overviewArtifactNote.textContent = `${rel}${goV}`;
  } else {
    dom.overviewArtifactCount.textContent = "-";
    dom.overviewArtifactNote.textContent = t("waitingArtifacts");
  }
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
  const labels = t("metaLabels");
  const dashMeta = (state.snapshotData && state.snapshotData.dashboard_meta) || {};
  const slsa = (state.snapshotData && state.snapshotData.slsa_attestations) || {};
  const flagship = (slsa.level_per_service && slsa.level_per_service[dashMeta.flagship_service || "user-service"]) || {};
  const goPinned = String(dashMeta.go_version_pinned || "");
  const goReq = String(dashMeta.go_version_required_fix || "");
  const goNeedsBump = goReq && goPinned && goReq !== goPinned;
  const goLine = goPinned ? `${goPinned}${goNeedsBump ? ` (need ${goReq})` : ""}` : "-";
  const builder = String(flagship.builder_id || "");
  const builderShort = builder ? (builder.length > 80 ? builder.slice(0, 77) + "..." : builder) : "-";
  const verifier = String(slsa.verifier || "-");

  const entries = [
    [labels[0], meta.runId, null],
    [labels[1], meta.context, null],
    [labels[2], meta.namespace, null],
    [labels[3], meta.signedDigest, null],
    [labels[4], meta.unsignedDigest, null],
    [labels[5], meta.sbomDigest, null],
    ["Go Version", goLine, goNeedsBump ? "fail" : null],
    ["Builder", builderShort, null],
    ["Verifier", verifier, null],
  ];

  dom.runMeta.innerHTML = "";
  for (const [label, value, badge] of entries) {
    const wrapper = document.createElement("div");
    const dt = document.createElement("dt");
    const dd = document.createElement("dd");
    dt.textContent = label;
    dd.textContent = value || "-";
    if (badge) {
      dd.classList.add(`meta-${badge}`);
    }
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

function formatVerdictLabel(verdict) {
  if (verdict === "PASS") return t("verdictMatched");
  if (verdict === "FAIL") return t("verdictMismatch");
  return verdict;
}

function signalClass(value) {
  const normalized = normalizeAdmissionDecision(value);
  if (normalized === "allow") return "signal-allowed";
  if (normalized === "deny") return "signal-denied";
  return "signal-neutral";
}

function matrixMissingCard(caseName) {
  return `
    <article class="matrix-card" data-missing="true" data-verdict="MISSING">
      <div class="matrix-head">
        <div>
          <div class="case-id">${caseName}</div>
          <div class="matrix-detail">${getCaseHint(caseName)}</div>
        </div>
        <span class="pill neutral">${t("missing")}</span>
      </div>
      <div class="matrix-detail">
        ${t("missingCaseText")}
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
    const actualRaw = String(item.actual || "UNKNOWN");
    const expectedRaw = String(item.expected || "UNKNOWN");
    const reasonRaw = String(item.reason || "-");
    const actual = localizeAdmissionValue(actualRaw);
    const expected = localizeAdmissionValue(expectedRaw);
    const reason = localizeReason(reasonRaw);
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
            <div class="matrix-detail">${getCaseHint(caseName)}</div>
          </div>
          <span class="pill ${verdictClass(verdict)}">${formatVerdictLabel(verdict)}</span>
        </div>
        <div class="matrix-detail">
          <strong>${t("expected")}:</strong> ${expected}<br />
          <strong>${t("actual")}:</strong> <span class="${signalClass(actual)}">${actual}</span><br />
          <strong>${t("reason")}:</strong> ${reason}
          <dl>
            <dt>include_sbom</dt><dd>${includeSbom}</dd>
            <dt>high_critical</dt><dd>${highCritical}</dd>
            <dt>apply_exit_code</dt><dd>${applyExit}</dd>
            <dt>wait_exit_code</dt><dd>${waitExit}</dd>
            <dt>image</dt><dd>${image}</dd>
          </dl>
        </div>
        <div class="matrix-actions">
          <span class="artifact-count">${t("artifacts")}: ${artifactCount}</span>
          <button type="button" data-action="show-artifacts" data-case="${caseName}">
            ${t("viewArtifacts")}
          </button>
        </div>
      </article>
    `;
  }

  dom.matrixGrid.innerHTML = html;
  for (const button of dom.matrixGrid.querySelectorAll("button[data-action='show-artifacts']")) {
    button.addEventListener("click", () => {
      const caseName = button.getAttribute("data-case") || "";
      if (state.latestRunData.mode === "snapshot") {
        renderSnapshotArtifacts(caseName);
      } else {
        renderArtifacts(basePath, caseName);
      }
    });
  }
}

function renderArtifacts(basePath, caseName) {
  state.selectedCaseName = caseName;
  const item = state.latestRunData.casesByName[caseName];
  dom.artifactList.innerHTML = "";

  if (!item || !item.artifacts || typeof item.artifacts !== "object") {
    dom.artifactContext.textContent = t("noArtifactMap", { caseName });
    return;
  }

  dom.artifactContext.textContent = t("showingArtifacts", { caseName });
  const keys = Object.keys(item.artifacts);
  if (keys.length === 0) {
    const li = document.createElement("li");
    li.textContent = t("noArtifactFiles");
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
    appendArtifactChips(li, relativePath);
    dom.artifactList.appendChild(li);
  }
}

function appendArtifactChips(li, name) {
  const n = String(name || "").toLowerCase();
  const chips = [];
  if (n.endsWith(".sig")) chips.push(["signed", "signed"]);
  if (n.endsWith(".att") || n.includes(".intoto.")) chips.push(["attested", "attested"]);
  if (n.endsWith(".spdx.json") || n.includes(".spdx.")) chips.push(["sbom", "sbom"]);
  if (n.includes(".grype.") || n.endsWith(".grype.json")) chips.push(["grype", "grype"]);
  for (const [label, cls] of chips) {
    const tag = document.createElement("span");
    tag.className = `artifact-chip artifact-chip-${cls}`;
    tag.textContent = label;
    li.appendChild(tag);
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
  dom.sbomTotal.textContent = t("totalPackages", { total });

  if (total === 0) {
    dom.sbomLegend.innerHTML = `<p>${t("noPackagesFound")}</p>`;
    dom.sbomTopList.innerHTML = `<li>${t("noDependencies")}</li>`;
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
  dom.sbomLegend.innerHTML = `<p>${t("sbomNotFound")}</p>`;
  dom.sbomTopList.innerHTML = `<li>${t("sbomMissing")}</li>`;
  dom.sbomTotal.textContent = t("totalPackages", { total: "-" });
  return "";
}

function sortGateFindings(rows) {
  const rank = { critical: 2, high: 1 };
  const isStdlib = (item) => {
    const cve = String((item && item.cve) || "");
    const pkg = String((item && item.package) || "");
    return /^GO-\d{4}-/.test(cve) || /^(net\/|crypto\/|encoding\/|stdlib)/.test(pkg);
  };
  return [...rows].sort((a, b) => {
    const sa = isStdlib(a) ? 1 : 0;
    const sb = isStdlib(b) ? 1 : 0;
    if (sb !== sa) return sb - sa;
    const ra = rank[String(a.severity || "").toLowerCase()] || 0;
    const rb = rank[String(b.severity || "").toLowerCase()] || 0;
    if (rb !== ra) return rb - ra;
    return String(a.cve || "").localeCompare(String(b.cve || ""));
  });
}

function renderGateFindings(rows, source) {
  const ordered = sortGateFindings(rows);
  dom.cveSummary.textContent = t("cveSourceSummary", { source, count: ordered.length });

  if (!ordered.length) {
    dom.cveTableBody.innerHTML = `<tr><td colspan="5">${t("cveNoFindings")}</td></tr>`;
    return;
  }

  dom.cveTableBody.innerHTML = ordered.map((item) => {
    const fixed = Array.isArray(item.fixed_versions) ? item.fixed_versions.join(", ") : String(item.fixed_versions || "-");
    const cve = String(item.cve || "-");
    const pkg = String(item.package || "-");
    const isStdlib = /^GO-\d{4}-/.test(cve) || /^(net\/|crypto\/|encoding\/|stdlib)/.test(pkg);
    const tag = isStdlib ? `<span class="cve-tag-stdlib">stdlib · blocking</span>` : "";
    return `<tr${isStdlib ? ' class="cve-row-stdlib"' : ""}>
      <td>${cve} ${tag}</td>
      <td>${String(item.severity || "-")}</td>
      <td>${pkg}</td>
      <td>${String(item.installed || "-")}</td>
      <td>${fixed || "-"}</td>
    </tr>`;
  }).join("");
}

function toFixableHighCriticalFromGrype(grypeJson) {
  const matches = Array.isArray(grypeJson.matches) ? grypeJson.matches : [];
  return matches.map((m) => {
    const vulnerability = m && m.vulnerability ? m.vulnerability : {};
    const artifact = m && m.artifact ? m.artifact : {};
    const severity = String(vulnerability.severity || "").toLowerCase();
    const fixState = String((vulnerability.fix && vulnerability.fix.state) || "unknown").toLowerCase();
    const fixedVersions = vulnerability.fix && Array.isArray(vulnerability.fix.versions)
      ? vulnerability.fix.versions
      : [];
    return {
      cve: String(vulnerability.id || ""),
      severity,
      package: String(artifact.name || ""),
      installed: String(artifact.version || ""),
      fix_state: fixState,
      fixed_versions: fixedVersions,
    };
  }).filter((item) =>
    (item.severity === "high" || item.severity === "critical")
    && item.fix_state !== "wont-fix"
    && item.fix_state !== "not-fixed"
    && item.fix_state !== "unknown"
  );
}

async function loadGateFindings() {
  for (const source of GATE_FINDINGS_SOURCES) {
    try {
      const data = await fetchJson(source);
      if (Array.isArray(data)) {
        renderGateFindings(data, source);
        return;
      }
      if (Array.isArray(data.findings)) {
        renderGateFindings(data.findings, source);
        return;
      }
      if (Array.isArray(data.matches)) {
        renderGateFindings(toFixableHighCriticalFromGrype(data), `${source} (derived)`);
        return;
      }
    } catch (error) {
      // Continue fallback chain.
    }
  }

  dom.cveSummary.textContent = t("cveSourceUnavailable");
  dom.cveTableBody.innerHTML = "<tr><td colspan=\"5\">-</td></tr>";
}

function toSnapshotRunKey(run) {
  return `gha:${String(run.workflow_key || "").trim()}#${String(run.run_number || "").trim()}`;
}

function toSnapshotRunMeta(run) {
  const workflowName = String(run.workflow_name || run.workflow_key || "workflow");
  const workflowKey = String(run.workflow_key || "workflow");
  const runNumber = String(run.run_number || "-");
  const conclusion = String(run.conclusion || run.status || "unknown");
  const createdAt = formatRunTimestamp(run.created_at);
  return {
    workflowLabel: workflowName,
    label: `${workflowName} #${runNumber} | ${conclusion} | ${createdAt}`,
    hint: `${workflowKey} #${runNumber} | ${conclusion} | run_id=${String(run.run_id || "-")}`,
    search: `${workflowName} ${workflowKey} #${runNumber} ${conclusion} ${createdAt} ${String(run.run_id || "")}`,
  };
}

function parseRunMetaFromSnapshot(run) {
  const matrixMeta = (run.matrix && run.matrix.metadata) || {};
  return {
    runId: `${String(run.workflow_name || run.workflow_key || "-")} #${String(run.run_number || "-")} (run_id: ${String(run.run_id || "-")})`,
    context: String(matrixMeta.context || run.workflow_key || "-"),
    namespace: String(matrixMeta.namespace || "-"),
    signedDigest: String(matrixMeta.signed_image_digest || "-"),
    unsignedDigest: String(matrixMeta.unsigned_image_digest || "-"),
    sbomDigest: String(matrixMeta.sbom_digest || "-"),
  };
}

function normalizeCasesFromSnapshot(run) {
  const matrix = run.matrix || {};
  const cases = Array.isArray(matrix.cases) ? matrix.cases : [];
  return byCaseName(cases);
}

function renderSnapshotArtifacts(caseName) {
  state.selectedCaseName = caseName;
  const item = state.latestRunData.casesByName[caseName];
  dom.artifactList.innerHTML = "";

  if (!item || !item.artifacts || typeof item.artifacts !== "object") {
    dom.artifactContext.textContent = t("noArtifactMap", { caseName });
    return;
  }

  dom.artifactContext.textContent = t("showingArtifacts", { caseName });
  const keys = Object.keys(item.artifacts);
  if (keys.length === 0) {
    const li = document.createElement("li");
    li.textContent = t("noArtifactFiles");
    dom.artifactList.appendChild(li);
    return;
  }

  const runHtmlUrl = state.latestRunData.runHtmlUrl || "#";
  for (const key of keys) {
    const relativePath = String(item.artifacts[key]);
    const li = document.createElement("li");
    const link = document.createElement("a");
    link.href = runHtmlUrl;
    link.target = "_blank";
    link.rel = "noopener noreferrer";
    link.textContent = `${key}: ${relativePath}`;
    li.appendChild(link);
    appendArtifactChips(li, relativePath);
    dom.artifactList.appendChild(li);
  }
}

function renderCveFromSnapshotRun(run) {
  const securityGate = run.security_gate || {};
  const findings = Array.isArray(securityGate.findings) ? securityGate.findings : [];
  const source = `snapshot/${String(run.workflow_key || "workflow")}#${String(run.run_number || "-")}`;
  if (findings.length > 0) {
    setCveSnapshotV2Mode(false);
    renderGateFindings(findings, source);
    return true;
  }
  const schemaVersion = Number((state.snapshotData && state.snapshotData.schema_version) || 0);
  if (schemaVersion === 2) {
    setCveSnapshotV2Mode(true);
    return true;
  }
  setCveSnapshotV2Mode(false);
  return false;
}

function setCveSnapshotV2Mode(active) {
  if (!dom.cveSnapshotNote || !dom.cveTableWrap) return;
  if (active) {
    dom.cveSnapshotNote.innerHTML = `${escapeHtml(t("cveSnapshotV2Note"))} <a href="#alert-cve">${escapeHtml(t("cveSnapshotV2JumpLink"))}</a>`;
    dom.cveSnapshotNote.classList.remove("hidden");
    dom.cveTableWrap.classList.add("hidden");
    if (dom.cveSummary) {
      dom.cveSummary.textContent = t("cveSnapshotV2Note");
    }
  } else {
    dom.cveSnapshotNote.classList.add("hidden");
    dom.cveSnapshotNote.innerHTML = "";
    dom.cveTableWrap.classList.remove("hidden");
  }
}

async function loadSnapshotRun(runId) {
  const run = state.snapshotRunById[runId];
  if (!run) {
    setStatus(t("statusSnapshotRunMissing", { runId }), "error");
    return;
  }

  const displayRunId = String(run.run_key || runId);
  setStatus(t("statusSnapshotLoading", { runId: displayRunId }), "info");

  const casesByName = normalizeCasesFromSnapshot(run);
  const summaryText = String((run.matrix && run.matrix.summary_text) || "");
  const regression = run.matrix ? run.matrix.regression : null;

  state.latestRunData = {
    mode: "snapshot",
    runId,
    basePath: "",
    runHtmlUrl: String(run.html_url || ""),
    casesByName,
    summaryText,
  };

  renderOverview(casesByName);
  renderRunMeta(parseRunMetaFromSnapshot(run));
  renderMatrix("", casesByName);
  dom.summaryPreview.textContent = summaryText.trim() || t("summaryEmpty");

  const firstSelectable = state.selectedCaseName && casesByName[state.selectedCaseName]
    ? state.selectedCaseName
    : CASE_ORDER.find((caseName) => Boolean(casesByName[caseName]));

  if (firstSelectable) {
    renderSnapshotArtifacts(firstSelectable);
  } else {
    state.selectedCaseName = "";
    dom.artifactContext.textContent = t("noKnownCases");
    dom.artifactList.innerHTML = `<li>${t("noArtifactsAvailable")}</li>`;
  }

  const renderedFromSnapshot = renderCveFromSnapshotRun(run);
  if (!renderedFromSnapshot) {
    // Fallback for runs where snapshot couldn't ingest artifact findings.
    void loadGateFindings();
  }

  if (regression && regression.verdict) {
    state.lastLoadStatus = "loaded";
    state.lastRegressionVerdict = String(regression.verdict);
    setStatus(t("statusSnapshotLoaded", { runId: displayRunId, result: regression.verdict }), "success");
  } else {
    state.lastLoadStatus = "regressionMissing";
    state.lastRegressionVerdict = "";
    if (run.evidence_unavailable) {
      setStatus(t("statusSnapshotPartial", { runId: displayRunId }), "warning");
    } else {
      setStatus(t("statusSnapshotRegressionMissing", { runId: displayRunId }), "warning");
    }
  }
}

async function tryLoadActionsSnapshot() {
  const snapshot = await fetchJson(ACTIONS_SNAPSHOT_PATH);
  const runs = Array.isArray(snapshot.runs) ? snapshot.runs : [];
  const workflows = Array.isArray(snapshot.workflows) ? snapshot.workflows : [];
  if (!runs.length && !workflows.length) {
    throw new Error("Snapshot has no runs");
  }

  const runIds = [];
  const runMetaById = {};
  const runById = {};
  const seenRunIds = new Set();

  const pushRun = (rawRun, workflowFallback) => {
    if (!rawRun || typeof rawRun !== "object") {
      return;
    }
    const run = {
      ...rawRun,
      workflow_key: rawRun.workflow_key || workflowFallback.workflow_key || "workflow",
      workflow_name: rawRun.workflow_name || workflowFallback.workflow_name || workflowFallback.workflow_key || "workflow",
    };
    const runId = toSnapshotRunKey(run);
    if (!isSnapshotRunId(runId) || seenRunIds.has(runId)) {
      return;
    }
    seenRunIds.add(runId);
    runIds.push(runId);
    runMetaById[runId] = toSnapshotRunMeta(run);
    runById[runId] = run;
  };

  for (const workflow of workflows) {
    const workflowRuns = Array.isArray(workflow.runs) ? [...workflow.runs] : [];
    workflowRuns.sort((a, b) => Number(b.run_number || 0) - Number(a.run_number || 0));
    for (const run of workflowRuns) {
      pushRun(run, workflow);
    }
  }

  if (!runIds.length) {
    const orderedRuns = [...runs].sort((a, b) => {
      const ta = Date.parse(String(a.created_at || "")) || 0;
      const tb = Date.parse(String(b.created_at || "")) || 0;
      if (tb !== ta) {
        return tb - ta;
      }
      return Number(b.run_number || 0) - Number(a.run_number || 0);
    });
    for (const run of orderedRuns) {
      pushRun(run, {});
    }
  }

  if (!runIds.length) {
    throw new Error("Snapshot has no valid run identifiers");
  }

  return {
    snapshot,
    runIds,
    runMetaById,
    runById,
  };
}

async function loadRun(runId) {
  if (state.dataMode === "snapshot" && isSnapshotRunId(runId)) {
    await loadSnapshotRun(runId);
    return;
  }

  if (!isValidRunId(runId)) {
    setStatus(t("statusInvalidRun"), "error");
    return;
  }

  if (!state.activeEvidenceBasePath) {
    setStatus(t("statusNoSource"), "error");
    return;
  }

  const basePath = `${state.activeEvidenceBasePath}${runId}`;
  setStatus(t("statusLoadingEvidence", { basePath }), "info");

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
      mode: "legacy",
      runId,
      basePath,
      runHtmlUrl: "",
      casesByName,
      summaryText,
    };

    renderOverview(casesByName);
    renderRunMeta(parseRunMetadata(summaryText));
    renderMatrix(basePath, casesByName);
    dom.summaryPreview.textContent = summaryText.trim() || t("summaryEmpty");

    const firstSelectable = state.selectedCaseName && casesByName[state.selectedCaseName]
      ? state.selectedCaseName
      : CASE_ORDER.find((caseName) => Boolean(casesByName[caseName]));

    if (firstSelectable) {
      renderArtifacts(basePath, firstSelectable);
    } else {
      state.selectedCaseName = "";
      dom.artifactContext.textContent = t("noKnownCases");
      dom.artifactList.innerHTML = `<li>${t("noArtifactsAvailable")}</li>`;
    }

    if (regression && regression.verdict) {
      state.lastLoadStatus = "loaded";
      state.lastRegressionVerdict = String(regression.verdict);
      setStatus(t("statusRunLoaded", { runId, result: regression.verdict }), "success");
    } else {
      state.lastLoadStatus = "regressionMissing";
      state.lastRegressionVerdict = "";
      setStatus(t("statusRegressionMissing", { runId }), "warning");
    }
  } catch (error) {
    state.lastLoadStatus = "loadFailed";
    state.lastRegressionVerdict = "";
    dom.matrixGrid.innerHTML = "";
    dom.artifactList.innerHTML = "";
    dom.summaryPreview.textContent = t("summaryNotLoaded");
    setStatus(t("statusLoadRunFailed", { runId }), "error");
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
  setComboDisplayValue(String((state.runMetaById[runId] || {}).label || runId));
  renderRunOptionList();

  if (opts.updateQuery) {
    updateUrlParams();
  }

  if (opts.close) {
    closeCombo();
  }

  if (opts.load) {
    await loadRun(runId);
  }
}

async function discoverAndLoadRuns(preferredRunId) {
  setStatus(t("statusScanningRuns"), "info");

  try {
    const snapshotBundle = await tryLoadActionsSnapshot();
    state.dataMode = "snapshot";
    state.activeEvidenceBasePath = "";
    state.snapshotData = snapshotBundle.snapshot;
    state.snapshotRunById = snapshotBundle.runById;
    renderPipelineStatus(snapshotBundle.snapshot);
    renderServiceCoverage(snapshotBundle.snapshot);
    renderDashboardMetaPanels(snapshotBundle.snapshot);
    renderRunOptions(snapshotBundle.runIds, snapshotBundle.runMetaById);

    let selectedRun = snapshotBundle.runIds[0];
    if (preferredRunId && snapshotBundle.runIds.includes(preferredRunId)) {
      selectedRun = preferredRunId;
    } else if (preferredRunId && !snapshotBundle.runIds.includes(preferredRunId)) {
      setStatus(t("statusSnapshotFallback", { preferredRunId, selectedRun }), "warning");
    }

    await chooseRun(selectedRun, { load: true, close: true, updateQuery: true });
    return;
  } catch (error) {
    // Snapshot unavailable; continue with legacy evidence directories.
  }

  state.dataMode = "legacy";
  state.snapshotData = null;
  state.snapshotRunById = {};

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
      if (preferredRunId && isLegacyRunId(preferredRunId) && runIds.includes(preferredRunId)) {
        selectedRun = preferredRunId;
      } else if (preferredRunId && isLegacyRunId(preferredRunId) && !runIds.includes(preferredRunId)) {
        setStatus(t("statusRunFallback", { preferredRunId, selectedRun }), "warning");
      }

      await chooseRun(selectedRun, { load: true, close: true, updateQuery: true });
      return;
    } catch (error) {
      // Try next source.
    }
  }

  state.activeEvidenceBasePath = "";
  renderRunOptions([]);
  renderOverview({});
  setStatus(t("statusNoRuns"), "error");
}

function rerenderUiForLanguage() {
  applyStaticTranslations();
  renderRunOptions(state.availableRunIds, state.runMetaById);
  if (state.snapshotData) {
    renderDashboardMetaPanels(state.snapshotData);
  }
  renderOverview(state.latestRunData.casesByName || {});

  if (state.dataMode === "snapshot" && isSnapshotRunId(state.selectedRunId) && state.snapshotRunById[state.selectedRunId]) {
    const run = state.snapshotRunById[state.selectedRunId];
    renderRunMeta(parseRunMetaFromSnapshot(run));
    renderMatrix("", state.latestRunData.casesByName || {});
    if (state.selectedCaseName) {
      renderSnapshotArtifacts(state.selectedCaseName);
    }
    dom.summaryPreview.textContent = state.latestRunData.summaryText.trim() || t("summaryEmpty");
    renderCveFromSnapshotRun(run);
  } else if (state.latestRunData.summaryText) {
    renderRunMeta(parseRunMetadata(state.latestRunData.summaryText));
    renderMatrix(state.latestRunData.basePath, state.latestRunData.casesByName || {});
    if (state.selectedCaseName) {
      renderArtifacts(state.latestRunData.basePath, state.selectedCaseName);
    }
    dom.summaryPreview.textContent = state.latestRunData.summaryText.trim() || t("summaryEmpty");
    void loadGateFindings();
  } else {
    dom.artifactContext.textContent = t("chooseCase");
    dom.summaryPreview.textContent = t("summaryNotLoaded");
    if (state.dataMode !== "snapshot") {
      void loadGateFindings();
    }
  }

  if (state.latestRunData.runId && state.lastLoadStatus === "loaded") {
    setStatus(
      t("statusRunLoaded", { runId: state.latestRunData.runId, result: state.lastRegressionVerdict || "N/A" }),
      "success"
    );
  } else if (state.latestRunData.runId && state.lastLoadStatus === "regressionMissing") {
    setStatus(t("statusRegressionMissing", { runId: state.latestRunData.runId }), "warning");
  } else if (state.latestRunData.runId && state.lastLoadStatus === "loadFailed") {
    setStatus(t("statusLoadRunFailed", { runId: state.latestRunData.runId }), "error");
  }

  void loadSbom();
  updateUrlParams();
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

function initLanguage() {
  const params = new URLSearchParams(window.location.search);
  const queryLang = String(params.get("lang") || "").trim().toLowerCase();
  const storedLang = String(localStorage.getItem("security-dashboard-lang") || "").trim().toLowerCase();
  state.language = queryLang === "vi" || (queryLang !== "en" && storedLang === "vi") ? "vi" : "en";
  localStorage.setItem("security-dashboard-lang", state.language);
  applyStaticTranslations();

  dom.langEn.addEventListener("click", () => {
    state.language = "en";
    localStorage.setItem("security-dashboard-lang", state.language);
    rerenderUiForLanguage();
  });

  dom.langVi.addEventListener("click", () => {
    state.language = "vi";
    localStorage.setItem("security-dashboard-lang", state.language);
    rerenderUiForLanguage();
  });
}

function init() {
  initLanguage();
  initComboEvents();

  dom.refreshRuns.addEventListener("click", () => {
    void discoverAndLoadRuns(state.selectedRunId);
  });

  const params = new URLSearchParams(window.location.search);
  const runFromQuery = String(params.get("run") || "").trim();

  dom.artifactContext.textContent = t("chooseCase");
  dom.summaryPreview.textContent = t("summaryNotLoaded");
  dom.overviewRunList.textContent = t("noRunsDetected");

  void discoverAndLoadRuns(runFromQuery);
  void loadSbom();
  void loadGateFindings();
}

init();


