param(
  [string]$Repo = "sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes",
  [string]$Branch = "test/ci-cve-signal",
  [int]$Iterations = 5
)

$ErrorActionPreference = "Stop"

function Get-PercentileValue {
  param([double[]]$Values,[double]$Percentile)
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
  if (($count % 2) -eq 1) { return [double]$sorted[[int]($count / 2)] }
  return ([double]$sorted[($count / 2) - 1] + [double]$sorted[($count / 2)]) / 2.0
}

$runIds = @()
for ($i = 1; $i -le $Iterations; $i++) {
  Write-Host "Dispatch iteration $i/$Iterations ..."
  gh workflow run ci-service.yml --repo $Repo --ref $Branch -f service=all -f mode=benchmark | Out-Null
  Start-Sleep -Seconds 2
  $run = gh run list --repo $Repo --workflow ci-service.yml --branch $Branch --limit 1 --json databaseId | ConvertFrom-Json
  $runId = $run[0].databaseId
  $runIds += $runId
  Write-Host "Run: $runId"
}

foreach ($runId in $runIds) {
  Write-Host "Waiting run $runId ..."
  gh run watch $runId --repo $Repo --interval 8 --exit-status | Out-Null
}

$targets = @("windows-gh-hosted-smoke", "windows-parity-smoke")
$rows = @()
foreach ($runId in $runIds) {
  $jobs = gh api repos/$Repo/actions/runs/$runId/jobs | ConvertFrom-Json
  foreach ($job in $jobs.jobs) {
    if ($job.name -notin $targets) { continue }
    $duration = ([DateTime]::Parse($job.completed_at).ToUniversalTime() - [DateTime]::Parse($job.started_at).ToUniversalTime()).TotalSeconds
    $rows += [PSCustomObject]@{
      run_id = $runId
      job_name = $job.name
      duration_sec = [Math]::Round($duration, 2)
      conclusion = $job.conclusion
      runner_name = $job.runner_name
    }
  }
}

$summary = foreach ($name in $targets) {
  $vals = @($rows | Where-Object { $_.job_name -eq $name } | ForEach-Object { [double]$_.duration_sec })
  [PSCustomObject]@{
    job_name = $name
    samples = $vals.Count
    median_sec = [Math]::Round((Get-MedianValue -Values $vals), 2)
    p95_sec = [Math]::Round((Get-PercentileValue -Values $vals -Percentile 95), 2)
    min_sec = [Math]::Round((($vals | Measure-Object -Minimum).Minimum), 2)
    max_sec = [Math]::Round((($vals | Measure-Object -Maximum).Maximum), 2)
  }
}

Write-Host ""
Write-Host "=== Raw Durations ==="
$rows | Sort-Object run_id,job_name | Format-Table -AutoSize
Write-Host ""
Write-Host "=== Summary ==="
$summary | Format-Table -AutoSize
