package oai

import (
	"fmt"
)

//InstallSnap is a wrapper function for installSnapCore
func InstallSnap(OaiObj Oai) {
	// Install Snap Core
	installSnapCore(OaiObj)
}

//InstallCN is a wrapper for installing OAI CN
func InstallCN(OaiObj Oai) {

	// Install oai-cn snap
	installOaicn(OaiObj)

}

// StartCN is a wrapper for configuring and starting OAI CN services
func StartCN(OaiObj Oai, CnAllInOneMode bool, build bool) {
	// Start HSS
	OaiObj.Logger.Print("Starting configuring HSS")
	fmt.Println("Starting configuring HSS")
	startHss(OaiObj, CnAllInOneMode, build)
	// Start MME
	OaiObj.Logger.Print("Starting configuring MME")
	fmt.Println("Starting configuring MME")
	startMme(OaiObj, CnAllInOneMode, build)
	// Start SPGW
	OaiObj.Logger.Print("Starting configuring SPGW")
	fmt.Println("Starting configuring SPGW")
	startSpgw(OaiObj, CnAllInOneMode)
}

// StartHSS is a wrapper for startHss
func StartHSS(OaiObj Oai, CnAllInOneMode bool, build bool) {
	// Start HSS
	startHss(OaiObj, CnAllInOneMode, build)
}

// StartMME is a wrapper for startMme
func StartMME(OaiObj Oai, CnAllInOneMode bool, build bool) {
	// Start Mme
	startMme(OaiObj, CnAllInOneMode, build)
}

// StartSPGW is a wrapper for startSpgw
func StartSPGW(OaiObj Oai, CnAllInOneMode bool) {
	// Start Mme
	startSpgw(OaiObj, CnAllInOneMode)
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
func StartENB(OaiObj Oai) {
	OaiObj.Logger.Print("Starting RAN")
	fmt.Println("Starting RAN")
	startENB(OaiObj)
	OaiObj.Logger.Print("RAN is started")
	fmt.Println("RAN is started")
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
