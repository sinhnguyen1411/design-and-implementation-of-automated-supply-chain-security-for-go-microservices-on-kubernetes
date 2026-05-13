param(
  [string]$Repo = "sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes",
  [string]$Branch = "test/ci-cve-signal",
  [int]$Iterations = 5,
  [string]$TargetService = "services/user-service",
  [string]$GoVersion = "1.25.10"
)

$ErrorActionPreference = "Stop"

if ($Iterations -lt 2) {
  throw "Iterations must be >= 2"
}

function Get-PercentileValue {
  param(
    [double[]]$Values,
    [double]$Percentile
  )
  $sorted = @($Values | Sort-Object)
  if ($sorted.Count -eq 0) { return [double]::NaN }
  $index = [Math]::Ceiling(($Percentile / 100.0) * $sorted.Count) - 1
  if ($index -lt 0) { $index = 0 }
  if ($index -ge $sorted.Count) { $index = $sorted.Count - 1 }
  return [double]$sorted[$index]
}

function Get-MedianValue {
  param([double[]]$Values)
  $sorted = @($Values | Sort-Object)
  $count = $sorted.Count
  if ($count -eq 0) { return [double]::NaN }
  if (($count % 2) -eq 1) {
    return [double]$sorted[[int]($count / 2)]
  }
  $a = [double]$sorted[($count / 2) - 1]
  $b = [double]$sorted[($count / 2)]
  return ($a + $b) / 2.0
}

$runIds = @()
for ($i = 1; $i -le $Iterations; $i++) {
  Write-Host "Dispatch iteration $i/$Iterations ..."
  gh workflow run runner-ab-benchmark.yml --repo $Repo --ref $Branch -f target_service=$TargetService -f go_version=$GoVersion | Out-Null
  Start-Sleep -Seconds 3
  $run = gh run list --repo $Repo --workflow runner-ab-benchmark.yml --branch $Branch --limit 1 --json databaseId,status,createdAt,headSha | ConvertFrom-Json
  $runId = $run[0].databaseId
  $runIds += $runId
  Write-Host "Run: $runId"
}

foreach ($runId in $runIds) {
  Write-Host "Waiting run $runId ..."
  gh run watch $runId --repo $Repo --interval 8 --exit-status | Out-Null
}

$records = @()
foreach ($runId in $runIds) {
  $jobs = gh api repos/$Repo/actions/runs/$runId/jobs --paginate | ConvertFrom-Json
  foreach ($job in $jobs.jobs) {
    if ($job.name -notin @("gh-hosted-windows", "self-hosted-windows-parity")) {
      continue
    }
    $start = [DateTime]::Parse($job.started_at).ToUniversalTime()
    $end = [DateTime]::Parse($job.completed_at).ToUniversalTime()
    $seconds = [Math]::Round(($end - $start).TotalSeconds, 2)
    $records += [PSCustomObject]@{
      run_id = $runId
      job_name = $job.name
      conclusion = $job.conclusion
      runner_name = $job.runner_name
      duration_sec = $seconds
      started_at = $job.started_at
      completed_at = $job.completed_at
    }
  }
}

$groups = $records | Group-Object job_name
$summary = @()
foreach ($g in $groups) {
  $vals = @($g.Group.duration_sec | ForEach-Object { [double]$_ })
  $summary += [PSCustomObject]@{
    job_name = $g.Name
    samples = $vals.Count
    median_sec = [Math]::Round((Get-MedianValue -Values $vals), 2)
    p95_sec = [Math]::Round((Get-PercentileValue -Values $vals -Percentile 95), 2)
    min_sec = [Math]::Round((($vals | Measure-Object -Minimum).Minimum), 2)
    max_sec = [Math]::Round((($vals | Measure-Object -Maximum).Maximum), 2)
  }
}

Write-Host ""
Write-Host "=== Raw Durations (sec) ==="
$records | Sort-Object run_id, job_name | Format-Table -AutoSize
Write-Host ""
Write-Host "=== Summary (Median/P95) ==="
$summary | Sort-Object job_name | Format-Table -AutoSize

