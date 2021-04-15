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
# file          main.go
# brief 		main file to create docker-hook, in order to install the required snaps and configure them correctly inside dockers
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package main

import (
	"flag"
	"fmt"
	"mosaic5g/docker-hook/internal/oai"
)

const (
	logPath          = "/root/hook.log"
	confPath         = "/root/config/conf.yaml"
	usersPath        = "/root/users.json"
	flexranStatsPath = "/root/flexran_stats.json"
)

func main() {
	// Initialize oai struct
	OaiObj := oai.Oai{}
	err := OaiObj.Init(logPath, confPath, usersPath, flexranStatsPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse input flags
	OaiObj.Logger.Print("Starting parsing input parameters")
	fmt.Println("Starting parsing input parameters")
	installCN := flag.Bool("installCN", false, "a bool to indicate whether to install the snap oai-cn (v1) or (oai-hss, oai-mme, oai-spgwc, and oai-spgwu) (v2), if it is true")
	installRanEnb := flag.Bool("installRAN", false, "a bool to indicate whether to deploy 4G eNB, if it is true")
	installRanGnb := flag.Bool("installRanGnb", false, "a bool to indicate whether to  deploy 5G gNB, if it is true")
	installHSS := flag.Bool("installHSS", false, "a bool to indicate whether to install the snap oai-cn (v1) oai-hss (v2), if it is true")
	installMME := flag.Bool("installMME", false, "a bool to indicate whether to install the snap oai-cn (v1) or oai-mme (v2), if it is true")
	installSPGW := flag.Bool("installSPGW", false, "a bool to indicate whether to install the snap oai-cn, if it is true")
	installSPGWC := flag.Bool("installSPGWC", false, "a bool to indicate whether to install the snap oai-spgwc, if it is true")
	installSPGWU := flag.Bool("installSPGWU", false, "a bool to indicate whether to install the snap oai-spgwu, if it is true")
	installFlexRAN := flag.Bool("installFlexRAN", false, "a bool to indicate whether to install the snap flexran, if it is true")
	//TODO: add the options for enb-sim, ue-sim
	//TODO: need to check if we need different options for enb-sim/ue-sim
	// maybe if enb-sim/ue-sim are inside a single snap, we might need only one option: SIM
	// installUESim := flag.Bool("installUESim", false, "a bool to indicate whether to install the snap ue-sim")
	installENBSim := flag.Bool("installENBSim", false, "a bool to indicate whether to install the snap enb-sim")
	installMEC := flag.Bool("installMEC", false, "a bool to indicate whether to install the snap ll-mec, if it is true")
	buildImage := flag.Bool("build", false, "a bool value to define that the current setup is to build the docker image.")
	var snapVersion string
	flag.StringVar(&snapVersion, "snapVersion", "v2", "a string value to specify the snap version that will be used to build the docker image. Valid values: v1, v2")
	flag.Parse()

	//Install snap core
	OaiObj.Logger.Print("Installing snap")
	fmt.Println("Installing snap")
	oai.InstallSnap(OaiObj)
	// Decide actions based on flags
	// CnAllInOneMode := true
	buildSnap := false
	if *buildImage {
		buildSnap = true
	}
	if *installCN {
		OaiObj.Logger.Print("Installing CN")
		fmt.Println("Installing CN")

		oai.InstallCN(OaiObj, buildSnap, snapVersion)

		OaiObj.Logger.Print("CN is installed")
		fmt.Println("CN is installed")

		OaiObj.Logger.Print("Starting CN")
		fmt.Println("Starting CN")

		oai.StartCN(OaiObj, buildSnap, snapVersion)

		OaiObj.Logger.Print("CN is started: exit")
		fmt.Println("CN is started: exit")
	} else if *installRanEnb || *installENBSim {
		// this if is ugly
		if *installRanEnb {
			OaiObj.Logger.Print("Installing RAN eNB")
			fmt.Println("Installing RAN eNB")
		} else {
			OaiObj.Logger.Print("Installing RAN eNB Sim")
			fmt.Println("Installing RAN eNB Sim")
		}
		oai.InstallRAN(OaiObj)
		if *installRanEnb {
			OaiObj.Logger.Print("Starting RAN eNB")
			fmt.Println("Starting RAN eNB")
		} else {
			OaiObj.Logger.Print("Starting RAN eNB Sim")
			fmt.Println("Starting RAN eNB Sim")
		}
		oai.StartENB(OaiObj, buildSnap)
		/*}else if *installUESim {
		OaiObj.Logger.Print("Installing UE-Sim")
		fmt.Println("Installing UE-Sim")

		oai.InstallUESim(OaiObj)

		OaiObj.Logger.Print("Starting UE-Sim")
		fmt.Println("Starting UE-Sim")

		oai.StartUESim(OaiObj, buildSnap) */
	} else if *installRanGnb {
		OaiObj.Logger.Print("Installing RAN gNB")
		fmt.Println("Installing RAN gNB")

		oai.InstallRanGnb(OaiObj)

		OaiObj.Logger.Print("Starting RAN gNB")
		fmt.Println("Starting RAN gNB")

		oai.StartGNB(OaiObj, buildSnap)
	} else if *installHSS {
		// CnAllInOneMode = false
		oai.InstallHSS(OaiObj, buildSnap, snapVersion)
		oai.StartHSS(OaiObj, buildSnap, snapVersion)
	} else if *installMME {
		// CnAllInOneMode = false
		oai.InstallMME(OaiObj, buildSnap, snapVersion)
		oai.StartMME(OaiObj, buildSnap, snapVersion)
	} else if *installSPGW {
		// CnAllInOneMode = false
		oai.InstallSPGW(OaiObj, buildSnap)
		oai.StartSPGW(OaiObj, buildSnap)
	} else if *installSPGWC {
		// CnAllInOneMode = false
		// Install SPGWC
		oai.InstallSPGWC(OaiObj)
		oai.StartSPGWCV2(OaiObj, buildSnap)
	} else if *installSPGWU {
		// CnAllInOneMode = false
		oai.InstallSPGWU(OaiObj)
		oai.StartSPGWUV2(OaiObj, buildSnap)
	} else if *installFlexRAN {
		oai.InstallFlexRAN(OaiObj)
		if buildSnap == false {
			oai.StartFlexRAN(OaiObj)
		}
	} else if *installMEC {
		oai.InstallMEC(OaiObj)
	} else {
		fmt.Println("This should only be executed in container!!")
		return
	}

	// Give a hello when program ends
	// Do not change the phrase "End of hook", it is used in docker-build
	OaiObj.Logger.Print("End of hook")
	OaiObj.Clean()
}
