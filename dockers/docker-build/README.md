# Container with Mosiac5g Snaps

This directory includes the build materials for building containers with mosaic5g snaps. These scripts allow you to

- Create docker containers that includes the mosaic5G snaps.
- Bring up a set of applications in one command

## WARNING: Read Before Using any Script

- This creates a set of containers with **security** options **disabled**, this is an unsupported setup, if you have multiple snap packages inside the same container they will be able to break out of the confinement and see each others data and processes. **Do not rely on security inside the container**.
- The scripts are tested and works fine in our environment. The models we used for testing is *GIGABYTE BRIX GB-BRi7-8550* and *Dell XPS-15*.  We ARE NOT sure they won't cause trouble in other environments. For more details, please read the known issue section.
- For the details of containers, please read the individual README in their foldes.

## Requirement

- Ubuntu 16.04/18.04
- Docker-ce 18.09
- Docker Compose
- Golang 1.10+ (Optional; If you want to rebuild the hook)

**Please make sure that docker can be run with non-root user.**

## Quick Start

### Run Pre-set Service

In **compose** directory, we provide docker-compose files that can bring up Mosaic5g services without configuring. Just `cd` to your desired service directory and run `docker-compose up -d`. For example, to start an OAI lte service,

1. Go to the lte folder `cd compose/lte`
2. Check if the parameters in `conf.yaml` meet your need
3. Run `docker-compose up -d`
4. The services will start running when ready

### Build From Source

In the build folder:

- To build oai-cn docker containers from source, with the tag mytest:
  - `./build.sh oai-cn mytest`
  - With default setting, you'll get an image **mosaic5gecosys/oaicn:mytest**.

- To build oai-ran docker containers from source, with the tag mytest:
  - `./build.sh oai-ran mytest`
  - With default setting, you'll get an image **mosaic5gecosys/oairan:mytest**

- To clean up unused containers & iamges:
  - `./build.sh clean_all`
  - This will clean the images and containers that are used for building

## Known Issues

- TOSHIBA PORTEGE Z30-C will freeze if running any docker container provided by this branch.
