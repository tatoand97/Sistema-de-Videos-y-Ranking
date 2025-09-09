@echo off
echo ========================================
echo GENERANDO REPORTE DE COBERTURA COMPLETO
echo ========================================
echo.

set REPORT_DIR=coverage_reports
if exist %REPORT_DIR% rmdir /s /q %REPORT_DIR%
mkdir %REPORT_DIR%

echo [1/8] Ejecutando pruebas API...
cd Api
go test ./... -coverprofile=../%REPORT_DIR%/api_coverage.out -covermode=atomic > ../%REPORT_DIR%/api_test_results.txt 2>&1
go tool cover -func=../%REPORT_DIR%/api_coverage.out > ../%REPORT_DIR%/api_coverage_summary.txt
go tool cover -html=../%REPORT_DIR%/api_coverage.out -o ../%REPORT_DIR%/api_coverage.html
cd ..

echo [2/8] Ejecutando pruebas Workers/shared...
cd Workers/shared
go test ./... -coverprofile=../../%REPORT_DIR%/shared_coverage.out -covermode=atomic > ../../%REPORT_DIR%/shared_test_results.txt 2>&1
if exist ../../%REPORT_DIR%/shared_coverage.out (
    go tool cover -func=../../%REPORT_DIR%/shared_coverage.out > ../../%REPORT_DIR%/shared_coverage_summary.txt
    go tool cover -html=../../%REPORT_DIR%/shared_coverage.out -o ../../%REPORT_DIR%/shared_coverage.html
)
cd ../..

echo [3/8] Ejecutando pruebas AudioRemoval...
cd Workers/AudioRemoval
go test ./... -coverprofile=../../%REPORT_DIR%/audioremoval_coverage.out -covermode=atomic > ../../%REPORT_DIR%/audioremoval_test_results.txt 2>&1
if exist ../../%REPORT_DIR%/audioremoval_coverage.out (
    go tool cover -func=../../%REPORT_DIR%/audioremoval_coverage.out > ../../%REPORT_DIR%/audioremoval_coverage_summary.txt
    go tool cover -html=../../%REPORT_DIR%/audioremoval_coverage.out -o ../../%REPORT_DIR%/audioremoval_coverage.html
)
cd ../..

echo [4/8] Ejecutando pruebas EditVideo...
cd Workers/EditVideo
go test ./... -coverprofile=../../%REPORT_DIR%/editvideo_coverage.out -covermode=atomic > ../../%REPORT_DIR%/editvideo_test_results.txt 2>&1
if exist ../../%REPORT_DIR%/editvideo_coverage.out (
    go tool cover -func=../../%REPORT_DIR%/editvideo_coverage.out > ../../%REPORT_DIR%/editvideo_coverage_summary.txt
    go tool cover -html=../../%REPORT_DIR%/editvideo_coverage.out -o ../../%REPORT_DIR%/editvideo_coverage.html
)
cd ../..

echo [5/8] Ejecutando pruebas GossipOpenClose...
cd Workers/gossipOpenClose
go test ./... -coverprofile=../../%REPORT_DIR%/gossipopenclose_coverage.out -covermode=atomic > ../../%REPORT_DIR%/gossipopenclose_test_results.txt 2>&1
if exist ../../%REPORT_DIR%/gossipopenclose_coverage.out (
    go tool cover -func=../../%REPORT_DIR%/gossipopenclose_coverage.out > ../../%REPORT_DIR%/gossipopenclose_coverage_summary.txt
    go tool cover -html=../../%REPORT_DIR%/gossipopenclose_coverage.out -o ../../%REPORT_DIR%/gossipopenclose_coverage.html
)
cd ../..

echo [6/8] Ejecutando pruebas StatesMachine...
cd Workers/StatesMachine
go test ./... -coverprofile=../../%REPORT_DIR%/statesmachine_coverage.out -covermode=atomic > ../../%REPORT_DIR%/statesmachine_test_results.txt 2>&1
if exist ../../%REPORT_DIR%/statesmachine_coverage.out (
    go tool cover -func=../../%REPORT_DIR%/statesmachine_coverage.out > ../../%REPORT_DIR%/statesmachine_coverage_summary.txt
    go tool cover -html=../../%REPORT_DIR%/statesmachine_coverage.out -o ../../%REPORT_DIR%/statesmachine_coverage.html
)
cd ../..

echo [7/8] Ejecutando pruebas TrimVideo...
cd Workers/TrimVideo
go test ./... -coverprofile=../../%REPORT_DIR%/trimvideo_coverage.out -covermode=atomic > ../../%REPORT_DIR%/trimvideo_test_results.txt 2>&1
if exist ../../%REPORT_DIR%/trimvideo_coverage.out (
    go tool cover -func=../../%REPORT_DIR%/trimvideo_coverage.out > ../../%REPORT_DIR%/trimvideo_coverage_summary.txt
    go tool cover -html=../../%REPORT_DIR%/trimvideo_coverage.out -o ../../%REPORT_DIR%/trimvideo_coverage.html
)
cd ../..

echo [8/8] Ejecutando pruebas Watermarking...
cd Workers/Watermarking
go test ./... -coverprofile=../../%REPORT_DIR%/watermarking_coverage.out -covermode=atomic > ../../%REPORT_DIR%/watermarking_test_results.txt 2>&1
if exist ../../%REPORT_DIR%/watermarking_coverage.out (
    go tool cover -func=../../%REPORT_DIR%/watermarking_coverage.out > ../../%REPORT_DIR%/watermarking_coverage_summary.txt
    go tool cover -html=../../%REPORT_DIR%/watermarking_coverage.out -o ../../%REPORT_DIR%/watermarking_coverage.html
)
cd ../..

echo.
echo ========================================
echo REPORTE GENERADO EXITOSAMENTE
echo ========================================
echo.
echo Archivos generados en: coverage_reports/
echo - *_coverage.html (reportes HTML por módulo)
echo - *_coverage_summary.txt (resúmenes por módulo)
echo - *_test_results.txt (resultados detallados)
echo.
echo Para ver los resúmenes:
dir coverage_reports\*_coverage_summary.txt
echo.
pause