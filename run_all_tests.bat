@echo off
echo ========================================
echo EJECUTANDO PRUEBAS UNITARIAS - SISTEMA DE VIDEOS Y RANKING
echo ========================================

echo.
echo [1/9] API - Pruebas Unitarias
echo ----------------------------------------
cd Api
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura API:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para API
)
cd ..

echo.
echo [2/9] AdminCache Worker - Pruebas Unitarias
echo ----------------------------------------
cd Workers\AdminCache
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura AdminCache:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para AdminCache
)
cd ..\..

echo.
echo [3/9] AudioRemoval Worker - Pruebas Unitarias
echo ----------------------------------------
cd Workers\AudioRemoval
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura AudioRemoval:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para AudioRemoval
)
cd ..\..

echo.
echo [4/9] EditVideo Worker - Pruebas Unitarias
echo ----------------------------------------
cd Workers\EditVideo
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura EditVideo:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para EditVideo
)
cd ..\..

echo.
echo [5/9] GossipOpenClose Worker - Pruebas Unitarias
echo ----------------------------------------
cd Workers\gossipOpenClose
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura GossipOpenClose:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para GossipOpenClose
)
cd ..\..

echo.
echo [6/9] Shared Workers - Pruebas Unitarias
echo ----------------------------------------
cd Workers\shared
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura Shared:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para Shared
)
cd ..\..

echo.
echo [7/9] StatesMachine Worker - Pruebas Unitarias
echo ----------------------------------------
cd Workers\StatesMachine
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura StatesMachine:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para StatesMachine
)
cd ..\..

echo.
echo [8/9] TrimVideo Worker - Pruebas Unitarias
echo ----------------------------------------
cd Workers\TrimVideo
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura TrimVideo:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para TrimVideo
)
cd ..\..

echo.
echo [9/9] Watermarking Worker - Pruebas Unitarias
echo ----------------------------------------
cd Workers\Watermarking
go test -coverprofile=coverage.out ./... 2>nul
if exist coverage.out (
    echo Cobertura Watermarking:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo No se pudo generar cobertura para Watermarking
)
cd ..\..

echo.
echo ========================================
echo RESUMEN COMPLETADO
echo ========================================