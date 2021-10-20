#! bin/bash


# check go version

function checkGoVerion() {
  useGoVersion="1120"
  goVersion=`go version |grep -Eo '([0-9])\.([0-9]{1,2})\.([0-9]*)' |awk -F '.' '{print $1 $2 $3}'`
  if [ ! -n "$goVersion" ] || [ "$goVersion" -le "$useGoVersion" ];then
    echo "go version < 1.12.0 pleace upgreate"
    echo "Installation package download address : https://golang.org/dl/ or https://golang.google.cn/dl/"
    echo "Installation tutorial address : https://www.runoob.com/go/go-environment.html"
    exit -1
  else
    initInstall
  fi
}


function initInstall() {
  go run main.go
}


checkGoVerion