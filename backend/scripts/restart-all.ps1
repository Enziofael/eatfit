# scripts/restart-all.ps1
# Полный перезапуск окружения: compose down, compose up, миграции

$ErrorActionPreference = "Stop"
$rootDir = Split-Path -Parent $PSScriptRoot

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Eatfit - Full Environment Restart" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# 1. Останавливаем и удаляем контейнеры
Write-Host "`n[1/4] Stopping containers..." -ForegroundColor Yellow
docker compose -f "$rootDir\docker-compose.yml" --env-file "$rootDir\.env" down -v

# 2. Запускаем контейнеры
Write-Host "`n[2/4] Starting containers..." -ForegroundColor Yellow
docker compose -f "$rootDir\docker-compose.yml" --env-file "$rootDir\.env" up -d

# 3. Ждём готовности сервисов
Write-Host "`n[3/4] Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Проверяем PostgreSQL
$maxRetries = 10
$retry = 0
while ($retry -lt $maxRetries) {
    try {
        docker exec eatfit-postgres pg_isready -U eatfit_user -d eatfit
        if ($LASTEXITCODE -eq 0) {
            Write-Host "PostgreSQL is ready!" -ForegroundColor Green
            break
        }
    } catch {
        # Игнорируем ошибки
    }
    $retry++
    Write-Host "Waiting for PostgreSQL... ($retry/$maxRetries)" -ForegroundColor Gray
    Start-Sleep -Seconds 2
}

# Проверяем Redis
$retry = 0
while ($retry -lt $maxRetries) {
    try {
        docker exec eatfit-redis redis-cli ping
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Redis is ready!" -ForegroundColor Green
            break
        }
    } catch {
        # Игнорируем ошибки
    }
    $retry++
    Write-Host "Waiting for Redis... ($retry/$maxRetries)" -ForegroundColor Gray
    Start-Sleep -Seconds 2
}

# 4. Запускаем миграции
Write-Host "`n[4/4] Running migrations..." -ForegroundColor Yellow
& "$PSScriptRoot\migrate.ps1"

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  All services are up and running!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Services:" -ForegroundColor White
Write-Host "  PostgreSQL:  localhost:5432" -ForegroundColor Gray
Write-Host "  Redis:       localhost:6379" -ForegroundColor Gray
Write-Host "  MailHog UI:  http://localhost:8025" -ForegroundColor Gray
Write-Host ""
Write-Host "Run backend:  cd backend && go run cmd/eatfit/main.go" -ForegroundColor White
Write-Host "Test gRPC:    grpcurl -plaintext localhost:50051 list" -ForegroundColor White
