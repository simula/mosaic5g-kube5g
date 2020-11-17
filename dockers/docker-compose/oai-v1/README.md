[<img align="center" width="40" height="40" src="https://mosaic5g.io/img/m5g-kube5g.png" />](https://gitlab.eurecom.fr/mosaic5g/kube5g/-/tree/develop/dockers/docker-compose)
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


For this tutorial we used the following:
- GigaByte Box with ubuntu 16.04 and 16 BG RAM
- Docker version: 19.03.5
- Docker-compose version: 1.8.0
- USRP: B210 mini
- Phone: Google pixel 2

The following figure illustrates the network that we will deploy, which concernes the example ```lte-all-in-one```. Note that the following instructions are also to be followed when deploying the other examples, namely ```lte-all-in-one-with-flexran``` and ```lte```. It is composed of i) mysql ii) oai-cn iii) oai-ran.

![Fig_4gNetwork](https://i.imgur.com/gFC5I8i.jpg)

### Step 1: Network configuration
- Check the config file ```oai-v1/lte-all-in-one/conf.yaml```:
  ```yaml
  oaiEnb:
    - mcc: "208"                 
      mnc: "95"
      mmeService: 
        description: |
          If 'ipv4' is provided, the enb connects to MME using the provied ip address. Otherwise, it will connect to the service provided by 'name'        
        name: "oaicn"
        snapVersion: "v1" # allowed values: v1, v2
        ipv4: ""
      flexRAN: false
      flexRANServiceName: "flexran"
      snap:
        description: |
          Define the version of the snap to be used. The field "refresh" if true, it will get the lates version of snap. Otherwise, the current version (supposing that it is already installed) will be kept
        name: "oai-ran"
        channel: "edge"
        devmode: true
        refresh: false 
        
      eutra_band:
        description: "Setting the LTE  EUTRA frequency band."
        default: "7"
      downlink_frequency:
        description: "Setting the downlink frequency band."
        default: "2685000000L"
      uplink_frequency_offset:
        description: "Setting the uplink frequency offset."
        default: "-120000000"
      N_RB_DL:
        description: |
          Setting the bandwidth in terms of number of available PRBS, 
          25 (5MHz), 50 (10MHz), and 100 (20MHz)."
        default: "25"
      parallel_config:
        description: |
          Defines the level of parallelism. There are three available values;
          'PARALLEL_SINGLE_THREAD', 'PARALLEL_RU_L1_SPLIT', or 'PARALLEL_RU_L1_TRX_SPLIT'"
        default: "PARALLEL_SINGLE_THREAD"
      max_rxgain:
        description: "defines the maximum Rx gain"  
        default: "110"

  oaiCn:
    v1:
    - realm: 
        description: "Realm of oai"
        default: "openair4G.eur"
      snap:
        description: |
          Define the version of the snap to be used
        name: "oai-cn"
        channel: "edge"
        devmode: true
        refresh: false
      oaiHss:
        databaseServiceName: "mysql"
      oaiMme:
        mcc: "208"  
        mnc: "95"
      oaiSpgw:
        dns: "192.168.1.1"
  ```

Where the important parameters to be modified:
- ```oaiEnb``` and ```oaiCn[v1][0][oaiHss]```: modify the ```mcc``` and ```mnc``` according to your deployment. 
- ```oaiEnb```: modify the following parameters: 
  - ```eutra_band```
  - ```downlink_frequency```
  - ```uplink_frequency_offset```
  
  Note that the you may also change following parameters if you want, but the eNB should work with the default values
  - ```N_RB_DL```
  - ```max_rxgain```
  - ```parallel_config```
  
- ```oaiCn[v1][0][oaiSpgw]```:
  - ```dns```: the dns of your network that you use to connect to the internet, and you can get it e.g. like the following:
    ```bash
    $ nmcli device show wlp2s0 | grep -i dns 
    IP4.DNS[1]:                             172.24.2.4
    ```
    where ```wlp2s0``` is the name of the netowkr interface that is connected to the network, which can be different for you. Then add the dns (which is ```172.24.2.4``` for our case) to the file ```conf.yaml```.

### Step 1: Network deployment
At this stage, we can now create the 4G network by running the following command:
```bash
cigarier@cigarier:~/mosaic5g/kube5g/dockers/docker-compose$ docker-compose -f oai-v1/lte-all-in-one/docker-compose.yaml up -d
```
After some time the USRP is on, and now you should be able to connect your UE to network. After that, you can stop the service by typing in the terminal ```docker-compose -f oai-v1/lte-all-in-one/docker-compose.yaml stop``` or you can delete the network by typing in the terminal ```docker-compose -f oai-v1/lte-all-in-one/docker-compose.yaml down```.


