#!/bin/sh
set -x
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
go build -o bin/gatekeeper
rm -rf ./logs/*
ps aux | grep gatekeeper | grep -v 'grep' | awk '{print $2}' | xargs kill

action=$1
case $action in
    "control" )
        ./bin/gatekeeper -c ./conf/dev/ -p $action
        ;;
    "proxy" )
        ./bin/gatekeeper -c ./conf/dev/ -p $action
        ;;
    * )
        echo "unknown command"
        exit 1
        ;;
esac