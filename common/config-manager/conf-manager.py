#!/usr/bin/python3
################################################################################
# Licensed to the Mosaic5G under one or more contributor license
# agreements. See the NOTICE file distributed with this
# work for additional information regarding copyright ownership.
# The Mosaic5G licenses this file to You under the
# Apache License, Version 2.0  (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#  
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
################################################################################
# file          conf_manager_operator.py
#
# brief         Configure all the Custom resourcs (CRs) for kube5g-operator according to your parameters, like dns of your network, lte band, etc. \
#               After configuring the crds according to your setup environment, you can start using kube5g-operator right away
#
# author        Osama Arouk (C) - 2020 arouk@eurecom.fr
#
# Dependencies  Here is the list of dependencies for this python script:
#               - ruamel.yaml==0.16.12 
#               - colorlog==4.6.2
#                   
#               To install these dependencies: 
#               1- sudo apt install python3-pip
#               2- pip3 install --upgrade pip
#               3- pip3 install ruamel.yaml==0.16.12 colorlog==4.6.2
import os, sys, subprocess, argparse, copy, logging
import ruamel.yaml
from colorlog import ColoredFormatter

#Logging
logger = logging.getLogger('conf.manager')
logger.setLevel(logging.DEBUG)
handler = logging.StreamHandler(sys.stdout)
handler.setLevel(logging.INFO)
# log_format = "%(asctime)s %(levelname)s %(name)s %(filename)s:%(lineno)s %(message)s"
# formatter = logging.Formatter(log_format, datefmt='%Y-%m-%dT%H:%M:%S')
LOGFORMAT = "%(asctime)s %(name)s:%(lineno)s %(log_color)s%(levelname)-3s%(reset)s | %(log_color)s%(message)s%(reset)s"
formatter = ColoredFormatter(LOGFORMAT)
handler.setFormatter(formatter)
logger.addHandler(handler)

# Define required directory
CURRENT_DIR = os.environ['PWD']
PWD = (CURRENT_DIR.split("/common/config-manager"))[0]
COMMON_DIR_CRD = "{}/openshift/kube5g-operator/deploy/crds".format(PWD)
COMMON_DIR_DOCKER = "{}/dockers/docker-compose".format(PWD)
COMMON_DIR_DOCKER_BUILD = "{}/dockers/docker-build".format(PWD)

