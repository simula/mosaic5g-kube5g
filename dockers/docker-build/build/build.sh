#! /bin/bash
# ################################################################################
# * Copyright 2016-2019 Eurecom and Mosaic5G Platforms Authors
# * Licensed to the Mosaic5G under one or more contributor license
# * agreements. See the NOTICE file distributed with this
# * work for additional information regarding copyright ownership.
# * The Mosaic5G licenses this file to You under the
# * Apache License, Version 2.0  (the "License");
# * you may not use this file except in compliance with the License.
# * You may obtain a copy of the License at
# *
# *      http://www.apache.org/licenses/LICENSE-2.0
# *
# * Unless required by applicable law or agreed to in writing, software
# * distributed under the License is distributed on an "AS IS" BASIS,
# * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# * See the License for the specific language governing permissions and
# * limitations under the License.
# ################################################################################
# #-------------------------------------------------------------------------------
# For more information about Mosaic5G:
#                                   admin@mosaic-5g.io
# file          build.sh
# brief 		Build docker images for mosaic5g snaps v1 and v2
# authors:
# 	- Osama Arouk (arouk@eurecom.fr)
# 	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
# *-------------------------------------------------------------------------------

REPO_NAME="mosaic5gecosys"  # dockerhub repository. Change it to your repository
TAG_BASE="base"             # The tag for the base image
BASE_CONTAINER="build_base" # The name of the temporary container
RELEASE_TAG="latest"        # Default release tag
SNAP_VERSION="v1"           #snap version: allowed values; v1, v2. For more info: https://gitlab.eurecom.fr/mosaic5g/mosaic5g/-/wikis/tutorials
DIR=""
DOCKER_HOOK_DIR="$HOME/go/src/mosaic5g/docker-hook/cmd/hook" # source-code of docker-hook, if you would like to build the docker-hook

# List of supported snaps
declare -a snap_list=("oai-ran oai-cn oai-hss oai-mme oai-spgw oai-spgwc oai-spgwu flexran")
declare -a snap_version_list=("v1 v2")


# contains(string, substring)
# Returns 0 if the specified string contains the specified substring, otherwise returns 1.
contains() {
    string="$1"
    substring="$2"
    if echo "$string" | $(type -p ggrep grep | head -1) -F -- "$substring" >/dev/null; then
        return 0    # $substring is in $string
    else
        return 1    # $substring is not in $string
    fi
}

# Build hook to update the change
build_hook(){
    export GOPATH=$HOME/go
    echo "build hook from source"
    CURRENT_DIR=`pwd`
    cd $DOCKER_HOOK_DIR
    go build -o hook main.go
    mv ./hook ${CURRENT_DIR}/
    cd $CURRENT_DIR
}

# Set variables
init() {
    TARGET="${REPO_NAME}/${TARGET_NAME}"
}

# Recreate base image
build_base(){
    cd ../${DIR}/
    cp ../build/hook ./
    cp ../build/conf.yaml ./
    docker build -t ${TARGET}:${TAG_BASE} --force-rm=true --rm=true .  |& tee build.log
    clean_up
}

