# kube5g: Cloud-Native Agile 5G Service Platforms

This repository includes the following

```bash
.
├── conf_global.yaml
├── conf-manager.py
└── README.md
```

- ```conf-manager.py```: ```config-manager``` takes ```conf_global.yaml``` as input (default mode), configure all the required configuraiton files for docker compose and custom resources (CRs) for kube5g-operator. More specifically, these are the files that will be configured by ```conf-manager.py```: 
    ```bash
    /home/cigarier/mosaic5g/kube5g/openshift/kube5g-operator/deploy/crds/cr-v1/lte-all-in-one/mosaic5g_v1alpha1_cr_v1_lte_all_in_one.yaml
    /home/cigarier/mosaic5g/kube5g/openshift/kube5g-operator/deploy/crds/cr-v1/lte/mosaic5g_v1alpha1_cr_v1_lte.yaml
    /home/cigarier/mosaic5g/kube5g/openshift/kube5g-operator/deploy/crds/cr-v2/lte-all-in-one/mosaic5g_v1alpha1_cr_v2_lte_all_in_one.yaml
    /home/cigarier/mosaic5g/kube5g/openshift/kube5g-operator/deploy/crds/cr-v2/lte/mosaic5g_v1alpha1_cr_v2_lte.yaml

    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v1/lte-all-in-one/conf.yaml
    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v1/lte/conf.yaml
    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v2/lte-all-in-one/conf.yaml
    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v2/lte/conf.yaml

    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v1/lte-all-in-one/docker-compose.yaml
    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v1/lte/docker-compose.yaml
    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v2/lte-all-in-one/docker-compose.yaml
    /home/cigarier/mosaic5g/kube5g/dockers/docker-compose/oai-v2/lte/docker-compose.yaml
    ```
    If you would lile to configure ```kube5g``` using ```conf_global.yaml```, all what you need is
    ```bash
    ./conf-manager.py
    ```
    or with the option ```./conf-manager.py -c file.yaml```.
    Note that ```file.yaml```can be only the file name of the file exists in the current directory, otherwise, the full path where the file exists is needed.
- ```conf_global.yaml```: It includes all the configurations required for both docker compose and kube5g-operator. For an easy deployment, you are required to verify only the following parameters:
    - ```mnc```: mobile network code, it is required for ```oaiEnb``` ```oaiCn``` and ```oaiMme```
    - ```mcc```: mobile country code, it is required for ```oaiEnb``` ```oaiCn``` and ```oaiMme```
    - the following parameters are required for ```oaiEnb```:
        - ```eutra_band```
        - ```downlink_frequency```
        - ```uplink_frequency_offset```    
        Note that you may also change following parameters if you want, but the eNB should work with the default values
        - ```N_RB_DL```
        - ```max_rxgain```
        - ```parallel_config```    
    - ```dns```: the dns of your network that you use to connect to the internet, and you can get it e.g. like the following:
        ```bash
        $ nmcli device show wlp2s0 | grep -i dns 
        IP4.DNS[1]:                             172.24.2.4
        ```
        where ```wlp2s0``` is the name of the netowkr interface that is connected to the network, which can be different for you. Then modify the dns ```dns```, which is ```172.24.2.4``` for our case.