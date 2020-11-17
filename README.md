# kube5g: Cloud-Native Agile 5G Service Platforms

This project includes the following
```bash
├── common
│   └── config-manager
├── dockers
│   ├── docker-build
│   ├── docker-compose
│   └── docker-hook
├── kubernetes
│   └── lte
├── openshift
│   └── kube5g-operator
└── README.md
```

- ```common```: It contains a set of common tools and scripts for kube5g project. Currently, it contains of:
    * ```config-manager```: It contains the config manager ```conf-manager.py``` and global configuration ```conf_global.yaml```. ```config-manager``` will configure automatically all the required configuraiton files for docker compose and custom resources (CRs) for kube5g-operator. Please refer to 
    ```common/config-manager/README.md``` for more information.
- ```dockers```: it includes a set of tools and scripts for building docker images of mosaic5g snaps ```v1``` and ```v2```, as well as deploying 4G/5G networks using docker and docker-compose
    * ```docker-hook```: It is to create the hook, which is the init for docker images. Through the hook, you can install the application (e.g., oai-cn snap), configure it correctly, and start it inside the dockers
    * ```docker-build```: It containes all what you need to build docker  images for Mosaic5G platforms, such as oai-ran, oai-cn, etc.
    * ```docker-compose```: This includes an easy-to-use docker-compose files along with the required configurations that you may need to change according to your setup.
- ```kubernetes```: Currently, it containes examples on how to deploy 4G/5G networks using kubernetes
- ```openshift```: Currently, it includes the ```kube5g-operator```, which is an orchestrator tool for managing 5G services in Kubernetes deployments. 
