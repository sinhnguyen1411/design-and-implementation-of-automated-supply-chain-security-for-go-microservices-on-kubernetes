param(
  [string]$Namespace = "stock-trading",
  [string]$KyvernoVersion = "v1.12.5",
  [string]$ImageTtl = "24h",
  [string]$CosignPassword = "local-demo-pass",
  [string]$DemoDir = ".demo",
  [switch]$ResetNamespace
)

$ErrorActionPreference = "Stop"

function Require-Cli($name) {
  if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
    throw "Required CLI not found: $name"
  }
}

function Write-Section([string]$message) {
  Write-Host ""
  Write-Host "== $message =="
}

function Invoke-YamlApply([string]$yaml) {
  $yaml | kubectl apply -f -
  if ($LASTEXITCODE -ne 0) {
    throw "kubectl apply failed"
  }
}

function Wait-ForDeployment([string]$name, [string]$ns, [int]$seconds = 180) {
  kubectl -n $ns rollout status "deploy/$name" --timeout="$($seconds)s"
  if ($LASTEXITCODE -ne 0) {
    throw "Deployment '$name' did not become ready in namespace '$ns'"
  }
}

function New-DemoImageRef([string]$ttl) {
  $suffix = [guid]::NewGuid().ToString("N").Substring(0, 12)
  return "ttl.sh/stock-trading-${suffix}:$ttl"
}

Require-Cli docker
Require-Cli kubectl
Require-Cli cosign
Require-Cli syft

Write-Section "Cluster Check"
$context = kubectl config current-context
if ($LASTEXITCODE -ne 0) {
  throw "kubectl is not configured for any cluster"
}
Write-Host "kubectl context: $context"
kubectl get nodes -o wide
if ($LASTEXITCODE -ne 0) {
  throw "Cluster is not reachable"
}

Write-Section "Workspace Prep"
New-Item -ItemType Directory -Force -Path $DemoDir | Out-Null
$cosignKeyPath = Join-Path $DemoDir "cosign.key"
$cosignPubPath = Join-Path $DemoDir "cosign.pub"
$sbomPath = Join-Path $DemoDir "sbom.spdx.json"
$provenancePath = Join-Path $DemoDir "provenance.json"
$policyPath = Join-Path $DemoDir "verify-local-demo-image.yaml"
$resourcesPath = Join-Path $DemoDir "local-demo-resources.yaml"
$deployPath = Join-Path $DemoDir "local-demo-deployment.yaml"

Write-Section "Build Image"
$localImage = "local-user-service:test"
docker build -t $localImage .
if ($LASTEXITCODE -ne 0) {
  throw "docker build failed"
}

Write-Section "Push Image"
$remoteImage = New-DemoImageRef $ImageTtl
docker tag $localImage $remoteImage
if ($LASTEXITCODE -ne 0) {
  throw "docker tag failed"
}
docker push $remoteImage
if ($LASTEXITCODE -ne 0) {
  throw "docker push failed"
}
$remoteRepo = $remoteImage.Split(":")[0]
$repoDigests = docker inspect --format='{{range .RepoDigests}}{{println .}}{{end}}' $remoteImage
if ($LASTEXITCODE -ne 0) {
  throw "docker inspect failed after push"
}
$remoteDigest = ($repoDigests -split "`r?`n" | Where-Object { $_ -like "$remoteRepo@*" } | Select-Object -First 1).Trim()
if (-not $remoteDigest) {
  throw "Unable to resolve pushed image digest"
}
Write-Host "Pushed image: $remoteDigest"

