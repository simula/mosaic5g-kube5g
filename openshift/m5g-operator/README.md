# Mosaic5g-Operator

This current branch for the objective of suporting more config parameters for both oai-ran and oai-cn, in addition to support the functional slit for oai-ran, i.e. CU-DU fuctional split
## Requirement for Development

- Operator SDK v0.7.0
- Golang 1.12+
- Docker 17.03+
- kubectl v1.11.3+
- a Kubernetes environment (microk8s or minikube, etc)
- Optional: [dep][https://golang.github.io/dep/docs/installation.html] version v0.5.0+.
- Optional: [delve](https://github.com/go-delve/delve/tree/master/Documentation/installation) version 1.2.0+ (for `up local --enable-delve`).


We will explain how to deploy and use mosaic5g-operator by the following example:
## Automate the deployment of OAI-LTE with mosaic5g-operator in K8S using microk8s

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
* Build kube5g using build_m5g script:
    ```bash
    $ ./build_m5g -k
    ```

* Install the requirements of kube5g 
    - Either using the build_m5g script:
        ```bash
        $ ./build_m5g -i
        ```
    - Or using the script of mosaic5g-operator:
        1. go to the following directory ```mosaic5g_DIR/kube5g/openshift/m5g-operator``` (it is assumed that you are already in the directory ```mosaic5g_DIR```)
            ```bash
            $ cd kube5g/openshift/m5g-operator
            ```
        2. Install the requirements of mosaic5g-operator by:
            ```bash
            $ ./m5goperator.sh -i
            ```
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
        IP4.DNS[1]:                             192.168.106.12
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
        upstreamNameservers: '["192.168.106.12", "8.8.8.8", "8.8.4.4"]'
        kind: ConfigMap
        metadata:
        annotations:
            kubectl.kubernetes.io/last-applied-configuration: |
            {"apiVersion":"v1","data":{"upstreamNameservers":"[\"192.168.106.12\", \"8.8.8.8\", \"8.8.4.4\"]"},"kind":"ConfigMap","metadata":{"annotations":{},"labels":{"addonmanager.kubernetes.io/mode":"EnsureExists"},"name":"kube-dns","namespace":"kube-system"}}
        creationTimestamp: "2020-03-24T08:36:56Z"
        labels:
            addonmanager.kubernetes.io/mode: EnsureExists
        name: kube-dns
        namespace: kube-system
        resourceVersion: "8703"
        selfLink: /api/v1/namespaces/kube-system/configmaps/kube-dns
        uid: 9f1c7876-6daa-11ea-93e1-ec21e5fc4532
    ```
* Run mosaic5g-operator as a pod in kubernetes
    - Go to ```mosaic5g_DIR/kube5g/openshift/m5g-operator```:
        
        ```bash
        $ cd kube5g/openshift/m5g-operator
        ```

    - You can discover the capability provided by the script ```m5g-operator.sh```
        
        ```bash
        $ ./m5goperator.sh 
            This program installs the requirements to run kubernets on one machine, 
            and run mosaic5g-operator as a pod inside kubernetes in order to manage
            the deployments and services of 4G/5G networks in cloud native environment.
            This program also allows to run mosaic5g-operator locally as Golang app.
            Options:
            -i | --install
                Install and run microk8s kubectl, then deploy operator on it"
            -n | --init
                Apply CRD to k8s cluster (Required for Operator)"
            -l | --local
                Run Operator as a Golang app at local"
            container [start|stop]
                Run Operator as a POD inside Kubernetes"
            -c | --clean 
                Remove CRD from cluster"
            -r | --remove
                remove the snap of kubectl and microk8s
            Usage:
            ./m5goperator.sh -i 
            ./m5goperator.sh container start
        ```
    
### Bring the network up
- Apply the Custom Resource Defintion (CRD) to k8s cluster
    
    ```bash
    $ ./m5goperator.sh -n
    ```

- Run mosaic5g-operator as a pod
    
    ```bash
    $ ./m5goperator.sh container start
    ```

* Now, Apply the custom resources to bring the LTE network up
    
    ```bash
    $ kubectl apply -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_cr.yaml
    ```

    After some time, the USRP will be on, and now you can connect the phone to the network.

### Bring the network down
- After that, you can bring the network down by:
    
    ```bash
    $ kubectl delete -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_cr.yaml
    ```

- If you wan to stop the mosaic5g-operator: 
    
    ```bash
    $ ./m5goperator.sh container stop
    ```

- If you want to remove the Custom Resource Defintion (CRD) from the k8s cluster
    
    ```bash
    $ ./m5goperator.sh -c
    ```

- If you want to clean your machine from the installed softwar:
    
    ```bash
    $ ./m5goperator.sh -r
    ```