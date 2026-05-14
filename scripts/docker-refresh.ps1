# docker-refresh.ps1

Write-Host "==================================="
Write-Host " DOCKER REFRESH " -BackgroundColor Cyan -ForegroundColor Black
Write-Host "-----------------------------------"

# Остановка и удаление контейнеров
Write-Host "[1/2]" -NoNewline -BackgroundColor Cyan
Write-Host " docker-compose down -v... " -ForegroundColor Cyan -NoNewline
docker-compose down -v *>$null

if ($LASTEXITCODE -eq 0) {
    Write-Host "OK" -ForegroundColor Green
}
else {
    Write-Host "ERROR" -BackgroundColor Red
}

# Запуск новых контейнеров
Write-Host "[2/2]" -NoNewline -BackgroundColor Cyan
Write-Host "  docker-compose up -d... " -ForegroundColor Cyan -NoNewline
docker-compose up -d *>$null

if ($LASTEXITCODE -eq 0) {
    Write-Host "OK" -ForegroundColor Green
}
else {
    Write-Host "ERROR" -BackgroundColor Red
}