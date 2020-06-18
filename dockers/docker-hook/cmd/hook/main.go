package main

// This APP is made for installing snaps in docker and
// handle the configurations
// Author: Osama Arouk, Kevin Hsi-Ping Hsu
import (
	"flag"
	"fmt"
	"mosaic5g/docker-hook/internal/oai"
)

const (
	logPath     = "/root/hook.log"
	confPath    = "/root/config/conf.yaml"
	ConfOaiPath = "/root/config/oai.yaml"
)

func main() {
	// Initialize oai struct
	OaiObj := oai.Oai{}
	err := OaiObj.Init(logPath, confPath, ConfOaiPath)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		OaiObj.Logger.Print("Init of OAI is successful")
		fmt.Println("Init of OAI is successful")
	}

	// Parse input flags
	OaiObj.Logger.Print("Starting parsing input parameters")
	fmt.Println("Starting parsing input parameters")
	installCN := flag.Bool("installCN", false, "a bool")
	installRAN := flag.Bool("installRAN", false, "a bool")
	installHSS := flag.Bool("installHSS", false, "a bool")
	installMME := flag.Bool("installMME", false, "a bool")
	installSPGW := flag.Bool("installSPGW", false, "a bool")
	installSPGWC := flag.Bool("installSPGWC", false, "a bool")
	installSPGWU := flag.Bool("installSPGWU", false, "a bool")
	installFlexRAN := flag.Bool("installFlexRAN", false, "a bool")
	installMEC := flag.Bool("installMEC", false, "a bool")
	buildImage := flag.Bool("build", false, "a bool value to define that the current setup is to build the docker image.")
	var snapVersion string
	flag.StringVar(&snapVersion, "snapVersion", "v2", "a string value to specify the snap version that will be used to build the docker image. Valid values: v1, v2")
	flag.Parse()
	//Install snap core
	OaiObj.Logger.Print("Installing snap")
	fmt.Println("Installing snap")
	oai.InstallSnap(OaiObj)

	// Decide actions based on flags
	CnAllInOneMode := true
	buildSnap := false
	if *buildImage {
		buildSnap = true
	}
	if *installCN {
		OaiObj.Logger.Print("Installing CN")
		fmt.Println("Installing CN")
		oai.InstallCN(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		OaiObj.Logger.Print("Starting CN")
		fmt.Println("Starting CN")
		oai.StartCN(OaiObj, CnAllInOneMode, buildSnap, snapVersion)
		OaiObj.Logger.Print("CN is started: exit")
	} else if *installRAN {
		OaiObj.Logger.Print("Installing RAN")
		fmt.Println("Installing RAN")
		oai.InstallRAN(OaiObj)
		// Define the functionality of the snap: oai-enb, oai-cu, oai-du, oai-rcc, oai-rru
		RanNodeFunction := OaiObj.Conf.NodeFunction.Default
		fmt.Println(RanNodeFunction)
		if (RanNodeFunction == "ENB") || (RanNodeFunction == "enb") {
			fmt.Println("Starting RAN eNB")
			oai.StartENB(OaiObj, snapVersion, buildSnap)
		} else if (RanNodeFunction == "CU") || (RanNodeFunction == "cu") {
			OaiObj.Logger.Print("Starting RAN CU")
			fmt.Println("Starting RAN CU")
			oai.StartCu(OaiObj)
		} else if (RanNodeFunction == "DU") || (RanNodeFunction == "du") {
			OaiObj.Logger.Print("Starting RAN DU")
			fmt.Println("Starting RAN DU")
			oai.StartDu(OaiObj)
		} else if (RanNodeFunction == "RRC") || (RanNodeFunction == "rrc") {
			OaiObj.Logger.Print("Starting RAN RRC")
			fmt.Println("Starting RAN RRC")
			oai.StartRrc(OaiObj)
		} else if (RanNodeFunction == "RRU") || (RanNodeFunction == "rru") {
			OaiObj.Logger.Print("Starting RAN RRU")
			fmt.Println("Starting RAN RRU")
			oai.StartRru(OaiObj)
		} else if (RanNodeFunction == "STOP") || (RanNodeFunction == "stop") {
			OaiObj.Logger.Print("Stopping RAN Service")
			fmt.Println("Stopping RAN Service")
			oai.StopRan(OaiObj)
		} else {
			OaiObj.Logger.Print("Error, unkown node function: ", RanNodeFunction, "Starting RAN eNB")
			fmt.Println("Error, unkown node function: ", RanNodeFunction, "Starting RAN eNB")
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
		oai.StartFlexRAN(OaiObj)
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
