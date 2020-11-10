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
# file          spgw.go
# brief 		configure the snap of oai-spgw v1, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package oai

import (
	"errors"
	"mosaic5g/docker-hook/internal/pkg/util"
	"strings"
)

// StartSpgw : Start SPGW as a daemon
func startSpgwV1(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-cn.spgw-status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-cn.spgw-conf-get"}, "/"))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")

	spgwConf := strings.Join([]string{confPath, "spgw.conf"}, "/")
	spgwBin := strings.Join([]string{snapBinaryPath, "oai-cn.spgw"}, "/")

	var defaultDNSIPV4Address string

	if CnAllInOneMode == true {
		defaultDNSIPV4Address = OaiObj.Conf.OaiCn.V1[0].OaiSpgw.DNS
	} else {
		defaultDNSIPV4Address = OaiObj.Conf.OaiSpgw.V1[0].DNS
	}

	// Init spgw
	OaiObj.Logger.Print("Init spgw")
	// if buildSnap == false {
	util.RunCmd(OaiObj.Logger, spgwBin+"-init")
	// }
	// Configure oai-spgw
	OaiObj.Logger.Print("Configure spgw.conf")

	// Get interface IP and outbound interface
	interfaceIP := util.GetOutboundIP()
	outInterface, _ := util.GetInterfaceByIP(interfaceIP)

	if CnAllInOneMode == true {
		// S-GW binded interface for S11 communication (GTPV2-C): interface name
		sedCommand := "s:SGW_INTERFACE_NAME_FOR_S11.*;:SGW_INTERFACE_NAME_FOR_S11              = \"" + "lo" + "\";:g"
		retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_INTERFACE_NAME_FOR_S11 in " + spgwConf + " failed")
		}
		// S-GW binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11.*;:SGW_IPV4_ADDRESS_FOR_S11                = \"" + "127.0.11.2" + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_IPV4_ADDRESS_FOR_S11 in " + spgwConf + " failed")
		}
	} else {
		// S-GW binded interface for S11 communication (GTPV2-C): interface name
		sedCommand := "s:SGW_INTERFACE_NAME_FOR_S11.*;:SGW_INTERFACE_NAME_FOR_S11              = \"" + outInterface + "\";:g"
		retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_INTERFACE_NAME_FOR_S11 in " + spgwConf + " failed")
		}
		// S-GW binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11.*;:SGW_IPV4_ADDRESS_FOR_S11                = \"" + interfaceIP + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_IPV4_ADDRESS_FOR_S11 in " + spgwConf + " failed")
		}
	}

	// S-GW binded interface for S1-U communication (GTPV1-U): interface name
	sedCommand := "s:SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP.*;:SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP    = \"" + outInterface + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP in " + spgwConf + " failed")
	}
	// S-GW binded interface for S1-U communication (GTPV1-U): ip address
	sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP.*;:SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP      = \"" + interfaceIP + "/24\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP in " + spgwConf + " failed")
	}

	// # P-GW binded interface for SGI (egress/ingress internet traffic): interface name
	sedCommand = "s:PGW_INTERFACE_NAME_FOR_SGI.*;:PGW_INTERFACE_NAME_FOR_SGI            = \"" + outInterface + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set PGW_INTERFACE_NAME_FOR_SGI in " + spgwConf + " failed")
	}

	// # DNS address communicated to UEs
	sedCommand = "s:DEFAULT_DNS_IPV4_ADDRESS.*;:DEFAULT_DNS_IPV4_ADDRESS     = \"" + defaultDNSIPV4Address + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set DEFAULT_DNS_IPV4_ADDRESS in " + spgwConf + " failed")
	}

	secondaryDNS := "8.8.4.4"
	sedCommand = "s:DEFAULT_DNS_SEC_IPV4_ADDRESS.*;:DEFAULT_DNS_SEC_IPV4_ADDRESS = \"" + secondaryDNS + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set DEFAULT_DNS_SEC_IPV4_ADDRESS in " + spgwConf + " failed")
	}

	// oai-cn.spgw-start
	if buildSnap == false {
		OaiObj.Logger.Print("start spgw as daemon")
		util.RunCmd(OaiObj.Logger, spgwBin+"-start")
	}
	return nil
}
