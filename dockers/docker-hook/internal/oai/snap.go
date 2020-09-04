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

// This file is for installing the snaps inside docker images
// Author: Osama Arouk, Kevin Hsi-Ping Hsu
*/

package oai

import (
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"os"
	"strings"
	"time"
)

// installSnapCore : Install Core
func installSnapCore(OaiObj Oai) {
	//Install Core
	util.PrintFunc(OaiObj.Logger, "Installing core")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "core")
	if err != nil {
		util.PrintFunc(OaiObj.Logger, err)
	}
	//Loop until package is installed
	if !ret {
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "install", "core", "--channel=edge")
		Snapfail := false
		for {
			OaiObj.Logger.Print("retStatus.Stderr=", retStatus.Stderr)

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
func installOaicn(OaiObj Oai, CnAllInOneMode bool, buildSnap bool, snapVersion string) {

	// Install oai-cn snap
	if snapVersion == "v1" {
		// realm := "openair4G.eur"
		realm := OaiObj.Conf.Realm.Default
		// realm := OaiObj.ConfOai.Realm.Default

		OaiObj.Logger.Print("the realm of OAI is: ", realm)
		fmt.Println("the realm of OAI is: ", realm)

		OaiObj.Logger.Print("Configure hostname before installing ")
		fmt.Println("Configure hostname before installing ")
		// Copy hosts
		util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")
		hostname, _ := os.Hostname()
		fullDomainName := "1s/^/127.0.0.1 " + hostname + "." + realm + " " + hostname + " mme\\n" +
			"127.0.0.1 " + hostname + "." + realm + " " + hostname + " hss \\n/"
		util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

		OaiObj.Logger.Print("hostname=", hostname)
		OaiObj.Logger.Print("fullDomainName=", fullDomainName)

		// Replace hosts
		util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")

		OaiObj.Logger.Print("Installing oai-cn")
		fmt.Println("Installing oai-cn")
		ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-cn")
		if err != nil {
			OaiObj.Logger.Print(err)
			fmt.Println("error=", err)
		}
		if !ret {
			util.RunCmd(OaiObj.Logger, "snap", "install", "oai-cn", "--channel=edge", "--devmode")
			OaiObj.Logger.Print("Installing oai-cn")
			fmt.Println("Installing oai-cn")
		}
	} else if snapVersion == "v2" {
		installOaiHssV2(OaiObj, CnAllInOneMode, buildSnap)
		installOaiMmeV2(OaiObj, CnAllInOneMode, buildSnap)
		installOaiSpgwcV2(OaiObj)
		installOaiSpgwuV2(OaiObj)
	} else {
		OaiObj.Logger.Print("Error while trying to install oai core entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of ", snapVersion, " are: v1, v2")
		os.Exit(1)
	}

}

// installOairan : Install oai-ran snap
func installOairan(OaiObj Oai) {
	// Install oai-ran snap
	util.PrintFunc(OaiObj.Logger, "Installing oai-ran")
	snapExist, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-ran")
	if err != nil {
		util.PrintFunc(OaiObj.Logger, err)
	}

	if !snapExist {

		if OaiObj.ConfOaiRan.OaianConf.Snap.Channel == "stable" {
			util.RunCmd(OaiObj.Logger, "snap", "install", "oai-ran")
			util.PrintFunc(OaiObj.Logger, "oairan stable is installed")
		} else {
			if OaiObj.ConfOaiRan.OaianConf.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-ran", "--channel=edge", "--devmode")
				util.PrintFunc(OaiObj.Logger, "oairan devmode is installed")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-ran", "--channel=edge", "--jailmode")
				util.PrintFunc(OaiObj.Logger, "oairan jailmode is installed")
			}
		}
	} else {
		if OaiObj.ConfOaiRan.OaianConf.Snap.Refresh == true {
			if OaiObj.ConfOaiRan.OaianConf.Snap.Channel == "stable" {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", "oai-ran")
				util.PrintFunc(OaiObj.Logger, "oairan stable is refresh")
			} else {
				if OaiObj.ConfOaiRan.OaianConf.Snap.Devmode == true {
					util.RunCmd(OaiObj.Logger, "snap", "refresh", "oai-ran", "--channel=edge", "--devmode")
					util.PrintFunc(OaiObj.Logger, "oairan devmode is refresh")
				} else {
					util.RunCmd(OaiObj.Logger, "snap", "refresh", "oai-ran", "--channel=edge", "--jailmode")
					util.PrintFunc(OaiObj.Logger, "oairan jailmode is refresh")
				}
			}

		}
	}
}

