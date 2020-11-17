# kube5g: Cloud-Native Agile 5G Service Platforms

This project includes the following
```bash
.
├── docker-build
│   ├── build
│   ├── flexran
│   ├── ll-mec
│   ├── oai-cn
│   ├── oai-cn-v2
│   ├── oai-hss
│   ├── oai-hss-v2
│   ├── oai-mme
│   ├── oai-mme-v2
│   ├── oai-ran
│   ├── oai-spgw
│   ├── oai-spgwc-v2
│   ├── oai-spgwu-v2
│   └── README.md
├── docker-compose
│   ├── oai-v1
│   ├── oai-v2
│   └── README.md
├── docker-hook
│   ├── cmd
│   ├── docker
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   └── README.md
└── README.md
```

* ```docker-hook```: It is to create the hook, which is the init for docker images. Through the hook, we can install the application (e.g., oai-cn snap), configure it correctly, and start it inside docker
* ```docker-build```: It containes all what you need to build docker images for Mosaic5G platforms, such as oai-ran, oai-cn, etc., for both versions of snaps ```v1``` and ```v2```
* ```docker-compose```: This includes an easy-to-use docker-compose files along with the required configurations that you may need to change (using config-manager that resides in common top directory) according to your setup. In ```docker-compose```, there are many examples forthe versions of snaps ```v1``` and ```v2```
