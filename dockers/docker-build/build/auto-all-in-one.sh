#! /bin/bash

rm -rf /home/cigarier/mosaic5g/kube5g/dockers/docker-build/build/hook
rm -rf /home/cigarier/go/src/docker-hook/cmd/hook/hook
go build -o  /home/cigarier/go/src/docker-hook/cmd/hook/hook /home/cigarier/go/src/docker-hook/cmd/hook/main.go
cp /home/cigarier/mosaic5g/kube5g/dockers/docker-hook/cmd/hook/hook /home/cigarier/mosaic5g/kube5g/dockers/docker-build/build/

docker images rm --force mosaic5gecosys/oaicn:mytest
docker images rm --force mosaic5gecosys/oairan:mytest

./build.sh oai-cn mytest
./build.sh oai-ran mytest

docker-compose -f /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/lte-all-in-one/docker-compose.yaml up -d