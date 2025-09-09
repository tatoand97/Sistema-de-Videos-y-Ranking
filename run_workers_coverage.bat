@echo off
echo ========================================
echo EJECUTANDO PRUEBAS DE WORKERS CON COBERTURA
echo ========================================
echo.

set REPORT_DIR=workers_coverage_reports
if exist %REPORT_DIR% rmdir /s /q %REPORT_DIR%
mkdir %REPORT_DIR%

echo [1/7] Ejecutando pruebas Workers/shared...
cd Workers\shared
go test ./... -coverprofile=..\..\%REPORT_DIR%\shared_coverage.out -covermode=atomic -v > ..\..\%REPORT_DIR%\shared_test_results.txt 2>&1
if exist ..\..\%REPORT_DIR%\shared_coverage.out (
    go tool cover -func=..\..\%REPORT_DIR%\shared_coverage.out > ..\..\%REPORT_DIR%\shared_coverage_summary.txt
    go tool cover -html=..\..\%REPORT_DIR%\shared_coverage.out -o ..\..\%REPORT_DIR%\shared_coverage.html
)
cd ..\..

echo [2/7] Ejecutando pruebas AudioRemoval...
cd Workers\AudioRemoval
go test ./... -coverprofile=..\..\%REPORT_DIR%\audioremoval_coverage.out -covermode=atomic -v > ..\..\%REPORT_DIR%\audioremoval_test_results.txt 2>&1
if exist ..\..\%REPORT_DIR%\audioremoval_coverage.out (
    go tool cover -func=..\..\%REPORT_DIR%\audioremoval_coverage.out > ..\..\%REPORT_DIR%\audioremoval_coverage_summary.txt
    go tool cover -html=..\..\%REPORT_DIR%\audioremoval_coverage.out -o ..\..\%REPORT_DIR%\audioremoval_coverage.html
)
cd ..\..

echo [3/7] Ejecutando pruebas EditVideo...
cd Workers\EditVideo
go test ./... -coverprofile=..\..\%REPORT_DIR%\editvideo_coverage.out -covermode=atomic -v > ..\..\%REPORT_DIR%\editvideo_test_results.txt 2>&1
if exist ..\..\%REPORT_DIR%\editvideo_coverage.out (
    go tool cover -func=..\..\%REPORT_DIR%\editvideo_coverage.out > ..\..\%REPORT_DIR%\editvideo_coverage_summary.txt
    go tool cover -html=..\..\%REPORT_DIR%\editvideo_coverage.out -o ..\..\%REPORT_DIR%\editvideo_coverage.html
)
cd ..\..

echo [4/7] Ejecutando pruebas GossipOpenClose...
cd Workers\gossipOpenClose
go test ./... -coverprofile=..\..\%REPORT_DIR%\gossipopenclose_coverage.out -covermode=atomic -v > ..\..\%REPORT_DIR%\gossipopenclose_test_results.txt 2>&1
if exist ..\..\%REPORT_DIR%\gossipopenclose_coverage.out (
    go tool cover -func=..\..\%REPORT_DIR%\gossipopenclose_coverage.out > ..\..\%REPORT_DIR%\gossipopenclose_coverage_summary.txt
    go tool cover -html=..\..\%REPORT_DIR%\gossipopenclose_coverage.out -o ..\..\%REPORT_DIR%\gossipopenclose_coverage.html
)
cd ..\..

echo [5/7] Ejecutando pruebas StatesMachine...
cd Workers\StatesMachine
go test ./... -coverprofile=..\..\%REPORT_DIR%\statesmachine_coverage.out -covermode=atomic -v > ..\..\%REPORT_DIR%\statesmachine_test_results.txt 2>&1
if exist ..\..\%REPORT_DIR%\statesmachine_coverage.out (
    go tool cover -func=..\..\%REPORT_DIR%\statesmachine_coverage.out > ..\..\%REPORT_DIR%\statesmachine_coverage_summary.txt
    go tool cover -html=..\..\%REPORT_DIR%\statesmachine_coverage.out -o ..\..\%REPORT_DIR%\statesmachine_coverage.html
)
cd ..\..

echo [6/7] Ejecutando pruebas TrimVideo...
cd Workers\TrimVideo
go test ./... -coverprofile=..\..\%REPORT_DIR%\trimvideo_coverage.out -covermode=atomic -v > ..\..\%REPORT_DIR%\trimvideo_test_results.txt 2>&1
if exist ..\..\%REPORT_DIR%\trimvideo_coverage.out (
    go tool cover -func=..\..\%REPORT_DIR%\trimvideo_coverage.out > ..\..\%REPORT_DIR%\trimvideo_coverage_summary.txt
    go tool cover -html=..\..\%REPORT_DIR%\trimvideo_coverage.out -o ..\..\%REPORT_DIR%\trimvideo_coverage.html
)
cd ..\..

echo [7/7] Ejecutando pruebas Watermarking...
cd Workers\Watermarking
go test ./... -coverprofile=..\..\%REPORT_DIR%\watermarking_coverage.out -covermode=atomic -v > ..\..\%REPORT_DIR%\watermarking_test_results.txt 2>&1
if exist ..\..\%REPORT_DIR%\watermarking_coverage.out (
    go tool cover -func=..\..\%REPORT_DIR%\watermarking_coverage.out > ..\..\%REPORT_DIR%\watermarking_coverage_summary.txt
    go tool cover -html=..\..\%REPORT_DIR%\watermarking_coverage.out -o ..\..\%REPORT_DIR%\watermarking_coverage.html
)
cd ..\..

echo.
echo ========================================
echo GENERANDO REPORTE CONSOLIDADO
echo ========================================

echo # REPORTE DE COBERTURA - WORKERS > %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
echo Generado: %date% %time% >> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
echo. >> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
echo ## RESUMEN POR WORKER >> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
echo. >> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md

for %%f in (%REPORT_DIR%\*_coverage_summary.txt) do (
    echo ### %%~nf >> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
    echo ```>> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
    type "%%f" >> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
    echo ```>> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
    echo. >> %REPORT_DIR%\WORKERS_COVERAGE_REPORT.md
)

echo.
echo ========================================
echo REPORTE GENERADO EXITOSAMENTE
echo ========================================
echo.
echo Archivos generados en: %REPORT_DIR%/
echo - WORKERS_COVERAGE_REPORT.md (resumen consolidado)
echo - *_coverage.html (reportes HTML por worker)
echo - *_coverage_summary.txt (resúmenes por worker)
echo - *_test_results.txt (resultados detallados)
echo.
pause