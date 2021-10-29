#!/bin/bash

function buildLinux() {
    echo =================================
    echo ==========Build Linux ======
    echo =================================
    CGO_ENABLED=0
    GOOS=linux
    GOARCH=amd64
    echo now the CGO_ENABLED:
    go env CGO_ENABLED
    echo now the GOOS:
    go env GOOS
    echo now the GOARCH:
    go env GOARCH
    go build -o bin/install_linux main.go
}


function buildMac() {
    echo =================================
    echo ==========Build Mac ======
    echo =================================
    CGO_ENABLED=0
    GOOS=darwin
    GOARCH=amd64
    echo now the CGO_ENABLED:
    go env CGO_ENABLED
    echo now the GOOS:
    go env GOOS
    echo now the GOARCH:
    go env GOARCH
    go build -o bin/install_mac main.go
}

function buildWindows() {
    echo =================================
    echo ==========Build Windows ======
    echo =================================
    CGO_ENABLED=1
    GOOS=windows
    GOARCH=amd64
    echo now the CGO_ENABLED:
    go env CGO_ENABLED
    echo now the GOOS:
    go env GOOS
    echo now the GOARCH:
    go env GOARCH
    go build -o bin/install_windows.exe main.go
}



buildLinux
buildMac
buildWindows
