/*
#!/usr/local/go/bin/go
################################################################################
* Copyright 2016-2019 Eurecom and Mosaic5G Platforms Authors
* Licensed to the Mosaic5G under one or more contributor license
* agreements. See the NOTICE file distributed with this
* work for additional information regarding copyright ownership.
* The Mosaic5G licenses this file to You under the
* Apache License, Version 2.0  (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
################################################################################
#-------------------------------------------------------------------------------
# For more information about Mosaic5G:
#                                   admin@mosaic-5g.io
# file          mme.go
# brief 		configure the snap of oai-mme v1, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package oai

import (
	"errors"
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"os"
	"strings"
	"time"
)

// StartMme : Start MME as a daemon
func startMmeV1(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring MME V1")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-cn.mme-status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-cn.mme-conf-get"}, "/"))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")

	mmeConf := strings.Join([]string{confPath, "mme.conf"}, "/")
	mmeFdConf := strings.Join([]string{confPath, "mme_fd.conf"}, "/")
	mmeBin := strings.Join([]string{snapBinaryPath, "oai-cn.mme"}, "/")

	OaiObj.Logger.Print("mmeConf=", mmeConf)
	fmt.Println("mmeConf=", mmeConf)

	OaiObj.Logger.Print("mmeFdConf=", mmeFdConf)
	fmt.Println("mmeFdConf=", mmeFdConf)

	OaiObj.Logger.Print("mmeBin=", mmeBin)
	fmt.Println("mmeBin=", mmeBin)

	hostname, _ := os.Hostname()

	// Init mme
	var hssServiceName, spgwServiceName, mncValue, mccValue, realm string
	if CnAllInOneMode == true {
		mncValue = OaiObj.Conf.OaiCn.V1[0].OaiMme.MNC
		mccValue = OaiObj.Conf.OaiCn.V1[0].OaiMme.MCC
		realm = OaiObj.Conf.OaiCn.V1[0].Realm.Default
	} else {
		hssServiceName = OaiObj.Conf.OaiMme.V1[0].HssServiceName
		spgwServiceName = OaiObj.Conf.OaiMme.V1[0].SpgwServiceName
		mncValue = OaiObj.Conf.OaiMme.V1[0].MNC
		mccValue = OaiObj.Conf.OaiMme.V1[0].MCC
		realm = OaiObj.Conf.OaiMme.V1[0].Realm.Default
	}

	OaiObj.Logger.Print("Init mme")
	retStatus = util.RunCmd(OaiObj.Logger, mmeBin+"-init")
	if retStatus.Exit != 0 {
		return errors.New("mme init failed ")
	}

	// Configure oai-mme
	OaiObj.Logger.Print("Configure mme.conf")

	// hostname
	sedCommand := "s:HSS_HOSTNAME.*;:HSS_HOSTNAME               = \"" + hostname + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set hss domain name in " + mmeConf + " failed")
	}
	// Replace GUMMEI
	OaiObj.Logger.Print("Replace MNC")
	sedCommand = "s/MNC=\"93\"/MNC=\"" + mncValue + "\\\"/g"
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set GUMMEI in " + mmeConf + " failed")
	}
	OaiObj.Logger.Print("Replace MCC")
	//Replace MCC
	sedCommand = "s:{MCC=\"208\":{MCC=\"" + mccValue + "\":g"
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set TAI in " + mmeConf + " failed")
	}

	// Get interface ip and replace the default one
	outInterfaceIP := util.GetOutboundIP()
	outInterface, _ := util.GetInterfaceByIP(outInterfaceIP)

	// MME binded interface for S1-C or S1-MME  communication (S1AP): interface name
	sedCommand = "s:MME_INTERFACE_NAME_FOR_S1_MME.*;:MME_INTERFACE_NAME_FOR_S1_MME         = \"" + outInterface + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MME_INTERFACE_NAME_FOR_S1_MME in " + mmeConf + " failed")
	}
	// MME binded interface for S1-C or S1-MME  communication (S1AP): ip address
	sedCommand = "s:MME_IPV4_ADDRESS_FOR_S1_MME.*;:MME_IPV4_ADDRESS_FOR_S1_MME           = \"" + outInterfaceIP + "/24\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MME_IPV4_ADDRESS_FOR_S1_MME in " + mmeConf + " failed")
	}
	if CnAllInOneMode == true {
		// MME binded interface for S11 communication (GTPV2-C): interface name
		sedCommand = "s:MME_INTERFACE_NAME_FOR_S11_MME.*;:MME_INTERFACE_NAME_FOR_S11_MME        = \"" + "lo" + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_INTERFACE_NAME_FOR_S11_MME in " + mmeConf + " failed")
		}
		// MME binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:MME_IPV4_ADDRESS_FOR_S11_MME.*;:MME_IPV4_ADDRESS_FOR_S11_MME          = \"" + "127.0.11.1" + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_IPV4_ADDRESS_FOR_S11_MME in " + mmeConf + " failed")
		}
	} else {
		// MME binded interface for S11 communication (GTPV2-C): interface name
		sedCommand = "s:MME_INTERFACE_NAME_FOR_S11_MME.*;:MME_INTERFACE_NAME_FOR_S11_MME        = \"" + outInterface + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_INTERFACE_NAME_FOR_S11_MME in " + mmeConf + " failed")
		}
		// MME binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:MME_IPV4_ADDRESS_FOR_S11_MME.*;:MME_IPV4_ADDRESS_FOR_S11_MME          = \"" + outInterfaceIP + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_IPV4_ADDRESS_FOR_S11_MME in " + mmeConf + " failed")
		}
	}

	var spgwIP string
	var err error
	spgwIP = "127.0.11.2"
	if (CnAllInOneMode != true) && (buildSnap != true) {
		spgwIP, err = util.GetIPFromDomain(OaiObj.Logger, spgwServiceName)
		for {
			if err != nil {
				OaiObj.Logger.Print(err)
			} else {
				hostNameSpgw, err := net.LookupHost(spgwIP)
				if len(hostNameSpgw) > 0 {
					break
				} else {
					OaiObj.Logger.Print(err)
				}
			}
			OaiObj.Logger.Print("Valid ip address for spgw not yet retreived")
			time.Sleep(1 * time.Second)
			spgwIP, err = util.GetIPFromDomain(OaiObj.Logger, spgwServiceName)
		}
	}

	//S-GW binded interface for S11 communication (GTPV2-C): ip address
	sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11.*;:SGW_IPV4_ADDRESS_FOR_S11          = \"" + spgwIP + "/8\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set SGW_IPV4_ADDRESS_FOR_S11 in " + mmeConf + " failed")
	}

	// Identity
	identity := hostname + "." + realm // use the Hostname we got before
	sedCommand = "s:Identity.*;:Identity = \"" + identity + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set Identity in " + mmeFdConf + " failed")
	}
	// Realm
	sedCommand = "s:Realm.*;:Realm = \"" + realm + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set Realm in " + mmeFdConf + " failed")
	}

	// Replace the hostname of Peer conectivity address
	sedCommand = "103s/ubuntu/" + hostname + "/g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set hostname in " + mmeFdConf + " failed")
	}
	// Get the IP address of oai-hss

	if CnAllInOneMode != true {
		hssIP, err := util.GetIPFromDomain(OaiObj.Logger, hssServiceName)
		if buildSnap == true {
			hssIP = "127.0.0.1"
		} else {
			for {
				if err != nil {
					OaiObj.Logger.Print(err)
				} else {
					hostNameHss, err := net.LookupHost(hssIP)
					if len(hostNameHss) > 0 {
						break
					} else {
						OaiObj.Logger.Print(err)
					}
				}
				OaiObj.Logger.Print("Valid ip address for oai-hss not yet retreived")
				time.Sleep(1 * time.Second)
				hssIP, err = util.GetIPFromDomain(OaiObj.Logger, hssServiceName)
			}
		}

		// replace the ip address of hss
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.1/"+hssIP+"/g", mmeFdConf)
		if retStatus.Exit != 0 {
			return errors.New("Set the ip address of oai-hss in " + mmeFdConf + " failed")
		}
	}

	// oai-cn.mme-start
	if buildSnap != true {
		OaiObj.Logger.Print("start mme as daemon")
		util.RunCmd(OaiObj.Logger, mmeBin+"-start")
	}
	return nil
}
