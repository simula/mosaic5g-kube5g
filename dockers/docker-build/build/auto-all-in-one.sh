#! /bin/bash

rm -rf $HOME/mosaic5g/kube5g/dockers/docker-build/build/hook
rm -rf $HOME/go/src/docker-hook/cmd/hook/hook
go build -o  $HOME/go/src/docker-hook/cmd/hook/hook $HOME/go/src/docker-hook/cmd/hook/main.go
cp $HOME/mosaic5g/kube5g/dockers/docker-hook/cmd/hook/hook $HOME/mosaic5g/kube5g/dockers/docker-build/build/

docker images rm --force mosaic5gecosys/oaicn:mytest
docker images rm --force mosaic5gecosys/oairan:mytest

./build.sh oai-cn mytest
./build.sh oai-ran mytest

docker-compose -f $HOME/mosaic5g/kube5g/dockers/docker-compose/lte-all-in-one/docker-compose.yaml up -d