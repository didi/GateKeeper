#!/bin/bash
#############################################
## main
## 以托管方式, 启动服务
## control.sh脚本, 必须实现start方法
#############################################
workspace=$(cd $(dirname $0) && pwd -P)
cd $workspace
app=gatekeeper

action=$1
case $action in
    "start" )
        ps aux | grep gatekeeper | grep -v 'grep' | awk '{print $2}' | xargs kill -9
        nohup ./bin/gatekeeper -config=./conf/dev/ >> /home/webroot/logs/gatekeeper/gatekeeper.log &
        ;;
    * )
        echo "unknown command"
        exit 1
        ;;
esac