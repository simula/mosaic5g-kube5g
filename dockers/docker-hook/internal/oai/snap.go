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
# file          snap.go
# brief 		install the required snaps: core, core18, etc
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
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
	snapName := "core"
	snapChannel := "edge"
	OaiObj.Logger.Print("Installing " + snapName)
	fmt.Println("Installing " + snapName)
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println(err)
	}
	//Loop until package is installed
	if !ret {
		// Install the snap
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+snapChannel)
		Snapfail := false
		for {
			OaiObj.Logger.Print("retStatus.Stderr=", retStatus.Stderr)
			fmt.Println("retStatus.Stderr=", retStatus.Stderr)

			if len(retStatus.Stderr) > 0 {
				if len(retStatus.Stderr[0]) > 0 {
					Snapfail = strings.Contains(retStatus.Stderr[0], "error")
				}
			}
			OaiObj.Logger.Print("Snapfail=", Snapfail)
			fmt.Println("Snapfail=", Snapfail)
			if Snapfail {
				OaiObj.Logger.Print("Wait for snapd being ready")
				fmt.Println("Wait for snapd being ready")
				time.Sleep(1 * time.Second)
				retStatus = util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+snapChannel)
			} else {
				OaiObj.Logger.Print("snapd is ready and " + snapName + " is installed")
				fmt.Println("snapd is ready and " + snapName + " is installed")
				break
			}
		}
	}

	//Install Core18
	snapName = "core18"
	snapChannel = "edge"
	OaiObj.Logger.Print("Installing " + snapName)
	fmt.Println("Installing " + snapName)
	ret, err = util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println(err)
	}
	//Loop until package is installed
	if !ret {
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+snapChannel)
		Snapfail := false
		for {
			OaiObj.Logger.Print("retStatus.Stderr=", retStatus.Stderr)
			fmt.Println("retStatus.Stderr=", retStatus.Stderr)

			if len(retStatus.Stderr) > 0 {
				if len(retStatus.Stderr[0]) > 0 {
					Snapfail = strings.Contains(retStatus.Stderr[0], "error")
				}
			}
			OaiObj.Logger.Print("Snapfail=", Snapfail)
			fmt.Println("Snapfail=", Snapfail)
			if Snapfail {
				OaiObj.Logger.Print("Wait for snapd being ready")
				fmt.Println("Wait for snapd being ready")
				time.Sleep(1 * time.Second)
				retStatus = util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+snapChannel)
			} else {
				OaiObj.Logger.Print("snapd is ready and " + snapName + " is installed")
				fmt.Println("snapd is ready and " + snapName + " is installed")
				break
			}
		}
	}

	// Install hello-world
	snapName = "hello-world"
	OaiObj.Logger.Print("Installing " + snapName)
	fmt.Println("Installing " + snapName)
	ret, err = util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print("err", err)
		fmt.Println("err", err)
	}
	if !ret {
		OaiObj.Logger.Print("installing " + snapName)
		fmt.Println("installing " + snapName)
		util.RunCmd(OaiObj.Logger, "snap", "install", snapName)
		OaiObj.Logger.Print(snapName + " is installed")
		fmt.Println(snapName + " is installed")
	}

}

// installOaicn : Install oai-cn snap
func installOaicn(OaiObj Oai, buildSnap bool, snapVersion string) {
	if snapVersion == "v1" {
		// get the realm of the network
		realm := OaiObj.Conf.OaiCn.V1[0].Realm.Default

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

		OaiObj.Logger.Print("hostname=", hostname, "\n fullDomainName= ", fullDomainName)
		fmt.Println("hostname=", hostname, "\n fullDomainName= ", fullDomainName)

		// Replace hosts
		util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")

		// checking if the snap cn is installed
		snapName := OaiObj.Conf.OaiCn.V1[0].Snap.Name
		OaiObj.Logger.Print("Installing " + snapName)
		fmt.Println("Installing " + snapName)

		ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
		if err != nil {
			OaiObj.Logger.Print(err)
			fmt.Println("error=", err)
		}
		if !ret {
			// Install the snap
			if OaiObj.Conf.OaiCn.V1[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V1[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V1[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V1[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V1[0].Snap.Channel)
			}
		} else {
			// Snap is already installed, refresh it if specified
			if OaiObj.Conf.OaiCn.V1[0].Snap.Refresh == true {
				if OaiObj.Conf.OaiCn.V1[0].Snap.Devmode == true {
					util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V1[0].Snap.Channel, "--devmode")
					OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V1[0].Snap.Channel + " in devmode")
				} else {
					util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V1[0].Snap.Channel)
					OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V1[0].Snap.Channel)
				}
			}
		}
	} else if snapVersion == "v2" {
		installOaiCnHssV2(OaiObj, buildSnap)
		installOaiCnMmeV2(OaiObj, buildSnap)
		installOaiCnSpgwcV2(OaiObj)
		installOaiCnSpgwuV2(OaiObj)
	} else {
		OaiObj.Logger.Print("Error while trying to install oai core entity: snap version", snapVersion, " is not recognized")
		OaiObj.Logger.Print("The allowed values of snapVersion are: v1, v2")
		os.Exit(1)
	}
}

