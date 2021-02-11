#! /bin/bash
# ################################################################################
# * Copyright 2019-2020 Eurecom and Mosaic5G Platforms Authors
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
# *-------------------------------------------------------------------------------
# file          snap-docker-build.jenkins
# brief         Build snap versions of Mosaic5G using different versions of ubuntu (e.g., 16.04 and 18.04) using
# contact       admin@mosaic-5g.io
# authors:
# 	- Robert Schmidt (robert.schmidt@eurecom.fr)
# *-------------------------------------------------------------------------------

host=localhost
port=9999

# get any temp file
tmpf=$(mktemp)

curl -sX GET http://$host:$port/capabilities > $tmpf

if [ $? -ne 0 ]; then
  echo "Could not check capabilities"
  exit 100
fi
# cat $tmpf |jq '.'

name=$(jq '.info.title' < $tmpf)
if [ "$name" != "\"FlexRAN NB API\"" ]; then
  echo "/capabilities broken: .info.title is $name"
  exit 101
fi

# verify that nothing is connected
curl -sX GET http://$host:$port/stats > $tmpf
# cat $tmpf 
# echo ""
date_time=$(jq '.date_time' < $tmpf)
eNB_config=$(jq '.eNB_config' < $tmpf)
eNB_config_n=$(echo $eNB_config | jq '. | length')
mac_stats=$(jq '.mac_stats' < $tmpf)
mac_stats_n=$(echo $mac_stats | jq '. | length')

if [ "$eNB_config" != "[]" ] || [ $eNB_config_n -ne 0 ]; then
  echo "/stats broken: found \"eNB_config: $eNB_config\", but should be [] and length 0!"
  exit 102
fi

if [ "$mac_stats" != "[]" ] || [ $mac_stats_n -ne 0 ]; then
  echo "/stats broken: found \"mac_stats: $mac_stats\", but should be [] and length 0!"
  exit 103
fi

rm $tmpf
echo "Unit Test SUCCESS: controller reports date_time $date_time"
exit 0




