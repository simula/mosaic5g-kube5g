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
# file          spgwcV2.go
# brief 		configure the snap of oai-spgwc v2, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package oai

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mosaic5g/docker-hook/internal/pkg/util"
	"os"
	"regexp"
	"strings"
	"time"
)

// initSpgwcV2 : Init SPGW
func initSpgwcV2(OaiObj Oai) error {
	return nil
}

// configSpgwcV2 : Config oai-spgw
func configSpgwcV2(OaiObj Oai) error {
	return nil
}

// StartSpgwcV2 : Start SPGW as a daemon
func startSpgwcV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring OAI-SPGWC V2")
	OaiObj.Logger.Print("Starting configuration of OAI-SPGWC V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-spgwc.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	spgwBin := strings.Join([]string{snapBinaryPath, "oai-spgwc"}, "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "conf-get"}, "."))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]

	spgwConf := strings.Join([]string{confPath, "spgwc.conf"}, "/")

	if buildSnap == false {
		// Init spgwc
		OaiObj.Logger.Print("Start Init of oai-spgwc")
		fmt.Println("Start Init of oai-spgwc")
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Init of oai-spgwc is failed")
				fmt.Println("Init of oai-spgwc is failed")
			} else {
				OaiObj.Logger.Print("Init of oai-spgwc is successful")
				fmt.Println("Init of oai-spgwc is successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Init of oai-spgwc")
			fmt.Println("Retrying to Init of oai-spgwc")
			retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		}

		// Configure oai-spgw
		OaiObj.Logger.Print("Configure of oai-spgwc")
		fmt.Println("Configure of oai-spgwc")

		// get the dns
		var DNSIPV4Address string
		var APINi string
		if CnAllInOneMode == true {
			DNSIPV4Address = OaiObj.Conf.OaiCn.V2[0].OaiSpgwc.DNS
			APINi = OaiObj.Conf.OaiCn.V2[0].ApnNi.Default
		} else {
			DNSIPV4Address = OaiObj.Conf.OaiSpgwc.V2[0].DNS
			APINi = OaiObj.Conf.OaiSpgwc.V2[0].ApnNi.Default
		}

		sedCommand := "s:DEFAULT_DNS_IPV4_ADDRESS.*;:DEFAULT_DNS_IPV4_ADDRESS     = \"" + DNSIPV4Address + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)

		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + DNSIPV4Address + " in " + spgwConf + " failed")
				fmt.Println("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + DNSIPV4Address + " in " + spgwConf + " failed")
			} else {
				OaiObj.Logger.Print("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + DNSIPV4Address + " in " + spgwConf + " successful")
				fmt.Println("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + DNSIPV4Address + " in " + spgwConf + " successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Set DEFAULT_DNS_IPV4_ADDRESS to the value " + DNSIPV4Address + " in " + spgwConf)
			fmt.Println("Retrying to Set DEFAULT_DNS_IPV4_ADDRESS to the value " + DNSIPV4Address + " in " + spgwConf)
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		}

		// retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/oai.openair5G.eur/"+APINi+"/g", spgwConf)
		// retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/oai.openair5G.eur/"+"oai.ipv4"+"/g", spgwConf)
		/*========================================= Open spgwConf as text-file =========================================*/
		spgwcConfInput, err := ioutil.ReadFile(spgwConf)
		if err != nil {
			OaiObj.Logger.Print(err)
			fmt.Println(err)
			os.Exit(1)
		}
		spgwcConfInputText := string(spgwcConfInput)

		// APN_LIST: APN_NI
		reqExpStr := `{APN_NI\s*=\s*".*";\s*PDN_TYPE\s*=\s*"IPv4";\s*IPV4_POOL\s*=\s*\d*;\s*IPV6_POOL\s*=\s*-*\d*}`
		apnListApnNi := regexp.MustCompile(reqExpStr)
		submatchall := apnListApnNi.FindAllString(spgwcConfInputText, -1)

		for _, element := range submatchall {
			apnListApnNi = regexp.MustCompile(element)
			break
		}
		fmt.Println(apnListApnNi)
		OaiObj.Logger.Print(apnListApnNi)

		newStr := `{APN_NI = "` + APINi + `"; PDN_TYPE = "IPv4"; IPV4_POOL = 0; IPV6_POOL = -1}`
		repStr := `${1}` + newStr + `$2`
		spgwcConfInputText = apnListApnNi.ReplaceAllString(spgwcConfInputText, repStr)

		fmt.Println(spgwcConfInputText)
		OaiObj.Logger.Print(spgwcConfInputText)

		spgwcConfOutput := []byte(spgwcConfInputText)

		if err = ioutil.WriteFile(spgwConf, spgwcConfOutput, 0666); err != nil {
			OaiObj.Logger.Print(err)
			fmt.Println(err)
			os.Exit(1)
		}
		//
		if CnAllInOneMode == false {
			// Get interface IP and outbound interface
			interfaceIP := util.GetOutboundIP()
			outInterface, _ := util.GetInterfaceByIP(interfaceIP)
			// INTERFACE_NAME of S11_CP
			sedCommand := "71s:\".*;:\"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set INTERFACE_NAME of S11_CP in " + spgwConf + " failed")
			}
			// IPV4_ADDRESS of S11_CP
			sedCommand = "72s:\".*;:\"" + interfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set IPV4_ADDRESS of S11_CP in " + spgwConf + " failed")
			}
			// INTERFACE_NAME of SX
			sedCommand = "161s:\".*;:\"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set INTERFACE_NAME of SX in " + spgwConf + " failed")
			}
			// IPV4_ADDRESS of SX
			sedCommand = "162s:\".*;:\"" + interfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set IPV4_ADDRESS of SX in " + spgwConf + " failed")
			}
		}

		// oai.spgwc-start
		// time.Sleep(10 * time.Second)
		OaiObj.Logger.Print("start spgwc as daemon")
		fmt.Println("start spgwc as daemon")
		retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
		time.Sleep(5 * time.Second)
		counter := 0
		maxCounter := 2
		for {
			time.Sleep(1 * time.Second)
			if len(retStatus.Stderr) == 0 {
				counter = counter + 1
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "status"}, "."))
				oaiSpgwcStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oaiSpgwcStatus, "inactive")
				if checkInactive != true {
					OaiObj.Logger.Print("Waiting to make sure that oai-spgwc is working properly")
					fmt.Println("Waiting to make sure that oai-spgwc is working properly")
					if counter >= maxCounter {
						break
					}
				} else {
					OaiObj.Logger.Print("oai-spgwc is in inactive status, restarting the service")
					fmt.Println("oai-spgwc is in inactive status, restarting the service")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "stop"}, "."))
					for {
						time.Sleep(1 * time.Second)
						retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "status"}, "."))
						oaiSpgwcStatus = strings.Join(retStatus.Stdout, " ")
						if strings.Contains(oaiSpgwcStatus, "disabled") && strings.Contains(oaiSpgwcStatus, "inactive") {
							break
						}
					}
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
					time.Sleep(5 * time.Second)
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start oai-spgwc failed, try again later")
				fmt.Println("Start oai-spgwc failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
				time.Sleep(5 * time.Second)
				counter = 0
			}
		}
	}

	fmt.Println("END of oai-spgwc configuring and starting")
	OaiObj.Logger.Print("END of oai-spgwc configuring and starting")
	return nil
}

// RestartSpgwcV2 : Restart SPGW as a daemon
func restartSpgwcV2(OaiObj Oai) error {
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

// stopSpgwcV2 : Stop SPGW as a daemon
func stopSpgwcV2(OaiObj Oai) error {
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
