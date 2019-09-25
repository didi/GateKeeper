#!/bin/bash

if [ $# -eq 1 ];then
    file_dir=$1
else
    echo input file_dir like service
    exit 1
fi

go test  -coverpkg github.com/didi/gatekeeper/${file_dir}/... -coverprofile=report/${file_dir}_coverage.out ./...
go tool cover -html=report/${file_dir}_coverage.out -o report/${file_dir}_coverage.html
open report/${file_dir}_coverage.html