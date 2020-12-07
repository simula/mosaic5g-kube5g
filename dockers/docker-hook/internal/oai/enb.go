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
	"regexp"
	"strconv"
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
			oaiRanDisabledInactive := 0
			for {
				time.Sleep(1 * time.Second)
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
				oairanStatus := strings.Join(retStatus.Stdout, " ")
				if strings.Contains(oairanStatus, "disabled") && strings.Contains(oairanStatus, "inactive") {
					oaiRanDisabledInactive = 1
					break
				}
			}
			if oaiRanDisabledInactive == 1 {
				break
			}
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

	// tx_gain
	sedCommand = "s:tx_gain.*;:tx_gain                                         = " + c.OaiEnb[0].TxGain.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// rx_gain
	sedCommand = "s:rx_gain.*;:rx_gain                                         = " + c.OaiEnb[0].RxGain.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// pusch_p0_Nominal
	sedCommand = "s:pusch_p0_Nominal.*;:pusch_p0_Nominal                                         = " + c.OaiEnb[0].PuschP0Nominal.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// pucch_p0_Nominal
	sedCommand = "s:pucch_p0_Nominal.*;:pucch_p0_Nominal                                         = " + c.OaiEnb[0].PucchP0Nominal.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// pdsch_referenceSignalPower
	sedCommand = "s:pdsch_referenceSignalPower.*;:pdsch_referenceSignalPower                                         = " + c.OaiEnb[0].PdschReferenceSignalPower.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// puSch10xSnr
	sedCommand = "s:puSch10xSnr.*;:puSch10xSnr                                         = " + c.OaiEnb[0].PuSch10xSnr.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// puCch10xSnr
	sedCommand = "s:puCch10xSnr.*;:puCch10xSnr                                         = " + c.OaiEnb[0].PuCch10xSnr.Default + ";:g"
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
		flexranIP, err = util.GetIPFromDomain(OaiObj.Logger, c.OaiEnb[0].FlexRANServiceName)
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
	mmeServiceName := OaiObj.Conf.OaiEnb[0].MmeService.Name
	mmeIPV4Customized := OaiObj.Conf.OaiEnb[0].MmeService.IPV4
	mmeSnapVersion := OaiObj.Conf.OaiEnb[0].MmeService.SnapVersion
	if buildSnap == false {
		mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, mmeServiceName)

		if (mmeServiceName == "") && (mmeIPV4Customized == "") {
			//skip configuring mme as there is nothing specified
			if mmeSnapVersion == "v1" {
				mmeIP = "127.0.1.10"
			} else {
				mmeIP = "127.0.1.1"
			}
		} else if mmeIPV4Customized != "" {
			mmeIP = mmeIPV4Customized
		} else if mmeServiceName != "" {
			mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, mmeServiceName)
			OaiObj.Logger.Print("err = ", err)
			OaiObj.Logger.Print("mmeIP = ", mmeIP)
			for {
				if err != nil {
					OaiObj.Logger.Print(err)
				} else {
					hostNameMme, err := net.LookupHost(mmeIP)
					OaiObj.Logger.Print("err = ", err)
					OaiObj.Logger.Print("hostNameMme = ", hostNameMme)
					if len(hostNameMme) > 0 {
						// time.Sleep(3 * time.Second)
						break
					} else {
						OaiObj.Logger.Print(err)
					}
				}
				OaiObj.Logger.Print("Valid ip address for oai-mme not get retreived")
				time.Sleep(1 * time.Second)
				mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, mmeServiceName)
			}
		} else {
			// both mmeServiceName and mmeIPV4Customized are defined, thus, the privilige is for mmeIPV4Customized
			mmeIP = mmeIPV4Customized
			OaiObj.Logger.Print("mmeIP = mmeIPV4Customized")
		}
		OaiObj.Logger.Print("mmeServiceName = ", mmeServiceName)
		OaiObj.Logger.Print("mmeIPV4Customized = ", mmeIPV4Customized)
		OaiObj.Logger.Print("mmeIP = ", mmeIP)

		fmt.Println("mmeServiceName = ", mmeServiceName)
		fmt.Println("mmeIPV4Customized = ", mmeIPV4Customized)
		fmt.Println("mmeIP = ", mmeIP)

		sedCommand = "175s:\".*;:\"" + mmeIP + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

		var counter int64
		var maxWaitTime int64
		if mmeIPV4Customized == "" {
			if mmeSnapVersion == "v1" {
				// maxWaitTime = 3
				maxWaitTime = 20
				OaiObj.Logger.Print("Waiting for " + strconv.FormatInt(maxWaitTime, 10) + " seconds until the service " + mmeServiceName + " becomes ready ...")
				fmt.Println("Waiting for " + strconv.FormatInt(maxWaitTime, 10) + " seconds until the service " + mmeServiceName + " becomes ready ...")
				time.Sleep(time.Duration(maxWaitTime) * time.Second)
				OaiObj.Logger.Print("Supposed that service " + mmeServiceName + " is active now")
				fmt.Println("Supposed that service " + mmeServiceName + " is active now")
			} else {
				// curl http://127.0.0.1:5552/mme/status
				urlMme := "http://" + mmeIP + ":5552/mme/status"

				counter = 0
				maxWaitTime = 90

				var mmeActiveTime int64
				var counterMmeActiveTime int64

				mmeActiveTime = 0
				counterMmeActiveTime = 1

				resp, err := http.Get(urlMme)
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
						if len(mmeStat) > 0 {
							if (mmeStat[0].Startup == "enabled") && (mmeStat[0].Current == "active") {
								mmeActiveTime++
								if mmeActiveTime >= counterMmeActiveTime {
									OaiObj.Logger.Print("The service " + mmeStat[0].Service + " is active")
									fmt.Println("The service " + mmeStat[0].Service + " is active")
									break
								} else {
									OaiObj.Logger.Print("Waiting time " + strconv.FormatInt(mmeActiveTime, 10) + "/" + strconv.FormatInt(counterMmeActiveTime, 10) + " seconds to make sure that the service " + mmeServiceName + " is active")
									fmt.Println("Waiting time " + strconv.FormatInt(mmeActiveTime, 10) + "/" + strconv.FormatInt(counterMmeActiveTime, 10) + " seconds to make sure that the service " + mmeServiceName + " is active")
								}
							}
						} else {
							mmeActiveTime = 0
							OaiObj.Logger.Print("The service " + mmeServiceName + " is NOT active yet, waiting ...")
							fmt.Println("The service " + mmeServiceName + " is NOT active yet, waiting ...")
						}
					}
					counter++
					if counter >= maxWaitTime {
						OaiObj.Logger.Print("Waiting for " + strconv.FormatInt(maxWaitTime, 10) + " seconds while the service " + mmeServiceName + " is not ready yet, exit...")
						fmt.Println("Waiting for " + strconv.FormatInt(maxWaitTime, 10) + " seconds while the service " + mmeServiceName + " is not ready yet, exit...")
						break
					}
					time.Sleep(1 * time.Second)
					resp, err = http.Get(urlMme)
				}
			}
		}
		OaiObj.Logger.Print("Start enb daemon")

		if mmeIPV4Customized == "" && mmeSnapVersion == "v2" {
			// curl http://127.0.0.1:5552/mme/journal
			urlMmeJournal := "http://" + mmeIP + ":5552/mme/journal"
			OaiObj.Logger.Print("urlMmeJournal=", urlMmeJournal)
			retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
			time.Sleep(5 * time.Second)
			var counterOairabActiveTime int64

			counter = 0
			counterOairabActiveTime = 30
			for {
				time.Sleep(1 * time.Second)
				if len(retStatus.Stderr) == 0 {
					counter++
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
					oairanStatus := strings.Join(retStatus.Stdout, " ")
					checkInactive := strings.Contains(oairanStatus, "inactive")
					if checkInactive != true {
						// check the journal of mme
						resp, err := http.Get(urlMmeJournal)
						// time.Sleep(3 * time.Second)
						for {
							time.Sleep(2 * time.Second)
							resp, err = http.Get(urlMmeJournal)
							if err != nil {
								OaiObj.Logger.Print(err)
							} else {
								defer resp.Body.Close()
								bodyBytes, _ := ioutil.ReadAll(resp.Body)
								mmeJournal := string(bodyBytes)
								connectedEnbStr := `Connected\s*eNBs\s*\|\s*[1-9][0-9]*\s*\|\s*\d*\s*\|\s*\d*\s*\|`
								OaiObj.Logger.Print("connectedEnbStr=", connectedEnbStr)
								reConnectedEnb := regexp.MustCompile(connectedEnbStr)
								submatchall := reConnectedEnb.FindAllString(mmeJournal, -1)
								if len(submatchall) > 1 {
									// enb is connected
									counter = counterOairabActiveTime
									OaiObj.Logger.Print("Found eNB connected to mme, exit...")
									break
								} else {
									// enb is NOT connected yet, restart the eNB
									OaiObj.Logger.Print("enb is in not connected to the mme/cn, restarting the service")
									util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-stop"}, "/"))
									// time.Sleep(5 * time.Second)
									for {
										time.Sleep(1 * time.Second)
										retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
										oairanStatus := strings.Join(retStatus.Stdout, " ")
										if strings.Contains(oairanStatus, "disabled") && strings.Contains(oairanStatus, "inactive") {
											break
										}
									}
									// time.Sleep(3 * time.Second)
									retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
									// time.Sleep(5 * time.Second)
									time.Sleep(5 * time.Second)
									counter = 0
									break
								}

							}

						}
						if counter >= counterOairabActiveTime {
							OaiObj.Logger.Print("enb is working, exit...")
							break
						}
					} else {
						OaiObj.Logger.Print("enb is in inactive status, restarting the service")
						util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-stop"}, "/"))
						// time.Sleep(5 * time.Second)
						for {
							time.Sleep(1 * time.Second)
							retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
							oairanStatus := strings.Join(retStatus.Stdout, " ")
							if strings.Contains(oairanStatus, "disabled") && strings.Contains(oairanStatus, "inactive") {
								break
							}
						}
						// time.Sleep(3 * time.Second)
						retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
						// time.Sleep(5 * time.Second)
						time.Sleep(5 * time.Second)
						counter = 0
					}
				} else {
					OaiObj.Logger.Print("Start enb failed, try again later")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
					time.Sleep(5 * time.Second)
					// time.Sleep(5 * time.Second)
					counter = 0
				}
			}
		} else {
			retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
			time.Sleep(5 * time.Second)
			var counterOairabActiveTime int64

			counter = 0
			counterOairabActiveTime = 30
			for {
				time.Sleep(1 * time.Second)
				if len(retStatus.Stderr) == 0 {
					counter++
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
					oairanStatus := strings.Join(retStatus.Stdout, " ")
					checkInactive := strings.Contains(oairanStatus, "inactive")
					if checkInactive != true {
						if counter >= counterOairabActiveTime {
							OaiObj.Logger.Print("enb is working, exit...")
							break
						}
					} else {
						OaiObj.Logger.Print("enb is in inactive status, restarting the service")
						util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-stop"}, "/"))
						for {
							time.Sleep(1 * time.Second)
							retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
							oairanStatus := strings.Join(retStatus.Stdout, " ")
							if strings.Contains(oairanStatus, "disabled") && strings.Contains(oairanStatus, "inactive") {
								break
							}
						}
						// time.Sleep(5 * time.Second)
						retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
						time.Sleep(5 * time.Second)
						counter = 0
					}
				} else {
					OaiObj.Logger.Print("Start enb failed, try again later")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
					time.Sleep(5 * time.Second)
					counter = 0
				}
			}
		}
	}
	OaiObj.Logger.Print("enb daemon Started")
	return nil
}
