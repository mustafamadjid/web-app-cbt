param(
    [string[]]$Packages = @("./..."),
    [string]$Tags = ""
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$envFile = Join-Path $PSScriptRoot "export_env.test.ps1"
if ([string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable("POSTGRES_TEST_DBURL")) -and (Test-Path $envFile)) {
    . $envFile
}

$testDbUrl = [Environment]::GetEnvironmentVariable("POSTGRES_TEST_DBURL")
if ([string]::IsNullOrWhiteSpace($testDbUrl)) {
    throw "POSTGRES_TEST_DBURL must be set before running integration tests."
}

$previousMainDbUrl = [Environment]::GetEnvironmentVariable("POSTGRES_DBURL")
$env:POSTGRES_DBURL = $testDbUrl.Trim()

try {
    Write-Host "Running integration tests against: $env:POSTGRES_DBURL"

    $goArgs = @("test")
    if (-not [string]::IsNullOrWhiteSpace($Tags)) {
        $goArgs += @("-tags", $Tags.Trim())
    }
    $goArgs += $Packages

    & go @goArgs
    if ($LASTEXITCODE -ne 0) {
        throw "go test failed with exit code $LASTEXITCODE."
    }
}
finally {
    if ([string]::IsNullOrWhiteSpace($previousMainDbUrl)) {
        Remove-Item Env:POSTGRES_DBURL -ErrorAction SilentlyContinue
    }
    else {
        $env:POSTGRES_DBURL = $previousMainDbUrl
    }
}
