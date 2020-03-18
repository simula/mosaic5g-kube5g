# How to deploy OAI-LTE using Kubernetes


In this tutorial, We demonstrate how to deploy OAI-LTE using Kubernetes (K8S).


## Requirements
The following requirements shall be met to sucessfully deploy 4G network using docker:
* One PC with Ubuntu 16.04/18.04 with (preferable) 8GB RAM
* Kubernetes deployment microk8s
* Snap enabled
* An USRP (as frontend) attached to the PC
* Commecial phone equipped with SIM card to be connected to the network


For this tutorial we use the following:
- GigaByte Box with ubuntu 16.04 and 16 BG RAM
- microk8s version: v1.14.10
- kubectl version: 1.17.3
- USRP: B210 mini
- Phone: Google pixel 2

The following figures illustrates the network that we will deploy. It is composed of i) mysql ii) oai-cn iii) oai-ran.

![Fig_oai_lte](https://i.imgur.com/wDSQiza.jpg)


### Kubernetes deployment
- Install microk8s
```bash=1
sudo snap install microk8s --channel=1.13/stable --classic
```
- Start microk8s 
```bash=2
microk8s.start
```
-  Enable dns
```bash=3
microk8s.enable dns
```

- Add the user to the group  ```microk8s``` to the permission
```bash=4
sudo usermod -a -G microk8s $USER
```

Note that you might need to change the dns address with ```microk8s.kubectl -n kube-system edit configmap/coredns``` in order to make it work in your environment. The default setting is 8.8.8.8 and 8.8.4.4

-  install kubectl
```bash=5
sudo snap install kubectl --classic
microk8s.kubectl config view --raw > $HOME/.kube/config
```
- Add --allow-privileged=true to both kubelet and kube-apiserver the microk8s, and then restart the services
```bash=6
# kubelet config
sudo vim /var/snap/microk8s/current/args/kubelet

#kube-apiserver config
sudo vim /var/snap/microk8s/current/args/kube-apiserver

# Restart services:
sudo systemctl restart snap.microk8s.daemon-kubelet.service
sudo systemctl restart snap.microk8s.daemon-apiserver.service
```
- Apply kubernetes services (it is supposed that you are in the root directory of mosaic5g):
```bash=15
kubectl apply -f kube5g/kubernetes/lte/mysql.yaml
kubectl apply -f kube5g/kubernetes/lte/oaicn.yaml
kubectl apply -f kube5g/kubernetes/lte/oairan.yaml
```
At this stage, the USRP should start functioning after certain time

In order to shutdown the network, type the following:
```bash
kubectl delete -f oai-cn/
kubectl delete -f oai-ran/
```
You may also want to stop microk8s by typing ```microk8s.stop``` in the terminal

