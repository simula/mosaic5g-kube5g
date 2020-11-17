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
#               To install these dependencies: pip3 install ruamel.yaml==0.16.12 colorlog==4.6.2
import os, sys, subprocess, argparse, copy, logging

## Install ruamel.yaml if it does not exist
try:
    import ruamel.yaml
except ImportError:
    subprocess.check_call([sys.executable, "-m", "pip", "install", 'uamel.yaml==0.16.12'])
finally:
    import ruamel.yaml
## Install colorlog if it does not exist
try:
    from colorlog import ColoredFormatter
except ImportError:
    subprocess.check_call([sys.executable, "-m", "pip", "install", 'colorlog==4.6.2'])
finally:
    from colorlog import ColoredFormatter

#Logging
logger = logging.getLogger('conf.manager')
logger.setLevel(logging.DEBUG)
handler = logging.StreamHandler(sys.stdout)
handler.setLevel(logging.INFO)
# log_format = "%(asctime)s %(levelname)s %(name)s %(filename)s:%(lineno)s %(message)s"
# formatter = logging.Formatter(log_format, datefmt='%Y-%m-%dT%H:%M:%S')
LOGFORMAT = "%(asctime)s %(name)s %(log_color)s%(levelname)-8s%(reset)s | %(log_color)s%(message)s%(reset)s"
formatter = ColoredFormatter(LOGFORMAT)
handler.setFormatter(formatter)
logger.addHandler(handler)

# Define required directory
CURRENT_DIR = os.environ['PWD']
PWD = (CURRENT_DIR.split("/common/config-manager"))[0]
COMMON_DIR_CRD = "{}/openshift/kube5g-operator/deploy/crds".format(PWD)
COMMON_DIR_DOCKER = "{}/dockers/docker-compose".format(PWD)

class ConfigManager(object):
    def __init__(self, global_config):
        self.config_file = global_config
        self.common_dir_crs = COMMON_DIR_CRD
        self.common_dir_docker = COMMON_DIR_DOCKER
        self.config_global_data = None
        self.docker_compose_data = None
        self.list_configured_files_crs = list()
        self.list_configured_files_docker = list()
        self.list_configured_files_docker_compose = list()
        
        self.yaml_config()
        self.open_config_global()
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
    
    # config docker compose and crs of kube5g-operators for lte-all-in-one
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
    # config docker compose and crs of kube5g-operators for lte
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
        if version == "v1":
            self.docker_compose_data["services"]["oaispgw"]["image"] = conf_docker_lte_data["oaiSpgw"][version][0]["oaiSpgwImage"]
        else:
            self.docker_compose_data["services"]["oaispgwc"]["image"] = conf_docker_lte_data["oaiSpgwc"][version][0]["oaiSpgwcImage"]
            self.docker_compose_data["services"]["oaispgwu"]["image"] = conf_docker_lte_data["oaiSpgwu"][version][0]["oaiSpgwuImage"]
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

    global_config = "{}/{}".format(CURRENT_DIR, "conf_global.yaml")
    parser = argparse.ArgumentParser(description='configure the custom resources for kube5g-operator')
    parser.add_argument('-c', '--global-config', metavar='[option]', action='store', type=str,
                        required=False, default='{}'.format(global_config), 
                        help="global configuration file of crs, default: {}".format(global_config))

    args = parser.parse_args()

    logger.info("global configuration file of crs: {}".format(args.global_config))
    conf_manager = ConfigManager(args.global_config)
    versions = ["v1", "v2"]
    for version in versions:
        conf_manager.config_lte_all_in_one(version)
        conf_manager.config_lte(version)
    logger.info("configuration of crs for the versions {} is successfully finished".format(versions))
    logger.info("Here is the list of configured files:")

    LOGFORMAT = "%(log_color)s%(message)s%(reset)s"
    formatter = ColoredFormatter(LOGFORMAT)
    handler.setFormatter(formatter)
    
    for file in conf_manager.list_configured_files_crs:
        logger.info(file)
    n = 0
    for file in conf_manager.list_configured_files_docker:
        if n== 0:
            message = "{} {}".format("\n", file)
            logger.info(message)
            n = 1
        else: 
            logger.info(file)
    n = 0
    for file in conf_manager.list_configured_files_docker_compose:
        if n== 0:
            message = "{} {}".format("\n", file)
            logger.info(message)
            n = 1
        else: 
            logger.info(file)
