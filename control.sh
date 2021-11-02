#!/bin/bash
#############################################
## main
## 以托管方式, 启动服务
## control.sh脚本, 必须实现start方法
#############################################

set -e 
workspace=$(cd $(dirname $0) && pwd -P)
cd $workspace
app=gatekeeper
cmd_run=" run -c conf/dev/ -p"

function check_go_version() {
  use_go_version="1120"
  go_version=`go version |grep -Eo '([0-9])\.([0-9]{1,2})\.([0-9]*)' |awk -F '.' '{print $1 $2 $3}'`
  if [ ! -n "$go_version" ] || [ "$go_version" -le "$use_go_version" ];then
    echo "go version < 1.12.0 pleace upgreate"
    echo "Installation package download address : https://golang.org/dl/ or https://golang.google.cn/dl/"
    echo "Installation tutorial address : https://www.runoob.com/go/go-environment.html"
    exit -1
  fi
}

function run() {
  panel_type=$1
  cmd_run="${cmd_run} ${panel_type}"

  # bin run gatekeeper
  if [ -f "${app}" ];then
    echo "run bin gatekeeper"
    chmod 755 ${app}
    eval ./${app}${cmd_run}
    return
  fi

  # go run gatekeeper
  if [ -f "main.go" ];then
    check_go_version
    echo "go run main.go ${cmd_run}"
    eval go run main.go ${cmd_run}
    return
  fi
  echo "not found run file"
}

function help() {
    echo -e "Control command manager"
    echo -e "Usage:"
    echo -e "\t[command]"
    echo -e "\n"
    echo -e "Available Commands:"
    echo -e "\tstart_proxy \t start gatekeeper proxy"
    echo -e "\tstart_control \t start gatekeeper control"
    echo -e "\tstart_both \t start gatekeeper control && proxy"
    echo -e "\n"
    echo -e "Flags:"
    echo -e "\t-h,\t--help\thelp for this command"
}

action=$1
case $action in
    -h|--help)
        help
        ;;
    "start_proxy" )
        run proxy
        ;;
    "start_control" )
        run control
        ;;
    "start_both" )
      run both
      ;;
    * )
        help
        ;;
esac
