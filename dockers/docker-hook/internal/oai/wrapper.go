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
# file          wrapper.go
# brief 		This is just wrapper of the functions in internal/oai
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package oai

import (
	"fmt"
	"os"
)

//InstallSnap is a wrapper function for installSnapCore
func InstallSnap(OaiObj Oai) {
	// Install Snap Core
	installSnapCore(OaiObj)
}

//InstallCN is a wrapper for installing OAI CN
func InstallCN(OaiObj Oai, buildSnap bool, snapVersion string) {

	// Install oai-cn snap
	installOaicn(OaiObj, buildSnap, snapVersion)

}

//InstallHSS is a wrapper for installing OAI HSS
func InstallHSS(OaiObj Oai, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		// Install oai-cn snap
		installOaiHssV1(OaiObj, buildSnap)
	} else {
		// Install oai-cn snap
		installOaiHssV2(OaiObj, buildSnap)
	}
}

//InstallMME is a wrapper for installing OAI MME
func InstallMME(OaiObj Oai, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		// Install oai-mme snap
		installOaiMmeV1(OaiObj, buildSnap)
	} else {
		// Install oai-mme snap
		installOaiMmeV2(OaiObj, buildSnap)
	}
}

//InstallSPGW is a wrapper for installing OAI SPGW V1
func InstallSPGW(OaiObj Oai, buildSnap bool) {
	// Install oai-spgw-v1 snap
	installOaiSpgwV1(OaiObj, buildSnap)
}

//InstallSPGWC is a wrapper for installing OAI MME
func InstallSPGWC(OaiObj Oai) {
	// Install oai-spgwc-v2 snap
	installOaiSpgwcV2(OaiObj)
}

//InstallSPGWU is a wrapper for installing OAI MME
func InstallSPGWU(OaiObj Oai) {
	// Install oai-spgwc-v2 snap
	installOaiSpgwuV2(OaiObj)
}

// StartCN is a wrapper for configuring and starting OAI CN services
func StartCN(OaiObj Oai, buildSnap bool, snapVersion string) {
	CnAllInOneMode := true
	if snapVersion == "v1" {
		// Start HSS
		OaiObj.Logger.Print("Starting configuring HSS")
		fmt.Println("Starting configuring HSS")
		startHssV1(OaiObj, CnAllInOneMode, buildSnap)
		// Start MME
		OaiObj.Logger.Print("Starting configuring MME")
		fmt.Println("Starting configuring MME")
		startMmeV1(OaiObj, CnAllInOneMode, buildSnap)
		// Start SPGW
		OaiObj.Logger.Print("Starting configuring SPGW")
		fmt.Println("Starting configuring SPGW")
		startSpgwV1(OaiObj, CnAllInOneMode, buildSnap)
	} else if snapVersion == "v2" {
		OaiObj.Logger.Print("Starting configuring HSS v2")
		fmt.Println("Starting configuring HSS v2")
		startHssV2(OaiObj, CnAllInOneMode, buildSnap)
		// time.Sleep(5 * time.Second)
		// Start MME
		OaiObj.Logger.Print("Starting configuring MME v2")
		fmt.Println("Starting configuring MME v2")
		startMmeV2(OaiObj, CnAllInOneMode, buildSnap)
		// time.Sleep(5 * time.Second)
		// Start SPGWC
		OaiObj.Logger.Print("Starting configuring SPGWC v2")
		fmt.Println("Starting configuring SPGWC v2")
		startSpgwcV2(OaiObj, CnAllInOneMode, buildSnap)
		// time.Sleep(5 * time.Second)
		// Start SPGWU
		OaiObj.Logger.Print("Starting configuring SPGWU v2")
		fmt.Println("Starting configuring SPGWU v2")
		startSpgwuV2(OaiObj, CnAllInOneMode, buildSnap)
		// time.Sleep(5 * time.Second)

	} else {
		OaiObj.Logger.Print("Error while trying to install oai core entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of ", snapVersion, " are: v1, v2")
		os.Exit(1)
	}

}

// StartHSS is a wrapper for startHssV1
func StartHSS(OaiObj Oai, buildSnap bool, snapVersion string) {
	CnAllInOneMode := false
	// Start HSS
	if snapVersion == "v1" {
		OaiObj.Logger.Print("Starting configuring HSS v1")
		fmt.Println("Starting configuring HSS v1")
		startHssV1(OaiObj, CnAllInOneMode, buildSnap)
	} else if snapVersion == "v2" {
		OaiObj.Logger.Print("Starting configuring HSS v2")
		fmt.Println("Starting configuring HSS v2")
		startHssV2(OaiObj, CnAllInOneMode, buildSnap)

	} else {
		OaiObj.Logger.Print("Error while trying to oai-hss entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of ", snapVersion, " are: v1, v2")
		os.Exit(1)
	}
}

