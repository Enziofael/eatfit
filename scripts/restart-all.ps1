# scripts/restart-all.ps1
# Полный перезапуск окружения: compose down, compose up, миграции

$ErrorActionPreference = "Stop"

# Определяем корень проекта (на уровень выше scripts)
$rootDir = Split-Path -Parent $PSScriptRoot

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Eatfit - Full Environment Restart" -ForegroundColor Cyan
Write-Host "  Root: $rootDir" -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Cyan

# Проверяем наличие .env
$envFile = Join-Path $rootDir ".env"
if (-not (Test-Path $envFile)) {
    Write-Host "ERROR: .env file not found at $envFile" -ForegroundColor Red
    Write-Host "Copy .env.example to .env and fill in the values" -ForegroundColor Yellow
    exit 1
}

# 1. Останавливаем и удаляем контейнеры
Write-Host "`n[1/4] Stopping containers..." -ForegroundColor Yellow
docker compose -f "$rootDir\docker-compose.yml" --env-file $envFile down

# 2. Запускаем контейнеры
Write-Host "`n[2/4] Starting containers..." -ForegroundColor Yellow
docker compose -f "$rootDir\docker-compose.yml" --env-file $envFile up -d

# 3. Ждём готовности сервисов
Write-Host "`n[3/4] Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# Проверяем PostgreSQL
$maxRetries = 15
$retry = 0
$pgReady = $false
while ($retry -lt $maxRetries) {
    try {
        $result = docker exec eatfit-postgres pg_isready -U eatfit_user -d eatfit 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "PostgreSQL is ready!" -ForegroundColor Green
            $pgReady = $true
            break
        }
    } catch {
        # Игнорируем ошибки
    }
    $retry++
    Write-Host "Waiting for PostgreSQL... ($retry/$maxRetries)" -ForegroundColor Gray
    Start-Sleep -Seconds 2
}

if (-not $pgReady) {
    Write-Host "ERROR: PostgreSQL failed to start" -ForegroundColor Red
    docker logs eatfit-postgres
    exit 1
}

# Проверяем Redis
$retry = 0
$redisReady = $false
while ($retry -lt $maxRetries) {
    try {
        # Если пароль пустой - ping без авторизации
        if ([string]::IsNullOrEmpty($env:REDIS_PASSWORD)) {
            $result = docker exec eatfit-redis redis-cli ping 2>&1
        } else {
            $result = docker exec eatfit-redis redis-cli -a $env:REDIS_PASSWORD ping 2>&1
        }
        if ($LASTEXITCODE -eq 0 -and $result -match "PONG") {
            Write-Host "Redis is ready!" -ForegroundColor Green
            $redisReady = $true
            break
        }
    } catch {
        # Игнорируем ошибки
    }
    $retry++
    Write-Host "Waiting for Redis... ($retry/$maxRetries)" -ForegroundColor Gray
    Start-Sleep -Seconds 2
}

if (-not $redisReady) {
    Write-Host "ERROR: Redis failed to start" -ForegroundColor Red
    docker logs eatfit-redis
    exit 1
}

# 4. Запускаем миграции
Write-Host "`n[4/4] Running migrations..." -ForegroundColor Yellow
& "$rootDir\backend\scripts\migrate.ps1"

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  All services are up and running!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Services:" -ForegroundColor White
Write-Host "  PostgreSQL:  localhost:5432" -ForegroundColor Gray
Write-Host "  Redis:       localhost:6379" -ForegroundColor Gray
Write-Host "  MailHog UI:  http://localhost:8025" -ForegroundColor Gray
Write-Host ""
Write-Host "Run backend:" -ForegroundColor White
Write-Host "  cd backend && go run cmd/eatfit/main.go" -ForegroundColor Gray
Write-Host ""
Write-Host "Test gRPC:" -ForegroundColor White
Write-Host "  grpcurl -plaintext localhost:50051 list" -ForegroundColor Gray