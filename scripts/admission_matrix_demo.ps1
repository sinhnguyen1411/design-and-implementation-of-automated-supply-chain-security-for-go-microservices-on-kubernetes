param(
  [string]$Context = "docker-desktop",
  [string]$Namespace = "stock-trading",
  [string]$ExportDir = ".demo/evidence",
  [string]$ImageTtl = "24h",
  [string]$CosignPassword = "local-demo-pass",
  [string]$DemoDir = ".demo",
  [string]$KyvernoVersion = "v1.12.5",
  [switch]$ResetNamespace,
  [switch]$SkipGoTest
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest
if ($PSVersionTable.PSVersion.Major -ge 7) {
  $PSNativeCommandUseErrorActionPreference = $false
}

function Require-Cli([string]$name) {
  if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
    throw "Required CLI not found: $name"
  }
}

function Write-Section([string]$message) {
  Write-Host ""
  Write-Host "== $message =="
}

function Ensure-Directory([string]$path) {
  New-Item -ItemType Directory -Force -Path $path | Out-Null
}

function Write-Text([string]$path, [string[]]$lines) {
  $lines | Set-Content -Path $path -Encoding UTF8
}

function Write-TextNoBom([string]$path, [string]$content) {
  $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
  [System.IO.File]::WriteAllText($path, $content, $utf8NoBom)
}

function New-DemoImageRef([string]$namePrefix, [string]$ttl) {
  $suffix = [guid]::NewGuid().ToString("N").Substring(0, 12)
  return "ttl.sh/$namePrefix-${suffix}:$ttl"
}

function Invoke-Kubectl([string[]]$KubectlArgs, [string]$OutFile, [switch]$AllowFailure) {
  $previousErrorAction = $ErrorActionPreference
  $result = $null
  $exitCode = 0
  try {
    $ErrorActionPreference = "Continue"
    $result = & kubectl @KubectlArgs 2>&1
    $exitCode = $LASTEXITCODE
  } finally {
    $ErrorActionPreference = $previousErrorAction
  }
  Write-Text -path $OutFile -lines @($result)
  if (-not $AllowFailure -and $exitCode -ne 0) {
    throw "kubectl failed: kubectl $($KubectlArgs -join ' ')"
  }
  return @{
    ExitCode = $exitCode
    Output = ($result -join [Environment]::NewLine)
  }
}

