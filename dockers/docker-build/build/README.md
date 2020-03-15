# Quick Start Commands

- To build oai-cn docker containers from source, with the tag mytest:
  - `./build.sh oai-cn mytest`
  - With default setting, you'll get an image **mosaic5gecosys/oaicn:mytest**.

- To build oai-ran docker containers from source, with the tag mytest:
  - `./build.sh oai-ran mytest`
  - With default setting, you'll get an image **mosaic5gecosys/oairan:mytest**

- In addition, `flexran` and `ll-mec` are for building *flexran* and *ll-mec* respectively.

- To clean up unused containers & iamges:
  - `./build.sh clean_all`
  - This will clean the images and containers that are used for building

## Build Info

Change them to meet your need, They're located at the beginning of the **build.sh**

```shell
...
REPO_NAME="mosaic5gecosys" # Change it to your repository
TARGET="${REPO_NAME}/${TARGET_NAME}" # The name of our image
TAG_BASE="base" # The tag for the base image
BASE_CONTAINER="build_base" # The name of the temporary container
RELEASE_TAG="latest" # Default release tag
...
```

If you have [source of hook](https://github.com/tig4605246/oai-snap-in-docker) in your GOPATH, you can run `./build.sh build-hook` to rebuild and update the hook in the build folder.
