<!-- # v1.1: Move config path to /root/config/conf.yaml. v1.2: A copy of 1.1. It's for upgrading demo. -->
<!-- 1.0: conf.yaml should be mounted at /root/conf.yaml ; 1.1: conf.yaml should be mounted at /root/config/conf.yaml ; 1.2: Image == 1.1, just a tag for demonstarting upgrade. -->


## WARNING: Read Before Using any Script 

This container is with **security** options **disabled**, this is an unsupported setup, if you have multiple snap packages inside the same container they will be able to break out of the confinement and see each others data and processes. **Do not rely on security inside the container**.

## Extra packages installed in this image
- apt-utils
- dnsutils
- net-tools
- iputils-ping
- vim
- virt-what

## version 1.0:


The following parameters are exposed to be configured in the current version of oaicn:
   
```yaml
mcc: "208"                 
mnc: "95"   

eutraBand: "7"             
downlinkFrequency: "2685000000L"    
uplinkFrequencyOffset: "-120000000"
NumberRbDl: "25"
MaxRxGain: "110"
ParallelConfig: "PARALLEL_SINGLE_THREAD"

hssDomainName: "oaicn"
mmeDomainName: "oaicn"
spgwDomainName: "oaicn"
mysqlDomainName: "mysql"
dns: "192.168.1.1"
```

- ```mcc```: Mobile Country Code
- ```mnc```: Mobile Network Code

- ```eutraBand```: EUTRA band, along with the downlink frequency ```downlinkFrequency```, and uplink frequency offset ```uplinkFrequencyOffset```
- ```NumberRbDl```: Number of Physical Resourc Blocks
- ```MaxRxGain```: Maximum Rx gain
- ```ParallelConfig```: three config for level of parallelism "PARALLEL_SINGLE_THREAD", "PARALLEL_RU_L1_SPLIT", or "PARALLEL_RU_L1_TRX_SPLIT"
- ```hssDomainName```: name of docker container for the service hss
- ```mmeDomainName```: name of docker container for the service mme
- ```spgwDomainName```: name of docker container for the service spgw
- ```mysqlDomainName```: name of docker container for the service mysql
Note that for the all-in-one deployment, the parameters ```hssDomainName```, ```mmeDomainName```, and ```spgwDomainName``` must be the same, which is the name of the docker running oai-cn.

For more information about these parameters and more, pleases visit: 
- [openairinterface5g](https://gitlab.eurecom.fr/oai/openairinterface5g)
- [mosaic5G](https://gitlab.eurecom.fr/mosaic5g/mosaic5g)


### Example Usage
1. Create docker-docmpose (```docker-compose.yaml```) file with the following content

    ```yaml
    version: '2'

    services:
      oaicn: # Domain name of container
        image: mosaic5gecosys/oaicn:1.0
        restart: always # Operation Policy
        container_name: oaicn # Name of the container
        hostname: ubuntu # hostname
        privileged: true # Give the container the permission to manipulate the host
        depends_on: # Before starting this container, what should be ready
          - "mysql"
        volumes: # Mounted from host
          - /sys/fs/cgroup:/sys/fs/cgroup:ro
          - /lib/modules:/lib/modules:ro
          - ./conf.yaml:/root/config/conf.yaml:ro
        networks: # the network this container attached to
          - oai

      oairan:
        image: mosaic5gecosys/oairan:1.0
        restart: always
        container_name: oairan
        hostname: oairan
        privileged: true
        depends_on:
          - "oaicn"
        volumes:
          - /sys/fs/cgroup:/sys/fs/cgroup:ro
          - /lib/modules:/lib/modules:ro
          - /dev/bus/usb:/dev/bus/usb
          - ./conf.yaml:/root/config/conf.yaml:ro
        networks:
          - oai

      mysql:
        image: mysql:5.6
        restart: always
        container_name: mysql
        environment: # shell variables
          - MYSQL_ROOT_PASSWORD=linux
        networks:
          - oai

    networks: # Define our network here
      oai:
        driver: bridge
        driver_opts:
          com.docker.network.driver.mtu: 9000 # Configure mtu
    ```
    
2. Create config file (```conf.yaml```) file with the following content

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
    dns: "192.168.1.1"

    flexRAN: false
    flexRANDomainName: "flexran"
    test: false
    ```

Please change the above parameters according to your setup.
3. deploy the network by typing in the terminal ```docker-compose up -d```

## version 1.1:


## version 1.2:

