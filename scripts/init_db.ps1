param(
  [string]$DBHost = "localhost",
  [int]$Port = 3306,
  [string]$User = "root",
  [string]$Password = "Ngdms1107#",
  [string]$Database = "stock",
  [string]$SchemaPath = "internal/adapters/database/schema_verification.sql"
)

function Require-Cli($name) {
  if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
    Write-Error "Required CLI not found: $name"
    exit 1
  }
}

function Try-LocalMySQL($mysqlArgs, [string]$stdin) {
  if (-not (Get-Command mysql -ErrorAction SilentlyContinue)) { return $false }
  try {
    if ($stdin) {
      $stdin | & mysql @mysqlArgs 2>$null
    } else {
      & mysql @mysqlArgs 2>$null
    }
    return $LASTEXITCODE -eq 0
  } catch {
    return $false
  }
}

function Normalize-HostForDocker($mysqlArgs) {
  $out = @()
  foreach ($a in $mysqlArgs) {
    if ($a -like '--host=*') {
      $v = $a.Substring(7)
      if ($v -eq 'localhost' -or $v -eq '127.0.0.1') {
        $out += '--host=host.docker.internal'
        continue
      }
    }
    $out += $a
  }
  return ,$out
}

function Run-DockerMySQL($mysqlArgs, [string]$stdin) {
  $norm = Normalize-HostForDocker $mysqlArgs
  $passwordArg = $null
  $mysqlCliArgs = @()
  foreach ($arg in $norm) {
    if ($arg -like '--password=*') {
      $passwordArg = $arg.Substring(11)
      continue
    }
    $mysqlCliArgs += $arg
  }
  $dockerArgs = @('run','--rm','-i')
  if ($null -ne $passwordArg) {
    $dockerArgs += @('-e', "MYSQL_PWD=$passwordArg")
  }
  $dockerArgs += @('mysql:8','mysql') + $mysqlCliArgs
  $display = $dockerArgs | ForEach-Object {
    if ($_ -like 'MYSQL_PWD=*') { 'MYSQL_PWD=****' } else { $_ }
  } | Out-String
  Write-Host "Running: docker $display"
  if ($stdin) {
    $stdin | & docker @dockerArgs
  } else {
    & docker @dockerArgs
  }
  return $LASTEXITCODE -eq 0
}

Require-Cli docker

if (!(Test-Path -Path $SchemaPath)) {
  Write-Error "Schema file not found: $SchemaPath"
  exit 1
}

Write-Host "Creating database '$Database' if not exists..."
$createSQL = "CREATE DATABASE IF NOT EXISTS ``$Database`` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
$args = @("--host=$DBHost","--port=$Port","--user=$User","--password=$Password","-e", $createSQL)
if (-not (Try-LocalMySQL $args $null)) {
  Write-Warning "Local mysql failed; falling back to docker mysql client..."
  if (-not (Run-DockerMySQL $args $null)) { exit 1 }
}

Write-Host "Applying schema from $SchemaPath..."
$schema = Get-Content -Raw -Path $SchemaPath
$args2 = @("--host=$DBHost","--port=$Port","--user=$User","--password=$Password", $Database)
if (-not (Try-LocalMySQL $args2 $schema)) {
  Write-Warning "Local mysql failed; falling back to docker mysql client..."
  if (-not (Run-DockerMySQL $args2 $schema)) { exit 1 }
}

Write-Host "Database initialized."
