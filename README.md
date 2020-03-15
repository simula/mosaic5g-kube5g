# kube5g: Cloud-Native Agile 5G Service Platforms

This project includes the following
```bash
├── dockers
│   ├── docker-build
│   ├── docker-compose
│   └── docker-hook
├── kubernetes
│   └── lte
├── openshift
│   └── m5g-operator
└── README.md
```

1. ```dockers```: it includes a set of tools and scripts for building docker images of OAI, as well as deploying 4G/5G networks using docker and docker-compose
    * ```docker-hook```: It is to create the hook, which is the init for docker images. Through the hook, we can install the application (e.g., oai-cn snap), configure it correctly, and start it
    * ```docker-build```: It containes all what you need to build docker  images for Mosaic5G platforms, such as oai-ran, oai-cn, etc.
    * ```docker-compose```: This includes an easy-to-use docker-compose files along with the required configurations that you may need to change according to your setup.
2. ```kubernetes```: Currently, it containes examples on how to deploy 4G/5G networks using kubernetes
2. ```openshift```: Currently, it includes the m5g-operator, which is an orchestrator tool for managing 5G services in Kubernetes deployments. 
