package main

// This APP is made for installing snaps in docker and
// handle the configurations
// Author: Kevin Hsi-Ping Hsu
import (
	"docker-hook/internal/oai"
	"flag"
	"fmt"
)

const (
	logPath  = "/root/hook.log"
	confPath = "/root/config/conf.yaml"
)

func main() {
	// Initialize oai struct
	OaiObj := oai.Oai{}
	err := OaiObj.Init(logPath, confPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse input flags
	installCN := flag.Bool("installCN", false, "a bool")
	installRAN := flag.Bool("installRAN", false, "a bool")
	installHSS := flag.Bool("installHSS", false, "a bool")
	installMME := flag.Bool("installMME", false, "a bool")
	installSPGW := flag.Bool("installSPGW", false, "a bool")
	installFlexRAN := flag.Bool("installFlexRAN", false, "a bool")
	installMEC := flag.Bool("installMEC", false, "a bool")
	flag.Parse()
	//Install snap core
	OaiObj.Logger.Print("Installing snap")
	fmt.Println("Installing snap")
	oai.InstallSnap(OaiObj)
	// Decide actions based on flags
	if *installCN {
		OaiObj.Logger.Print("Installing CN")
		fmt.Println("Installing CN")
		oai.InstallCN(OaiObj)
		OaiObj.Logger.Print("Starting CN")
		fmt.Println("Starting CN")
		oai.StartCN(OaiObj)
	} else if *installRAN {
		OaiObj.Logger.Print("Installing RAN")
		fmt.Println("Installing RAN")
		oai.InstallRAN(OaiObj)
		OaiObj.Logger.Print("Starting RAN")
		fmt.Println("Starting RAN")
		oai.StartENB(OaiObj)
	} else if *installHSS {
		oai.InstallCN(OaiObj)
		oai.StartHSS(OaiObj)
	} else if *installMME {
		oai.InstallCN(OaiObj)
		oai.StartMME(OaiObj)
	} else if *installSPGW {
		oai.InstallCN(OaiObj)
		oai.StartSPGW(OaiObj)
	} else if *installFlexRAN {
		oai.InstallFlexRAN(OaiObj)
		oai.StartFlexRAN(OaiObj)
	} else if *installMEC {
		oai.InstallMEC(OaiObj)
	} else {
		fmt.Println("This should only be executed in container!!")
		return
	}

	// Give a hello when program ends
	OaiObj.Logger.Print("End of hook")
	OaiObj.Clean()
}