class ConfigManager(object):
    def __init__(self, conf_short, conf_global_default, config_short, entity_test_cicd):
        self.config_file_short = conf_short
        self.config_file = conf_global_default
        self.common_dir_crs = COMMON_DIR_CRD
        self.common_dir_docker = COMMON_DIR_DOCKER
        self.common_dir_docker_build = COMMON_DIR_DOCKER_BUILD
        self.config_global_data = None
        self.docker_compose_data = None
        self.list_configured_files_crs = list()
        self.list_configured_files_docker = list()
        self.list_configured_files_docker_compose = list()
        
        self.yaml_config()
        self.open_config_global()
        if config_short:
            # open short config file and configure global long file
            self.open_short_config_and_configure_config_global()
        if entity_test_cicd != "none":
            # configure the config file of docker build for testing purposes in ci/cd jenkins
            self.config_ci_cd_test(entity_test_cicd)
    def yaml_config(self):
        ruamel.yaml.YAML().indent(mapping=4, sequence=6, offset=3)
    def open_config_global(self):
        logger.info("getting the global configuration from {}".format(self.config_file))
        try:
            with open(self.config_file) as file_conf_global:
                self.config_global_data = ruamel.yaml.round_trip_load(file_conf_global, preserve_quotes=True)
            logger.info("global configuration is successfully retrieved")
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
    def open_short_config_and_configure_config_global(self):
        logger.info("getting the short configuration from {}".format(self.config_file_short))
        try:
            with open(self.config_file_short) as file_conf_short:
                config_short_data = ruamel.yaml.round_trip_load(file_conf_short, preserve_quotes=True)
            # mcc
            self.config_global_data["spec"]["oaiEnb"][0]["mcc"] = \
                    config_short_data["common"][0]["mcc"]
            # mnc
            self.config_global_data["spec"]["oaiEnb"][0]["mnc"] = \
                    config_short_data["common"][0]["mnc"]
            # eutra_band
            self.config_global_data["spec"]["oaiEnb"][0]["eutra_band"]["default"] = \
                    config_short_data["oaiEnb"][0]["eutra_band"]["default"]
            # downlink_frequency
            self.config_global_data["spec"]["oaiEnb"][0]["downlink_frequency"]["default"] = \
                    config_short_data["oaiEnb"][0]["downlink_frequency"]["default"]
            # uplink_frequency_offset
            self.config_global_data["spec"]["oaiEnb"][0]["uplink_frequency_offset"]["default"] = \
                    config_short_data["oaiEnb"][0]["uplink_frequency_offset"]["default"]
            # N_RB_DL
            self.config_global_data["spec"]["oaiEnb"][0]["N_RB_DL"]["default"] = \
                    config_short_data["oaiEnb"][0]["N_RB_DL"]["default"]
            # tx_gain
            self.config_global_data["spec"]["oaiEnb"][0]["tx_gain"]["default"] = \
                    config_short_data["oaiEnb"][0]["tx_gain"]["default"]
            # rx_gain
            self.config_global_data["spec"]["oaiEnb"][0]["rx_gain"]["default"] = \
                    config_short_data["oaiEnb"][0]["rx_gain"]["default"]
            # pusch_p0_Nominal
            self.config_global_data["spec"]["oaiEnb"][0]["pusch_p0_Nominal"]["default"] = \
                    config_short_data["oaiEnb"][0]["pusch_p0_Nominal"]["default"]
            # pucch_p0_Nominal
            self.config_global_data["spec"]["oaiEnb"][0]["pucch_p0_Nominal"]["default"] = \
                    config_short_data["oaiEnb"][0]["pucch_p0_Nominal"]["default"]
            # pdsch_referenceSignalPower
            self.config_global_data["spec"]["oaiEnb"][0]["pdsch_referenceSignalPower"]["default"] = \
                    config_short_data["oaiEnb"][0]["pdsch_referenceSignalPower"]["default"]
            # puSch10xSnr
            self.config_global_data["spec"]["oaiEnb"][0]["puSch10xSnr"]["default"] = \
                    config_short_data["oaiEnb"][0]["puSch10xSnr"]["default"]
            # puCch10xSnr
            self.config_global_data["spec"]["oaiEnb"][0]["puCch10xSnr"]["default"] = \
                    config_short_data["oaiEnb"][0]["puCch10xSnr"]["default"]
            # parallel_config
            self.config_global_data["spec"]["oaiEnb"][0]["parallel_config"]["default"] = \
                    config_short_data["oaiEnb"][0]["parallel_config"]["default"]
            # max_rxgain
            self.config_global_data["spec"]["oaiEnb"][0]["max_rxgain"]["default"] = \
                    config_short_data["oaiEnb"][0]["max_rxgain"]["default"]
            """ CN: DNS """
            # oaiCn v1
            self.config_global_data["spec"]["oaiCn"]["v1"][0]["oaiSpgw"]["dns"] = \
                config_short_data["oaiCn"][0]["dns"]
            # oaiCn v2
            self.config_global_data["spec"]["oaiCn"]["v2"][0]["oaiSpgwc"]["dns"] = \
                            config_short_data["oaiCn"][0]["dns"]
            # oaiSpgw v1
            self.config_global_data["spec"]["oaiSpgw"]["v1"][0]["dns"] = \
                            config_short_data["oaiCn"][0]["dns"]
            # oaiSpgwc v2
            self.config_global_data["spec"]["oaiSpgwc"]["v2"][0]["dns"] = \
                            config_short_data["oaiCn"][0]["dns"]
            """ CN: APN v2 """
            self.config_global_data["spec"]["oaiCn"]["v2"][0]["APN_NI"]["default"] = \
                            config_short_data["oaiCn"][0]["v2"]["APN_NI"]["default"]
            # oaiHss v2
            self.config_global_data["spec"]["oaiHss"]["v2"][0]["APN_NI"]["default"] = \
                            config_short_data["oaiCn"][0]["v2"]["APN_NI"]["default"]
            # oaiSpgwc v2
            self.config_global_data["spec"]["oaiSpgwc"]["v2"][0]["APN_NI"]["default"] = \
                            config_short_data["oaiCn"][0]["v2"]["APN_NI"]["default"]
                            
            logger.info("configuration is successfully retrieved from short configuration file")
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
    def open_docker_compose(self, file):
        logger.info("getting docker compose {}".format(file))
        self.docker_compose_data = None
        try:
            with open(file) as file_docker_compose:
                self.docker_compose_data = ruamel.yaml.round_trip_load(file_docker_compose, preserve_quotes=True)
            logger.info("docker compose data is successfully retrieved")
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
    
    # change the concerned parameters for testing purpose in ci/cd jenkins
    def config_ci_cd_test(self, entity_test_cicd):
        def config_cicd_test_oairan(self):
            # oaiEnb
            self.config_global_data['spec']["oaiEnb"][0]["oaiEnbImage"]             = "mosaic5gecosys/oairan:v1.test"
            self.config_global_data['spec']["oaiEnb"][0]["snap"]["channel"]         = "edge/ci"
            self.config_global_data['spec']["oaiEnb"][0]["snap"]["devmode"]         = True
            self.config_global_data['spec']["oaiEnb"][0]["snap"]["refresh"]         = True
        def config_cicd_test_flexran(self):
            # flexran
            self.config_global_data['spec']["flexran"][0]["flexranImage"]           = "mosaic5gecosys/flexran:v1.test"
            self.config_global_data['spec']["flexran"][0]["snap"]["channel"]        = "edge/ci"
            self.config_global_data['spec']["flexran"][0]["snap"]["devmode"]        = True
            self.config_global_data['spec']["flexran"][0]["snap"]["refresh"]        = True
        def config_cicd_test_oaiCn(self):
            # oaiCn v1
            self.config_global_data['spec']["oaiCn"]["v1"][0]["oaiCnImage"]         = "mosaic5gecosys/oaicn:v1.test"
            self.config_global_data['spec']["oaiCn"]["v1"][0]["snap"]["channel"]    = "edge/ci"
            self.config_global_data['spec']["oaiCn"]["v1"][0]["snap"]["devmode"]    = True
            self.config_global_data['spec']["oaiCn"]["v1"][0]["snap"]["refresh"]    = True
            # oaiCn v2
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiCnImage"]         = "mosaic5gecosys/oaicn:v2.test"
            # oaiCn:oaiHss v2
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiHss"]["snap"]["channel"]    = "edge/ci"
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiHss"]["snap"]["devmode"]    = True
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiHss"]["snap"]["refresh"]    = True
            # oaiCn:oaiMme v2
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiMme"]["snap"]["channel"]    = "edge/ci"
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiMme"]["snap"]["devmode"]    = True
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiMme"]["snap"]["refresh"]    = True
            # oaiCn:oaiSpgwc v2
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiSpgwc"]["snap"]["channel"]    = "edge/ci"
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiSpgwc"]["snap"]["devmode"]    = True
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiSpgwc"]["snap"]["refresh"]    = True
            # oaiCn:oaiSpgwu v2
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiSpgwu"]["snap"]["channel"]    = "edge/ci"
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiSpgwu"]["snap"]["devmode"]    = True
            self.config_global_data['spec']["oaiCn"]["v2"][0]["oaiSpgwu"]["snap"]["refresh"]    = True


            # oaiHss
        def config_cicd_test_oaiHss(self):
            # oaiHss v1
            self.config_global_data['spec']["oaiHss"]["v1"][0]["oaiHssImage"]       = "mosaic5gecosys/oaihss:v1.test"
            self.config_global_data['spec']["oaiHss"]["v1"][0]["snap"]["channel"]   = "edge/ci"
            self.config_global_data['spec']["oaiHss"]["v1"][0]["snap"]["devmode"]   = True
            self.config_global_data['spec']["oaiHss"]["v1"][0]["snap"]["refresh"]   = True
            # oaiHss v2
            self.config_global_data['spec']["oaiHss"]["v2"][0]["oaiHssImage"]       = "mosaic5gecosys/oaihss:v2.test"
            self.config_global_data['spec']["oaiHss"]["v2"][0]["snap"]["channel"]   = "edge/ci"
            self.config_global_data['spec']["oaiHss"]["v2"][0]["snap"]["devmode"]   = True
            self.config_global_data['spec']["oaiHss"]["v2"][0]["snap"]["refresh"]   = True
        def config_cicd_test_oaiMme(self):
            # oaiMme v1
            self.config_global_data['spec']["oaiMme"]["v1"][0]["oaiMmeImage"] = "mosaic5gecosys/oaimme:v1.test"
            self.config_global_data['spec']["oaiMme"]["v1"][0]["snap"]["channel"]   = "edge/ci"
            self.config_global_data['spec']["oaiMme"]["v1"][0]["snap"]["devmode"]   = True
            self.config_global_data['spec']["oaiMme"]["v1"][0]["snap"]["refresh"]   = True
            # oaiMme v2
            self.config_global_data['spec']["oaiMme"]["v2"][0]["oaiMmeImage"] = "mosaic5gecosys/oaimme:v2.test"
            self.config_global_data['spec']["oaiMme"]["v2"][0]["snap"]["channel"]   = "edge/ci"
            self.config_global_data['spec']["oaiMme"]["v2"][0]["snap"]["devmode"]   = True
            self.config_global_data['spec']["oaiMme"]["v2"][0]["snap"]["refresh"]   = True
        def config_cicd_test_oaiSpgw(self):
            # oaiSpgw v1
            self.config_global_data['spec']["oaiSpgw"]["v1"][0]["oaiSpgwImage"] = "mosaic5gecosys/oaispgw:v1.test"
            self.config_global_data['spec']["oaiSpgw"]["v1"][0]["snap"]["channel"]   = "edge/ci"
            self.config_global_data['spec']["oaiSpgw"]["v1"][0]["snap"]["devmode"]   = True
            self.config_global_data['spec']["oaiSpgw"]["v1"][0]["snap"]["refresh"]   = True
        def config_cicd_test_oaiSpgwc(self):
            # oaiSpgwc v2
            self.config_global_data['spec']["oaiSpgwc"]["v2"][0]["oaiSpgwcImage"] = "mosaic5gecosys/oaispgwc:v2.test"
            self.config_global_data['spec']["oaiSpgwc"]["v2"][0]["snap"]["channel"]   = "edge/ci"
            self.config_global_data['spec']["oaiSpgwc"]["v2"][0]["snap"]["devmode"]   = True
            self.config_global_data['spec']["oaiSpgwc"]["v2"][0]["snap"]["refresh"]   = True
        def config_cicd_test_oaiSpgwu(self):
            # oaiSpgwu v2
            self.config_global_data['spec']["oaiSpgwu"]["v2"][0]["oaiSpgwuImage"] = "mosaic5gecosys/oaispgwu:v2.test"
            self.config_global_data['spec']["oaiSpgwu"]["v2"][0]["snap"]["channel"]   = "edge/ci"
            self.config_global_data['spec']["oaiSpgwu"]["v2"][0]["snap"]["devmode"]   = True
            self.config_global_data['spec']["oaiSpgwu"]["v2"][0]["snap"]["refresh"]   = True
        if entity_test_cicd == "all":
            config_cicd_test_oairan(self)
            config_cicd_test_oaiHss(self)
            config_cicd_test_oaiSpgw(self)
            config_cicd_test_oaiSpgwc(self)
            config_cicd_test_oaiSpgwu(self)
            config_cicd_test_oaiMme(self)
            config_cicd_test_oaiCn(self)
            config_cicd_test_flexran(self)
        elif entity_test_cicd == "oai-ran":
            config_cicd_test_oairan(self)
        elif entity_test_cicd == "oai-hss":
            config_cicd_test_oaiHss(self)
        elif entity_test_cicd == "oai-spgw":
            config_cicd_test_oaiSpgw(self)
        elif entity_test_cicd == "oai-spgwc":
            config_cicd_test_oaiSpgwc(self)
        elif entity_test_cicd == "oai-spgwu":
            config_cicd_test_oaiSpgwu(self)
        elif entity_test_cicd == "oai-mme":
            config_cicd_test_oaiMme(self)
        elif entity_test_cicd == "oai-cn":
            config_cicd_test_oaiCn(self)
        elif entity_test_cicd == "flexran":
            config_cicd_test_flexran(self)
        else:
            logger.error("The entity {} is not supported yet, exit ...".format(entity_test_cicd))
            exit(0)
    ##################################################################
    # change the parameters for testing purpose in ci/cd jenkins
    def config_docker_build(self):
        logger.debug("configuring docker-build config fule")
        conf_file_out_docker_build = "{}/build/conf.yaml".format(self.common_dir_docker_build)
        """                    docker-build config file                        """
        ## Docker config
        # getting the config from global
        conf_docker_lte_all_in_one_data = copy.deepcopy(self.config_global_data['spec'])

        # getting names of docker and other values for the concerned entities
        # oaiEnb
        conf_docker_lte_all_in_one_data["oaiEnb"][0]["flexRAN"] = False

        # remove un-necessary parameters for docker
        k8s_param = ["k8sGlobalNamespace", "database"]
        for key in k8s_param:
            try:
                del conf_docker_lte_all_in_one_data[key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data))

        oaienb_param = ["oaiEnbSize", "oaiEnbImage"]
        oaicn_param = ["oaiCnSize", "oaiCnImage"]
        flexran_param = ["flexranSize", "flexranImage"]
        oaiHss_param = ["oaiHssSize", "oaiHssImage"]
        oaiMme_param = ["oaiMmeSize", "oaiMmeImage"]
        oaiSpgw_param = ["oaiSpgwSize", "oaiSpgwImage"]
        oaiSpgwc_param = ["oaiSpgwcSize", "oaiSpgwcImage"]
        oaiSpgwu_param = ["oaiSpgwuSize", "oaiSpgwuImage"]
        k8s_param = ["k8sDeploymentName", "k8sServiceName", "k8sLabelSelector", "k8sEntityNamespace", "k8sNodeSelector"]
        versions = ["v1", "v2"]
        # remove un-necessary parameters from oaiEnb for docker
        oaienb_param_total = oaienb_param + k8s_param
        for key in oaienb_param_total:
            try:
                del conf_docker_lte_all_in_one_data["oaiEnb"][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiEnb"][0]))
        # oaimme
        conf_docker_lte_all_in_one_data["oaiEnb"][0]["mmeService"]["name"] = "oaimme"
        conf_docker_lte_all_in_one_data["oaiEnb"][0]["mmeService"]["snapVersion"] = "v1"
        conf_docker_lte_all_in_one_data["oaiEnb"][0]["mmeService"]["ipv4"] = ""
        # remove un-necessary parameters from oaiCn for docker
        oaicn_param_total = oaicn_param + k8s_param
        for key in oaicn_param_total:
            for version in versions:
                try:
                    del conf_docker_lte_all_in_one_data["oaiCn"][version][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiCn"][version][0]))

        # remove un-necessary parameters from flexran for docker
        flexran_param_total = flexran_param + k8s_param
        for key in flexran_param_total:
            try:
                del conf_docker_lte_all_in_one_data["flexran"][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["flexran"][0]))
        # remove un-necessary parameters from oaiHss for docker
        oaiHss_param_total = oaiHss_param + k8s_param
        for key in oaiHss_param_total:
            for version in versions:
                try:
                    del conf_docker_lte_all_in_one_data["oaiHss"][version][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiHss"][version][0]))
        conf_docker_lte_all_in_one_data["oaiHss"]["v1"][0]["databaseServiceName"] = "mysql"
        conf_docker_lte_all_in_one_data["oaiHss"]["v1"][0]["mmeServiceName"] = "oaimme"
        conf_docker_lte_all_in_one_data["oaiHss"]["v2"][0]["databaseServiceName"] = "cassandra"
        conf_docker_lte_all_in_one_data["oaiHss"]["v2"][0]["hssServiceName"] = "oaihss"
        conf_docker_lte_all_in_one_data["oaiHss"]["v2"][0]["mmeServiceName"] = "oaimme"
        conf_docker_lte_all_in_one_data["oaiHss"]["v2"][0]["spgwcServiceName"] = "oaispgwc"

        # remove un-necessary parameters from oaiMme for docker
        oaiMme_param_total = oaiMme_param + k8s_param
        for key in oaiMme_param_total:
            for version in versions:
                try:
                    del conf_docker_lte_all_in_one_data["oaiMme"][version][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiMme"][version][0]))
        conf_docker_lte_all_in_one_data["oaiMme"]["v1"][0]["hssServiceName"] = "oaihss"
        conf_docker_lte_all_in_one_data["oaiMme"]["v1"][0]["spgwServiceName"] = "oaispgw"
        conf_docker_lte_all_in_one_data["oaiMme"]["v2"][0]["hssServiceName"] = "oaihss"
        conf_docker_lte_all_in_one_data["oaiMme"]["v2"][0]["spgwcServiceName"] = "oaispgwc"
        # remove un-necessary parameters from oaiSpgw for docker
        oaiSpgw_param_total = oaiSpgw_param + k8s_param
        for key in oaiSpgw_param_total:
            for version in versions:
                try:
                    del conf_docker_lte_all_in_one_data["oaiSpgw"]["v1"][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiSpgw"]["v1"][0]))
        conf_docker_lte_all_in_one_data["oaiSpgw"]["v1"][0]["hssServiceName"] = "oaihss"
        conf_docker_lte_all_in_one_data["oaiSpgw"]["v1"][0]["mmeServiceName"] = "oaimme"
        # remove un-necessary parameters from oaiSpgwc for docker
        oaiSpgwc_param_total = oaiSpgwc_param + k8s_param
        for key in oaiSpgwc_param_total:
            for version in versions:
                try:
                    del conf_docker_lte_all_in_one_data["oaiSpgwc"]["v2"][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiSpgwc"]["v2"][0]))
        # remove un-necessary parameters from oaiSpgwu for docker
        oaiSpgwu_param_total = oaiSpgwu_param + k8s_param
        for key in oaiSpgwu_param_total:
            for version in versions:
                try:
                    del conf_docker_lte_all_in_one_data["oaiSpgwu"]["v2"][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiSpgwu"]["v2"][0]))
        conf_docker_lte_all_in_one_data["oaiSpgwu"]["v2"][0]["spgwcServiceName"] = "oaispgwc"
        # write config of docker-compose
        try:
            with open(conf_file_out_docker_build, 'w') as outfile:
                ruamel.yaml.dump(conf_docker_lte_all_in_one_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_docker.append(conf_file_out_docker_build)
            logger.debug("configuration for docker lte_all_in_one is successfully written in {}".format(conf_file_out_docker_build))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
    ##################################################################
    # config docker compose and crs of kube5g-operator for lte-all-in-one
    def config_lte_all_in_one(self, version):
        logger.debug("configuring crs of lte_all_in_one for the version {}".format(version))
        alter_version = "v1" if version == "v2" else "v2"
        alter_database_type = "cassandra" if version == "v1" else "mysql"
        database_type = "mysql" if version == "v1" else "cassandra"
        conf_file_out_crs = "{}/cr-{}/lte-all-in-one/mosaic5g_v1alpha1_cr_{}_lte_all_in_one.yaml".format(self.common_dir_crs, version, version)
        conf_file_out_docker = "{}/oai-{}/lte-all-in-one/conf.yaml".format(self.common_dir_docker, version)
        docker_compose_file = "{}/oai-{}/lte-all-in-one/docker-compose.yaml".format(self.common_dir_docker, version)
        docker_compose_file_output = "{}/oai-{}/lte-all-in-one/docker-compose.yaml".format(self.common_dir_docker, version)
        """                    crs v1/v2 lte-all-in-one                        """
        # change mme config of oaiEnb
        conf_crs_lte_all_in_one_data = copy.deepcopy(self.config_global_data)
        conf_crs_lte_all_in_one_data["spec"]["oaiEnb"][0]["mmeService"]["snapVersion"] = version
        conf_crs_lte_all_in_one_data["spec"]["oaiEnb"][0]["mmeService"]["name"] = \
                conf_crs_lte_all_in_one_data["spec"]["oaiCn"][version][0]["k8sServiceName"]
        # delete "oaiCn"
        try:
            del conf_crs_lte_all_in_one_data["spec"]["oaiCn"][alter_version]
        except:
            logger.debug("the key {} does not exist in {}, skipping".format(alter_version, conf_crs_lte_all_in_one_data["spec"]["oaiCn"]))
        
        # delete"flexran", "llmec", "oaiHss", "oaiMme", "oaiSpgw", "oaiSpgwc", "oaiSpgwu"        
        oaispgwu_param = ["flexran", "llmec", "oaiHss", "oaiMme", "oaiSpgw", "oaiSpgwc", "oaiSpgwu"]
        for key in oaispgwu_param:
            try:
                del conf_crs_lte_all_in_one_data["spec"][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_crs_lte_all_in_one_data["spec"]))

        # getting the right database
        for item in range(len(conf_crs_lte_all_in_one_data["spec"]["database"])):
            if conf_crs_lte_all_in_one_data["spec"]["database"][item]["databaseType"] == alter_database_type:
                del conf_crs_lte_all_in_one_data["spec"]["database"][item]
                break
        try:
            with open(conf_file_out_crs, 'w') as outfile:
                ruamel.yaml.dump(conf_crs_lte_all_in_one_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_crs.append(conf_file_out_crs)
            logger.debug("configuration of crs for lte_all_in_one {} is successfully written in {}".format(version, conf_file_out_crs))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)

        """                    docker v1/v2 lte-all-in-one                        """
        ## Docker config
        # getting the config from global
        conf_docker_lte_all_in_one_data = copy.deepcopy(conf_crs_lte_all_in_one_data['spec'])
        self.open_docker_compose(docker_compose_file)
        self.docker_compose_data["services"]["oairan"]["image"] = conf_docker_lte_all_in_one_data["oaiEnb"][0]["oaiEnbImage"]
        self.docker_compose_data["services"]["oaicn"]["image"] = conf_docker_lte_all_in_one_data["oaiCn"][version][0]["oaiCnImage"]
        self.docker_compose_data["services"][database_type]["image"] = conf_docker_lte_all_in_one_data["database"][0]["databaseImage"]

        # getting names of docker and other values for the concerned entities
        # oaiEnb
        conf_docker_lte_all_in_one_data["oaiEnb"][0]["mmeService"]["snapVersion"] = version
        conf_docker_lte_all_in_one_data["oaiEnb"][0]["mmeService"]["name"] = self.docker_compose_data["services"]["oaicn"]["container_name"]
        conf_docker_lte_all_in_one_data["oaiEnb"][0]["flexRAN"] = False
        # oaiCn
        conf_docker_lte_all_in_one_data["oaiCn"][version][0]["oaiHss"]["databaseServiceName"] = self.docker_compose_data["services"][database_type]["container_name"]

        # remove un-necessary parameters for docker
        k8s_param = ["k8sGlobalNamespace", "database"]
        for key in k8s_param:
            try:
                del conf_docker_lte_all_in_one_data[key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data))

        oaienb_param = ["oaiEnbSize", "oaiEnbImage"]
        oaicn_param = ["oaiCnSize", "oaiCnImage"]
        k8s_param = ["k8sDeploymentName", "k8sServiceName", "k8sLabelSelector", "k8sEntityNamespace", "k8sNodeSelector"]
        # remove un-necessary parameters from oaiEnb for docker
        oaienb_param_total = oaienb_param + k8s_param
        for key in oaienb_param_total:
            try:
                del conf_docker_lte_all_in_one_data["oaiEnb"][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiEnb"][0]))
        # remove un-necessary parameters from oaiCn for docker
        oaicn_param_total = oaicn_param + k8s_param
        for key in oaicn_param_total:
            try:
                del conf_docker_lte_all_in_one_data["oaiCn"][version][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_data["oaiCn"][version][0]))

        # write config of docker-compose
        try:
            with open(conf_file_out_docker, 'w') as outfile:
                ruamel.yaml.dump(conf_docker_lte_all_in_one_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_docker.append(conf_file_out_docker)
            logger.debug("configuration for docker lte_all_in_one is successfully written in {}".format(conf_file_out_docker))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
        # write docker-compose
        try:
            with open(docker_compose_file_output, 'w') as outfile:
                ruamel.yaml.dump(self.docker_compose_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_docker_compose.append(docker_compose_file_output)
            logger.debug("configuration for docker lte_all_in_one is successfully written in {}".format(docker_compose_file_output))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
    # config docker compose and crs of kube5g-operator for lte-all-in-one-with-flexran
    def config_lte_all_in_one_flexran(self, version):
        logger.debug("configuring crs of lte_all_in_one with flexran for the version {}".format(version))
        alter_version = "v1" if version == "v2" else "v2"
        alter_database_type = "cassandra" if version == "v1" else "mysql"
        database_type = "mysql" if version == "v1" else "cassandra"
        conf_file_out_crs = "{}/cr-{}/lte-all-in-one-with-flexran/mosaic5g_v1alpha1_cr_{}_lte_all_in_one_flexran.yaml".format(self.common_dir_crs, version, version)
        
        conf_file_out_docker = "{}/oai-{}/lte-all-in-one-with-flexran/conf.yaml".format(self.common_dir_docker, version)
        docker_compose_file = "{}/oai-{}/lte-all-in-one-with-flexran/docker-compose.yaml".format(self.common_dir_docker, version)
        docker_compose_file_output = "{}/oai-{}/lte-all-in-one-with-flexran/docker-compose.yaml".format(self.common_dir_docker, version)
        """                    crs v1/v2 lte-all-in-one-with-flexran                        """
        # change mme config of oaiEnb
        conf_crs_lte_all_in_one_data = copy.deepcopy(self.config_global_data)
        conf_crs_lte_all_in_one_data["spec"]["oaiEnb"][0]["mmeService"]["snapVersion"] = version
        conf_crs_lte_all_in_one_data["spec"]["oaiEnb"][0]["mmeService"]["name"] = \
                conf_crs_lte_all_in_one_data["spec"]["oaiCn"][version][0]["k8sServiceName"]
        conf_crs_lte_all_in_one_data["spec"]["oaiEnb"][0]["flexRAN"] = True
        # delete "oaiCn" 
        try:
            del conf_crs_lte_all_in_one_data["spec"]["oaiCn"][alter_version]
        except:
            logger.debug("the key {} does not exist in {}, skipping".format(alter_version, conf_crs_lte_all_in_one_data["spec"]["oaiCn"]))
        
        # delete "llmec", "oaiHss", "oaiMme", "oaiSpgw", "oaiSpgwc", "oaiSpgwu"        
        oaispgwu_param = ["llmec", "oaiHss", "oaiMme", "oaiSpgw", "oaiSpgwc", "oaiSpgwu"]
        for key in oaispgwu_param:
            try:
                del conf_crs_lte_all_in_one_data["spec"][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_crs_lte_all_in_one_data["spec"]))

        # getting the right database
        for item in range(len(conf_crs_lte_all_in_one_data["spec"]["database"])):
            if conf_crs_lte_all_in_one_data["spec"]["database"][item]["databaseType"] == alter_database_type:
                del conf_crs_lte_all_in_one_data["spec"]["database"][item]
                break
        # This will be skipped, as it is not supported yet
        skip = True 
        if not skip:
            try:
                with open(conf_file_out_crs, 'w') as outfile:
                    ruamel.yaml.dump(conf_crs_lte_all_in_one_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
                self.list_configured_files_crs.append(conf_file_out_crs)
                logger.debug("configuration of crs for lte_all_in_one {} is successfully written in {}".format(version, conf_file_out_crs))
            except Exception as ex:
                message = "Error while trying to open the file: {}".format(ex) 
                logger.error(message)
                exit(0) 

        """                    docker v1/v2 lte-all-in-one-with-flexran                        """
        ## Docker config
        # getting the config from global 
        conf_docker_lte_all_in_one_with_flexran_data = copy.deepcopy(conf_crs_lte_all_in_one_data['spec'])
        self.open_docker_compose(docker_compose_file)
        self.docker_compose_data["services"]["oairan"]["image"] = conf_docker_lte_all_in_one_with_flexran_data["oaiEnb"][0]["oaiEnbImage"]
        self.docker_compose_data["services"]["oaicn"]["image"] = conf_docker_lte_all_in_one_with_flexran_data["oaiCn"][version][0]["oaiCnImage"]
        self.docker_compose_data["services"]["flexran"]["image"] = conf_docker_lte_all_in_one_with_flexran_data["flexran"][0]["flexranImage"]
        self.docker_compose_data["services"][database_type]["image"] = conf_docker_lte_all_in_one_with_flexran_data["database"][0]["databaseImage"]

        # getting names of docker and other values for the concerned entities
        # oaiEnb
        conf_docker_lte_all_in_one_with_flexran_data["oaiEnb"][0]["mmeService"]["snapVersion"] = version
        conf_docker_lte_all_in_one_with_flexran_data["oaiEnb"][0]["mmeService"]["name"] = self.docker_compose_data["services"]["oaicn"]["container_name"]
        conf_docker_lte_all_in_one_with_flexran_data["oaiEnb"][0]["flexRAN"] = True
        conf_docker_lte_all_in_one_with_flexran_data["oaiEnb"][0]["flexRANServiceName"] = self.docker_compose_data["services"]["flexran"]["container_name"]
        # oaiCn
        conf_docker_lte_all_in_one_with_flexran_data["oaiCn"][version][0]["oaiHss"]["databaseServiceName"] = self.docker_compose_data["services"][database_type]["container_name"]

        # remove un-necessary parameters for docker
        k8s_param = ["k8sGlobalNamespace", "database"]
        for key in k8s_param:
            try:
                del conf_docker_lte_all_in_one_with_flexran_data[key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_with_flexran_data))

        oaienb_param = ["oaiEnbSize", "oaiEnbImage"]
        oaicn_param = ["oaiCnSize", "oaiCnImage"]
        k8s_param = ["k8sDeploymentName", "k8sServiceName", "k8sLabelSelector", "k8sEntityNamespace", "k8sNodeSelector"]
        # remove un-necessary parameters from oaiEnb for docker
        oaienb_param_total = oaienb_param + k8s_param
        for key in oaienb_param_total:
            try:
                del conf_docker_lte_all_in_one_with_flexran_data["oaiEnb"][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_with_flexran_data["oaiEnb"][0]))

        # remove un-necessary parameters from flexran for docker
        flexran_param = ["flexranSize", "flexranImage"]
        flexran_param_total = flexran_param + k8s_param
        for key in flexran_param_total:
            try:
                del conf_docker_lte_all_in_one_with_flexran_data["flexran"][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_with_flexran_data["flexran"][0]))

        # remove un-necessary parameters from oaiCn for docker
        oaicn_param_total = oaicn_param + k8s_param
        for key in oaicn_param_total:
            try:
                del conf_docker_lte_all_in_one_with_flexran_data["oaiCn"][version][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_all_in_one_with_flexran_data["oaiCn"][version][0]))

        # write config of docker-compose
        try:
            with open(conf_file_out_docker, 'w') as outfile:
                ruamel.yaml.dump(conf_docker_lte_all_in_one_with_flexran_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_docker.append(conf_file_out_docker)
            logger.debug("configuration for docker lte_all_in_one is successfully written in {}".format(conf_file_out_docker))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
        # write docker-compose
        try:
            with open(docker_compose_file_output, 'w') as outfile:
                ruamel.yaml.dump(self.docker_compose_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_docker_compose.append(docker_compose_file_output)
            logger.debug("configuration for docker lte_all_in_one is successfully written in {}".format(docker_compose_file_output))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
    # config docker compose and crs of kube5g-operator for lte
    def config_lte(self, version):
        logger.debug("configuring crs of lte for the version {}".format(version))
        alter_version = "v1" if version == "v2" else "v2"
        alter_database_type = "cassandra" if version == "v1" else "mysql"
        database_type = "mysql" if version == "v1" else "cassandra"
        conf_file_out_crs = "{}/cr-{}/lte/mosaic5g_v1alpha1_cr_{}_lte.yaml".format(self.common_dir_crs, version, version)
        conf_file_out_docker = "{}/oai-{}/lte/conf.yaml".format(self.common_dir_docker, version)
        docker_compose_file = "{}/oai-{}/lte/docker-compose.yaml".format(self.common_dir_docker, version)
        docker_compose_file_output = "{}/oai-{}/lte/docker-compose.yaml".format(self.common_dir_docker, version)
        """                    crs v1/v2 lte                    """
        # change mme config of oaiEnb
        conf_crs_lte_data = copy.deepcopy(self.config_global_data)
        conf_crs_lte_data["spec"]["oaiEnb"][0]["mmeService"]["snapVersion"] = version
        conf_crs_lte_data["spec"]["oaiEnb"][0]["mmeService"]["name"] = \
                conf_crs_lte_data["spec"]["oaiMme"][version][0]["k8sServiceName"]
        # delete flexran, llmec, oaiCn
        entity_param = ["flexran", "llmec", "oaiCn"]
        for key in entity_param:
            try:
                del conf_crs_lte_data["spec"][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_crs_lte_data["spec"]))
        # delete "oaiHss", "oaiMme"
        entity_param = ["oaiHss", "oaiMme"]
        if alter_version == "v1":
            entity_param.append("oaiSpgw")
        else:
            entity_param.append("oaiSpgwc")
            entity_param.append("oaiSpgwu")
        for key in entity_param:
            try:
                del conf_crs_lte_data["spec"][key][alter_version]
            except:
                logger.debug("the key [{}][{}] does not exist in {}, skipping".format(key, alter_version, conf_crs_lte_data["spec"]))
        # delete "oaiSpgw" for v2, and ("oaiSpgwc" and "oaiSpgwu") for v1
        entity_param = list()
        if alter_version == "v1":
            entity_param.append("oaiSpgw")
        else:
            entity_param.append("oaiSpgwc")
            entity_param.append("oaiSpgwu")
        for key in entity_param:
            try:
                del conf_crs_lte_data["spec"][key]
            except:
                logger.debug("the key [{}][{}] does not exist in {}, skipping".format(key, alter_version, conf_crs_lte_data["spec"]))

        for item in range(len(conf_crs_lte_data["spec"]["database"])):
            if conf_crs_lte_data["spec"]["database"][item]["databaseType"] == alter_database_type:
                del conf_crs_lte_data["spec"]["database"][item]
                break
        try:
            with open(conf_file_out_crs, 'w') as outfile:
                ruamel.yaml.dump(conf_crs_lte_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_crs.append(conf_file_out_crs)
            logger.debug("configuration of crs for lte {} is successfully written in {}".format(version, conf_file_out_crs))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
        
        """                    docker v1/v2 lte                        """
        ## Docker config
        # getting the config from global
        conf_docker_lte_data = copy.deepcopy(conf_crs_lte_data['spec'])
        self.open_docker_compose(docker_compose_file)
        self.docker_compose_data["services"]["oairan"]["image"] = conf_docker_lte_data["oaiEnb"][0]["oaiEnbImage"]
        self.docker_compose_data["services"]["oaihss"]["image"] = conf_docker_lte_data["oaiHss"][version][0]["oaiHssImage"]
        self.docker_compose_data["services"]["oaimme"]["image"] = conf_docker_lte_data["oaiMme"][version][0]["oaiMmeImage"]
        # getting names of docker and other values for the concerned entities
        ## oaiEnb
        conf_docker_lte_data["oaiEnb"][0]["mmeService"]["snapVersion"] = version
        conf_docker_lte_data["oaiEnb"][0]["mmeService"]["name"] = self.docker_compose_data["services"]["oaimme"]["container_name"]
        conf_docker_lte_data["oaiEnb"][0]["flexRAN"] = False
        ## oaiHss
        conf_docker_lte_data["oaiHss"][version][0]["databaseServiceName"] = self.docker_compose_data["services"][database_type]["container_name"]
        conf_docker_lte_data["oaiHss"][version][0]["hssServiceName"] = self.docker_compose_data["services"]["oaihss"]["container_name"]
        conf_docker_lte_data["oaiHss"][version][0]["mmeServiceName"] = self.docker_compose_data["services"]["oaimme"]["container_name"]
        ## oaiMme
        conf_docker_lte_data["oaiMme"][version][0]["hssServiceName"] = self.docker_compose_data["services"]["oaihss"]["container_name"]

        if version == "v1":
            self.docker_compose_data["services"]["oaispgw"]["image"] = conf_docker_lte_data["oaiSpgw"][version][0]["oaiSpgwImage"]

            # getting name of docker of oaispgw for the oaiMme entity
            conf_docker_lte_data["oaiMme"][version][0]["spgwServiceName"] = self.docker_compose_data["services"]["oaispgw"]["container_name"]
            # getting name of docker of oaihss and oaimme for the oaiSpgw entity
            conf_docker_lte_data["oaiSpgw"][version][0]["hssServiceName"] = self.docker_compose_data["services"]["oaihss"]["container_name"]
            conf_docker_lte_data["oaiSpgw"][version][0]["mmeServiceName"] = self.docker_compose_data["services"]["oaimme"]["container_name"]
        else:
            self.docker_compose_data["services"]["oaispgwc"]["image"] = conf_docker_lte_data["oaiSpgwc"][version][0]["oaiSpgwcImage"]
            self.docker_compose_data["services"]["oaispgwu"]["image"] = conf_docker_lte_data["oaiSpgwu"][version][0]["oaiSpgwuImage"]
            # getting name of docker of oaispgwc for the oaiHss entity
            conf_docker_lte_data["oaiHss"][version][0]["spgwcServiceName"] = self.docker_compose_data["services"]["oaispgwc"]["container_name"]
            # getting name of docker of oaispgwc for the oaiMme entity
            conf_docker_lte_data["oaiMme"][version][0]["spgwcServiceName"] = self.docker_compose_data["services"]["oaispgwc"]["container_name"]

            # getting name of docker of oaispgwc for the oaispgwu entity
            conf_docker_lte_data["oaiSpgwu"][version][0]["spgwcServiceName"] = self.docker_compose_data["services"]["oaispgwc"]["container_name"]


        self.docker_compose_data["services"][database_type]["image"] = conf_docker_lte_data["database"][0]["databaseImage"]

        # remove un-necessary parameters for docker
        k8s_param = ["k8sGlobalNamespace", "database"]
        for key in k8s_param:
            try:
                del conf_docker_lte_data[key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_data))

        oaienb_param = ["oaiEnbSize", "oaiEnbImage"]
        oaihss_param = ["oaiHssSize", "oaiHssImage"]
        oaimme_param = ["oaiMmeSize", "oaiMmeImage"]
        oaispgw_param = ["oaiSpgwSize", "oaiSpgwImage"]
        oaispgwc_param = ["oaiSpgwcSize", "oaiSpgwcImage"]
        oaispgwu_param = ["oaiSpgwuSize", "oaiSpgwuImage"]
        k8s_param = ["k8sDeploymentName", "k8sServiceName", "k8sLabelSelector", "k8sEntityNamespace", "k8sNodeSelector"]
        # oaiEnb
        oaienb_param_total = oaienb_param + k8s_param
        for key in oaienb_param_total:
            try:
                del conf_docker_lte_data["oaiEnb"][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_data["oaiEnb"][0]))
        # oaiHss
        oaihss_param_total = oaihss_param + k8s_param
        for key in oaihss_param_total:
            try:
                del conf_docker_lte_data["oaiHss"][version][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_data["oaiHss"][version][0]))
        # oaiMme
        oaimme_param_total = oaimme_param + k8s_param
        for key in oaimme_param_total:
            try:
                del conf_docker_lte_data["oaiMme"][version][0][key]
            except:
                logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_data["oaiMme"][version][0]))
        
        if version == "v1":
            # oaiSpgw
            oaispgw_param_total = oaispgw_param + k8s_param
            for key in oaispgw_param_total:
                try:
                    del conf_docker_lte_data["oaiSpgw"][version][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_data["oaiSpgw"][version][0]))
        else:
            # oaiSpgwc
            oaispgwc_param_total = oaispgwc_param + k8s_param
            for key in oaispgwc_param_total:
                try:
                    del conf_docker_lte_data["oaiSpgwc"][version][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_data["oaiSpgwc"][version][0]))
            # oaiSpgwu
            oaispgwu_param_total = oaispgwu_param + k8s_param
            for key in oaispgwu_param_total:
                try:
                    del conf_docker_lte_data["oaiSpgwu"][version][0][key]
                except:
                    logger.debug("the key {} does not exist in {}, skipping".format(key, conf_docker_lte_data["oaiSpgwu"][version][0]))

        # write config of docker-compose
        try:
            with open(conf_file_out_docker, 'w') as outfile:
                ruamel.yaml.dump(conf_docker_lte_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_docker.append(conf_file_out_docker)
            logger.debug("configuration for docker lte_all_in_one is successfully written in {}".format(conf_file_out_docker))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
        # write docker-compose
        try:
            with open(docker_compose_file_output, 'w') as outfile:
                ruamel.yaml.dump(self.docker_compose_data, outfile, Dumper=ruamel.yaml.RoundTripDumper)
            self.list_configured_files_docker_compose.append(docker_compose_file_output)
            logger.debug("configuration for docker lte_all_in_one is successfully written in {}".format(docker_compose_file_output))
        except Exception as ex:
            message = "Error while trying to open the file: {}".format(ex) 
            logger.error(message)
            exit(0)
        
if __name__ == "__main__":
    conf_global_default = "{}/{}".format(CURRENT_DIR, "conf_global_default.yaml")
    conf_short_default = "{}/{}".format(CURRENT_DIR, "conf_short_default.yaml")
    parser = argparse.ArgumentParser(description='configure the custom resources for kube5g-operator')
    parser.add_argument('-g', '--conf-global', metavar='[option]', action='store', type=str,
                        required=False, default='{}'.format(conf_global_default), 
                        help="global configuration file, default: {}".format(conf_global_default))
    parser.add_argument('-s', '--conf-short', metavar='[option]', action='store', type=str,
                            required=False, default='{}'.format(conf_short_default), 
                            help="short configuration file, default: {}".format(conf_short_default))
    parser.add_argument('-t', '--test', metavar='[option]', action='store', type=str,
                        required=False, default='none', 
                        choices=("none", "all", "oai-ran", "oai-hss", "oai-spgw", "oai-spgwc", "oai-spgwu", "oai-mme", "oai-cn", "flexran"),
                            help="THis is intended ONLY for testing in CI/CD. Specify the entity to be tested"
                            "this will replace the tag of the concerned docker images by v1.test (for snap version v1) and v2.test (for snap version v2)")
    
    args = parser.parse_args()

    config_short = False
    if(args.conf_short):
        logger.info("short configuration file: {}".format(args.conf_short))
        config_short = True
    else:
        logger.info("global configuration file: {}".format(args.conf_global))
    
    
    conf_manager = ConfigManager(args.conf_short, args.conf_global, config_short, args.test)
    versions = ["v1", "v2"]
    for version in versions:
        conf_manager.config_lte_all_in_one(version)
        conf_manager.config_lte(version)
    conf_manager.config_lte_all_in_one_flexran("v1")
    conf_manager.config_docker_build()
    logger.info("configuration of crs for the versions {} is successfully finished".format(versions))
    logger.info("Here is the list of configured files:")

    LOGFORMAT = "%(log_color)s%(message)s%(reset)s"
    formatter = ColoredFormatter(LOGFORMAT)
    handler.setFormatter(formatter)
    
    logger.info("Custom Resources (CRs) files:")
    for file in conf_manager.list_configured_files_crs:
        logger.info(file)

    logger.info("\ndocker config files:")
    for file in conf_manager.list_configured_files_docker:
        logger.info(file)
    
    logger.info("\ndocker-compose files:")
    for file in conf_manager.list_configured_files_docker_compose:
        logger.info(file)
