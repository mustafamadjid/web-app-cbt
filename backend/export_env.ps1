$env:ENVIRONMENT="prod"
$env:BASE_URL="http://localhost:8080"
$env:UPLOAD_DIR="./public/uploads"
$env:POSTGRES_DBURL="postgres://postgres:mustafa100304@localhost:5432/smafi_cbt_db?sslmode=disable"
$env:ISSUER="web-app-cbt"
$env:ACCESS_TOKEN_SECRET="Nitqv9-DFKwT7MxXF077oHc8A5OFFt3BDV7qr3fE0voJLpWrDXkF6FlUlderC0UpaZtZPWYQz5OZE5RQku1gQA"
$env:REFRESH_TOKEN_SECRET="8BofFULi2o9sB0M3_YnzknZB6SdlzIPFMtA8FCiCxrc23MW2m_upgwcDtDVNE4kNfw2xr6sPAWDxgDeHc57w0Q"
$env:TRUSTED_ORIGINS="http://localhost:3000,http://localhost:5173"

# . ./export_env.ps1