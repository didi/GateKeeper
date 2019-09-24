#!/bin/bash
if [ $# -eq 1 ];then
    http_port=$1
else
    echo need http_port like ./reload.sh 8081
    exit 1
fi

# 重新载入配置文件
curl "http://127.0.0.1:${http_port}/reload"
echo ""