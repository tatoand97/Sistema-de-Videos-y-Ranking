@echo off
echo === COBERTURA DE CODIGO POR PROYECTO ===
echo.

echo [API]
cd Api
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..

echo.
echo [AdminCache]
cd Workers\AdminCache
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..

echo.
echo [AudioRemoval]
cd Workers\AudioRemoval
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..

echo.
echo [EditVideo]
cd Workers\EditVideo
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..

echo.
echo [GossipOpenClose]
cd Workers\gossipOpenClose
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..

echo.
echo [StatesMachine]
cd Workers\StatesMachine
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..

echo.
echo [TrimVideo]
cd Workers\TrimVideo
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..

echo.
echo [Watermarking]
cd Workers\Watermarking
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..

echo.
echo [Shared]
cd Workers\shared
go test -coverprofile=coverage.out ./... 2>nul && go tool cover -func=coverage.out | findstr "total:" || echo ERROR: Fallo en compilacion
cd ..\..