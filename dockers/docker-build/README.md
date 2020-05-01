# Container with Mosiac5g Snaps

This directory includes the build materials for building containers with mosaic5g snaps.
## WARNING: Read Before Using any Script

- This creates a set of containers with **security** options **disabled**, this is an unsupported setup, if you have multiple snap packages inside the same container they will be able to break out of the confinement and see each others data and processes. **Do not rely on security inside the container**.
<!-- - The scripts are tested and works fine in our environment. The models we used for testing is *GIGABYTE BRIX GB-BRi7-8550* and *Dell XPS-15*.  We ARE NOT sure they won't cause trouble in other environments. For more details, please read the known issue section.
- For the details of containers, please read the individual README in their foldes. -->

## Requirement
- Ubuntu 16.04/18.04
- Docker-ce 18.09
- Docker Compose
- Golang 1.10+ (Optional; If you want to rebuild the hook)

**Please make sure that docker can be run with non-root user.**

**Please note that for building the docker containers, we created the hook using the projet ```mosaic5g/kube5g/dockers/docker-hook```.**

## Quick Start

### Run Pre-set Service
<!-- 
In **compose** directory, we provide docker-compose files that can bring up Mosaic5g services without configuring. Just `cd` to your desired service directory and run `docker-compose up -d`. For example, to start an OAI lte service,

1. Go to the lte folder `cd compose/lte`
2. Check if the parameters in `conf.yaml` meet your need
3. Run `docker-compose up -d`
4. The services will start running when ready -->

### How to build docker images for mosaic5G snaps?
- Go to the build folder:
  ```bash
  $ cd build
  ```
  Among the files exist there, you will find the binary file ```hook```. This hook is generate using the project ```mosaic5g/kube5g/dockers/docker-hook``` Plese, check this project for more details.
- Check the capabilities of the provided script:
  ```bash
  $./build.sh 

  Description:
  This Script will remove the old docker snap image and build a new one
  Usage:
          ./build.sh [oai-cn|oai-hss|oai-mme|oai-spgw|oai-ran|flexran|ll-mec] [release tag(default is latest)]
  Example:
          ./build.sh oai-cn mytest
  ```
- You may change the repository name and other parameters that exist in the beggening of the script. For example the default repository is ```mosaic5gecosys```

- To build for example oai-cn docker containe, with the tag ```mytest```:
  - `./build.sh oai-cn mytest`
  - With default setting, you'll get an image ```mosaic5gecosys/oaicn:mytest```. You can check that by typing in the terminal ``` docker images```
<!-- ## Known Issues

- TOSHIBA PORTEGE Z30-C will freeze if running any docker container provided by this branch. -->


**After building dokcer images for mosaic5G snaps, you may proceed now with the examples that are provided using dokcer-compose. There examples are in the project ```mosaic5g/kube5g/dockers/docker-compose```**.