// installOairan : Install oai-ran snap
func installOairan(OaiObj Oai) {
	// Install oai-ran snap
	OaiObj.Logger.Print("Installing oai-ran")
	fmt.Println("Installing oai-ran")

	snapName := OaiObj.Conf.OaiEnb[0].Snap.Name
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		OaiObj.Logger.Print("err", err)
		fmt.Println("err", err)
	}

	if !ret {
		// Install the snap

		if OaiObj.Conf.OaiEnb[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiEnb[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiEnb[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiEnb[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiEnb[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiEnb[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiEnb[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiEnb[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiEnb[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiEnb[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiEnb[0].Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiEnb[0].Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiEnb[0].Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}

	//Wait a moment, cn is not ready yet !
	OaiObj.Logger.Print("Wait 15 seconds... OK now cn should be ready")
	fmt.Println("Wait 15 seconds... OK now cn should be ready")
	// time.Sleep(15 * time.Second)

}

// installOaiHssV2 : Install oai-hss v2 snap for all-in-one mode
func installOaiCnHssV2(OaiObj Oai, buildSnap bool) {
	// get the snap name
	snapName := OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Name
	// get the realm of the network
	realm := OaiObj.Conf.OaiCn.V2[0].Realm.Default

	// Install oai-hss v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

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
		retStatus := util.RunCmd(OaiObj.Logger, "test", "-f", "./hosts_original")
		if retStatus.Exit != 0 {
			fmt.Println("File does not exist")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_original")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")
		} else {
			fmt.Println("File ./hosts_original already exist")
			util.RunCmd(OaiObj.Logger, "cp", "./hosts_original", "./hosts_new")
		}

		// the installation is all-in-one mode
		mmeIP := "127.0.0.1"

		hostname, _ := os.Hostname()
		fullDomainName := "1s/^/" + mmeIP + " " + hostname + "." + realm + " " + hostname + " mme\\n" +
			"127.0.0.1 " + hostname + "." + realm + " " + hostname + " hss \\n/"
		util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

		OaiObj.Logger.Print("hostname=", hostname)
		OaiObj.Logger.Print("fullDomainName=", fullDomainName)
		// Replace hosts
		util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")
	}

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Refresh == true {
			if OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiCn.V2[0].OaiHss.Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}
}

// installOaiHssV1 : Install oai-hss v1 snap for disaggregated mode
func installOaiHssV1(OaiObj Oai, buildSnap bool) {
	// get the snap to be installed
	snapName := OaiObj.Conf.OaiHss.V1[0].Snap.Name
	// get the realm of the network
	realm := OaiObj.Conf.OaiHss.V1[0].Realm.Default

	// Install oai-hss v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v1")
	fmt.Println("Installing " + snapName + " v1")

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
		retStatus := util.RunCmd(OaiObj.Logger, "test", "-f", "./hosts_original")
		if retStatus.Exit != 0 {
			fmt.Println("File does not exist")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_original")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")
		} else {
			fmt.Println("File ./hosts_original already exist")
			util.RunCmd(OaiObj.Logger, "cp", "./hosts_original", "./hosts_new")
		}

		mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiHss.V1[0].MmeServiceName)

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
			mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiHss.V1[0].MmeServiceName)
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

	// checking if the snap cn is installed
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}
	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiHss.V1[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiHss.V1[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiHss.V1[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiHss.V1[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiHss.V1[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiHss.V1[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiHss.V1[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiHss.V1[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiHss.V1[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiHss.V1[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiHss.V1[0].Snap.Channel)
			}
		}
	}
}

// installOaiHssV2 : Install oai-hss v2 snap for disaggregated mode
func installOaiHssV2(OaiObj Oai, buildSnap bool) {
	snapName := OaiObj.Conf.OaiHss.V2[0].Snap.Name
	realm := OaiObj.Conf.OaiHss.V2[0].Realm.Default

	// Install oai-hss v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

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
		retStatus := util.RunCmd(OaiObj.Logger, "test", "-f", "./hosts_original")
		if retStatus.Exit != 0 {
			fmt.Println("File does not exist")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_original")
			util.RunCmd(OaiObj.Logger, "cp", "/etc/hosts", "./hosts_new")
		} else {
			fmt.Println("File ./hosts_original already exist")
			util.RunCmd(OaiObj.Logger, "cp", "./hosts_original", "./hosts_new")
		}

		mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiHss.V2[0].MmeServiceName)

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
			mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiHss.V2[0].MmeServiceName)
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
	// check if the snap of oai-hss is already installed
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiHss.V2[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiHss.V2[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiHss.V2[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiHss.V2[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiHss.V2[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiHss.V2[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiHss.V2[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiHss.V2[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiHss.V2[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiHss.V2[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiHss.V2[0].Snap.Channel)
			}
		}
	}

	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiHss.V2[0].Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiHss.V2[0].Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}

}

// installOaiCnMmeV2 : Install oai-mme v2 snap for all-in-one mode
func installOaiCnMmeV2(OaiObj Oai, buildSnap bool) {
	snapName := OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Name
	realm := OaiObj.Conf.OaiCn.V2[0].Realm.Default

	// Install oai-mme v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

	OaiObj.Logger.Print("the realm of OAI is: ", realm)
	fmt.Println("the realm of OAI is: ", realm)

	OaiObj.Logger.Print("Configure hostname before installing ")
	fmt.Println("Configure hostname before installing ")

	// Install oai-mme v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel, "--devmode")
			OaiObj.Logger.Print("oaimme is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel)
			OaiObj.Logger.Print("oaimme is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Refresh == true {
			if OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel, "--devmode")
				OaiObj.Logger.Print("oaimme is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel)
				OaiObj.Logger.Print("oaimme is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiCn.V2[0].OaiMme.Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}
}

// installOaiMmeV1 : Install oai-mme v1 snap in dissaggregated mode
func installOaiMmeV1(OaiObj Oai, buildSnap bool) {
	snapName := OaiObj.Conf.OaiMme.V1[0].Snap.Name
	realm := OaiObj.Conf.OaiMme.V1[0].Realm.Default

	// Install oai-mme v1 snap
	OaiObj.Logger.Print("Installing " + snapName + " v1")
	fmt.Println("Installing " + snapName + " v1")

	OaiObj.Logger.Print("the realm of OAI is: ", realm)
	fmt.Println("the realm of OAI is: ", realm)

	OaiObj.Logger.Print("Configure hostname before installing ")
	fmt.Println("Configure hostname before installing ")

	// Copy hosts
	// if CnAllInOneMode == false {
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

	hssIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiMme.V1[0].HssServiceName)
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
			hssIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiMme.V1[0].HssServiceName)
		}
	}

	hostname, _ := os.Hostname()
	fullDomainName := "1s/^/" + "127.0.0.1" + " " + hostname + "." + realm + " " + hostname + " mme\\n" +
		hssIP + " " + hostname + "." + realm + " " + hostname + " hss \\n/"
	util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

	OaiObj.Logger.Print("hostname=", hostname)
	OaiObj.Logger.Print("fullDomainName=", fullDomainName)
	// Replace hosts
	util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")

	// }
	// Install oai-mme v1 snap
	OaiObj.Logger.Print("Installing " + snapName + " v1")
	fmt.Println("Installing " + snapName + " v1")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiMme.V1[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiMme.V1[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiMme.V1[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiMme.V1[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiMme.V1[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiMme.V1[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiMme.V1[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiMme.V1[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiMme.V1[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiMme.V1[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiMme.V1[0].Snap.Channel)
			}
		}
	}
}

// installOaiMmeV2 : Install oai-mme v2 snap in dissaggregated mode
func installOaiMmeV2(OaiObj Oai, buildSnap bool) {
	// func installOaiMmeV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) {
	snapName := OaiObj.Conf.OaiMme.V2[0].Snap.Name
	realm := OaiObj.Conf.OaiMme.V2[0].Realm.Default

	// Install oai-mme v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

	OaiObj.Logger.Print("the realm of OAI is: ", realm)
	fmt.Println("the realm of OAI is: ", realm)

	OaiObj.Logger.Print("Configure hostname before installing ")
	fmt.Println("Configure hostname before installing ")

	// Copy hosts
	// if CnAllInOneMode == false {
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

	hssIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiMme.V2[0].HssServiceName)
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
			hssIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiMme.V2[0].HssServiceName)
		}
	}

	hostname, _ := os.Hostname()
	fullDomainName := "1s/^/" + "127.0.0.1" + " " + hostname + "." + realm + " " + hostname + " mme\\n" +
		hssIP + " " + hostname + "." + realm + " " + hostname + " hss \\n/"
	util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

	OaiObj.Logger.Print("hostname=", hostname)
	OaiObj.Logger.Print("fullDomainName=", fullDomainName)
	// Replace hosts
	util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")

	// }
	// Install oai-mme v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiMme.V2[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiMme.V2[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiMme.V2[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiMme.V2[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiMme.V2[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiMme.V2[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiMme.V2[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiMme.V2[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiMme.V2[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiMme.V2[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiMme.V2[0].Snap.Channel)
			}
		}
	}

	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiMme.V2[0].Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiMme.V2[0].Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}
}

// installOaiSpgwV1 : Install oai-spgw v1 snap in dissagregated mode
func installOaiSpgwV1(OaiObj Oai, buildSnap bool) {
	// get the snap to be installed
	snapName := OaiObj.Conf.OaiSpgw.V1[0].Snap.Name
	// get the realm of the network
	realm := OaiObj.Conf.OaiSpgw.V1[0].Realm.Default

	// Install oai-hss v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v1")
	fmt.Println("Installing " + snapName + " v1")

	OaiObj.Logger.Print("the realm of OAI is: ", realm)
	fmt.Println("the realm of OAI is: ", realm)

	OaiObj.Logger.Print("Configure hostname before installing ")
	fmt.Println("Configure hostname before installing ")

	// Copy hosts
	// if CnAllInOneMode == false {
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

	hssIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiSpgw.V1[0].HssServiceName)
	mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiSpgw.V1[0].MmeServiceName)
	if buildSnap == true {
		hssIP = "127.0.0.1"
		mmeIP = "127.0.0.1"
	} else {
		// get the hssIP
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
			hssIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiSpgw.V1[0].HssServiceName)
		}
		// get the mmeIP
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
			mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiSpgw.V1[0].MmeServiceName)
		}

	}

	hostname, _ := os.Hostname()
	fullDomainName := "1s/^/" + mmeIP + " " + hostname + "." + realm + " " + hostname + " mme\\n" +
		hssIP + " " + hostname + "." + realm + " " + hostname + " hss \\n/"
	util.RunCmd(OaiObj.Logger, "sed", "-i", fullDomainName, "./hosts_new")

	OaiObj.Logger.Print("hostname=", hostname)
	OaiObj.Logger.Print("fullDomainName=", fullDomainName)
	// Replace hosts
	util.RunCmd(OaiObj.Logger, "cp", "-f", "./hosts_new", "/etc/hosts")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)

	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiSpgw.V1[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiSpgw.V1[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiSpgw.V1[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiSpgw.V1[0].Snap.Channel)
			}
		}
	}

}

// installOaiCnSpgwcV2 : Install oai-spgwc v2 snap in all-in-one mode
func installOaiCnSpgwcV2(OaiObj Oai) {
	snapName := OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Name

	// Install oai-spgwc v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)

	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Refresh == true {
			if OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}
}

// installOaiSpgwcV2 : Install oai-spgwc v2 snap in dissagregated mode
func installOaiSpgwcV2(OaiObj Oai) {
	snapName := OaiObj.Conf.OaiSpgwc.V2[0].Snap.Name

	// Install oai-spgwc v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)

	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiSpgwc.V2[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiSpgwc.V2[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiSpgwc.V2[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiSpgwc.V2[0].Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiSpgwc.V2[0].Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiSpgwc.V2[0].Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}

}

// installOaiCnSpgwuV2 : Install oai-spgwu v2 snap in all-in-one mode
func installOaiCnSpgwuV2(OaiObj Oai) {
	snapName := OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Name

	// Install oai-spgwu v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Refresh == true {
			if OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiCn.V2[0].OaiSpgwu.Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}

}

// installOaiSpgwuV2 : Install oai-spgwu v2 snap in dissagregated mode
func installOaiSpgwuV2(OaiObj Oai) {
	snapName := OaiObj.Conf.OaiSpgwu.V2[0].Snap.Name

	// Install oai-spgwu v2 snap
	OaiObj.Logger.Print("Installing " + snapName + " v2")
	fmt.Println("Installing " + snapName + " v2")

	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println("error=", err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.OaiSpgwu.V2[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.OaiSpgwu.V2[0].Snap.Refresh == true {
			if OaiObj.Conf.OaiSpgwu.V2[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.OaiSpgwu.V2[0].Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.OaiSpgwu.V2[0].Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.OaiSpgwu.V2[0].Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}

}

// installFlexRAN : Install FlexRAN snap
func installFlexRAN(OaiObj Oai) {
	snapName := OaiObj.Conf.Flexran[0].Snap.Name

	// Install FlexRAN snap
	OaiObj.Logger.Print("Installing " + snapName)
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
		fmt.Println(err)
	}

	if !ret {
		// Install the snap

		if OaiObj.Conf.Flexran[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.Flexran[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.Flexran[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.Flexran[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.Flexran[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.Flexran[0].Snap.Refresh == true {
			if OaiObj.Conf.Flexran[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.Flexran[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.Flexran[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.Flexran[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.Flexran[0].Snap.Channel)
			}
		}
	}
	// enable the plugs
	var flexranPlugs []string = []string{"network", "network-control", "log-observe", "process-control", "cpu-control", "network-observe", "network-bind"}
	var permission string
	if len(OaiObj.Conf.Flexran[0].Snap.Plugs) > 0 {
		flexranPlugs = OaiObj.Conf.Flexran[0].Snap.Plugs
	}
	for i := 0; i < len(flexranPlugs); i++ {
		permission = snapName + ":" + flexranPlugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}

	//Wait a moment, cn is not ready yet !
	OaiObj.Logger.Print("Wait 5 seconds... OK now " + snapName + " should be ready")
	time.Sleep(5 * time.Second)

}

// installMEC : Install LL-MEC snap
func installMEC(OaiObj Oai) {
	snapName := OaiObj.Conf.LlMec[0].Snap.Name
	// Install LL-MEC snap
	OaiObj.Logger.Print("Installing " + snapName)
	ret, err := util.CheckSnapPackageExist(OaiObj.Logger, snapName)
	if err != nil {
		OaiObj.Logger.Print(err)
	}

	if !ret {
		// Install the snap
		if OaiObj.Conf.LlMec[0].Snap.Devmode == true {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.LlMec[0].Snap.Channel, "--devmode")
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.LlMec[0].Snap.Channel + " in devmode")
		} else {
			util.RunCmd(OaiObj.Logger, "snap", "install", snapName, "--channel="+OaiObj.Conf.LlMec[0].Snap.Channel)
			OaiObj.Logger.Print(snapName + " is installed from the channel " + OaiObj.Conf.LlMec[0].Snap.Channel)
		}
	} else {
		// Snap is already installed, refresh it if specified
		if OaiObj.Conf.LlMec[0].Snap.Refresh == true {
			if OaiObj.Conf.LlMec[0].Snap.Devmode == true {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.LlMec[0].Snap.Channel, "--devmode")
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.LlMec[0].Snap.Channel + " in devmode")
			} else {
				util.RunCmd(OaiObj.Logger, "snap", "refresh", snapName, "--channel="+OaiObj.Conf.LlMec[0].Snap.Channel)
				OaiObj.Logger.Print(snapName + " is refreshed from the channel " + OaiObj.Conf.LlMec[0].Snap.Channel)
			}
		}
	}
	// enable the plugs
	var permission string
	for i := 0; i < len(OaiObj.Conf.LlMec[0].Snap.Plugs); i++ {
		permission = snapName + ":" + OaiObj.Conf.LlMec[0].Snap.Plugs[i]
		OaiObj.Logger.Print("giving the permission; " + permission)
		retStatus := util.RunCmd(OaiObj.Logger, "snap", "connect", permission)
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error while giving the permission "+permission+"\n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Successfully giving the permission "+permission+": \n", retStatus.Stdout)
		}
	}
	//Wait a moment, cn is not ready yet !
	OaiObj.Logger.Print("Wait 5 seconds... OK now " + snapName + " should be ready")
	time.Sleep(5 * time.Second)

}