function Build-DeploymentYaml(
  [string]$NamespaceName,
  [string]$Image,
  [string]$SbomDigest,
  [bool]$IncludeSbom,
  [string]$HighCritical,
  [string]$CaseName
) {
  $sbomLine = ""
  if ($IncludeSbom) {
    $sbomLine = "        security.stock-trading.dev/sbom-digest: `"$SbomDigest`""
  }

  return @"
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: $NamespaceName
  labels:
    app.kubernetes.io/name: user-service
    app.kubernetes.io/component: backend
    app.kubernetes.io/part-of: stock-trading
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: user-service
  template:
    metadata:
      labels:
        app.kubernetes.io/name: user-service
        app.kubernetes.io/component: backend
        app.kubernetes.io/part-of: stock-trading
      annotations:
        kubectl.kubernetes.io/default-container: user-service
        security.grype.io/high_critical: "$HighCritical"
$sbomLine
        security.stock-trading.dev/matrix-case: "$CaseName"
    spec:
      serviceAccountName: user-service
      automountServiceAccountToken: false
      securityContext:
        fsGroup: 65532
        runAsNonRoot: true
        seccompProfile:
          type: RuntimeDefault
      containers:
        - name: user-service
          image: $Image
          imagePullPolicy: Always
          args: ["server", "--config", "/etc/stock-trading/config/config.yaml"]
          env:
            - name: AUTH__ACCESS_TOKEN_SECRET
              valueFrom:
                secretKeyRef:
                  name: user-service-secrets
                  key: auth-access-token-secret
            - name: AUTH__REFRESH_TOKEN_SECRET
              valueFrom:
                secretKeyRef:
                  name: user-service-secrets
                  key: auth-refresh-token-secret
          ports:
            - name: http
              containerPort: 18080
            - name: grpc
              containerPort: 19090
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          securityContext:
            allowPrivilegeEscalation: false
            runAsNonRoot: true
            runAsUser: 65532
            runAsGroup: 65532
            readOnlyRootFilesystem: true
            capabilities:
              drop: ["ALL"]
          readinessProbe:
            tcpSocket:
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
          livenessProbe:
            tcpSocket:
              port: http
            initialDelaySeconds: 15
            periodSeconds: 15
          startupProbe:
            tcpSocket:
              port: http
            failureThreshold: 18
            periodSeconds: 5
          volumeMounts:
            - name: config
              mountPath: /etc/stock-trading/config
              readOnly: true
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: config
          configMap:
            name: user-service-config
            items:
              - key: config.yaml
                path: config.yaml
        - name: tmp
          emptyDir: {}
"@
}

function Get-KyvernoControllerName() {
  $deploys = & kubectl -n kyverno get deploy -o name 2>$null
  if ($LASTEXITCODE -ne 0 -or -not $deploys) {
    return $null
  }

  $preferred = $deploys | Where-Object { $_ -match "kyverno-admission-controller" } | Select-Object -First 1
  if ($preferred) {
    return $preferred
  }

  return ($deploys | Select-Object -First 1)
}

function Collect-CaseEvidence([string]$CaseDir, [string]$NamespaceName) {
  Invoke-Kubectl -KubectlArgs @("-n", $NamespaceName, "get", "events", "--sort-by=.lastTimestamp") -OutFile (Join-Path $CaseDir "events.txt") | Out-Null
  Invoke-Kubectl -KubectlArgs @("-n", $NamespaceName, "get", "deploy,rs,pods", "-o", "wide") -OutFile (Join-Path $CaseDir "workloads.txt") | Out-Null
  Invoke-Kubectl -KubectlArgs @("-n", $NamespaceName, "describe", "deployment", "user-service") -OutFile (Join-Path $CaseDir "describe-deployment.txt") -AllowFailure | Out-Null

  $rsNames = (& kubectl -n $NamespaceName get rs -l app.kubernetes.io/name=user-service -o name 2>$null)
  if ($LASTEXITCODE -eq 0 -and $rsNames) {
    $rsDescribe = @()
    foreach ($rs in $rsNames) {
      $rsDescribe += "===== $rs ====="
      $rsDescribe += (& kubectl -n $NamespaceName describe $rs 2>&1)
      $rsDescribe += ""
    }
    Write-Text -path (Join-Path $CaseDir "describe-replicasets.txt") -lines $rsDescribe
  } else {
    Write-Text -path (Join-Path $CaseDir "describe-replicasets.txt") -lines @("No ReplicaSet found.")
  }

  $podNames = (& kubectl -n $NamespaceName get pods -l app.kubernetes.io/name=user-service -o name 2>$null)
  if ($LASTEXITCODE -eq 0 -and $podNames) {
    $podDescribe = @()
    foreach ($pod in $podNames) {
      $podDescribe += "===== $pod ====="
      $podDescribe += (& kubectl -n $NamespaceName describe $pod 2>&1)
      $podDescribe += ""
    }
    Write-Text -path (Join-Path $CaseDir "describe-pods.txt") -lines $podDescribe
  } else {
    Write-Text -path (Join-Path $CaseDir "describe-pods.txt") -lines @("No Pod found.")
  }

  $kyvernoDeploy = Get-KyvernoControllerName
  if ($kyvernoDeploy) {
    Invoke-Kubectl -KubectlArgs @("-n", "kyverno", "logs", $kyvernoDeploy, "-c", "kyverno", "--tail=400") -OutFile (Join-Path $CaseDir "kyverno-logs.txt") -AllowFailure | Out-Null
  } else {
    Write-Text -path (Join-Path $CaseDir "kyverno-logs.txt") -lines @("Kyverno deployment not found.")
  }
}

function Get-CaseVerdict([string]$Expected, [int]$ApplyExitCode, [int]$WaitExitCode, [string]$CaseDir) {
  $applyOutput = Get-Content -Path (Join-Path $CaseDir "kubectl-apply.txt") -Raw
  $events = Get-Content -Path (Join-Path $CaseDir "events.txt") -Raw
  $describeDeploy = Get-Content -Path (Join-Path $CaseDir "describe-deployment.txt") -Raw
  $describeRs = Get-Content -Path (Join-Path $CaseDir "describe-replicasets.txt") -Raw
  $combined = "$applyOutput`n$events`n$describeDeploy`n$describeRs"
  $denyRegex = "(?i)(denied the request|rejected|failedcreate|forbidden|policy violation)"
  $hasDenyEvidence = [regex]::IsMatch($combined, $denyRegex)

  if ($Expected -eq "Allowed") {
    if ($ApplyExitCode -ne 0) {
      return @{
        Actual = "DeniedAtApply"
        Verdict = "FAIL"
        Reason = "Expected Allowed but apply was denied."
      }
    }
    if ($WaitExitCode -eq 0) {
      return @{
        Actual = "Allowed"
        Verdict = "PASS"
        Reason = "Deployment became Available."
      }
    }
    return @{
      Actual = "DeniedOrUnavailable"
      Verdict = "FAIL"
      Reason = "Expected Allowed but deployment did not become Available."
    }
  }

  if ($ApplyExitCode -ne 0 -and $hasDenyEvidence) {
    return @{
      Actual = "Denied"
      Verdict = "PASS"
      Reason = "Admission deny evidence detected at apply phase."
    }
  }

  if ($WaitExitCode -ne 0 -and $hasDenyEvidence) {
    return @{
      Actual = "Denied"
      Verdict = "PASS"
      Reason = "Admission deny evidence detected in events/ReplicaSet describe."
    }
  }

  return @{
    Actual = "UnknownOrAllowed"
    Verdict = "FAIL"
    Reason = "Expected Denied but no strong deny signal was found."
  }
}

Require-Cli docker
Require-Cli kubectl
Require-Cli cosign
Require-Cli syft
Require-Cli go

$contextNames = & kubectl config get-contexts -o name
if ($LASTEXITCODE -ne 0) {
  throw "Unable to list kubectl contexts."
}
if (-not ($contextNames -contains $Context)) {
  throw "Kubernetes context '$Context' not found."
}

Write-Section "Use Kubernetes Context"
$currentContext = & kubectl config current-context
if ($currentContext -ne $Context) {
  & kubectl config use-context $Context | Out-Null
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to switch kubectl context to '$Context'."
  }
}
Write-Host "Active context: $Context"

