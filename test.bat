@echo off
cd Api && go test -coverprofile=coverage.out -coverpkg=./internal/... ./tests/... >nul 2>&1 && echo API: && go tool cover -func=coverage.out | findstr "total:" || echo API: 0.0%%
cd ..
cd Workers\AdminCache && go test -coverprofile=coverage.out ./... >nul 2>&1 && echo AdminCache: && go tool cover -func=coverage.out | findstr "total:" || echo AdminCache: 0.0%%
cd ..\..
cd Workers\AudioRemoval && go test -coverprofile=coverage.out -coverpkg=./internal/... ./tests/... >nul 2>&1 && echo AudioRemoval: && go tool cover -func=coverage.out | findstr "total:" || echo AudioRemoval: 0.0%%
cd ..\..
cd Workers\EditVideo && go test -coverprofile=coverage.out -coverpkg=./internal/... ./tests/... >nul 2>&1 && echo EditVideo: && go tool cover -func=coverage.out | findstr "total:" || echo EditVideo: 0.0%%
cd ..\..
cd Workers\gossipOpenClose && go test -coverprofile=coverage.out ./... >nul 2>&1 && echo GossipOpenClose: && go tool cover -func=coverage.out | findstr "total:" || echo GossipOpenClose: 0.0%%
cd ..\..
cd Workers\shared && go test -coverprofile=coverage.out ./... >nul 2>&1 && echo Shared: && go tool cover -func=coverage.out | findstr "total:" || echo Shared: 0.0%%
cd ..\..
cd Workers\StatesMachine && go test -coverprofile=coverage.out ./... >nul 2>&1 && echo StatesMachine: && go tool cover -func=coverage.out | findstr "total:" || echo StatesMachine: 0.0%%
cd ..\..
cd Workers\TrimVideo && go test -coverprofile=coverage.out ./... >nul 2>&1 && echo TrimVideo: && go tool cover -func=coverage.out | findstr "total:" || echo TrimVideo: 0.0%%
cd ..\..
cd Workers\Watermarking && go test -coverprofile=coverage.out ./... >nul 2>&1 && echo Watermarking: && go tool cover -func=coverage.out | findstr "total:" || echo Watermarking: 0.0%%