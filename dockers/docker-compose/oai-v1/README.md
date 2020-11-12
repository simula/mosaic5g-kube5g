# How to deploy 4G network using docker compose

In the current directory, we provided three examples to run mosaic5g docker containers
- lte-all-in-one: this directory contains docker-compose, along with the configuration file. This docker compose will create a networks that is composed of i) mysql, ii) oaicn, iii) oairan
- lte-all-in-one-with-flexran: it is similar to ```lte-all-in-one```, but with flexran as ran controller additionally
- lte: Different from the first example, this dokcer-compose will create three docker containers for the core network (cn), namely, one for ```oaihss```, one for ```oaimme```, and one for ```oaispgw```

In this tutorial, We demonstrate how to deploy 4G network using docker and docker-compose (the first example explained above).


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

The following figure illustrates the network that we will deploy, which concernes the example ```lte-all-in-one```. Note that the following instructions are also to be followed when deploying the other examples, namely ```lte-all-in-one-with-flexran``` and ```lte```. It is composed of i) mysql ii) oai-cn iii) oai-ran.

![Fig_4gNetwork](https://i.imgur.com/gFC5I8i.jpg)


### Step 1: Network configuration
- Check the config file ```conf.yaml```:
  ```yaml
  mcc: "208"                 
  mnc: "95"   

  eutraBand: "7"             
  downlinkFrequency: "2660000000L"    
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
cigarier@cigarier:~/mosaic5g/kube5g/dockers/docker-compose$ docker-compose -f lte-all-in-one/docker-compose.yaml up -d
```
After some time (generally less than one minute) the USRP is on, and now you should be able to connect your UE to network. After that, you can stop the service by typing in the terminal ```docker-compose -f lte-all-in-one/docker-compose.yaml stop``` or you can delete the network by typing in the terminal ```docker-compose -f lte-all-in-one/docker-compose.yaml down```.