Write-Section "Pre-check Cluster"
& kubectl get nodes -o wide
if ($LASTEXITCODE -ne 0) {
  throw "Cluster is not reachable on context '$Context'."
}

if (-not $SkipGoTest) {
  Write-Section "Pre-check Go Tests"
  & go test ./...
  if ($LASTEXITCODE -ne 0) {
    throw "go test failed."
  }
}

$runId = Get-Date -Format "yyyyMMdd-HHmmss"
$runDir = Join-Path $ExportDir $runId
Ensure-Directory $runDir
Ensure-Directory $DemoDir

$cosignKeyPath = Join-Path $DemoDir "cosign.key"
$cosignPubPath = Join-Path $DemoDir "cosign.pub"
$sbomPath = Join-Path $DemoDir "sbom.spdx.json"
$provenancePath = Join-Path $DemoDir "provenance.json"
$baseResourcesPath = Join-Path $runDir "base-resources.yaml"
$localPolicyPath = Join-Path $runDir "verify-local-matrix-image.yaml"

Write-Section "Build Local Image Once"
$localImage = "local-user-service:matrix"
& docker build -t $localImage .
if ($LASTEXITCODE -ne 0) {
  throw "docker build failed."
}

Write-Section "Prepare Signed And Unsigned Images"
$signedImage = New-DemoImageRef -namePrefix "stock-trading-matrix-signed" -ttl $ImageTtl
$unsignedImage = New-DemoImageRef -namePrefix "stock-trading-matrix-unsigned" -ttl $ImageTtl

& docker tag $localImage $signedImage
if ($LASTEXITCODE -ne 0) { throw "docker tag signed image failed." }
& docker push $signedImage
if ($LASTEXITCODE -ne 0) { throw "docker push signed image failed." }

& docker tag $localImage $unsignedImage
if ($LASTEXITCODE -ne 0) { throw "docker tag unsigned image failed." }
& docker push $unsignedImage
if ($LASTEXITCODE -ne 0) { throw "docker push unsigned image failed." }