// installOaiHssV2 : Install oai-hss v2 snap
func installOaiHssV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) {
	// Install oai-hss v2 snap
	OaiObj.Logger.Print("Installing oai-hss v2")
	fmt.Println("Installing oai-hss v2")

	realm := OaiObj.Conf.Realm.Default
	OaiObj.Logger.Print("the realm of OAI is: ", realm)
	fmt.Println("the realm of OAI is: ", realm)

	OaiObj.Logger.Print("Configure hostname before installing ")
	fmt.Println("Configure hostname before installing ")

	// Copy hosts
	if buildSnap == true {
		retStatus := util.RunCmd(OaiObj.Logger, "test", "-f", "./hosts_original")
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("File does not exist")
			fmt.Println("File does not exist")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_original")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")

		} else {
			OaiObj.Logger.Print("File ./hosts_original already exist")
			fmt.Println("File ./hosts_original already exist")
			util.RunCmd(OaiObj.Logger, "cp", "./hosts_original", "./hosts_new")
		}
		hostname, _ := os.Hostname()
		fullDomainName := "1s/^/127.0.0.1 " + hostname + "." + realm + " " + hostname + " mme\\n" +
			"127.0.0.1 " + hostname + "." + realm + " " + hostname + " hss \\n/"
		util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

		OaiObj.Logger.Print("hostname=", hostname)
		OaiObj.Logger.Print("fullDomainName=", fullDomainName)
		// Replace hosts
		util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")

	} else {
		/////////////////////////////////////////////
		retStatus := util.RunCmd(OaiObj.Logger, "test", "-f", "./hosts_original")
		if retStatus.Exit != 0 {
			fmt.Println("File does not exist")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_original")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")
		} else {
			fmt.Println("File ./hosts_original already exist")
			util.RunCmd(OaiObj.Logger, "cp", "./hosts_original", "./hosts_new")
		}
		/////////////////////////////////////////////
		mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.MmeDomainName)
		if CnAllInOneMode == true {
			mmeIP = "127.0.0.1"
		} else {
			for {
				if err != nil {
					OaiObj.Logger.Print(err)
				} else {
					hostNameMme, err := net.LookupHost(mmeIP)
					if len(hostNameMme) > 0 {
						break
					} else {
						OaiObj.Logger.Print(err)
					}
				}
				OaiObj.Logger.Print("Valid ip address for mme not yet retreived")
				time.Sleep(1 * time.Second)
				mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.MmeDomainName)
			}
		}
		hostname, _ := os.Hostname()
		fullDomainName := "1s/^/" + mmeIP + " " + hostname + "." + realm + " " + hostname + " mme\\n" +
			"127.0.0.1 " + hostname + "." + realm + " " + hostname + " hss \\n/"
		util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

		OaiObj.Logger.Print("hostname=", hostname)
		OaiObj.Logger.Print("fullDomainName=", fullDomainName)
		// Replace hosts
		util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")
	}
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-hss")
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	if !ret {
		if OaiObj.Conf.Snap.Channel == "stable" {
			util.RunCmd(OaiObj.Logger, "snap", "install", "oai-hss")
			OaiObj.Logger.Print("oai-hss stable is installed")
			fmt.Println("oai-hss stable is installed")
		} else {
			if OaiObj.Conf.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-hss", "--channel=edge", "--devmode")
				OaiObj.Logger.Print("oai-hss devmode is installed")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-hss", "--channel=edge", "--jailmode")
				OaiObj.Logger.Print("oai-hss jailmode is installed")
			}
		}
	}
}

