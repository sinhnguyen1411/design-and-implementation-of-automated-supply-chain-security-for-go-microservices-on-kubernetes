param(
  [string]$Repo = "",
  [bool]$InstallMissing = $true,
  [bool]$SetupAuth = $true,
  [bool]$NonInteractive = $false,
  [bool]$SkipWorkflowChecks = $false
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

function Add-Result(
  [System.Collections.Generic.List[object]]$Results,
  [string]$CheckName,
  [string]$Status,
  [string]$Code,
  [string]$Detail
) {
  $Results.Add([pscustomobject]@{
      check  = $CheckName
      status = $Status
      code   = $Code
      detail = $Detail
    })
}

function Invoke-QuickCheck(
  [System.Collections.Generic.List[object]]$Results,
  [string]$CheckName,
  [scriptblock]$Action
) {
  try {
    & $Action | Out-Null
    Add-Result -Results $Results -CheckName $CheckName -Status "READY" -Code "OK" -Detail ""
    return $true
  } catch {
    Add-Result -Results $Results -CheckName $CheckName -Status "NOT_READY" -Code "ERR" -Detail $_.Exception.Message
    return $false
  }
}

function Get-RepoFromOrigin {
  $originUrl = (& git remote get-url origin 2>$null).Trim()
  if (-not $originUrl) {
    throw "Unable to resolve origin remote URL. Provide -Repo owner/repo."
  }

  $pattern = 'github\.com[:/](?<owner>[^/]+)/(?<repo>[^/.]+)(?:\.git)?$'
  $match = [regex]::Match($originUrl, $pattern, [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
  if (-not $match.Success) {
    throw "Origin does not point to github.com: $originUrl"
  }

  return "$($match.Groups['owner'].Value)/$($match.Groups['repo'].Value)"
}

function Resolve-Repo([string]$RepoInput) {
  if ($RepoInput) {
    return $RepoInput
  }
  return Get-RepoFromOrigin
}

function Require-Command([string]$Name) {
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    throw "Missing command: $Name"
  }
}

function Install-GhCli {
  Require-Command -Name "winget"
  Write-Host "Installing GitHub CLI with winget..."
  & winget install --id GitHub.cli --exact --silent --accept-package-agreements --accept-source-agreements
  if ($LASTEXITCODE -ne 0) {
    throw "winget failed to install GitHub CLI. Try elevated PowerShell and run again."
  }
}

function Test-GhAuth {
  & cmd /c "gh auth status -h github.com 1>nul 2>nul"
  return ($LASTEXITCODE -eq 0)
}

function ConvertTo-PlainText([securestring]$SecureValue) {
  $ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($SecureValue)
  try {
    return [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr)
  } finally {
    [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr)
  }
}

function Get-PatToken([bool]$UseNonInteractive) {
  if ($UseNonInteractive) {
    $token = $env:GH_TOKEN
    if (-not $token) {
      throw "NonInteractive mode requires GH_TOKEN environment variable."
    }
    return $token.Trim()
  }

  $secure = Read-Host "Enter GitHub PAT (input hidden)" -AsSecureString
  $token = ConvertTo-PlainText -SecureValue $secure
  if (-not $token) {
    throw "PAT cannot be empty."
  }
  return $token.Trim()
}

function Login-Gh([bool]$UseNonInteractive) {
  $token = Get-PatToken -UseNonInteractive $UseNonInteractive
  try {
    $token | & gh auth login --hostname github.com --with-token
    if ($LASTEXITCODE -ne 0) {
      throw "gh auth login failed."
    }

    & gh auth setup-git
    if ($LASTEXITCODE -ne 0) {
      throw "gh auth setup-git failed."
    }
  } finally {
    Remove-Variable token -ErrorAction SilentlyContinue
  }
}

function Assert-GhRepoAccess([string]$RepoName) {
  & gh api "repos/$RepoName" 1>$null 2>$null
  if ($LASTEXITCODE -ne 0) {
    throw "Unable to access repo via gh api: $RepoName"
  }
}

function Assert-GhWorkflowAccess([string]$RepoName) {
  & gh workflow list --repo $RepoName 1>$null 2>$null
  if ($LASTEXITCODE -ne 0) {
    throw "Unable to list workflows for repo: $RepoName"
  }
}

$results = New-Object 'System.Collections.Generic.List[object]'
$resolvedRepo = Resolve-Repo -RepoInput $Repo

Write-Section "Preflight checks"
Invoke-QuickCheck -Results $results -CheckName "git-installed" -Action { Require-Command -Name "git" } | Out-Null

$hasGh = Invoke-QuickCheck -Results $results -CheckName "gh-installed" -Action { Require-Command -Name "gh" }
Invoke-QuickCheck -Results $results -CheckName "winget-installed" -Action { Require-Command -Name "winget" } | Out-Null

Invoke-QuickCheck -Results $results -CheckName "github-desktop-installed" -Action {
  if (-not (Test-Path "$env:LOCALAPPDATA\GitHubDesktop\GitHubDesktop.exe")) {
    throw "GitHubDesktop.exe not found in LOCALAPPDATA."
  }
} | Out-Null

Invoke-QuickCheck -Results $results -CheckName "origin-reachable" -Action { & git ls-remote --heads origin 1>$null 2>$null } | Out-Null

if (-not $hasGh -and $InstallMissing) {
  Write-Section "Auto install"
  Install-GhCli
  Invoke-QuickCheck -Results $results -CheckName "gh-installed-after-install" -Action { Require-Command -Name "gh" } | Out-Null
} elseif (-not $hasGh -and -not $InstallMissing) {
  throw "GitHub CLI missing and InstallMissing=false. Install gh and retry."
}

Write-Section "Authentication"
$alreadyAuth = Test-GhAuth
if ($alreadyAuth) {
  Add-Result -Results $results -CheckName "gh-auth-status" -Status "READY" -Code "OK" -Detail "Already authenticated."
  Write-Host "gh is already authenticated for github.com"
} elseif ($SetupAuth) {
  Login-Gh -UseNonInteractive $NonInteractive
  if (-not (Test-GhAuth)) {
    throw "gh authentication did not succeed."
  }
  Add-Result -Results $results -CheckName "gh-auth-status" -Status "READY" -Code "OK" -Detail "Authenticated during bootstrap."
} else {
  Add-Result -Results $results -CheckName "gh-auth-status" -Status "NOT_READY" -Code "AUTH_REQUIRED" -Detail "Set SetupAuth=true to authenticate."
}

Write-Section "Automation self-check"
Invoke-QuickCheck -Results $results -CheckName "gh-api-user" -Action { & gh api user 1>$null 2>$null; if ($LASTEXITCODE -ne 0) { throw "gh api user failed." } } | Out-Null
Invoke-QuickCheck -Results $results -CheckName "gh-api-repo" -Action { Assert-GhRepoAccess -RepoName $resolvedRepo } | Out-Null

if (-not $SkipWorkflowChecks) {
  Invoke-QuickCheck -Results $results -CheckName "gh-workflow-list" -Action { Assert-GhWorkflowAccess -RepoName $resolvedRepo } | Out-Null
} else {
  Add-Result -Results $results -CheckName "gh-workflow-list" -Status "SKIPPED" -Code "SKIP" -Detail "SkipWorkflowChecks=true"
}

Write-Section "Result"
$results | Format-Table -AutoSize

$failed = $results | Where-Object { $_.status -eq "NOT_READY" }
if ($failed.Count -gt 0) {
  throw "Bootstrap finished with NOT_READY checks. Review table above."
}

Write-Host ""
Write-Host "Bootstrap completed successfully for repo: $resolvedRepo"
