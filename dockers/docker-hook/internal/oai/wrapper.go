package oai

import (
	"fmt"
	"os"
	"time"
)

//InstallSnap is a wrapper function for installSnapCore
func InstallSnap(OaiObj Oai) {
	// Install Snap Core
	installSnapCore(OaiObj)
}

//InstallCN is a wrapper for installing OAI CN
func InstallCN(OaiObj Oai, CnAllInOneMode bool, buildSnap bool, snapVersion string) {

	// Install oai-cn snap
	installOaicn(OaiObj, CnAllInOneMode, buildSnap, snapVersion)

}

//InstallHSS is a wrapper for installing OAI HSS
func InstallHSS(OaiObj Oai, CnAllInOneMode bool, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		// Install oai-cn snap
		installOaicn(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
	} else {
		// Install oai-cn snap
		installOaiHssV2(OaiObj, CnAllInOneMode, buildSnap)
	}
}

//InstallMME is a wrapper for installing OAI MME
func InstallMME(OaiObj Oai, CnAllInOneMode bool, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		// Install oai-mme snap
		installOaicn(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
	} else {
		// Install oai-mme snap
		installOaiMmeV2(OaiObj, CnAllInOneMode, buildSnap)
	}
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

// if snapVersion == "v1" {

// // InitCN is a wrapper for initing OAI CN services
// func InitCN(OaiObj Oai) {
// 	// Init HSS
// 	OaiObj.Logger.Print("Initing  HSS")
// 	fmt.Println("Initing HSS")
// 	initHss(OaiObj)
// 	// Init MME
// 	OaiObj.Logger.Print("Initing MME")
// 	fmt.Println("Initing MME")
// 	initMme(OaiObj)
// 	// Init SPGW
// 	OaiObj.Logger.Print("Initing SPGW")
// 	fmt.Println("Initing SPGW")
// 	initSpgw(OaiObj)
// }

// StartCN is a wrapper for configuring and starting OAI CN services
func StartCN(OaiObj Oai, CnAllInOneMode bool, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		// Start HSS
		OaiObj.Logger.Print("Starting configuring HSS")
		fmt.Println("Starting configuring HSS")
		startHss(OaiObj, CnAllInOneMode, buildSnap)
		// Start MME
		OaiObj.Logger.Print("Starting configuring MME")
		fmt.Println("Starting configuring MME")
		startMme(OaiObj, CnAllInOneMode, buildSnap)
		// Start SPGW
		OaiObj.Logger.Print("Starting configuring SPGW")
		fmt.Println("Starting configuring SPGW")
		startSpgw(OaiObj, CnAllInOneMode, buildSnap)
	} else if snapVersion == "v2" {
		/////
		// Start HSS
		// time.Sleep(20 * time.Second)
		OaiObj.Logger.Print("Starting configuring HSS v2")
		fmt.Println("Starting configuring HSS v2")
		startHssV2(OaiObj, CnAllInOneMode, buildSnap)
		time.Sleep(5 * time.Second)
		// Start MME
		OaiObj.Logger.Print("Starting configuring MME v2")
		fmt.Println("Starting configuring MME v2")
		startMmeV2(OaiObj, CnAllInOneMode, buildSnap)
		time.Sleep(5 * time.Second)
		// Start SPGWC
		OaiObj.Logger.Print("Starting configuring SPGWC v2")
		fmt.Println("Starting configuring SPGWC v2")
		startSpgwcV2(OaiObj, CnAllInOneMode, buildSnap)
		time.Sleep(5 * time.Second)
		// Start SPGWU
		OaiObj.Logger.Print("Starting configuring SPGWU v2")
		fmt.Println("Starting configuring SPGWU v2")
		startSpgwuV2(OaiObj, CnAllInOneMode, buildSnap)
		time.Sleep(5 * time.Second)

		// Start and block HSS
		OaiObj.Logger.Print("Starting and block configuring HSS v2")
		fmt.Println("Starting  and block configuring HSS v2")
		startAndBlockHssV2(OaiObj, CnAllInOneMode, buildSnap)
		// time.Sleep(10 * time.Second)
		////////////
	} else {
		OaiObj.Logger.Print("Error while trying to install oai core entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of ", snapVersion, " are: v1, v2")
		os.Exit(1)
	}
}

// // InitHSS is a wrapper for initHss
// func InitHSS(OaiObj Oai) {
// 	// Init HSS
// 	initHss(OaiObj)
// }

