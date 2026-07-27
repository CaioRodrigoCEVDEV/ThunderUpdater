@echo off
chcp 65001 >nul

setlocal

set PROJECT_DIR=%~dp0
cd /d "%PROJECT_DIR%"

where goversioninfo >nul 2>nul
if %errorlevel% neq 0 (
    echo Instalando goversioninfo...
    go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
    for /f "tokens=*" %%i in ('go env GOPATH') do set "PATH=%%i\bin;%PATH%"
) else (
    echo goversioninfo encontrado.
)

if not exist dist mkdir dist

echo ========================================
echo   ThunderUpdaterGO - Build All
echo ========================================
echo.

echo === Build Windows x86 (32-bit) ===
echo Gerando resource.syso para 386...
goversioninfo -64=false -icon=assets\icon.ico -o cmd\thunderupdater\resource_windows_386.syso cmd\thunderupdater\versioninfo.json
if %errorlevel% neq 0 (
    echo ERRO: Falha ao gerar resource.syso para 386
    exit /b %errorlevel%
)

set GOARCH=386
go build -o dist\ThunderUpdater-x86.exe .\cmd\thunderupdater
set BUILD_ERR=%errorlevel%
del /f cmd\thunderupdater\resource_windows_386.syso >nul 2>&1
if %BUILD_ERR% neq 0 (
    echo ERRO: Falha ao compilar x86
    exit /b %BUILD_ERR%
)
echo   dist\ThunderUpdater-x86.exe
echo.

echo === Build Windows x64 (64-bit) ===
echo Gerando resource.syso para amd64...
goversioninfo -icon=assets\icon.ico -o cmd\thunderupdater\resource_windows_amd64.syso cmd\thunderupdater\versioninfo.json
if %errorlevel% neq 0 (
    echo ERRO: Falha ao gerar resource.syso para amd64
    exit /b %errorlevel%
)

set GOARCH=amd64
go build -o dist\ThunderUpdater-x64.exe .\cmd\thunderupdater
set BUILD_ERR=%errorlevel%
del /f cmd\thunderupdater\resource_windows_amd64.syso >nul 2>&1
if %BUILD_ERR% neq 0 (
    echo ERRO: Falha ao compilar x64
    exit /b %BUILD_ERR%
)
echo   dist\ThunderUpdater-x64.exe
echo.

echo Build concluído com sucesso!

endlocal
