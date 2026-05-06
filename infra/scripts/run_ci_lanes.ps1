param(
  [string]$Repo = "",
  [string]$Branch = "main",
  [bool]$Wait = $true,
  [bool]$StopOnFail = $true
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest
if ($PSVersionTable.PSVersion.Major -ge 7) {
  $PSNativeCommandUseErrorActionPreference = $false
}

function Write-Section([string]$Message) {
  Write-Host ""
  Write-Host "== $Message =="
}

function Require-Cli([string]$Name) {
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    throw "Required CLI not found: $Name"
  }
}

function Resolve-Repo([string]$RepoInput) {
  if ($RepoInput) {
    return $RepoInput
  }
  $origin = (& git remote get-url origin 2>$null).Trim()
  if (-not $origin) {
    throw "Cannot resolve origin remote. Provide -Repo owner/repo."
  }
  $pattern = 'github\.com[:/](?<owner>[^/]+)/(?<repo>[^/.]+)(?:\.git)?$'
  $match = [regex]::Match($origin, $pattern, [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
  if (-not $match.Success) {
    throw "Origin does not point to github.com: $origin"
  }
  return "$($match.Groups['owner'].Value)/$($match.Groups['repo'].Value)"
}

function Get-WorkflowRunId([string]$RepoName, [string]$WorkflowName, [string]$Ref, [datetime]$NotBeforeUtc) {
  $query = "repos/$RepoName/actions/workflows/$WorkflowName/runs?branch=$Ref&per_page=1"
  $json = & gh api $query --jq '.workflow_runs[0] | {id,created_at}'
  if ($LASTEXITCODE -ne 0 -or -not $json) {
    return $null
  }
  $obj = $json | ConvertFrom-Json
  if (-not $obj) {
    return $null
  }
  $created = [datetime]::Parse($obj.created_at).ToUniversalTime()
  if ($created -lt $NotBeforeUtc) {
    return $null
  }
  return [int64]$obj.id
}

function Wait-ForRunCreation([string]$RepoName, [string]$WorkflowName, [string]$Ref, [datetime]$NotBeforeUtc) {
  for ($i = 0; $i -lt 30; $i++) {
    $id = Get-WorkflowRunId -RepoName $RepoName -WorkflowName $WorkflowName -Ref $Ref -NotBeforeUtc $NotBeforeUtc
    if ($id) {
      return $id
    }
    Start-Sleep -Seconds 2
  }
  throw "No run detected for workflow '$WorkflowName' on branch '$Ref' after dispatch."
}

function Get-RunStatus([string]$RepoName, [int64]$RunId) {
  $json = & gh api "repos/$RepoName/actions/runs/$RunId" --jq '{id,status,conclusion,html_url,name,head_sha,run_number}'
  if ($LASTEXITCODE -ne 0 -or -not $json) {
    throw "Unable to query run status for run id: $RunId"
  }
  return ($json | ConvertFrom-Json)
}

function Wait-ForRunCompletion([string]$RepoName, [int64]$RunId) {
  while ($true) {
    $run = Get-RunStatus -RepoName $RepoName -RunId $RunId
    Write-Host ("[{0}] status={1} conclusion={2}" -f $run.name, $run.status, $run.conclusion)
    if ($run.status -eq "completed") {
      return $run
    }
    Start-Sleep -Seconds 10
  }
}

function Dispatch-Workflow([string]$RepoName, [string]$WorkflowName, [string]$Ref, [hashtable]$Inputs) {
  $args = @("workflow", "run", $WorkflowName, "--repo", $RepoName, "--ref", $Ref)
  foreach ($key in $Inputs.Keys) {
    $args += @("-f", "$key=$($Inputs[$key])")
  }
  & gh @args
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to dispatch workflow: $WorkflowName"
  }
}

Require-Cli -Name "gh"
Require-Cli -Name "git"

& cmd /c "gh auth status -h github.com 1>nul 2>nul"
if ($LASTEXITCODE -ne 0) {
  throw "gh is not authenticated. Run bootstrap_github_automation.ps1 first."
}

$resolvedRepo = Resolve-Repo -RepoInput $Repo
$summary = New-Object 'System.Collections.Generic.List[object]'

$lanes = @(
  @{ key = "ci-service"; workflow = "ci-service.yml"; inputs = @{ service = "all" } },
  @{ key = "admission-lab"; workflow = "admission-matrix-evidence.yml"; inputs = @{} },
  @{ key = "onboarding-lab"; workflow = "service-scs-matrix-evidence.yml"; inputs = @{} },
  @{ key = "dashboard-data-sync"; workflow = "dashboard-data-sync.yml"; inputs = @{} }
)

foreach ($lane in $lanes) {
  Write-Section ("Dispatch lane: {0}" -f $lane.key)
  $dispatchTime = [datetime]::UtcNow.AddSeconds(-2)
  Dispatch-Workflow -RepoName $resolvedRepo -WorkflowName $lane.workflow -Ref $Branch -Inputs $lane.inputs
  $runId = Wait-ForRunCreation -RepoName $resolvedRepo -WorkflowName $lane.workflow -Ref $Branch -NotBeforeUtc $dispatchTime
  Write-Host ("Created run id: {0}" -f $runId)

  if (-not $Wait) {
    $summary.Add([pscustomobject]@{
        lane       = $lane.key
        run_id     = $runId
        conclusion = "queued"
        url        = "https://github.com/$resolvedRepo/actions/runs/$runId"
      })
    continue
  }

  $run = Wait-ForRunCompletion -RepoName $resolvedRepo -RunId $runId
  $summary.Add([pscustomobject]@{
      lane       = $lane.key
      run_id     = $run.id
      conclusion = $run.conclusion
      url        = $run.html_url
    })

  if ($StopOnFail -and $run.conclusion -ne "success") {
    Write-Host ("Lane failed and StopOnFail=true: {0}" -f $lane.key)
    break
  }
}

Write-Section "Run summary"
$summary | Format-Table -AutoSize

$failedLane = $summary | Where-Object { $_.conclusion -notin @("success", "queued") } | Select-Object -First 1
if ($failedLane -and $StopOnFail) {
  throw ("CI lane sequence stopped on failure: {0} (run {1})" -f $failedLane.lane, $failedLane.run_id)
}