// StartHSS is a wrapper for startHss
func StartHSS(OaiObj Oai, CnAllInOneMode bool, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		OaiObj.Logger.Print("Starting configuring HSS v1")
		fmt.Println("Starting configuring HSS v1")
		startHss(OaiObj, CnAllInOneMode, buildSnap)
	} else if snapVersion == "v2" {
		OaiObj.Logger.Print("Starting configuring HSS v2")
		fmt.Println("Starting configuring HSS v2")
		startHssV2(OaiObj, CnAllInOneMode, buildSnap)

		// Start and block HSS
		// time.Sleep(3 * time.Second)
		OaiObj.Logger.Print("Starting and block configuring HSS v2")
		fmt.Println("Starting  and block configuring HSS v2")
		startAndBlockHssV2(OaiObj, CnAllInOneMode, buildSnap)

	} else {
		OaiObj.Logger.Print("Error while trying to oai-hss entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of ", snapVersion, " are: v1, v2")
		os.Exit(1)
	}
}

// // InitMME is a wrapper for initHss
// func InitMME(OaiObj Oai) {
// 	// Init MME
// 	initMme(OaiObj)
// }

// StartMME is a wrapper for startMme
func StartMME(OaiObj Oai, CnAllInOneMode bool, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		OaiObj.Logger.Print("Starting configuring MME v1")
		fmt.Println("Starting configuring MME v1")
		startMme(OaiObj, CnAllInOneMode, buildSnap)
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

// StartSPGW is a wrapper for startSpgw
func StartSPGW(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) {
	// Start Mme
	OaiObj.Logger.Print("Starting configuring SPGW v1")
	fmt.Println("Starting configuring SPGW v1")
	startSpgw(OaiObj, CnAllInOneMode, buildSnap)
}

// StartSPGWCV2 is a wrapper for startSpgw
func StartSPGWCV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) {
	OaiObj.Logger.Print("Starting configuring SPGWC v2")
	fmt.Println("Starting configuring SPGWC v2")
	startSpgwcV2(OaiObj, CnAllInOneMode, buildSnap)
}

// StartSPGWUV2 is a wrapper for startSpgw
func StartSPGWUV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) {
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

//StartENB is a wrapper for configuring and starting OAI RAN services
func StartENB(OaiObj Oai, snapVersion string, build bool) {
	if snapVersion == "v1" {

		OaiObj.Logger.Print("Starting RAN for core network v1")
		fmt.Println("Starting RAN for core network v1")
		startENB(OaiObj, build)
		OaiObj.Logger.Print("RAN for core network v1 is started")
		fmt.Println("RAN for core network v1 is started")
	} else if snapVersion == "v2" {
		OaiObj.Logger.Print("Starting RAN for core network v2")
		fmt.Println("Starting RAN for core network v2")
		startENBV2(OaiObj, build)
		OaiObj.Logger.Print("RAN for core network v2 is started")
		fmt.Println("RAN for core network v2 is started")
	} else {
		OaiObj.Logger.Print("Error while trying to install oai core entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of ", snapVersion, " are: v1, v2")
		os.Exit(1)
	}
}

//StartCu is a wrapper for configuring and starting OAI CU service
func StartCu(OaiObj Oai) {
	OaiObj.Logger.Print("Starting CU")
	fmt.Println("Starting CU")
	startCu(OaiObj)
	OaiObj.Logger.Print("CU is started")
	fmt.Println("CU is started")
}

//StartDu is a wrapper for configuring and starting OAI CU service
func StartDu(OaiObj Oai) {
	OaiObj.Logger.Print("Starting DU")
	fmt.Println("Starting DU")
	startDu(OaiObj)
	OaiObj.Logger.Print("DU is started")
	fmt.Println("DU is started")
}

//StartRrc is a wrapper for configuring and starting OAI CU service
func StartRrc(OaiObj Oai) {
	OaiObj.Logger.Print("Starting RRC")
	fmt.Println("Starting RRC")
	startRcc(OaiObj)
	OaiObj.Logger.Print("RRC is started")
	fmt.Println("RRC is started")
}

//StartRru is a wrapper for configuring and starting OAI CU service
func StartRru(OaiObj Oai) {
	OaiObj.Logger.Print("Starting RRU")
	fmt.Println("Starting RRU")
	startRru(OaiObj)
	OaiObj.Logger.Print("RRU is started")
	fmt.Println("RRU is started")
}

//StopRan is a wrapper to stop RAN service
func StopRan(OaiObj Oai) {
	OaiObj.Logger.Print("Stopping RAN Service")
	fmt.Println("Stopping RAN Service")
	stopRan(OaiObj)
	OaiObj.Logger.Print("RAN Service is stopped")
	fmt.Println("RAN Service is stopped")
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
