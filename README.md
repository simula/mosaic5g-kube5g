# kube5g: Cloud-Native Agile 5G Service Platforms

## Quick Start
1. install the required dependencies of kube5g:
    ```bash
    ./build_kube5g -i
    ```
2. modify ```common/config-manager/conf_short_default.yaml``` according to your setup
3. run the configuration manager
    ```bash
     ./build_kube5g -c
    ```
4. to deploy network in kubernetes using kube5g-operator
    ```bash
    cd openshift/kube5g-operator/
    ```
    * apply the custom resource definition (crd) of the network
        ```bash
        ./k5goperator.sh -n
        ```
    * stat kube5g-operator as pod in kubernetes
        ```bash
        ./k5goperator.sh container start
        ```
    * deploy network with monolithic RAN and CN
        ```bash
        ./k5goperator.sh deploy v1 all-in-one
        ```
    
5. to deploy network using docker
    ```bash
    cd dockers/docker-compose/oai-v1/lte-all-in-one
    ```
    * deploy network with monolithic RAN and CN
        ```bash
        docker-compose up -d
        ```
    

This project includes the following
```bash
├── build_kube5g
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
- ```build_kube5g```: It helps in building kube5g, by installig the required dependencies ```./build_kube5g -i```, optional dependencies (intended for development) ```./build_kube5g -I```. It can also help in configuring kube5G ```./build_kube5g -c``` using th defaul short version of configuration file ```common/config-manager/conf_short_default.yaml``` and the global configuration file ```common/config-manager/conf_global_default.yaml```. For more information and options of ```build_kube5g```, type ```./build_kube5g -h```

- ```common```: It contains a set of common tools and scripts for kube5g project. Currently, it contains of:
    * ```config-manager```: It contains the config manager ```conf-manager.py``` and global configuration ```conf_global.yaml```. ```config-manager``` will configure automatically all the required configuraiton files for docker compose and custom resources (CRs) for kube5g-operator. Please refer to 
    ```common/config-manager/README.md``` for more information.
- ```dockers```: it includes a set of tools and scripts for building docker images of mosaic5g snaps ```v1``` and ```v2```, as well as deploying 4G/5G networks using docker and docker-compose
    * ```docker-hook```: It is to create the hook, which is the init for docker images. Through the hook, you can install the application (e.g., oai-cn snap), configure it correctly, and start it inside the dockers
    * ```docker-build```: It containes all what you need to build docker  images for Mosaic5G platforms, such as oai-ran, oai-cn, etc.
    * ```docker-compose```: This includes an easy-to-use docker-compose files along with the required configurations that you may need to change according to your setup.
- ```kubernetes```: Currently, it containes examples on how to deploy 4G/5G networks using kubernetes
- ```openshift```: Currently, it includes the ```kube5g-operator```, which is an orchestrator tool for managing 5G services in Kubernetes deployments. 
