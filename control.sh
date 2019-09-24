#!/bin/bash
#############################################
## main
## 以托管方式, 启动服务
## control.sh脚本, 必须实现start方法
#############################################
workspace=$(cd $(dirname $0) && pwd -P)
cd $workspace
module=`cat modulename`
app=$module

action=$1
case $action in
    "start" )
        exec "./bin/${app}" --config=./conf/dev/
        ;;
    * )
        echo "unknown command"
        exit 1
        ;;
esac