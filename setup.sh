#!/bin/bash



# init color
RED='\E[1;31m'
GREEN='\E[1;32m'
YELOW='\E[1;33m'
BLUE='\E[1;34m'
PINK='\E[1;35m'
GREEN_SHAN='\E[5;32;49;1m'
RES='\E[0m'

# set msg
function info_msg() {
    message=$1
    printf "${GREEN}[INFO]${RES}\t${message}\n"
}

function warn_msg() {
    message=$1
    printf "${YELOW}[WARNING]${RES}\t${message}\n"
}

function error_msg() {
    message=$1
    printf "${RED}[ERROR]${RES}\t${message}\n"
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
release_url="https://download.fastgit.org/didi/GateKeeper/releases/download/v${version}/${release_file}"


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
  cd ${gatekeeper_install_dir} &&  chmod 755 install && ./install
  if [ $? != 0 ]; then
    error_msg "init conf error"
    exit -1
  fi

  printf "you can sh ${gatekeeper_dir}/control.sh [ ${YELOW}start_proxy || start_control${RES}]\n"
  printf "demo: [ ${GREEN_SHAN}cd ${install_dir} && sh control.sh start_proxy${RES} ] start gatekeeper proxy\n"
  printf "or you can run gatekeeper binary file [ ${YELOW}gatekeeper${RES} ]\n"
  printf "demo: [ ${GREEN_SHAN}cd ${install_dir} && ./gatekeeper run -c conf/dev/ -p proxy${RES} ]"
}

setup




