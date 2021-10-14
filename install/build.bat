@echo off
color 0d

echo =================================
echo ==========Build Linux ======
echo =================================
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
echo now the CGO_ENABLED:
 go env CGO_ENABLED
echo now the GOOS:
 go env GOOS
echo now the GOARCH:
 go env GOARCH
go build -o bin/install_linux main.go

echo =================================
echo ==========Build Mac ======
echo =================================
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
echo now the CGO_ENABLED:
 go env CGO_ENABLED
echo now the GOOS:
 go env GOOS
echo now the GOARCH:
 go env GOARCH
go build -o bin/install_mac main.go

echo =================================
echo ==========Build Windows ======
echo =================================
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
echo now the CGO_ENABLED:
 go env CGO_ENABLED
echo now the GOOS:
 go env GOOS
echo now the GOARCH:
 go env GOARCH
go build -o bin/install_windows.exe main.go


