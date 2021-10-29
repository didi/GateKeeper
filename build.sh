#! bin/bash

function build_linux() {
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

    go build -o "${build_path}/${bin_file}" main.go
}


function build_mac() {
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

    go build -o "${build_path}/${bin_file}" main.go
}

function build_windows() {
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

    bin_file="${bin_file}.exe"
    go build -o "${build_path}/${bin_file}" main.go
}


function build() {

    # get build os
    build_os

    if [ -d "build" ];then
      rm -rf build
    fi

    build_workspace="${workspace}/build"
    build_dir_name="gatekeeper-${version}-"

#    build_dir_name="${workspace}/build/"

    for os in ${arr_build_os[*]}
    do
      build_path="${build_workspace}/${build_dir_name}${os}"
      echo $build_path

      echo "create build dir ${build_path}"
      mkdir -p "${build_path}/install"
      cp -r conf    "${build_path}"
      cp -r dist    "${build_path}"
      cp control.sh "${build_path}"
      cp gatekeeper.sql "${build_path}"

      echo ==========Build Gatekeeper======
      bin_file="gatekeeper"
      build_${os}


      echo ==========Build Gatekeeper install======
      cd install
      bin_file="install"
      build_path="${build_path}/install"
      build_${os}
      cd -

      echo "pack ${build_dir_name}${os}"
      cd ${build_workspace}
      tar -zcf "${build_dir_name}${os}.tar.gz" ${build_dir_name}${os}
      cd -
    done
}

function build_os() {
  arr_build_os=("windows" "mac" "linux")
}

function main() {

  # input git hub release version(1.0.1)
  echo -n  "Input gatekeeper release version(1.0.1):"
  read -a version

  version=`echo $version |grep -Eo '^[0-9]+.[0-9]+.[0-9]+$'`
  if [ ! -n "$version" ];then
     echo "input version error demo: (1.0.1)"
     exit -1
  fi
  #set -x


  set -e

  # set go proxy
  export GO111MODULE=on && export GOSUMDB=off
  export GOPROXY=https://goproxy.cn


  # set output dir path
  workspace=`pwd`

  # build bin
  build
}

main

