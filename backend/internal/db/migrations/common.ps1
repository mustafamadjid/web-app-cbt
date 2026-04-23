$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$script:MigrationsDir = $PSScriptRoot

function Get-DatabaseUrl {
    param(
        [string]$DbUrl,
        [Parameter(Mandatory = $true)]
        [string]$EnvVarName
    )

    if (-not [string]::IsNullOrWhiteSpace($DbUrl)) {
        return $DbUrl.Trim()
    }

    $envValue = [Environment]::GetEnvironmentVariable($EnvVarName)
    if ([string]::IsNullOrWhiteSpace($envValue)) {
        throw "Database URL not provided. Pass -DbUrl or set $EnvVarName."
    }

    return $envValue.Trim()
}

function Assert-TestDatabaseUrlIsSafe {
    param(
        [Parameter(Mandatory = $true)]
        [string]$TestDbUrl,
        [string]$MainDbUrl
    )

    if ([string]::IsNullOrWhiteSpace($TestDbUrl)) {
        throw "POSTGRES_TEST_DBURL must be set before running test database scripts."
    }

    if (-not [string]::IsNullOrWhiteSpace($MainDbUrl) -and $TestDbUrl.Trim() -eq $MainDbUrl.Trim()) {
        throw "POSTGRES_TEST_DBURL must not be the same as POSTGRES_DBURL."
    }
}

function Get-GooseExecutable {
    $goose = Get-Command goose -ErrorAction SilentlyContinue
    if (-not $goose) {
        throw "goose command not found in PATH."
    }

    return $goose.Source
}

function Invoke-GooseCommand {
    param(
        [Parameter(Mandatory = $true)]
        [string]$DbUrl,
        [Parameter(Mandatory = $true)]
        [string[]]$Commands
    )

    $goose = Get-GooseExecutable

    Write-Host "Using database: $DbUrl"

    foreach ($command in $Commands) {
        Write-Host "Running: goose -dir `"$script:MigrationsDir`" postgres <db-url> $command"
        & $goose -dir $script:MigrationsDir postgres $DbUrl $command
        if ($LASTEXITCODE -ne 0) {
            throw "goose command failed with exit code $LASTEXITCODE while running '$command'."
        }
    }
}
