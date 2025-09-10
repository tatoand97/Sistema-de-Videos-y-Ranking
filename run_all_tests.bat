@echo off
echo ========================================
echo EJECUTANDO PRUEBAS UNITARIAS DEL PROYECTO
echo ========================================

echo.
echo [1/7] API - Pruebas unitarias
echo ----------------------------------------
cd Api
go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./...
if %errorlevel% equ 0 (
    echo API - Cobertura:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo API - ERROR en las pruebas
)
cd ..

echo.
echo [2/7] StatesMachine - Pruebas unitarias
echo ----------------------------------------
cd Workers\StatesMachine
go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./...
if %errorlevel% equ 0 (
    echo StatesMachine - Cobertura:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo StatesMachine - ERROR en las pruebas
)
cd ..\..

echo.
echo [3/7] AudioRemoval - Pruebas unitarias
echo ----------------------------------------
cd Workers\AudioRemoval
go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./...
if %errorlevel% equ 0 (
    echo AudioRemoval - Cobertura:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo AudioRemoval - ERROR en las pruebas
)
cd ..\..

echo.
echo [4/7] EditVideo - Pruebas unitarias
echo ----------------------------------------
cd Workers\EditVideo
go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./...
if %errorlevel% equ 0 (
    echo EditVideo - Cobertura:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo EditVideo - ERROR en las pruebas
)
cd ..\..

echo.
echo [5/7] GossipOpenClose - Pruebas unitarias
echo ----------------------------------------
cd Workers\gossipOpenClose
go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./...
if %errorlevel% equ 0 (
    echo GossipOpenClose - Cobertura:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo GossipOpenClose - ERROR en las pruebas
)
cd ..\..

echo.
echo [6/7] TrimVideo - Pruebas unitarias
echo ----------------------------------------
cd Workers\TrimVideo
go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./...
if %errorlevel% equ 0 (
    echo TrimVideo - Cobertura:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo TrimVideo - ERROR en las pruebas
)
cd ..\..

echo.
echo [7/7] Watermarking - Pruebas unitarias
echo ----------------------------------------
cd Workers\Watermarking
go test -v ./tests/unit/... -coverprofile=coverage.out -coverpkg=./...
if %errorlevel% equ 0 (
    echo Watermarking - Cobertura:
    go tool cover -func=coverage.out | findstr "total:"
) else (
    echo Watermarking - ERROR en las pruebas
)
cd ..\..

echo.
echo ========================================
echo RESUMEN COMPLETADO
echo ========================================