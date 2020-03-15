package oai

import (
	"docker-hook/internal/pkg/util"
	"fmt"
	"os"
	"strings"
	"time"
)

// installSnapCore : Install Core
func installSnapCore(OaiObj Oai) {
	//Install Core
	OaiObj.Logger.Print("Installing core")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "core")
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	//Loop until package is installed
	if !ret {
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "install", "core", "--channel=edge")
		Snapfail := false
		for {
			OaiObj.Logger.Print("loop for")
			OaiObj.Logger.Print("retStatus.Stderr=", retStatus.Stderr)
			//OaiObj.Logger.Print("retStatus.Stderr[0]=", retStatus.Stderr[0])
			//OaiObj.Logger.Print("len(retStatus.Stderr[0])=", len(retStatus.Stderr[0]))

			if len(retStatus.Stderr) > 0 {
				if len(retStatus.Stderr[0]) > 0 {
					Snapfail = strings.Contains(retStatus.Stderr[0], "error")
				}
			}
			OaiObj.Logger.Print("Snapfail=", Snapfail)
			if Snapfail {
				OaiObj.Logger.Print("Wait for snapd being ready")
				time.Sleep(1 * time.Second)
				retStatus = util.RunCmd(OaiObj.Logger, "snap", "install", "core", "--channel=edge")
			} else {
				OaiObj.Logger.Print("snapd is ready and core is installed")
				break
			}
		}
	}

	// Install hello-world
	OaiObj.Logger.Print("Installing hello-world")
	ret, err = util.CheckSnapPackageExist(OaiObj.Logger, "hello-world")
	if err != nil {
		OaiObj.Logger.Print("err", err)
	}
	if !ret {
		OaiObj.Logger.Print("installing hello-world")
		util.RunCmd(OaiObj.Logger, "snap", "install", "hello-world")
		OaiObj.Logger.Print("hello-world is installed")
	}

}

// installOaicn : Install oai-cn snap
func installOaicn(OaiObj Oai) {
	OaiObj.Logger.Print("Configure hostname before installing ")
	fmt.Println("Configure hostname before installing ")
	// Copy hosts
	util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")
	hostname, _ := os.Hostname()
	fullDomainName := "1s/^/127.0.0.1 " + hostname + ".openair4G.eur " + hostname + " hss\\n127.0.0.1 " + hostname + ".openair4G.eur " + hostname + " mme \\n/"
	util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

	fmt.Println("hostname=", hostname)
	fmt.Println("fullDomainName=", fullDomainName)
	// Replace hosts
	util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")
	// Install oai-cn snap
	OaiObj.Logger.Print("Installing oai-cn")
	fmt.Println("Installing oai-cn")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-cn")
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}
	if !ret {
		util.RunCmd(OaiObj.Logger, "snap", "install", "oai-cn", "--channel=edge", "--devmode")
		fmt.Println("snap install oai-cn")
	}

}

// installOairan : Install oai-ran snap
func installOairan(OaiObj Oai) {
	// Install oai-ran snap
	OaiObj.Logger.Print("Installing oai-ran")
	fmt.Println("Installing oai-ran")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-ran")
	if err != nil {
		OaiObj.Logger.Print(err)
		OaiObj.Logger.Print("err", err)
		fmt.Println("err", err)
	}
	if !ret {
		OaiObj.Logger.Print("installing oairan devmode")
		fmt.Println("installing oairan devmode")
		util.RunCmd(OaiObj.Logger, "snap", "install", "oai-ran", "--channel=edge", "--devmode")
		OaiObj.Logger.Print("oairan devmode is installed")
		fmt.Println("oairan devmode is installed")
	}
	//Wait a moment, cn is not ready yet !
	OaiObj.Logger.Print("Wait 15 seconds... OK now cn should be ready")
	fmt.Println("Wait 15 seconds... OK now cn should be ready")
	time.Sleep(15 * time.Second)

}

// installFlexRAN : Install FlexRAN snap
func installFlexRAN(OaiObj Oai) {
	// Install FlexRAN snap
	OaiObj.Logger.Print("Installing FlexRAN")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "flexran")
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	if !ret {
		util.RunCmd(OaiObj.Logger, "snap", "install", "flexran", "--channel=edge", "--devmode")
	}
	//Wait a moment, cn is not ready yet !
	OaiObj.Logger.Print("Wait 5 seconds... OK now flexran should be ready")
	time.Sleep(5 * time.Second)

}

// installMEC : Install LL-MEC snap
func installMEC(OaiObj Oai) {
	// Install LL-MEC snap
	OaiObj.Logger.Print("Installing LL-MEC")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "ll-mec")
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	if !ret {
		util.RunCmd(OaiObj.Logger, "snap", "install", "ll-mec", "--channel=edge", "--devmode")
	}
	//Wait a moment, cn is not ready yet !
	OaiObj.Logger.Print("Wait 5 seconds... OK now ll-mec should be ready")
	time.Sleep(5 * time.Second)

}