$signedRepo = $signedImage.Split(":")[0]
$unsignedRepo = $unsignedImage.Split(":")[0]

$signedRepoDigests = & docker inspect --format='{{range .RepoDigests}}{{println .}}{{end}}' $signedImage
if ($LASTEXITCODE -ne 0) { throw "docker inspect signed image failed." }
$signedDigest = (($signedRepoDigests -split "`r?`n") | Where-Object { $_ -like "$signedRepo@*" } | Select-Object -First 1).Trim()
if (-not $signedDigest) { throw "Unable to resolve signed image digest." }

$unsignedRepoDigests = & docker inspect --format='{{range .RepoDigests}}{{println .}}{{end}}' $unsignedImage
if ($LASTEXITCODE -ne 0) { throw "docker inspect unsigned image failed." }
$unsignedDigest = (($unsignedRepoDigests -split "`r?`n") | Where-Object { $_ -like "$unsignedRepo@*" } | Select-Object -First 1).Trim()
if (-not $unsignedDigest) { throw "Unable to resolve unsigned image digest." }

Write-Host "Signed digest  : $signedDigest"
Write-Host "Unsigned digest: $unsignedDigest"

Write-Section "Generate SBOM + Sign + Attest (Signed Image Only)"
$env:COSIGN_PASSWORD = $CosignPassword
if (!(Test-Path $cosignKeyPath) -or !(Test-Path $cosignPubPath)) {
  & cosign generate-key-pair --output-key-prefix ($cosignKeyPath -replace '\.key$', '')
  if ($LASTEXITCODE -ne 0) { throw "cosign key generation failed." }
}

& syft $signedDigest -o "spdx-json=$sbomPath"
if ($LASTEXITCODE -ne 0) { throw "syft SBOM generation failed." }
$sbomDigest = (Get-FileHash $sbomPath -Algorithm SHA256).Hash

$provenanceJson = @"
{
  "buildType": "https://slsa.dev/provenance/v0.2",
  "builder": { "id": "local-admission-matrix" },
  "invocation": {
    "configSource": {
      "uri": "https://github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes",
      "digest": { "gitCommit": "local-matrix-demo" }
    }
  }
}
"@
Write-TextNoBom -path $provenancePath -content $provenanceJson

& cosign sign --yes --key $cosignKeyPath $signedDigest
if ($LASTEXITCODE -ne 0) { throw "cosign sign failed." }
& cosign attest --yes --key $cosignKeyPath --predicate $provenancePath --type slsaprovenance $signedDigest
if ($LASTEXITCODE -ne 0) { throw "cosign attest failed." }
& cosign verify --key $cosignPubPath $signedDigest | Out-Null
if ($LASTEXITCODE -ne 0) { throw "cosign verify failed." }
& cosign verify-attestation --key $cosignPubPath --type slsaprovenance $signedDigest | Out-Null
if ($LASTEXITCODE -ne 0) { throw "cosign verify-attestation failed." }

Write-Section "Ensure Kyverno"
$hasClusterPolicyApi = & kubectl api-resources | Select-String -Pattern "^clusterpolicies\s"
if (-not $hasClusterPolicyApi) {
  & kubectl apply --server-side -f "https://github.com/kyverno/kyverno/releases/download/$KyvernoVersion/install.yaml"
  if ($LASTEXITCODE -ne 0) { throw "Kyverno install failed." }
  & kubectl -n kyverno rollout status deploy/kyverno-admission-controller --timeout=240s
  if ($LASTEXITCODE -ne 0) { throw "Kyverno admission controller did not become ready." }
}

Write-Section "Apply Repository Policies"
& kubectl apply -k deploy/policies/kyverno
if ($LASTEXITCODE -ne 0) { throw "Applying repository Kyverno policies failed." }

Write-Section "Apply Local Matrix Signature + Attestation Policy"
$publicKeyBlock = Get-Content -Path $cosignPubPath -Raw
$indentedKeyVerify = ($publicKeyBlock.TrimEnd() -split "`r?`n" | ForEach-Object { "                      $_" }) -join "`n"
$indentedKeyAttestation = ($publicKeyBlock.TrimEnd() -split "`r?`n" | ForEach-Object { "                          $_" }) -join "`n"
@"
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: verify-local-matrix-images
spec:
  validationFailureAction: Enforce
  background: true
  rules:
    - name: verify-local-matrix-signature-and-attestation
      match:
        any:
          - resources:
              kinds: ["Pod"]
      verifyImages:
        - imageReferences:
            - "$signedRepo*"
            - "$unsignedRepo*"
          attestors:
            - entries:
                - keys:
                    publicKeys: |
