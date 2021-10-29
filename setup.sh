#!/bin/bash

# set msg
function info_msg() {
    message=$1
    printf "\033[;32m[INFO]\033[0m\t${message}\n"
}

function warn_msg() {
    message=$1
    printf "\033[;33m[WARNING]\033[0m\t${message}\n"
}

function error_msg() {
    message=$1
    printf "\033[;31m[ERROR]\033[0m\t${message}\n"
}

# check os
uname=`uname -s`
if [[ $uname = "Darwin" ]];then
  gatekeeper_os="mac"
elif [[ "$uname" =~ "MINGW" ]]; then
  gatekeeper_os="windows"
else
  gatekeeper_os="linux"
fi

if [[ $? != 0 ]];then
  error_msg "get os error"
  exit -1
fi

info_msg "you system: ${gatekeeper_os}"


# pack setup url
version="1.0.0"
gatekeeper_dir="gatekeeper-${version}-${gatekeeper_os}"
pageage_type="tar.gz"
release_file="${gatekeeper_dir}.${pageage_type}"
release_url="https://github.com/didi/GateKeeper/releases/download/v${version}/${release_file}"


function setup() {
  # get gatekeeper from github release
  info_msg "get: ${gatekeeper_dir} from ${release_url}"
  eval "curl -L '${release_url}' -o ${release_file}"
  if [[ $? != 0 ]];then
    error_msg "get gatekeeper from [${release_url}] error"
    exit -1
  fi


  # unpack gatekeeper to gatekeeper_dir
  info_msg "unpack ${gatekeeper_dir}"
  tar -xf ${release_file}
  if [[ $? != 0 ]];then
    error_msg "unpack ${release_file} error"
    exit -1
  fi

  # remove release_file
  info_msg "remove: ${release_file}"
  rm ${release_file}
  install
}


function install() {
  # init gatekeeper
  workspace=`pwd`
  install_dir="${workspace}/${gatekeeper_dir}"

  info_msg "gatekeeper install dir: ${install_dir}"

  cd ${install_dir} && chmod 755 gatekeeper

  # init gatekeeper conf
  gatekeeper_install_dir="${install_dir}/install"
  cd ${gatekeeper_install_dir} &&  chmod 755 install.* && ./install

  printf "you can sh ${gatekeeper_dir}/control.sh [start_proxy || start_control]"
  printf "demo: [ cd ${install_dir} && sh control.sh start_proxy ] start gatekeeper proxy"
  printf "or you can run gatekeeper binary file [ gatekeeper ]"
  printf "demo: [ cd ${install_dir} && ./gatekeeper run -c conf/dev/ -p proxy]"
}

setup




