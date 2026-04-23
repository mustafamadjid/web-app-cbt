param(
    [string]$DbUrl
)

$ErrorActionPreference = "Stop"

. "$PSScriptRoot/common.ps1"

$targetDbUrl = Get-DatabaseUrl -DbUrl $DbUrl -EnvVarName "POSTGRES_TEST_DBURL"
$mainDbUrl = [Environment]::GetEnvironmentVariable("POSTGRES_DBURL")

Assert-TestDatabaseUrlIsSafe -TestDbUrl $targetDbUrl -MainDbUrl $mainDbUrl
Invoke-GooseCommand -DbUrl $targetDbUrl -Commands @("reset", "up")