$indentedKeyVerify
          attestations:
            - type: https://slsa.dev/provenance/v0.2
              attestors:
                - entries:
                    - keys:
                        publicKeys: |
$indentedKeyAttestation
"@ | Set-Content -Path $localPolicyPath -Encoding UTF8

& kubectl delete clusterpolicy verify-local-matrix-images --ignore-not-found=true | Out-Null
& kubectl apply -f $localPolicyPath
if ($LASTEXITCODE -ne 0) { throw "Applying local matrix policy failed." }

Write-Section "Apply Base Namespace Resources"
if ($ResetNamespace) {
  & kubectl delete namespace $Namespace --ignore-not-found=true --timeout=120s | Out-Null
  Start-Sleep -Seconds 3
}

@"
apiVersion: v1
kind: Namespace
metadata:
  name: $Namespace
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: user-service
  namespace: $Namespace
  labels:
    app.kubernetes.io/name: user-service
    app.kubernetes.io/component: backend
    app.kubernetes.io/part-of: stock-trading
---
apiVersion: v1
kind: Secret
metadata:
  name: user-service-secrets
  namespace: $Namespace
  labels:
    app.kubernetes.io/name: user-service
    app.kubernetes.io/component: backend
    app.kubernetes.io/part-of: stock-trading
stringData:
  db-password: "unused"
  auth-access-token-secret: "local-demo-access-secret-0123456789"
  auth-refresh-token-secret: "local-demo-refresh-secret-0123456789"
  smtp-username: ""
  smtp-password: ""
type: Opaque
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: user-service-config
  namespace: $Namespace
  labels:
    app.kubernetes.io/name: user-service
    app.kubernetes.io/component: backend
    app.kubernetes.io/part-of: stock-trading
data:
  config.yaml: |
    env: local
    grpc:
      host: 0.0.0.0
      port: 19090
    http:
      host: 0.0.0.0
      port: 18080
    db:
      host: 127.0.0.1
      port: 3306
      user: root
      password: ""
      name: stock
    auth:
      access_token_secret: "local-demo-access-secret-0123456789"
      access_token_ttl_minutes: 15
      refresh_token_secret: "local-demo-refresh-secret-0123456789"
      refresh_token_ttl_minutes: 4320
      issuer: stock-trading-be
      audience: stock-trading-clients
    notification:
      kafka:
        brokers: []
        topic: ""
        group_id: email-service
      email:
        provider: noop
        smtp:
          host: localhost
          port: 1025
          username: ""
          password: ""
          from: no-reply@example.com
          use_tls: false
        verification_url_base: "http://127.0.0.1:18080/users/verify?token="
    verification:
      token_ttl_hours: 24
      resend_cooldown_seconds: 60
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: $Namespace
  labels:
    app.kubernetes.io/name: user-service
    app.kubernetes.io/component: backend
    app.kubernetes.io/part-of: stock-trading
spec:
  selector:
    app.kubernetes.io/name: user-service
  ports:
    - name: http
      port: 18080
      targetPort: http
    - name: grpc
      port: 19090
      targetPort: grpc
"@ | Set-Content -Path $baseResourcesPath -Encoding UTF8

& kubectl apply -f $baseResourcesPath
if ($LASTEXITCODE -ne 0) { throw "Applying base resources failed." }

Write-Section "Run Admission Matrix"
$cases = @(
  [ordered]@{
    Name = "VALID_ALLOW"
    Expected = "Allowed"
    Image = $signedDigest
    IncludeSbom = $true
    HighCritical = "0"
  },
  [ordered]@{
    Name = "NEG_UNSIGNED_DENY"
    Expected = "Denied"
    Image = $unsignedDigest
    IncludeSbom = $true
    HighCritical = "0"
  },
  [ordered]@{
    Name = "NEG_MISSING_SBOM_DENY"
    Expected = "Denied"
    Image = $signedDigest
    IncludeSbom = $false
    HighCritical = "0"
  },
  [ordered]@{
    Name = "NEG_CVE_THRESHOLD_DENY"
    Expected = "Denied"
    Image = $signedDigest
    IncludeSbom = $true
    HighCritical = "2"
  }
)