// StartMME is a wrapper for startMmeV1
func StartMME(OaiObj Oai, buildSnap bool, snapVersion string) {
	CnAllInOneMode := false
	// Start Mme
	if snapVersion == "v1" {
		OaiObj.Logger.Print("Starting configuring MME v1")
		fmt.Println("Starting configuring MME v1")
		startMmeV1(OaiObj, CnAllInOneMode, buildSnap)
	} else if snapVersion == "v2" {
		OaiObj.Logger.Print("Starting configuring MME v2")
		fmt.Println("Starting configuring MME v2")
		startMmeV2(OaiObj, CnAllInOneMode, buildSnap)

	} else {
		OaiObj.Logger.Print("Error while trying to oai-mme entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of ", snapVersion, " are: v1, v2")
		os.Exit(1)
	}
}

// StartSPGW is a wrapper for startSpgwV1
func StartSPGW(OaiObj Oai, buildSnap bool) {
	CnAllInOneMode := false
	// Start Mme
	OaiObj.Logger.Print("Starting configuring SPGW v1")
	fmt.Println("Starting configuring SPGW v1")
	startSpgwV1(OaiObj, CnAllInOneMode, buildSnap)
}

// StartSPGWCV2 is a wrapper for startSpgwV1
func StartSPGWCV2(OaiObj Oai, buildSnap bool) {
	CnAllInOneMode := false
	OaiObj.Logger.Print("Starting configuring SPGWC v2")
	fmt.Println("Starting configuring SPGWC v2")
	startSpgwcV2(OaiObj, CnAllInOneMode, buildSnap)
}

// StartSPGWUV2 is a wrapper for startSpgwV1
func StartSPGWUV2(OaiObj Oai, buildSnap bool) {
	CnAllInOneMode := false
	OaiObj.Logger.Print("Starting configuring SPGWU v2")
	fmt.Println("Starting configuring SPGWU v2")
	startSpgwuV2(OaiObj, CnAllInOneMode, buildSnap)
}

//InstallRAN is a wrapper for installing OAI RAN
func InstallRAN(OaiObj Oai) {

	// Install oai-ran snap
	OaiObj.Logger.Print("Installing RAN")
	fmt.Println("Installing RAN")
	installOairan(OaiObj)
	OaiObj.Logger.Print("RAN is installed")
	fmt.Println("RAN RAN is installed")
}

//InstallRanGnb is a wrapper for installing OAI RAN
func InstallRanGnb(OaiObj Oai) {

	// Install oai-ran snap
	OaiObj.Logger.Print("Installing RAN gNB")
	fmt.Println("Installing RAN gNB")
	installOairanGnb(OaiObj)
	OaiObj.Logger.Print("RAN gNB is installed")
	fmt.Println("RAN gNB is installed")
}

//StartENB is a wrapper for configuring and starting OAI RAN services as 4G eNB
func StartENB(OaiObj Oai, buildSnap bool) {
	nodeFunction := OaiObj.Conf.OaiEnb[0].NodeFunction
	if nodeFunction == "enb" || nodeFunction == "enb-sim" || nodeFunction == "enbsim" {
		OaiObj.Logger.Print("Starting 4G RAN")
		fmt.Println("Starting 4G RAN")
		if nodeFunction != "enb" {
			OaiObj.Logger.Print("Starting a simulation")
			fmt.Println("Starting a simulation")
		}
		startENB(OaiObj, buildSnap)
		OaiObj.Logger.Print("4G RAN is started")
		fmt.Println("4G RAN is started")
	} else {
		OaiObj.Logger.Print("Error while trying to start the RAN: the functionality of RAN as ", nodeFunction, " is not recognized")
		OaiObj.Logger.Print("Starting the RAN in the defaut mode as 4G RAN")
		startENB(OaiObj, buildSnap)
		OaiObj.Logger.Print("4G RAN is started")
		fmt.Println("4G RAN is started")
	}

}

//StartGNB is a wrapper for configuring and starting OAI RAN services as 5G gNB
func StartGNB(OaiObj Oai, buildSnap bool) {
	nodeFunction := OaiObj.Conf.OaiGnb[0].NodeFunction
	if nodeFunction == "gnb" {
		OaiObj.Logger.Print("Starting 5G RAN")
		fmt.Println("Starting 5G RAN")
		// startGNB(OaiObj, buildSnap) // this will be used later
		startGNB(OaiObj, buildSnap)
		OaiObj.Logger.Print("5G RAN is started")
		fmt.Println("5G RAN is started")
	} else {
		OaiObj.Logger.Print("Error while trying to start the RAN: the functionality of RAN as ", nodeFunction, " is not recognized")
		OaiObj.Logger.Print("Starting the RAN in the defaut mode as 5G RAN")
		startGNB(OaiObj, buildSnap)
		OaiObj.Logger.Print("5G RAN is started")
		fmt.Println("5G RAN is started")
	}

}

//InstallFlexRAN is a wrapper for installing FlexRAN
func InstallFlexRAN(OaiObj Oai) {

	// Install flexran snap
	installFlexRAN(OaiObj)
}

//StartFlexRAN is a wrapper for installing FlexRAN
func StartFlexRAN(OaiObj Oai) {

	// start FlexRAN
	startFlexRAN(OaiObj)
}

//InstallMEC is a wrapper for installing LL-MEC
func InstallMEC(OaiObj Oai) {

	// Install ll-mec snap
	installMEC(OaiObj)
}
