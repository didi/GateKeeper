#!/bin/bash
export PATH=$GOROOT/bin:$PATH
workspace=$(cd $(dirname $0) && pwd -P)
cd $workspace

## const
app=`cat modulename`
function vertifyBuildPath() {
    export GO11MODULE=on
    export GO111MODULE=on
    export GOPROXY=https://goproxy.io
}

## function
function build() {
    # 进行编译
    go build -o ./bin/$app main.go
    local sc=$?
    if [ $sc -ne 0 ];then
        ## 编译失败, 退出码为 非0
        echo "$app build error"
        exit $sc
    else
        echo -e "$app build ok"
    fi
}

function make_output() {
    # 新建output目录
    local output="./output"
    rm -rf $output &>/dev/null
    mkdir -p $output/bin &>/dev/null

    # 填充output目录, output 内的内容 即为 线上部署内容
    (
        cp -rf ./bin/$app $output/bin/ &&
        cp -rf ./conf $output/ &&
        cp -rf ./control.sh $output/ &&
        cp -rf ./modulename $output &&
        cp -rf ./tmpl $output &&
        echo -e "make output ok"
    ) || { echo -e "make output error"; rm -rf "./output"; exit 2; } # 填充output目录失败后, 退出码为 非0
}

##########################################
## main
## 其中,
##      1.进行编译
##      2.生成部署包output
##########################################

vertifyBuildPath

# 1.进行编译
build

# 2.生成部署包output
make_output

# 编译成功
echo -e "build done"
exit 0