# Build the target image
build_target(){
    list_include_item  "$snap_list" $1
    [[ $? -ne 0 ]] && echo "Error: Snap name \"$1\" not recognized" && echo "Allowed values are: $snap_list" && return $?
    list_include_item  "$snap_version_list" $SNAP_VERSION
    [[ $? -ne 0 ]] && echo "Error: Snap version \"${SNAP_VERSION}\" not recognized" && echo "Allowed values are: $snap_version_list" && return $?
    
    init
    build_base $1
    docker run --name=${BASE_CONTAINER} -ti --privileged -v /proc:/writable-proc -v /sys/fs/cgroup:/sys/fs/cgroup:ro -v /lib/modules:/lib/modules:ro -h ubuntu -d ${TARGET}:${TAG_BASE}
    RET=1
    echo "Waiting for the snaps to be installed insider dockers..."
    while  [ ${RET} -ne 0 ] ;
    do
        sleep 3
        LIST=`docker exec ${BASE_CONTAINER} snap list`
        if [ "${SNAP_VERSION}" = "v1" ] ; then
            if [ "${1}" = "oai-hss" ] || [ "${1}" = "oai-mme" ] || [ "${1}" = "oai-spgw" ] ; then
                echo "Waiting for snap oai-cn to be installed..."
                contains "${LIST}" "oai-cn"
            else
                echo "Waiting for snap ${1} to be installed..."
                contains "${LIST}" "${1}"
            fi
        elif [ "${SNAP_VERSION}" = "v2" ] ; then
            
            if [ "${1}" = "oai-hss" ] || [ "${1}" = "oai-mme" ] || [ "${1}" = "oai-spgwc" ] || [ "${1}" = "oai-spgwu" ] || [ "${1}" = "oai-ran" ] ; then
                echo "Waiting for snap ${1} to be installed..."
                contains "${LIST}" "${1}"
            elif [ "${1}" = "oai-cn" ] ; then
                echo "Waiting for snaps oai-hss, oai-mme, oai-spgwc, oai-spgwu to be installed..."
                (contains "${LIST}" "oai-hss") && (contains "${LIST}" "oai-mme") && (contains "${LIST}" "oai-spgwc") && (contains "${LIST}" "oai-spgwu")
            else
                echo "Error: the snap ${1} is not recognized, exit ..."
                exit 0
            fi
        else
            echo "Error: the snap version ${SNAP_VERSION} is not supported, exit ..."
            exit 0
        fi
        RET=$?
        
    done

    # Wait until the hook inside docker image finish
    cmd="tail -n 1 /root/hook.log"
    EXIT_STAT="End of hook"
    RET=1
    while  [ ${RET} -ne 0 ] ;
    do
        echo "Waiting until the docker-hook inside docker image finish ..."
        sleep 1
        DOCK_EXEC_OUT=$(docker exec ${BASE_CONTAINER} $cmd); echo $DOCK_EXEC_OUT
        if [[ $DOCK_EXEC_OUT == *"$EXIT_STAT"* ]]; then
            RET=0
        fi
    done
    
    
    echo "copying init_deploy.sh to docker"
    cmd="mv /root/init_deploy.sh /root/init.sh"
    docker exec ${BASE_CONTAINER} $cmd

    docker commit ${BASE_CONTAINER} ${TARGET}:${RELEASE_TAG}
    docker stop ${BASE_CONTAINER}
    docker container rm ${BASE_CONTAINER} -f
    docker image prune -f
    echo "Now ${TARGET}:${RELEASE_TAG} is ready"
    echo "All done, please use docker push ${TARGET}:${RELEASE_TAG} to push image to your repository"
}

function list_include_item {
  local list="$1"
  local item="$2"
  if [[ $list =~ (^|[[:space:]])"$item"($|[[:space:]]) ]] ; then
    # yes, list include item
    result=0
  else
    result=1
  fi
  return $result
}

clean_up(){
    rm hook
    rm conf.yaml
}

clean_all(){
    docker stop ${BASE_CONTAINER}
    docker container rm ${BASE_CONTAINER} -f
    docker image prune -f
}

main() {
    if [ "${2}" != "" ] ; then
        RELEASE_TAG=${2}
    fi
    if [ "${3}" != "" ] ; then
        SNAP_VERSION=${3}
    fi
    case ${1} in
        oai-cn)
            if [ "${SNAP_VERSION}" = "v1" ] ; then
                DIR="oai-cn"
            else
                DIR="oai-cn-v2"
            fi
            TARGET_NAME="oaicn"
            
            build_target ${1}
        ;;
        oai-hss)
            if [ "${SNAP_VERSION}" = "v1" ] ; then
                DIR="oai-hss"
            else
                DIR="oai-hss-v2"
            fi
            TARGET_NAME="oaihss"
            build_target ${1}
        ;;
        oai-mme)
            if [ "${SNAP_VERSION}" = "v1" ] ; then
                DIR="oai-mme"
            else
                DIR="oai-mme-v2"
            fi
            TARGET_NAME="oaimme"
            build_target ${1}
        ;;
        oai-spgw)
            DIR="oai-spgw"
            TARGET_NAME="oaispgw"
            build_target ${1}
        ;;
        oai-spgwc)
            DIR="oai-spgwc-v2"
            TARGET_NAME="oaispgwc"
            build_target ${1}
        ;;
        oai-spgwu)
            DIR="oai-spgwu-v2"
            TARGET_NAME="oaispgwu"
            build_target ${1}
        ;;
        oai-ran)
            DIR="oai-ran"
            TARGET_NAME="oairan"
            build_target ${1}
        ;;
        oai-gnb)
            DIR="oai-gnb"
            TARGET_NAME="oaignb"
            build_target "oai-ran" # oai-ran enb and oai-ran gnb have the same snap that is oai-ran
        ;;
        flexran)
            DIR="flexran"
            TARGET_NAME="flexran"
            build_target ${1}
        ;;
        ll-mec)
            DIR="ll-mec"
            TARGET_NAME="llmec"
            build_target ${1}
        ;;
        build-hook)
            build_hook
        ;;
        clean-all)
            clean_all
        ;;
        stop)
            stop
        ;;
        *)
            echo '
Description:
This Script will remove the old docker snap image and build a new one
Usage:
        ./build.sh [oai-cn|oai-hss|oai-mme|oai-spgw|oai-ran|flexran|ll-mec] [release tag(default is latest)] [snap version(default is v1. alowed values: v1, v2)]
Example:
        ./build.sh oai-cn mytestv1 v1
'
            exit 0
        ;;
    esac
    
    
}
main "$@"
