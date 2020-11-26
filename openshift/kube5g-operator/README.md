# Kube5g-Operator

This current branch for the objective of suporting more config parameters for both oai-ran and oai-cn, in addition to support the functional slit for oai-ran, i.e. CU-DU fuctional split

## Requirement for Development, along with the versions that are used
- Operator SDK v0.18.1
    ```bash
    jenkins@cigarier:~$ operator-sdk version
    operator-sdk version: "v0.18.1", commit: "7bf7b6886d647dc202525daec16fab67dcc52a3d", kubernetes version: "v1.18.2", go version: "go1.13.6 linux/amd64"
    ```
- Golang 1.12+
    ```bash
    jenkins@cigarier:~$ go version
    go version go1.14.4 linux/amd64
    ```
- Docker 17.03+
    ```bash
    jenkins@cigarier:~$ docker version
    Client: Docker Engine - Community
    Version:           19.03.13
    API version:       1.40
    Go version:        go1.13.15
    Git commit:        4484c46d9d
    Built:             Wed Sep 16 17:02:36 2020
    OS/Arch:           linux/amd64
    Experimental:      false

    Server: Docker Engine - Community
    Engine:
    Version:          19.03.13
    API version:      1.40 (minimum version 1.12)
    Go version:       go1.13.15
    Git commit:       4484c46d9d
    Built:            Wed Sep 16 17:01:06 2020
    OS/Arch:          linux/amd64
    Experimental:     false
    containerd:
    Version:          1.3.7
    GitCommit:        8fba4e9a7d01810a393d5d25a3621dc101981175
    runc:
    Version:          1.0.0-rc10
    GitCommit:        dc9208a3303feef5b3839f4323d9beb36df0a9dd
    docker-init:
    Version:          0.18.0
    GitCommit:        fec3683
    ```
- kubectl v1.11.3+
    ```bash
    jenkins@cigarier:~$ kubectl version
    Client Version: version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.6", GitCommit:"dff82dc0de47299ab66c83c626e08b245ab19037", GitTreeState:"clean", BuildDate:"2020-07-15T16:58:53Z", GoVersion:"go1.13.9", Compiler:"gc", Platform:"linux/amd64"}
    Server Version: version.Info{Major:"1", Minor:"18", GitVersion:"v1.18.8", GitCommit:"9f2892aab98fe339f3bd70e3c470144299398ace", GitTreeState:"clean", BuildDate:"2020-08-13T16:04:18Z", GoVersion:"go1.13.15", Compiler:"gc", Platform:"linux/amd64"}
    ```
- a Kubernetes environment (microk8s or minikube, etc)



We will explain how to deploy and use kube5g-operator by the following example:
## Automate the deployment of OAI-LTE with kube5g-operator in K8S using microk8s

## Setup requirements
The following requirements shall be met to sucessfully deploy 4G network using docker:
* One PC with Ubuntu 16.04/18.04 with (preferable) 8GB RAM
* Kubernetes deployment microk8s
* Snap enabled
* An USRP (as frontend) attached to the PC
* Commecial phone equipped with SIM card to be connected to the network


For this tutorial we use the following:
- GigaByte Box with ubuntu 16.04 and 16 BG RAM
- microk8s version: v1.14
- kubectl version: 1.17.3
- USRP: B210 mini
- Phone: Google pixel 2

The following figure illustrates the network that we will deploy. It is composed of i) mysql ii) oai-cn iii) oai-ran.

![Fig_oai_lte](https://i.imgur.com/wDSQiza.jpg)

## kube5g setup
* Install the requirements of kube5g 
    - Either using the build_m5g that resides in [mosaic5g](https://gitlab.eurecom.fr/mosaic5g/mosaic5g) project:
        ```bash
        $ ./build_m5g -i
        ```
    - Or using the script of kube5g-operator:
        1. go to the following directory ```kube5g/openshift/kube5g-operator```
            ```bash
            $ cd kube5g/openshift/kube5g-operator
            ```
        2. Install the requirements of kube5g-operator by:
            ```bash
            $ ./k5goperator.sh -i
            ```
            You can discover the capabilities of ```k5goperator.sh``` by ```./k5goperator.sh -h```
* Make sure that everything is fine by checking the status of ```microk8s```
    
    ```bash
    $ microk8s.status 
        microk8s is running
        addons:
        rbac: disabled
        ingress: disabled
        dns: enabled
        metrics-server: disabled
        linkerd: disabled
        prometheus: disabled
        istio: disabled
        jaeger: disabled
        fluentd: disabled
        gpu: disabled
        storage: disabled
        dashboard: disabled
        registry: disabled
    ```

* Add the dns to the configuration of kubectl
    
    - get the dns of your network:
        
        ```bash
        $ nmcli device show enp0s31f6 |grep -i dns
        IP4.DNS[1]:                             192.168.1.1
        ```
        
        where ```enp0s31f6``` is the interface name connected to the internet. 
    - Add the dns of your network to the config of kubectl
        
        ```bash
        microk8s.kubectl -n kube-system edit configmap/kube-dns
        ```
        and then add the dns. Here is an example:

        ```bash
        # Please edit the object below. Lines beginning with a '#' will be ignored,
        # and an empty file will abort the edit. If an error occurs while saving this file will be
        # reopened with the relevant failures.
        #
        apiVersion: v1
        data:
        upstreamNameservers: '["192.168.1.1", "8.8.8.8", "8.8.4.4"]'
        kind: ConfigMap
        metadata:
        annotations:
            kubectl.kubernetes.io/last-applied-configuration: |
            {"apiVersion":"v1","data":{"upstreamNameservers":"[\"192.168.1.1\", \"8.8.8.8\", \"8.8.4.4\"]"},"kind":"ConfigMap","metadata":{"annotations":{},"labels":{"addonmanager.kubernetes.io/mode":"EnsureExists"},"name":"kube-dns","namespace":"kube-system"}}
        creationTimestamp: "2020-03-24T08:36:56Z"
        labels:
            addonmanager.kubernetes.io/mode: EnsureExists
        name: kube-dns
        namespace: kube-system
        resourceVersion: "8703"
        selfLink: /api/v1/namespaces/kube-system/configmaps/kube-dns
        uid: 9f1c7876-6daa-11ea-93e1-ec21e5fc4532
    ```
### Bring the network up
- Apply the Custom Resource Defintion (CRD) to k8s cluster
    
    ```bash
    $ ./k5goperator.sh -n
    ```

- Run kube5g-operator as a pod
    
    ```bash
    $ ./k5goperator.sh container start
    ```

* Now, Apply the custom resources to bring the LTE network up
    
    ```bash
    $ ./k5goperator.sh deploy v1 all-in-one
    ```
    Where this will deploy lte network with ``àll-in-one``` mode for the core network, where ```v1```indicates to the version of the snaps used insider the pods. For more information on the snap version, plese check [here](https://gitlab.eurecom.fr/mosaic5g/mosaic5g/-/wikis/tutorials).

    After some time, the USRP will be on, and now you can connect the phone to the network.

### Bring the network down
- After that, you can bring the network down by:
    
    ```bash
    $ ./k5goperator.sh -d
    ```

- If you wan to stop the kube5g-operator: 
    
    ```bash
    $ ./k5goperator.sh container stop
    ```

- If you want to remove the Custom Resource Defintion (CRD) from the k8s cluster
    
    ```bash
    $ ./k5goperator.sh -c
    ```

- If you want to clean your machine from the installed softwar:
    
    ```bash
    $ ./k5goperator.sh -r
    ```