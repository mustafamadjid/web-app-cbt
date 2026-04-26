param(
    [string[]]$Packages = @("./..."),
    [string]$Tags = ""
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

function Get-Percent {
    param(
        [int]$Count,
        [int]$Total
    )

    if ($Total -le 0) {
        return "0.00"
    }

    return ([math]::Round(($Count * 100.0) / $Total, 2)).ToString("0.00")
}

function TryParse-TestEvent {
    param(
        [string]$Line,
        [ref]$Event
    )

    try {
        $Event.Value = $Line | ConvertFrom-Json -ErrorAction Stop
        return $true
    }
    catch {
        $Event.Value = $null
        return $false
    }
}

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

    $goArgs = @("test", "-json")
    if (-not [string]::IsNullOrWhiteSpace($Tags)) {
        $goArgs += @("-tags", $Tags.Trim())
    }
    $goArgs += $Packages

    $testResults = @{}

    & go @goArgs 2>&1 | ForEach-Object {
        $line = $_.ToString()
        $event = $null

        if (-not (TryParse-TestEvent -Line $line -Event ([ref]$event))) {
            Write-Host $line
            return
        }

        if (-not [string]::IsNullOrWhiteSpace($event.Output)) {
            $message = $event.Output.TrimEnd("`r", "`n")
            if (-not [string]::IsNullOrWhiteSpace($message)) {
                Write-Host $message
            }
        }

        if ([string]::IsNullOrWhiteSpace($event.Test)) {
            return
        }

        switch ($event.Action) {
            "pass" {
                $testResults[$event.Test] = "pass"
            }
            "fail" {
                $testResults[$event.Test] = "fail"
            }
            "skip" {
                if (-not $testResults.ContainsKey($event.Test)) {
                    $testResults[$event.Test] = "skip"
                }
            }
        }
    }

    $exitCode = $LASTEXITCODE
    $totalTests = $testResults.Count
    $passedTests = @($testResults.GetEnumerator() | Where-Object { $_.Value -eq "pass" }).Count
    $failedTests = @($testResults.GetEnumerator() | Where-Object { $_.Value -eq "fail" }).Count
    $skippedTests = @($testResults.GetEnumerator() | Where-Object { $_.Value -eq "skip" }).Count

    Write-Host ""
    Write-Host "Integration test summary"
    Write-Host "Total tests : $totalTests"
    Write-Host "Passed      : $passedTests ($(Get-Percent -Count $passedTests -Total $totalTests)%)"
    Write-Host "Failed      : $failedTests ($(Get-Percent -Count $failedTests -Total $totalTests)%)"
    Write-Host "Skipped     : $skippedTests ($(Get-Percent -Count $skippedTests -Total $totalTests)%)"

    if ($exitCode -ne 0) {
        throw "go test failed with exit code $exitCode."
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
