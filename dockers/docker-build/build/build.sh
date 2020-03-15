#! /bin/bash
################################################################################
# Licensed to the Mosaic5G under one or more contributor license
# agreements. See the NOTICE file distributed with this
# work for additional information regarding copyright ownership.
# The Mosaic5G licenses this file to You under the
# Apache License, Version 2.0  (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#  
#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
# -------------------------------------------------------------------------------
#   For more information about the:
#       
#
#
################################################################################
# file          build_snap_docker.sh
# brief         Build and renew the local image  
# author        Kevin Hsu (C) - 2019 hsuh@eurecom.fr

# Information of the image
REPO_NAME="mosaic5gecosys" # Change it to your repository
TARGET="${REPO_NAME}/${TARGET_NAME}" # The name of our image
TAG_BASE="base" # The tag for the base image
BASE_CONTAINER="build_base" # The name of the temporary container
RELEASE_TAG="latest" # Default release tag
DIR=""

# contains(string, substring)
#
# Returns 0 if the specified string contains the specified substring,
# otherwise returns 1.
contains() {
    string="$1"
    substring="$2"

    if echo "$string" | $(type -p ggrep grep | head -1) -F -- "$substring" >/dev/null; then
        return 0    # $substring is in $string
    else
        return 1    # $substring is not in $string
    fi
}

# Rebuild hook to update the change
build_hook(){
    echo "build hook from source"
    NOW=`pwd`
    cd ${GOPATH}/src/oai-snap-in-docker/cmd/hook/
    go build 
    mv ./hook ${NOW}/
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
    init
    build_base
    docker run --name=${BASE_CONTAINER} -ti --privileged -v /proc:/writable-proc -v /sys/fs/cgroup:/sys/fs/cgroup:ro -v /lib/modules:/lib/modules:ro -h ubuntu -d ${TARGET}:${TAG_BASE}
    RET=1
    echo "Installing snaps..."
    while  [ ${RET} -ne 0 ] ;
    do
        sleep 5
        LIST=`docker exec ${BASE_CONTAINER} snap list`
        echo "Waiting for snap to be installed..."
        contains "${LIST}" "${1}"
        RET=$?
        
    done
    sleep 5
    docker commit ${BASE_CONTAINER} ${TARGET}:${RELEASE_TAG}
    docker stop ${BASE_CONTAINER}
    docker container rm ${BASE_CONTAINER} -f
    docker image prune -f
    echo "Now ${TARGET}:${RELEASE_TAG} is ready"
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
    RELEASE_TAG=${2}
    case ${1} in
        oai-cn)
            DIR="oai-cn"
            TARGET_NAME="oaicn"
            build_target ${1}
        ;;
        oai-ran)
            DIR="oai-ran"
            TARGET_NAME="oairan"
            build_target ${1}
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
            exit 0
        ;;
        clean-all)
            clean_all
        ;;
        stop)
            stop
        ;;
        *)
            echo "Description:"
            echo "This Script will remove the old docker snap image and build a new one"
            echo "tested with 16.04 Ubuntu"
            echo "./build_snap_docker.sh [oai-cn|oai-ran|flexran|ll-mec] [release tag(default is latest)]"
            echo "Example: ./build_snap_docker.sh oai-cn mytest"
            exit 0
        ;;
    esac
    echo "All done, please use docker push [IMAGE NAME]:[TAG] to push image to your repository"

}
main ${1} ${2}