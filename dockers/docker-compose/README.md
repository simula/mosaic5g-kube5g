---
title: 'How to deploy 4G network using docker compose'
---
###### tags: `mosaic5g`

# How to deploy 4G network using docker compose

In this tutorial, We demonstrate how to deploy 4G network using docker and docker-compose.


## Requirements
The following requirements shall be met to sucessfully deploy 4G network using docker:
* One PC with Ubuntu 16.04/18.04 with (preferable) 8GB RAM
* Snap enabled
* docker-compose
* An USRP (as frontend) attached to the PC
* Commecial phone equipped with SIM card to be connected to the network


For this tutorial we use the following:
- GigaByte Box with ubuntu 16.04 and 16 BG RAM
- Docker version: 19.03.5
- Docker-compose version: 1.8.0
- USRP: B210 mini
- Phone: Google pixel 2

The following figures illustrates the network that we will deploy, which concerne the example ```lte```. Note that the following instructions are also to be followed when deployign the example ```lte-with-flexran```. It is composed of i) mysql ii) oai-cn iii) oai-ran.

![Fig_4gNetwork](https://i.imgur.com/gFC5I8i.jpg)


### Step 1: Network configuration
- Check the config file ```conf.yaml```:
  ```yaml
  mcc: "208"                 
  mnc: "95"   

  eutraBand: "7"             
  downlinkFrequency: "2685000000L"    
  uplinkFrequencyOffset: "-120000000"

  NumberRbDl: "25"
  MaxRxGain: "110"
  ParallelConfig: "PARALLEL_SINGLE_THREAD"

  configurationPathofCN: "/var/snap/oai-cn/current/"
  configurationPathofRAN: "/var/snap/oai-ran/current/"
  snapBinaryPath: "/snap/bin/"
  hssDomainName: "oaicn"
  mmeDomainName: "oaicn"
  spgwDomainName: "oaicn"
  mysqlDomainName: "mysql"
  dns: "172.24.2.4"

  flexRAN: false
  flexRANDomainName: "flexran"
  test: false
  ```

and modify the ```mcc``` and ```mnc``` according to your deployment. Note that the following parameters are related to oai-ran and they need to be modified based on your equipment: ```NumberRbDl```, ```downlinkFrequency```, ```uplinkFrequencyOffset```, ```uplinkFrequencyOffset```, ```MaxRxGain```, and ```ParallelConfig```. Additionally, you need also to change the ```dns``` according to the network that you use to connect to the internet, and you can get it e.g. like the following:
```bash
$ nmcli device show wlp2s0 | grep -i dns 
IP4.DNS[1]:                             172.24.2.4
```
where ```wlp2s0``` is the name of the netowkr interface that is connected to the network, which can be different for you. Then add the dns (which is ```172.24.2.4``` for our case) to the file ```conf.yaml```.

### Step 1: Network deployment
At this stage, we can now create the 4G network by running the following command:
```bash
cigarier@cigarier:~/mosaic5g/kube5g/dockers/docker-compose$ docker-compose -f lte/docker-compose.yaml up -d
```
After some time (generally less than one minute) the USRP is on, and now you should be able to connect your UE to network. After that, you can stop the service by typing in the terminal ```docker-compose -f lte/docker-compose.yaml stop``` or you can delete the network by typing in the terminal ```docker-compose -f lte/docker-compose.yaml down```.


