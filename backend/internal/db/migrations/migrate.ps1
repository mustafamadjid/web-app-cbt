param(
    [string]$DbUrl
)

$ErrorActionPreference = "Stop"

. "$PSScriptRoot/common.ps1"

$targetDbUrl = Get-DatabaseUrl -DbUrl $DbUrl -EnvVarName "POSTGRES_DBURL"

Invoke-GooseCommand -DbUrl $targetDbUrl -Commands @("up")