$results = @()

foreach ($case in $cases) {
  $caseName = [string]$case.Name
  $caseDir = Join-Path $runDir $caseName
  Ensure-Directory $caseDir

  Write-Section "Case $caseName"
  Write-Host "Expected: $($case.Expected)"

  & kubectl -n $Namespace delete deploy user-service --ignore-not-found=true | Out-Null
  Start-Sleep -Seconds 4

  $deployYaml = Build-DeploymentYaml `
    -NamespaceName $Namespace `
    -Image ([string]$case.Image) `
    -SbomDigest $sbomDigest `
    -IncludeSbom ([bool]$case.IncludeSbom) `
    -HighCritical ([string]$case.HighCritical) `
    -CaseName $caseName

  $caseDeployPath = Join-Path $caseDir "deployment.yaml"
  $deployYaml | Set-Content -Path $caseDeployPath -Encoding UTF8

  $allowApplyFailure = ([string]$case.Expected -eq "Denied")
  $apply = Invoke-Kubectl -KubectlArgs @("-n", $Namespace, "apply", "-f", $caseDeployPath) -OutFile (Join-Path $caseDir "kubectl-apply.txt") -AllowFailure:$allowApplyFailure

  if ($apply.ExitCode -eq 0) {
    $wait = Invoke-Kubectl -KubectlArgs @("-n", $Namespace, "wait", "--for=condition=Available", "deployment/user-service", "--timeout=35s") -OutFile (Join-Path $caseDir "kubectl-wait.txt") -AllowFailure
  } else {
    Write-Text -path (Join-Path $caseDir "kubectl-wait.txt") -lines @("Skipped wait because apply failed with exit code $($apply.ExitCode).")
    $wait = @{
      ExitCode = -1
      Output = "Skipped wait because apply failed."
    }
  }

  Start-Sleep -Seconds 4
  Collect-CaseEvidence -CaseDir $caseDir -NamespaceName $Namespace
  $verdictResult = Get-CaseVerdict -Expected ([string]$case.Expected) -ApplyExitCode $apply.ExitCode -WaitExitCode $wait.ExitCode -CaseDir $caseDir

  $results += [pscustomobject]@{
    case = $caseName
    expected = [string]$case.Expected
    actual = [string]$verdictResult.Actual
    verdict = [string]$verdictResult.Verdict
    reason = [string]$verdictResult.Reason
    image = [string]$case.Image
    include_sbom = [bool]$case.IncludeSbom
    high_critical = [string]$case.HighCritical
    apply_exit_code = [int]$apply.ExitCode
    wait_exit_code = [int]$wait.ExitCode
    artifacts = [pscustomobject]@{
      apply = "$caseName/kubectl-apply.txt"
      wait = "$caseName/kubectl-wait.txt"
      events = "$caseName/events.txt"
      workloads = "$caseName/workloads.txt"
      describe_deployment = "$caseName/describe-deployment.txt"
      describe_replicasets = "$caseName/describe-replicasets.txt"
      describe_pods = "$caseName/describe-pods.txt"
      kyverno_logs = "$caseName/kyverno-logs.txt"
    }
  }
}

$summaryPath = Join-Path $runDir "matrix-summary.md"
$indexPath = Join-Path $runDir "matrix-index.json"
$regressionPath = Join-Path $runDir "regression-valid-allow.json"

Write-Section "Regression Re-check VALID_ALLOW"
$regressionCaseName = "VALID_ALLOW_RECHECK"
$regressionCaseDir = Join-Path $runDir $regressionCaseName
Ensure-Directory $regressionCaseDir

& kubectl -n $Namespace delete deploy user-service --ignore-not-found=true | Out-Null
Start-Sleep -Seconds 4

$regressionDeployYaml = Build-DeploymentYaml `
  -NamespaceName $Namespace `
  -Image $signedDigest `
  -SbomDigest $sbomDigest `
  -IncludeSbom $true `
  -HighCritical "0" `
  -CaseName $regressionCaseName

$regressionDeployPath = Join-Path $regressionCaseDir "deployment.yaml"
$regressionDeployYaml | Set-Content -Path $regressionDeployPath -Encoding UTF8

