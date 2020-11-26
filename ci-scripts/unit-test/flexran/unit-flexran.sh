#! /bin/bash

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




