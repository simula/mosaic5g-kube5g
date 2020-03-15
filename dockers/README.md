# kube5g: Cloud-Native Agile 5G Service Platforms

This project includes the following
```bash
.
├── docker-build
│   ├── build
│   ├── flexran
│   ├── ll-mec
│   ├── oai-cn
│   ├── oai-ran
│   └── README.md
├── docker-compose
│   ├── lte
│   └── lte-with-flexran
├── docker-hook
│   ├── cmd
│   ├── docker
│   ├── internal
│   └── README.md
└── README.md
```

* ```docker-hook```: It is to create the hook, which is the init for docker images. Through the hook, we can install the application (e.g., oai-cn snap), configure it correctly, and start it
* ```docker-build```: It containes all what you need to build docker  images for Mosaic5G platforms, such as oai-ran, oai-cn, etc.
* ```docker-compose```: This includes an easy-to-use docker-compose files along with the required configurations that you may need to change according to your setup.