$regressionApply = Invoke-Kubectl -KubectlArgs @("-n", $Namespace, "apply", "-f", $regressionDeployPath) -OutFile (Join-Path $regressionCaseDir "kubectl-apply.txt")

if ($regressionApply.ExitCode -eq 0) {
  $regressionWait = Invoke-Kubectl -KubectlArgs @("-n", $Namespace, "wait", "--for=condition=Available", "deployment/user-service", "--timeout=35s") -OutFile (Join-Path $regressionCaseDir "kubectl-wait.txt") -AllowFailure
} else {
  Write-Text -path (Join-Path $regressionCaseDir "kubectl-wait.txt") -lines @("Skipped wait because apply failed with exit code $($regressionApply.ExitCode).")
  $regressionWait = @{
    ExitCode = -1
    Output = "Skipped wait because apply failed."
  }
}

Start-Sleep -Seconds 4
Collect-CaseEvidence -CaseDir $regressionCaseDir -NamespaceName $Namespace
$regressionVerdict = Get-CaseVerdict -Expected "Allowed" -ApplyExitCode $regressionApply.ExitCode -WaitExitCode $regressionWait.ExitCode -CaseDir $regressionCaseDir

$regressionResult = [pscustomobject]@{
  case = $regressionCaseName
  expected = "Allowed"
  actual = [string]$regressionVerdict.Actual
  verdict = [string]$regressionVerdict.Verdict
  reason = [string]$regressionVerdict.Reason
  image = $signedDigest
  include_sbom = $true
  high_critical = "0"
  apply_exit_code = [int]$regressionApply.ExitCode
  wait_exit_code = [int]$regressionWait.ExitCode
  artifacts = [pscustomobject]@{
    apply = "$regressionCaseName/kubectl-apply.txt"
    wait = "$regressionCaseName/kubectl-wait.txt"
    events = "$regressionCaseName/events.txt"
    workloads = "$regressionCaseName/workloads.txt"
    describe_deployment = "$regressionCaseName/describe-deployment.txt"
    describe_replicasets = "$regressionCaseName/describe-replicasets.txt"
    describe_pods = "$regressionCaseName/describe-pods.txt"
    kyverno_logs = "$regressionCaseName/kyverno-logs.txt"
  }
}

$summaryLines = @(
  "# Admission Matrix Summary",
  "",
  "- Run ID: $runId",
  "- Kubernetes context: $Context",
  "- Namespace: $Namespace",
  "- Signed image digest: $signedDigest",
  "- Unsigned image digest: $unsignedDigest",
  "- SBOM digest: $sbomDigest",
  "",
  "| Case | Expected | Actual | Verdict | Reason |",
  "|---|---|---|---|---|"
)

foreach ($result in $results) {
  $summaryLines += "| $($result.case) | $($result.expected) | $($result.actual) | $($result.verdict) | $($result.reason) |"
}

$summaryLines += ""
$summaryLines += "## Regression Re-check"
$summaryLines += ""
$summaryLines += "| Check | Expected | Actual | Verdict | Reason |"
$summaryLines += "|---|---|---|---|---|"
$summaryLines += "| $($regressionResult.case) | $($regressionResult.expected) | $($regressionResult.actual) | $($regressionResult.verdict) | $($regressionResult.reason) |"

Write-Text -path $summaryPath -lines $summaryLines
($results | ConvertTo-Json -Depth 8) | Set-Content -Path $indexPath -Encoding UTF8
$regressionResult | ConvertTo-Json -Depth 8 | Set-Content -Path $regressionPath -Encoding UTF8

$failed = @($results | Where-Object { $_.verdict -ne "PASS" })
if ($regressionResult.verdict -ne "PASS") {
  $failed += $regressionResult
}
$failedCount = @($failed).Count

Write-Section "Matrix Summary"
Get-Content -Path $summaryPath
Write-Host ""
Write-Host "Evidence directory: $runDir"
Write-Host "JSON index: $indexPath"
Write-Host "Regression result: $regressionPath"

if ($failedCount -gt 0) {
  throw "Admission matrix has failing cases. Check $summaryPath."
}

Write-Host ""
Write-Host "Admission matrix completed successfully."
