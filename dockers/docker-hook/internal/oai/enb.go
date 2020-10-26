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
# file          enb.go
# brief 		configure the snap of oai-ran, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/
package oai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mosaic5g/docker-hook/internal/pkg/common"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"net/http"
	"strings"
	"time"
)

func startENB(OaiObj Oai, buildSnap bool) error {
	// get the configuration
	c := OaiObj.Conf
	// config filename of the snap
	confFileName := "enb.band7.tm1.50PRB.usrpb210.conf"

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-ran.enb-status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")

	// Stop oai-enb
	OaiObj.Logger.Print("Stop enb daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-stop"}, "/"))
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Stop oai-enb failed, try again later")
		time.Sleep(1 * time.Second)
	}

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-conf-get"}, "/"))

	s = strings.Split(retStatus.Stdout[0], "/")
	enbConf := strings.Join(s[0:len(s)-1], "/")
	enbConf = strings.Join([]string{enbConf, confFileName}, "/")
	OaiObj.Logger.Print("enbConf=", enbConf)
	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-conf-set"}, "/"), enbConf)

	// Replace MCC
	sedCommand := "s/mcc =.[^;]*/mcc = " + c.OaiEnb[0].MCC + "/g"
	OaiObj.Logger.Print("Replace MCC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MCC in " + enbConf + " failed")
	}

	//Replace MNC
	sedCommand = "s/mnc =.[^;]*/mnc = " + c.OaiEnb[0].MNC + "/g"
	OaiObj.Logger.Print("Replace MNC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MNC in " + enbConf + " failed")
	}

	//eutra_band
	sedCommand = "s:eutra_band.*;:eutra_band                                      = " + c.OaiEnb[0].EutraBand.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// downlink_frequency
	sedCommand = "s:downlink_frequency.*;:downlink_frequency                              = " + c.OaiEnb[0].DownlinkFrequency.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// uplink_frequency_offset
	sedCommand = "s:uplink_frequency_offset.*;:uplink_frequency_offset                         = " + c.OaiEnb[0].UplinkFrequencyOffset.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// N_RB_DL
	sedCommand = "s:N_RB_DL.*;:N_RB_DL                                         = " + c.OaiEnb[0].NumberRbDl.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Get Outbound IP and Interface name
	outIP := util.GetOutboundIP()
	outInterface, err := util.GetInterfaceByIP(outIP)
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	OaiObj.Logger.Print("Outbound Interface and IP is ", outInterface, " ", outIP)
	// Replace interface
	sedCommand = "s:ENB_INTERFACE_NAME_FOR_S1_MME.*;:ENB_INTERFACE_NAME_FOR_S1_MME            = \"" + outInterface + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	sedCommand = "s:ENB_INTERFACE_NAME_FOR_S1U.*;:ENB_INTERFACE_NAME_FOR_S1U               = \"" + outInterface + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Replace enb IP
	sedCommand = "s:ENB_IPV4_ADDRESS_FOR_S1_MME.*;:ENB_IPV4_ADDRESS_FOR_S1_MME              = \"" + outIP + "/23\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	sedCommand = "s:ENB_IPV4_ADDRESS_FOR_S1U.*;:ENB_IPV4_ADDRESS_FOR_S1U                 = \"" + outIP + "/23\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	sedCommand = "s:ENB_IPV4_ADDRESS_FOR_X2C.*;:ENB_IPV4_ADDRESS_FOR_X2C                 = \"" + outIP + "/24\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Set up FlexRAN
	if (OaiObj.Conf.OaiEnb[0].FlexRAN == true) && (buildSnap == false) {
		// Get flexRAN ip
		var flexranIP string
		OaiObj.Logger.Print("Configure FlexRAN Parameters")
		flexranIP, err = util.GetIPFromDomain(OaiObj.Logger, c.OaiEnb[0].FlexranServiceName)
		if err != nil {
			OaiObj.Logger.Print(err)
			OaiObj.Logger.Print("Getting IP of FlexRAN failed, try again later")
		}
		sedCommand = "s:FLEXRAN_ENABLED.*;:FLEXRAN_ENABLED=        \"yes\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
		sedCommand = "s:FLEXRAN_INTERFACE_NAME.*;:FLEXRAN_INTERFACE_NAME= \"eth0\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
		sedCommand = "s:FLEXRAN_IPV4_ADDRESS.*;:FLEXRAN_IPV4_ADDRESS   = \"" + flexranIP + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	} else {
		OaiObj.Logger.Print("Disable FlexRAN Feature")
		sedCommand = "s:FLEXRAN_ENABLED.*;:FLEXRAN_ENABLED=        \"no\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	}

	// parallel_config
	sedCommand = "s:parallel_config.*;:parallel_config    = \"" + c.OaiEnb[0].ParallelConfig.Default + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// max_rxgain
	sedCommand = "s:max_rxgain.*;:max_rxgain     = " + c.OaiEnb[0].MaxRxGain.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Get the IP address of oai-mme
	if buildSnap == false {
		mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiEnb[0].MmeService.Name)
		for {
			if err != nil {
				OaiObj.Logger.Print(err)
			} else {
				hostNameMme, err := net.LookupHost(mmeIP)
				if len(hostNameMme) > 0 {
					// time.Sleep(3 * time.Second)
					break
				} else {
					OaiObj.Logger.Print(err)
				}
			}
			OaiObj.Logger.Print("Valid ip address for oai-hss not get retreived")
			time.Sleep(1 * time.Second)
			mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiEnb[0].MmeService.Name)
		}
		sedCommand = "175s:\".*;:\"" + mmeIP + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

		// curl 172.20.0.5:5552/mme/status
		urlMme := "http://" + mmeIP + ":5552/mme/status"
		resp, err := http.Get(urlMme)
		counter := 0
		mmeActiveTime := 0

		maxWaitTime := 200
		counterMmeActiveTime := 30
		for {
			if err != nil {
				OaiObj.Logger.Print(err)
			} else {
				defer resp.Body.Close()
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				bodyString := string(bodyBytes)

				var mmeStat []common.MmeStatus
				json.Unmarshal([]byte(bodyString), &mmeStat)

				OaiObj.Logger.Print(mmeStat)
				fmt.Println(mmeStat)
				/*
					mmeStat=
					[
						{
							"service": "oai-mme.mmed",
							"startup": "enabled",
							"current": "active",
							"notes": "-"
						}
					]
				*/
				if (mmeStat[0].Startup == "enabled") && (mmeStat[0].Current == "active") {
					mmeActiveTime++
					if mmeActiveTime >= counterMmeActiveTime {
						OaiObj.Logger.Print("The service " + mmeStat[0].Service + " is active")
						fmt.Println("The service " + mmeStat[0].Service + " is active")
						break
					} else {
						OaiObj.Logger.Print("Waiting for " + string(counterMmeActiveTime) + " seconds to make sure that the service " + OaiObj.Conf.OaiEnb[0].MmeService.Name + " is active")
						fmt.Println("Waiting for " + string(counterMmeActiveTime) + " seconds to make sure that the service " + OaiObj.Conf.OaiEnb[0].MmeService.Name + " is active")
					}
				} else {
					mmeActiveTime = 0
					OaiObj.Logger.Print("The service " + OaiObj.Conf.OaiEnb[0].MmeService.Name + " is NOT active yet, waiting ...")
					fmt.Println("The service " + OaiObj.Conf.OaiEnb[0].MmeService.Name + " is NOT active yet, waiting ...")
				}
			}
			counter++
			if counter >= maxWaitTime {
				OaiObj.Logger.Print("Waiting for " + string(counter) + " seconds while the service " + OaiObj.Conf.OaiEnb[0].MmeService.Name + " is not ready yet, exit...")
				fmt.Println("Waiting for " + string(counter) + " seconds while the service " + OaiObj.Conf.OaiEnb[0].MmeService.Name + " is not ready yet, exit...")
			}
			time.Sleep(1 * time.Second)
			resp, err = http.Get(urlMme)
		}

		OaiObj.Logger.Print("Start enb daemon")

		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
		counter = 0
		for {
			if len(retStatus.Stderr) == 0 {
				time.Sleep(5 * time.Second)
				counter++
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
				oairanStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oairanStatus, "inactive")
				if checkInactive != true {
					if counter >= 30 {
						break
					}
				} else {
					OaiObj.Logger.Print("enb is in inactive status, restarting the service")
					util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-stop"}, "/"))
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start enb failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
				counter = 0
			}
		}
	}
	OaiObj.Logger.Print("enb daemon Started")
	return nil
}
