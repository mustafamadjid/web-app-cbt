[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$expectedScenarios = 495
$backendDirectory = Join-Path $PSScriptRoot "backend"

if (-not (Test-Path -LiteralPath (Join-Path $backendDirectory "go.mod"))) {
    throw "File backend/go.mod tidak ditemukan dari lokasi skrip: $PSScriptRoot"
}

Push-Location $backendDirectory

try {
    $packages = @(
        go list "./internal/core/service/..." |
            Where-Object { $_ -notmatch "/integration_test$" }
    )

    if ($LASTEXITCODE -ne 0) {
        throw "Gagal mendapatkan daftar package service."
    }

    if ($packages.Count -eq 0) {
        throw "Tidak ada package unit test service yang ditemukan."
    }

    Write-Host "Menjalankan unit test service dan menghitung $expectedScenarios skenario terdaftar..."

    $jsonLines = @(
        # Semua unit test perlu dijalankan karena beberapa skenario target
        # berada di dalam subtest bertingkat.
        go test -count=1 -json $packages
    )
    $testExitCode = $LASTEXITCODE

    $events = @(
        $jsonLines |
            Where-Object { $_ } |
            ConvertFrom-Json
    )

    $scenarios = @(
        $events |
            Where-Object {
                $_.Test -match '/(Path_[0-9]+_->_|P[0-9]+_->_|Branch_[0-9]+_->_)' -and
                $_.Action -in @("pass", "fail", "skip")
            }
    )

    $failedScenarios = @(
        $scenarios |
            Where-Object Action -eq "fail" |
            Sort-Object Package, Test
    )

    $passed = @($scenarios | Where-Object Action -eq "pass").Count
    $failed = $failedScenarios.Count
    $skipped = @($scenarios | Where-Object Action -eq "skip").Count
    $total = $scenarios.Count

    $percentage = if ($expectedScenarios -gt 0) {
        [math]::Round(($passed / $expectedScenarios) * 100, 2)
    }
    else {
        0
    }

    Write-Host ""
    Write-Host "Hasil unit test service"
    Write-Host "Passed     : $passed"
    Write-Host "Failed     : $failed"
    Write-Host "Skipped    : $skipped"
    Write-Host "Terdeteksi : $total"
    Write-Host "Target     : $expectedScenarios"
    Write-Host "Persentase : $percentage%"

    if ($failedScenarios.Count -gt 0) {
        Write-Host ""
        Write-Host "Skenario yang gagal:"

        foreach ($scenario in $failedScenarios) {
            Write-Host "- $($scenario.Package) :: $($scenario.Test)"
        }
    }

    if ($total -ne $expectedScenarios) {
        Write-Warning "Jumlah skenario tidak sesuai: ditemukan $total dari target $expectedScenarios."
    }

    if ($testExitCode -ne 0 -or $failed -gt 0 -or $total -ne $expectedScenarios) {
        exit 1
    }

    exit 0
}
finally {
    Pop-Location
}
