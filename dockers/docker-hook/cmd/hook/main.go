/*
# Copyright (c) 2020 Eurecom
################################################################################
# Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The OpenAirInterface Software Alliance licenses this file to You under
# the Apache License, Version 2.0  (the "License"); you may not use this file
# except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#-------------------------------------------------------------------------------
# For more information about the OpenAirInterface (OAI) Software Alliance:
#      contact@openairinterface.org
################################################################################

// This hook is made for installing and configuring snaps inside docker
// Author: Osama Arouk (arouk@eurcom.fr), Kevin Hsi-Ping Hsu
*/
package main

import (
	"flag"
	"fmt"
	"mosaic5g/docker-hook/internal/oai"
	"mosaic5g/docker-hook/internal/pkg/util"
)

const (
	logPath  = "/root/hook.log"
	confPath = "/root/config/conf.yaml"

	oaicnLogPathV1  = "/root/hook-oaicn-v1.log"
	oaicnConfPathV1 = "/root/config/conf-oaicn-v1.yaml"

	oaicnLogPathV2  = "/root/hook-oaicn-v2.log"
	oaicnConfPathV2 = "/root/config/conf-oaicn-v2.yaml"

	oairanLogPathV2  = "/root/hook-oairan-v2.log"
	oairanConfPathV2 = "/root/config/conf-oairan-v2.yaml"
)

