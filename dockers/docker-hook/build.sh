
#!/bin/bash
################################################################################
# Licensed to the Mosaic5G under one or more contributor license
# agreements. See the NOTICE file distributed with this
# work for additional information regarding copyright ownership.
# The Mosaic5G licenses this file to You under the
# Apache License, Version 2.0  (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#  
#    	http://www.apache.org/licenses/LICENSE-2.0
  
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
# -------------------------------------------------------------------------------
#   For more information about the Mosaic5G:
#   	contact@mosaic-5g.io
#
#
################################################################################
# file build_m5g
# brief  Install dependencies to develop docker-hook in golang
# author  Osama Arouk

export DEBIAN_FRONTEND=noninteractive
current_DIR=$PWD
cd ../../../
mosaic5g_DIR=$PWD

# make sure that wget is installed
sudo apt install wget -y

# download golang binary
wget https://dl.google.com/go/go1.15.1.linux-amd64.tar.gz

# Exctract it to /usr/local:
sudo tar -C /usr/local -xzf go1.15.1.linux-amd64.tar.gz

# Add /usr/local/go/bin to the PATH environment variable:
echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
echo 'export GOPATH=$HOME/go' >> $HOME/.profile
source $HOME/.profile


###############docker-hook
mkdir $HOME/go/src/mosaic5g/
export GOPATH=$HOME/go/src
ln -s $mosaic5g_DIR/kube5g/dockers/docker-hook $HOME/go/src/mosaic5g/