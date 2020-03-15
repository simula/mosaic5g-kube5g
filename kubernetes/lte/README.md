# How to deploy OAI-LTE using Kubernetes

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
sudo snap install microk8s --classic
```
- Start 
```bash=2
Start microk8s
```
-  Enable dns working
```bash=3
microk8s.enable dns
```
-  install kubectl
```bash=4
sudo snap install kubectl
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
- clone this repo, then apply the manifest in oai-cn and oai-ran (in order).
```bash=15
git clone https://github.com/tig4605246/docker-oai.git
cd docker-oai
kubectl apply -f oai-cn/
kubectl apply -f oai-ran/
```
At this stage, the USRP should start functioning after certain time

In order to shutdown the network, type the following:
```bash
kubectl delete -f oai-cn/
kubectl delete -f oai-ran/
```
You may also want to stop microk8s by typing ```microk8s.stop``` in the terminal