// installOaiMmeV2 : Install oai-hss v2 snap
func installOaiMmeV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) {
	// Install oai-hss v2 snap
	OaiObj.Logger.Print("Installing oai-mme v2")
	fmt.Println("Installing oai-mme v2")

	realm := OaiObj.Conf.Realm.Default
	OaiObj.Logger.Print("the realm of OAI is: ", realm)
	fmt.Println("the realm of OAI is: ", realm)

	OaiObj.Logger.Print("Configure hostname before installing ")
	fmt.Println("Configure hostname before installing ")

	// Copy hosts
	if CnAllInOneMode == false {
		retStatus := util.RunCmd(OaiObj.Logger, "test", "-f", "./hosts_original")
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_original")
			fmt.Println("File does not exist")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_original")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")
		} else {
			OaiObj.Logger.Print("File ./hosts_original already exist")
			fmt.Println("File ./hosts_original already exist")
			util.RunCmd(OaiObj.Logger, "cp", "./hosts_original", "./hosts_new")
		}
		/////////////////////////////////////////////
		hssIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.HssDomainName)
		if buildSnap == true {
			hssIP = "127.0.0.1"
		} else {
			for {
				if err != nil {
					OaiObj.Logger.Print(err)
					fmt.Println(err)
				} else {
					hostNameHss, err := net.LookupHost(hssIP)
					if len(hostNameHss) > 0 {
						break
					} else {
						OaiObj.Logger.Print(err)
						fmt.Println(err)
					}
				}
				OaiObj.Logger.Print("Valid ip address for hss not yet retreived")
				fmt.Println("Valid ip address for hss not yet retreived")
				time.Sleep(1 * time.Second)
				hssIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.HssDomainName)
			}
		}
		/////////////////////////////////////////////
		hostname, _ := os.Hostname()
		fullDomainName := "1s/^/" + "127.0.0.1" + " " + hostname + "." + realm + " " + hostname + " mme\\n" +
			hssIP + " " + hostname + "." + realm + " " + hostname + " hss \\n/"
		util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

		OaiObj.Logger.Print("hostname=", hostname)
		OaiObj.Logger.Print("fullDomainName=", fullDomainName)
		// Replace hosts
		util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")

	}
	// Install oai-mme v2 snap
	OaiObj.Logger.Print("Installing oai-mme v2")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-mme")
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	if !ret {
		if OaiObj.Conf.Snap.Channel == "stable" {
			util.RunCmd(OaiObj.Logger, "snap", "install", "oai-mme")
			OaiObj.Logger.Print("oai-mme stable is installed")
			fmt.Println("oai-mme stable is installed")
		} else {
			if OaiObj.Conf.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-mme", "--channel=edge", "--devmode")
				OaiObj.Logger.Print("oai-mme devmode is installed")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-mme", "--channel=edge", "--jailmode")
				OaiObj.Logger.Print("oai-mme jailmode is installed")
			}
		}
	}
}

// installOaiSpgwcV2 : Install oai-hss v2 snap
func installOaiSpgwcV2(OaiObj Oai) {
	// Install oai-spgwc v2 snap
	OaiObj.Logger.Print("Installing oai-spgwc v2")
	fmt.Println("Installing oai-spgwc v2")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-spgwc")
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	if !ret {
		if OaiObj.Conf.Snap.Channel == "stable" {
			util.RunCmd(OaiObj.Logger, "snap", "install", "oai-spgwc")
			OaiObj.Logger.Print("oai-spgwc stable is installed")
			fmt.Println("oai-spgwc stable is installed")
		} else {
			if OaiObj.Conf.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-spgwc", "--channel=edge", "--devmode")
				OaiObj.Logger.Print("oai-spgwc devmode is installed")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-spgwc", "--channel=edge", "--jailmode")
				OaiObj.Logger.Print("oai-spgwc jailmode is installed")
			}
		}
	}
}

// installOaiSpgwuV2 : Install oai-spgwu v2 snap
func installOaiSpgwuV2(OaiObj Oai) {
	// Install oai-spgwu v2 snap
	OaiObj.Logger.Print("Installing oai-spgwu v2")
	fmt.Println("Installing oai-spgwu v2")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, "oai-spgwu")
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	if !ret {
		if OaiObj.Conf.Snap.Channel == "stable" {
			util.RunCmd(OaiObj.Logger, "snap", "install", "oai-spgwu")
			OaiObj.Logger.Print("oai-spgwu stable is installed")
			fmt.Println("oai-spgwu stable is installed")
		} else {
			if OaiObj.Conf.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-spgwu", "--channel=edge", "--devmode")
				OaiObj.Logger.Print("oai-spgwu devmode is installed")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "install", "oai-spgwu", "--channel=edge", "--jailmode")
				OaiObj.Logger.Print("oai-spgwu jailmode is installed")
			}
		}
	}
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
	waitTime := 5
	OaiObj.Logger.Print("Wait " + string(waitTime) + " seconds... OK now flexran should be ready")
	time.Sleep(time.Duration(waitTime) * time.Second)

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