func main() {

	// Parse input flags
	fmt.Println("Starting parsing input parameters")
	installCN := flag.Bool("installCN", false, "Bool value to define that the hook will install and configure oai-cn inside the docker image")
	installRAN := flag.Bool("installRAN", false, "Bool value to define that the hook will install and configure oai-ran inside the docker image")
	installHSS := flag.Bool("installHSS", false, "Bool value to define that the hook will install and configure oai-hss inside the docker image")
	installMME := flag.Bool("installMME", false, "Bool value to define that the hook will install and configure oai-mme inside the docker image")
	installSPGW := flag.Bool("installSPGW", false, "Bool value to define that the hook will install and configure oai-spgw inside the docker image")
	installSPGWC := flag.Bool("installSPGWC", false, "Bool value to define that the hook will install and configure oai-spgwc inside the docker image")
	installSPGWU := flag.Bool("installSPGWU", false, "Bool value to define that the hook will install and configure oai-spgwu inside the docker image")
	installFlexRAN := flag.Bool("installFlexRAN", false, "Bool value to define that the hook will install and configure flexran inside the docker image")
	installMEC := flag.Bool("installMEC", false, "Bool value to define that the hook will install and configure ll-mec inside the docker image")
	buildImage := flag.Bool("build", false, "a bool value to define that the current setup is to build the docker image.")
	var snapVersion string
	flag.StringVar(&snapVersion, "snapVersion", "v2", "a string value to specify the snap version that will be used to build the docker image. Valid values: v1, v2")
	flag.Parse()

	// Decide actions based on flags
	CnAllInOneMode := true
	buildSnap := false
	if *buildImage {
		buildSnap = true
	}
	var OaiObj oai.Oai
	if *installCN {
		// Initialize oai struct
		OaiObj = oaiInit("cn")

		util.PrintFunc(OaiObj.Logger, "Installing CN")
		oai.InstallCN(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		util.PrintFunc(OaiObj.Logger, "Starting CN")
		oai.StartCN(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		util.PrintFunc(OaiObj.Logger, "CN is started: exit")
	} else if *installRAN {
		// Initialize oai struct
		OaiObj = oaiInit("ran")

		util.PrintFunc(OaiObj.Logger, "Installing RAN")
		oai.InstallRAN(OaiObj)
		// Define the functionality of the snap: oai-enb, oai-cu, oai-du, oai-rcc, oai-rru
		RanNodeFunction := OaiObj.ConfOaiRan.OaiRanConf.ComponentCarriers.NodeFunction
		fmt.Println(RanNodeFunction)
		if (RanNodeFunction == "ENB") || (RanNodeFunction == "enb") {
			util.PrintFunc(OaiObj.Logger, "Starting RAN ENB")
			oai.StartENB(OaiObj, snapVersion, buildSnap)
			util.PrintFunc(OaiObj.Logger, "RAN ENB Started: exit")
		} else if (RanNodeFunction == "CU") || (RanNodeFunction == "cu") {
			util.PrintFunc(OaiObj.Logger, "Starting RAN CU")
			oai.StartCu(OaiObj)
		} else if (RanNodeFunction == "DU") || (RanNodeFunction == "du") {
			util.PrintFunc(OaiObj.Logger, "Starting RAN DU")
			oai.StartDu(OaiObj)
		} else if (RanNodeFunction == "RRC") || (RanNodeFunction == "rrc") {
			util.PrintFunc(OaiObj.Logger, "Starting RAN RRC")
			oai.StartRrc(OaiObj)
		} else if (RanNodeFunction == "RRU") || (RanNodeFunction == "rru") {
			util.PrintFunc(OaiObj.Logger, "Starting RAN RRU")
			oai.StartRru(OaiObj)
		} else if (RanNodeFunction == "STOP") || (RanNodeFunction == "stop") {
			util.PrintFunc(OaiObj.Logger, "Stopping RAN Service")
			oai.StopRan(OaiObj)
		} else {
			util.PrintFunc(OaiObj.Logger, "Error, unkown node function: ", RanNodeFunction, "Starting RAN eNB")
			oai.StartENB(OaiObj, snapVersion, buildSnap)
		}
		//////////
		// OaiObj.Logger.Print("Starting RAN")
		// fmt.Println("Starting RAN")
		// oai.StartENB(OaiObj)
	} else if *installHSS {
		CnAllInOneMode = false
		oai.InstallHSS(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		oai.StartHSS(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
	} else if *installMME {
		CnAllInOneMode = false
		// // Start; This is for testing
		// CnAllInOneMode = true
		// oai.InstallHSS(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		// oai.StartHSS(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		// // End; This is for testing
		oai.InstallMME(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		oai.StartMME(OaiObj, CnAllInOneMode, buildSnap, snapVersion)

	} else if *installSPGW {
		CnAllInOneMode = false
		oai.InstallCN(OaiObj, CnAllInOneMode, buildSnap, "v1")
		// oai.InstallCN(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		oai.StartSPGW(OaiObj, CnAllInOneMode, buildSnap)
	} else if *installSPGWC {
		CnAllInOneMode = false
		// Install SPGWC
		oai.InstallSPGWC(OaiObj)
		oai.StartSPGWCV2(OaiObj, CnAllInOneMode, buildSnap)
	} else if *installSPGWU {
		CnAllInOneMode = false
		// // Start; This is for testing
		// CnAllInOneMode = true
		// oai.InstallSPGWC(OaiObj)
		// oai.StartSPGWCV2(OaiObj, CnAllInOneMode, buildSnap)
		// // End; This is for testing
		// Install SPGWU
		oai.InstallSPGWU(OaiObj)
		oai.StartSPGWUV2(OaiObj, CnAllInOneMode, buildSnap)
	} else if *installFlexRAN {
		oai.InstallFlexRAN(OaiObj)
		oai.StartFlexRAN(OaiObj, buildSnap)
	} else if *installMEC {
		oai.InstallMEC(OaiObj)
	} else {
		fmt.Println("This should only be executed in container!!")
		return
	}
	OaiObj.Logger.Print("CnAllInOneMode=", CnAllInOneMode)
	OaiObj.Logger.Print("buildSnap=", buildSnap)
	OaiObj.Logger.Print("snapVersion=", snapVersion)

	fmt.Print("CnAllInOneMode=", CnAllInOneMode)
	fmt.Print("buildSnap=", buildSnap)
	fmt.Print("snapVersion=", snapVersion)
	// Give a hello when program ends
	OaiObj.Logger.Print("End of hook")
	OaiObj.Clean()
}

func oaiInit(entity string) oai.Oai {
	// Initialize oai struct
	OaiObj := oai.Oai{}
	OaiObj.Init(entity)
	util.PrintFunc(OaiObj.Logger, "Init of OAI is successful")

	//Install snap core
	util.PrintFunc(OaiObj.Logger, "Installing snap")
	oai.InstallSnap(OaiObj)
	return OaiObj
}
