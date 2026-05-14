# backend/scripts/migrate.ps1
# Применение миграций к PostgreSQL через Docker контейнер

$ErrorActionPreference = "Stop"

# Загружаем переменные из .env
$rootDir = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$envFile = Join-Path $rootDir ".env"

if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
}

$DB_USER = if ($env:POSTGRES_USER) { $env:POSTGRES_USER } else { "eatfit_user" }
$DB_NAME = if ($env:POSTGRES_DB) { $env:POSTGRES_DB } else { "eatfit" }

Write-Host "Running migrations in container..." -ForegroundColor Yellow

# Папка с миграциями
$migrationsDir = Join-Path $PSScriptRoot "..\migrations"

if (-not (Test-Path $migrationsDir)) {
    Write-Host "Migrations directory not found: $migrationsDir" -ForegroundColor Red
    exit 1
}

# Применяем каждый SQL файл
Get-ChildItem $migrationsDir -Filter *.sql | Sort-Object Name | ForEach-Object {
    $fileName = $_.Name
    Write-Host "  Executing: $fileName" -ForegroundColor Gray
    
    # Копируем содержимое SQL файла в контейнер и выполняем
    Get-Content $_.FullName -Raw | docker exec -i eatfit-postgres psql -U $DB_USER -d $DB_NAME -v ON_ERROR_STOP=1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Migration failed: $fileName" -ForegroundColor Red
        exit 1
    }
}

Write-Host "Migrations completed successfully!" -ForegroundColor Green