#! /bin/bash

rm -rf /home/cigarier/mosaic5g/kube5g/dockers/docker-build/build/hook
rm -rf /home/cigarier/go/src/docker-hook/cmd/hook/hook
go build -o  /home/cigarier/go/src/docker-hook/cmd/hook/hook /home/cigarier/go/src/docker-hook/cmd/hook/main.go
cp /home/cigarier/mosaic5g/kube5g/dockers/docker-hook/cmd/hook/hook /home/cigarier/mosaic5g/kube5g/dockers/docker-build/build/

docker images rm --force mosaic5gecosys/oaihss:mytest
docker images rm --force mosaic5gecosys/oaimme:mytest
docker images rm --force mosaic5gecosys/oaispgw:mytest
docker images rm --force mosaic5gecosys/oairan:mytest


./build.sh oai-hss mytest
./build.sh oai-mme mytest
./build.sh oai-spgw mytest
./build.sh oai-ran mytest

docker-compose -f /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/lte/docker-compose.yaml up -d