Write-Section "Generate Signing Materials"
$env:COSIGN_PASSWORD = $CosignPassword
if (!(Test-Path $cosignKeyPath) -or !(Test-Path $cosignPubPath)) {
  cosign generate-key-pair --output-key-prefix ($cosignKeyPath -replace '\.key$','')
  if ($LASTEXITCODE -ne 0) {
    throw "cosign key generation failed"
  }
}
syft $remoteDigest -o "spdx-json=$sbomPath"
if ($LASTEXITCODE -ne 0) {
  throw "syft SBOM generation failed"
}
$sbomDigest = (Get-FileHash $sbomPath -Algorithm SHA256).Hash
@"
{
  "buildType": "https://slsa.dev/provenance/v0.2",
  "builder": { "id": "local-demo" },
  "invocation": {
    "configSource": {
      "uri": "https://github.com/sinhnguyen1411/stock-trading-be",
      "digest": { "gitCommit": "local-demo" }
    }
  }
}
"@ | Set-Content -Path $provenancePath

Write-Section "Sign And Attest"
cosign sign --yes --key $cosignKeyPath $remoteDigest
if ($LASTEXITCODE -ne 0) {
  throw "cosign sign failed"
}
cosign attest --yes --key $cosignKeyPath --predicate $provenancePath --type slsaprovenance $remoteDigest
if ($LASTEXITCODE -ne 0) {
  throw "cosign attest failed"
}
cosign verify --key $cosignPubPath $remoteDigest | Out-Null
if ($LASTEXITCODE -ne 0) {
  throw "cosign verify failed"
}

Write-Section "Ensure Kyverno"
$hasKyverno = kubectl api-resources | Select-String -Pattern "^clusterpolicies\s"
if (-not $hasKyverno) {
  kubectl apply --server-side -f "https://github.com/kyverno/kyverno/releases/download/$KyvernoVersion/install.yaml"
  if ($LASTEXITCODE -ne 0) {
    throw "Kyverno install failed"
  }
  kubectl -n kyverno rollout status deploy/kyverno-admission-controller --timeout=240s
  if ($LASTEXITCODE -ne 0) {
    throw "Kyverno admission controller did not become ready"
  }
}

Write-Section "Apply Repo Policies"
kubectl apply -k deploy/policies/kyverno
if ($LASTEXITCODE -ne 0) {
  throw "Applying repository Kyverno policies failed"
}

Write-Section "Apply Local Demo Policy"
$publicKeyBlock = Get-Content -Path $cosignPubPath -Raw
$indentedKey = ($publicKeyBlock.TrimEnd() -split "`r?`n" | ForEach-Object { "                      $_" }) -join "`n"
@"
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: verify-local-demo-image
spec:
  validationFailureAction: Enforce
  background: true
  rules:
    - name: verify-local-demo-signature
      match:
        any:
          - resources:
              kinds: ["Pod"]
      verifyImages:
        - imageReferences:
            - "$($remoteDigest.Split('@')[0])*"
          attestors:
            - entries:
                - keys:
                    publicKeys: |
$indentedKey
"@ | Set-Content -Path $policyPath
kubectl delete clusterpolicy verify-local-demo-image --ignore-not-found=true | Out-Null
kubectl apply -f $policyPath
if ($LASTEXITCODE -ne 0) {
  throw "Applying local demo Kyverno policy failed"
}

Write-Section "Deploy Local Demo"
if ($ResetNamespace) {
  kubectl delete namespace $Namespace --ignore-not-found=true --timeout=120s | Out-Null
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
"@ | Set-Content -Path $resourcesPath
Invoke-YamlApply (Get-Content -Raw $resourcesPath)

@"
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: $Namespace
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
        security.grype.io/high_critical: "0"
        security.stock-trading.dev/sbom-digest: "$sbomDigest"
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
          image: $remoteDigest
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
"@ | Set-Content -Path $deployPath
kubectl delete deploy user-service -n $Namespace --ignore-not-found=true | Out-Null
Invoke-YamlApply (Get-Content -Raw $deployPath)
Wait-ForDeployment "user-service" $Namespace 180

Write-Section "Summary"
kubectl get deploy,pods -n $Namespace -o wide
kubectl logs -n $Namespace deploy/user-service --tail=40
Write-Host ""
Write-Host "Image digest: $remoteDigest"
Write-Host "SBOM digest : $sbomDigest"
Write-Host "Public key  : $cosignPubPath"
