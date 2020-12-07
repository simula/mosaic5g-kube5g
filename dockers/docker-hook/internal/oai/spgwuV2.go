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
# file          spgwuV2.go
# brief 		configure the snap of oai-spgwu v2, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package oai

import (
	"errors"
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"strings"
	"time"
)

// initSpgwuV2 : Init SPGW
func initSpgwuV2(OaiObj Oai) error {
	return nil
}

// configSpgwuV2 : Config oai-spgw
func configSpgwuV2(OaiObj Oai) error {
	return nil
}

// StartSpgwuV2 : Start SPGW as a daemon
func startSpgwuV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring OAI-SPGWu V2")
	OaiObj.Logger.Print("Starting configuration of OAI-SPGWu V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-spgwu.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	spgwBin := strings.Join([]string{snapBinaryPath, "oai-spgwu"}, "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "conf-get"}, "."))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]

	spgwConf := strings.Join([]string{confPath, "spgwu.conf"}, "/")

	if buildSnap == false {
		// Init spgwu
		OaiObj.Logger.Print("Start Init of oai-spgwu")
		fmt.Println("Start Init of oai-spgwu")
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Init of oai-spgwu is failed")
				fmt.Println("Init of oai-spgwu is failed")
			} else {
				OaiObj.Logger.Print("Init of oai-spgwu is successful")
				fmt.Println("Init of oai-spgwu is successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Init of oai-spgwu")
			fmt.Println("Retrying to Init of oai-spgwu")
			retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		}

		// Configure oai-spgw
		OaiObj.Logger.Print("Configure of oai-spgwu")
		fmt.Println("Configure of oai-spgwu")

		// Get interface IP and outbound interface
		interfaceIP := util.GetOutboundIP()
		outInterface, _ := util.GetInterfaceByIP(interfaceIP)
		// INTERFACE_NAME of S1U_S12_S4_UP
		sedCommand := "59s:\".*;:\"" + outInterface + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		// IPV4_ADDRESS of S1U_S12_S4_UP
		sedCommand = "60s:\".*;:\"" + interfaceIP + "/24\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)

		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/wlp2s0/"+outInterface+"/g", spgwConf)
		if CnAllInOneMode == false {
			// Get interface IP and outbound interface
			interfaceIP := util.GetOutboundIP()
			outInterface, _ := util.GetInterfaceByIP(interfaceIP)
			// INTERFACE_NAME of SX
			sedCommand := "72s:\".*;:\"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set INTERFACE_NAME of SX in " + spgwConf + " failed")
			}
			// IPV4_ADDRESS of SX
			sedCommand = "73s:\".*;:\"" + interfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set IPV4_ADDRESS of SX in " + spgwConf + " failed")
			}

			spgwcIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiSpgwu.V2[0].SpgwcServiceName)

			if buildSnap == true {
				spgwcIP = "127.0.12.1"
			} else {
				for {
					if err != nil {
						OaiObj.Logger.Print(err)
					} else {
						hostNameSpgwc, err := net.LookupHost(spgwcIP)
						if len(hostNameSpgwc) > 0 {
							break
						} else {
							OaiObj.Logger.Print(err)
						}
					}
					OaiObj.Logger.Print("Valid ip address for oai-spgwc not yet retreived")
					time.Sleep(1 * time.Second)
					spgwcIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiSpgwu.V2[0].SpgwcServiceName)
				}
			}
			sedCommand = "s:IPV4_ADDRESS=\"127.0.12.1.*;:IPV4_ADDRESS=\"" + spgwcIP + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				fmt.Println("Set IPV4_ADDRESS in " + spgwConf + " failed")
			}

		}
		// oai.spgwu-start
		// time.Sleep(10 * time.Second)
		OaiObj.Logger.Print("start spgwu as daemon")
		fmt.Println("start spgwu as daemon")

		retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
		time.Sleep(5 * time.Second)
		counter := 0
		maxCounter := 2
		for {
			time.Sleep(1 * time.Second)
			if len(retStatus.Stderr) == 0 {
				counter = counter + 1
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "status"}, "."))
				oaiSpgwuStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oaiSpgwuStatus, "inactive")
				if checkInactive != true {
					OaiObj.Logger.Print("Waiting to make sure that oai-spgwu is working properly")
					fmt.Println("Waiting to make sure that oai-spgwu is working properly")
					if counter >= maxCounter {
						break
					}
				} else {
					OaiObj.Logger.Print("oai-spgwu is in inactive status, restarting the service")
					fmt.Println("oai-spgwu is in inactive status, restarting the service")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "stop"}, "."))
					for {
						time.Sleep(1 * time.Second)
						retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "status"}, "."))
						oaiSpgwuStatus = strings.Join(retStatus.Stdout, " ")
						if strings.Contains(oaiSpgwuStatus, "disabled") && strings.Contains(oaiSpgwuStatus, "inactive") {
							break
						}
					}
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
					time.Sleep(5 * time.Second)
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start oai-spgwu failed, try again later")
				fmt.Println("Start oai-spgwu failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
				time.Sleep(5 * time.Second)
				counter = 0
			}
		}

	}
	fmt.Println("END of oai-spgwu configuring and starting")
	OaiObj.Logger.Print("END of oai-spgwu configuring and starting")
	return nil
}

// RestartSpgwuV2 : Restart SPGW as a daemon
func restartSpgwuV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Restart oai-spgw daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.spgw-restart")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Restart oai-spgw failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-spgw is successfully restarted")
	return nil
}

// stopSpgwuV2 : Stop SPGW as a daemon
func stopSpgwuV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Stop oai-spgw daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.spgw-stop")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Stop oai-spgw failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-spgw is successfully stopped")
	return nil
}
