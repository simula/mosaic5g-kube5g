# Docker Container with OAI-CN Snap

## WARNING: Read Before Using any Script 

- This creates a set of containers with **security** options **disabled**, this is an unsupported setup, if you have multiple snap packages inside the same container they will be able to break out of the confinement and see each others data and processes. **Do not rely on security inside the container**.
- The scripts are tested and works fine in our environment. We ARE NOT sure they won't cause trouble in other environments. For more details, please read the known issue section. The models we used for testing is *GIGABYTE BRIX GB-BRi7-8550* and *Dell XPS-15*

## Extra packages installed in this image

- apt-utils
- dnsutils
- net-tools
- iputils-ping
- vim
