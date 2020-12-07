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
# file          mme.go
# brief 		configure the snap of oai-mme v2, and start it
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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StartMmeV2 : Start MME as a daemon
func startMmeV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring OAI-MME V2")
	OaiObj.Logger.Print("Starting configuration of OAI-MME V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-mme.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	mmeBin := strings.Join([]string{snapBinaryPath, "oai-mme"}, "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "conf-get"}, "."))
	fmt.Println(retStatus)
	OaiObj.Logger.Print(retStatus)
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]
	mmeConf := strings.Join([]string{confPath, "mme.conf"}, "/")
	mmeFdConf := strings.Join([]string{confPath, "mme_fd.conf"}, "/")

	// get the dns
	var mcc string
	var mnc string
	if CnAllInOneMode == true {
		mcc = OaiObj.Conf.OaiCn.V2[0].OaiMme.MCC
		mnc = OaiObj.Conf.OaiCn.V2[0].OaiMme.MNC
	} else {
		mcc = OaiObj.Conf.OaiMme.V2[0].MCC
		mnc = OaiObj.Conf.OaiMme.V2[0].MNC
	}

	if buildSnap == false {
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "init"}, "."))
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Init of oai-mme is failed")
				fmt.Println("Init of oai-mme is failed")
			} else {
				OaiObj.Logger.Print("Init of oai-mme is successful")
				fmt.Println("Init of oai-mme is successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Init of oai-mme")
			fmt.Println("Retrying to Init of oai-mme")
			retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "init"}, "."))
		}

		// Get interface ip and replace the default one
		outInterfaceIP := util.GetOutboundIP()
		outInterface, _ := util.GetInterfaceByIP(outInterfaceIP)

		// MME_INTERFACE_NAME_FOR_S1_MME
		sedCommand := "s:MME_INTERFACE_NAME_FOR_S1_MME.*;:MME_INTERFACE_NAME_FOR_S1_MME               = \"" + outInterface + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " failed")
				fmt.Println("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " failed")
			} else {
				OaiObj.Logger.Print("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " successful")
				fmt.Println("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf)
			fmt.Println("Retrying to Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf)
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		}

		// MME_IPV4_ADDRESS_FOR_S1_MME
		sedCommand = "s:MME_IPV4_ADDRESS_FOR_S1_MME.*;:MME_IPV4_ADDRESS_FOR_S1_MME          = \"" + outInterfaceIP + "/24\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)

		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf + " failed")
				fmt.Println("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf + " failed")
			} else {
				OaiObj.Logger.Print("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + " in " + mmeConf + " successful")
				fmt.Println("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf + " successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf)
			fmt.Println("Retrying to Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf)
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		}

		/*========================================= Open mmeConf as text-file =========================================*/
		mmeConfInput, err := ioutil.ReadFile(mmeConf)
		if err != nil {
			OaiObj.Logger.Print(err)
			fmt.Println(err)
			os.Exit(1)
		}
		mmeConfInputText := string(mmeConfInput)

		// GUMMEI_LIST: MCC and MCC
		gummeiListMccMnc := regexp.MustCompile(`{MCC="\d{3}"\s.*;\s.*MNC="\d{2}";\s.*MME_GID="\d{1}"\s.*;\s.*MME_CODE="\d{1}";\s.*}`)
		submatchall := gummeiListMccMnc.FindAllString(mmeConfInputText, -1)
		for _, element := range submatchall {
			gummeiListMccMnc = regexp.MustCompile(element)
			break
		}
		fmt.Println(gummeiListMccMnc)
		OaiObj.Logger.Print(gummeiListMccMnc)

		newStrExpr := `{MCC="` + mcc + `" ; MNC="` + mnc + `"; MME_GID="4" ; MME_CODE="1"; }`
		repStr := `${1}` + newStrExpr + `$2`
		mmeConfInputText = gummeiListMccMnc.ReplaceAllString(mmeConfInputText, repStr)

		// TAI_LIST: MCC and MCC
		taiListMccMnc := regexp.MustCompile(`{MCC="\d{3}" ; MNC="\d{2}";  TAC = "\d{1}"; }`)
		submatchall = taiListMccMnc.FindAllString(mmeConfInputText, 3)
		for _, element := range submatchall {
			taiListMccMnc = regexp.MustCompile(element)
		}
		fmt.Println(taiListMccMnc)
		OaiObj.Logger.Print(taiListMccMnc)

		newStrExpr = `{MCC="` + mcc + `" ; MNC="` + mnc + `";  TAC = "1"; }`
		repStr = `${1}` + newStrExpr + `$2`
		mmeConfInputText = taiListMccMnc.ReplaceAllString(mmeConfInputText, repStr)

		// WRR_LIST_SELECTION: MCC and MCC
		if len(mnc) == 2 {
			mnc = "0" + mnc
		}
		regExpStr := `{ID="` + `tac-lb01.tac-hb00.tac.epc.mnc` + `\d+` + `.mcc` + `\d+` + `.3gppnetwork.org"`
		wrrListSelection := regexp.MustCompile(regExpStr)

		submatchall = wrrListSelection.FindAllString(mmeConfInputText, 3)
		for _, element := range submatchall {
			wrrListSelection = regexp.MustCompile(element)
		}
		fmt.Println(wrrListSelection)
		OaiObj.Logger.Print(wrrListSelection)

		newStrExpr = `{ID="` + `tac-lb01.tac-hb00.tac.epc.mnc` + mnc + `.mcc` + mcc + `.3gppnetwork.org"`
		repStr = `${1}` + newStrExpr + `$2`
		mmeConfInputText = wrrListSelection.ReplaceAllString(mmeConfInputText, repStr)

		mmeConfOutput := []byte(mmeConfInputText)
		if err = ioutil.WriteFile(mmeConf, mmeConfOutput, 0666); err != nil {
			OaiObj.Logger.Print(err)
			fmt.Println(err)
			os.Exit(1)
		}
		//
		hssIP := "127.0.0.1"
		hssServiceName := "oai-hss"
		if CnAllInOneMode == false {
			hssServiceName = OaiObj.Conf.OaiMme.V2[0].HssServiceName

			outInterfaceIP := util.GetOutboundIP()
			outInterface, _ := util.GetInterfaceByIP(outInterfaceIP)

			sedCommand := "s:MME_INTERFACE_NAME_FOR_S11.*;:MME_INTERFACE_NAME_FOR_S11 = \"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
			if retStatus.Exit != 0 {
				return errors.New("Set MME_INTERFACE_NAME_FOR_S11 in " + mmeConf + " failed")
			}

			sedCommand = "s:MME_IPV4_ADDRESS_FOR_S11.*;:MME_IPV4_ADDRESS_FOR_S11          = \"" + outInterfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
			if retStatus.Exit != 0 {
				return errors.New("Set MME_IPV4_ADDRESS_FOR_S11 in " + mmeConf + " failed")
			}

			sedCommand = "s:ListenOn.*;:ListenOn = \"" + outInterfaceIP + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
			if retStatus.Exit != 0 {
				return errors.New("Set ListenOn in " + mmeFdConf + " failed")
			}
			hssIP, err := util.GetIPFromDomain(OaiObj.Logger, hssServiceName)
			for {
				if err != nil {
					OaiObj.Logger.Print(err)
				} else {
					hostNameHss, err := net.LookupHost(hssIP)
					if len(hostNameHss) > 0 {
						break
					} else {
						OaiObj.Logger.Print(err)
					}
				}
				OaiObj.Logger.Print("Valid ip address for oai-hss not yet retreived")
				time.Sleep(1 * time.Second)
				hssIP, err = util.GetIPFromDomain(OaiObj.Logger, hssServiceName)
			}

			spgwcIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiMme.V2[0].SpgwcServiceName)
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
				spgwcIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiMme.V2[0].SpgwcServiceName)
			}

			// replace the ip address of hss
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.10/"+hssIP+"/g", mmeFdConf)
			if retStatus.Exit != 0 {
				return errors.New("Set the ip address of oai-hss in " + mmeFdConf + " failed")
			}

			// SGW_IPV4_ADDRESS_FOR_S11: this value was in the old config of oai-mme
			// replace SGW_IPV4_ADDRESS_FOR_S11
			sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11=\"127.0.11.2.*;:SGW_IPV4_ADDRESS_FOR_S11=\"" + spgwcIP + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print(errors.New("Set the ip address of oai-spgwc SGW_IPV4_ADDRESS_FOR_S11 in " + mmeConf + " failed"))
			}

			// SGW_IP_ADDRESS_FOR_S11; this is the new value in the config of oai-mme
			// replace SGW_IP_ADDRESS_FOR_S11
			sedCommand = "s:SGW_IP_ADDRESS_FOR_S11=\"127.0.11.2.*;:SGW_IP_ADDRESS_FOR_S11=\"" + spgwcIP + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
			if retStatus.Exit != 0 {
				return errors.New("Set the ip address of oai-spgwc SGW_IP_ADDRESS_FOR_S11 in " + mmeConf + " failed")
			}

		}

		if CnAllInOneMode == true {
			// deuplicated from oai-hss
			var APINi string
			APINi = OaiObj.Conf.OaiCn.V2[0].ApnNi.Default
			cassandraIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiCn.V2[0].OaiHss.DatabaseServiceName)
			OaiObj.Logger.Print("cassandraIP " + cassandraIP)
			for {
				if err != nil {
					OaiObj.Logger.Print(err)
				} else {
					hostNameCassandra, err := net.LookupHost(cassandraIP)

					if len(hostNameCassandra) > 0 {
						break
					} else {
						OaiObj.Logger.Print(err)
					}
				}
				OaiObj.Logger.Print("Valid ip address for mysql not yet retreived")
				time.Sleep(1 * time.Second)
				cassandraIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.OaiCn.V2[0].OaiHss.DatabaseServiceName)
			}
			// deuplicated from oai-hss
			// sudo oai-hss.add-mme -i ubuntu.openair5G.eur -C 172.18.0.2
			OaiObj.Logger.Print("Adding oai-mme to Cassanra DB ")
			retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
			for {
				if retStatus.Exit != 0 {
					OaiObj.Logger.Print("Adding the mme to hss database failed")
					fmt.Println("Adding the mme to hss database failed")
				} else {
					OaiObj.Logger.Print("Adding the mme to hss database was successful")
					fmt.Println("Adding the mme to hss database was successful")
					break
				}
				time.Sleep(1 * time.Second)
				OaiObj.Logger.Print("Retrying to add oai-mme to hss database")
				fmt.Println("Retrying to add oai-mme to hss database")
				retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
			}
			// oai-hss.add-users -I208950000000001-208950000000010 -a oai.ipv4 -C 172.18.0.2
			OaiObj.Logger.Print("Adding users to Cassanra DB ")
			retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-users", "-I", "208950000000001-208950000000010", "-a", APINi, "-C", cassandraIP)
			for {
				if retStatus.Exit != 0 {
					OaiObj.Logger.Print("Adding users to hss database failed")
					fmt.Println("Adding users to hss database failed")
					// return errors.New("Adding users to hss database failed")
				} else {
					OaiObj.Logger.Print("Adding users to hss database was successful")
					fmt.Println("Adding users to hss database was successful")
					break
				}
				time.Sleep(1 * time.Second)
				OaiObj.Logger.Print("Retrying to add users to hss database")
				fmt.Println("Retrying to add users to hss database")
				retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-users", "-I", "208950000000001-208950000000010", "-a", APINi, "-C", cassandraIP)
			}
		}

		// curl http://127.0.0.1:5551/hss/status
		urlHssStatus := "http://" + hssIP + ":5552/hss/status"
		var counter int64
		counter = 0

		var hssActiveTime int64
		var counterHssActiveTime int64
		var maxWaitTimeHssStatus int64

		hssActiveTime = 0
		counterHssActiveTime = 3
		maxWaitTimeHssStatus = 15
		resp, err := http.Get(urlHssStatus)
		for {
			if err != nil {
				OaiObj.Logger.Print(err)
			} else {
				defer resp.Body.Close()
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				bodyString := string(bodyBytes)

				var hssStat []common.CnEntityV2Status
				json.Unmarshal([]byte(bodyString), &hssStat)

				OaiObj.Logger.Print(hssStat)
				fmt.Println(hssStat)
				/*
					hssStat=
					[
						{
							"service": "oai-mme.mmed",
							"startup": "enabled",
							"current": "active",
							"notes": "-"
						}
					]
				*/
				if len(hssStat) > 0 {
					if (hssStat[0].Startup == "enabled") && (hssStat[0].Current == "active") {
						hssActiveTime++
						if hssActiveTime >= counterHssActiveTime {
							OaiObj.Logger.Print("The service " + hssStat[0].Service + " is active")
							fmt.Println("The service " + hssStat[0].Service + " is active")
							break
						} else {
							OaiObj.Logger.Print("Waiting time " + strconv.FormatInt(hssActiveTime, 10) + "/" + strconv.FormatInt(counterHssActiveTime, 10) + " seconds to make sure that the service " + hssServiceName + " is active")
							fmt.Println("Waiting time " + strconv.FormatInt(hssActiveTime, 10) + "/" + strconv.FormatInt(counterHssActiveTime, 10) + " seconds to make sure that the service " + hssServiceName + " is active")
						}
					}
				} else {
					hssActiveTime = 0
					OaiObj.Logger.Print("The service " + hssServiceName + " is NOT active yet, waiting ...")
					fmt.Println("The service " + hssServiceName + " is NOT active yet, waiting ...")
				}
			}
			counter++
			if counter >= maxWaitTimeHssStatus {
				OaiObj.Logger.Print("Waiting for " + strconv.FormatInt(maxWaitTimeHssStatus, 10) + " seconds while the service " + hssServiceName + " is not ready yet, exit...")
				fmt.Println("Waiting for " + strconv.FormatInt(maxWaitTimeHssStatus, 10) + " seconds while the service " + hssServiceName + " is not ready yet, exit...")
				break
			}
			time.Sleep(1 * time.Second)
			resp, err = http.Get(urlHssStatus)
		}
		// oai-mme.start
		// time.Sleep(10 * time.Second)
		retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "start"}, "."))
		time.Sleep(5 * time.Second)
		var maxCounterHssStatus int64

		// curl http://127.0.0.1:5551/hss/journal
		urlHssJournal := "http://" + hssIP + ":5551/hss/journal"
		OaiObj.Logger.Print("urlHssJournal=", urlHssJournal)

		counter = 0
		maxCounterHssStatus = 3
		for {
			time.Sleep(1 * time.Second)
			if len(retStatus.Stderr) == 0 {
				counter = counter + 1
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "status"}, "."))
				oaiMmeStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oaiMmeStatus, "inactive")
				if checkInactive != true {
					OaiObj.Logger.Print("Waiting to make sure that oai-mme is working properly")
					fmt.Println("Waiting to make sure that oai-mme is working properly")
					resp, err := http.Get(urlHssJournal)
					time.Sleep(1 * time.Second)
					for {
						if err != nil {
							OaiObj.Logger.Print(err)
						} else {
							defer resp.Body.Close()
							bodyBytes, _ := ioutil.ReadAll(resp.Body)
							hssJournal := string(bodyBytes)
							ClosedToOpenStateStr := `NOTI\s*'STATE_CLOSED'\s*->\s*'STATE_OPEN'\s*`
							ClosedToOpenStateStr2 := `\s*->\s*'STATE_OPEN'\s*`
							OaiObj.Logger.Print("ClosedToOpenStateStr=", ClosedToOpenStateStr)
							reClosedToOpenState := regexp.MustCompile(ClosedToOpenStateStr)
							reClosedToOpenState2 := regexp.MustCompile(ClosedToOpenStateStr2)

							submatchall := reClosedToOpenState.FindAllString(hssJournal, -1)
							submatchall2 := reClosedToOpenState2.FindAllString(hssJournal, -1)
							if (len(submatchall) >= 1) || (len(submatchall2) >= 1) {
								// oai-hss is in STATE_OPEN
								counter = maxCounterHssStatus
								OaiObj.Logger.Print("oai-hss switched from STATE_CLOSED to STATE_OPEN , exit...")
								break
							} else {
								// oai-hss is not in STATE_OPEN
								OaiObj.Logger.Print("oai-hss is not in STATE_OPEN, restarting the service")
								retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "stop"}, "."))
								for {
									retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "status"}, "."))
									oaiMmeStatus = strings.Join(retStatus.Stdout, " ")
									if strings.Contains(oaiMmeStatus, "disabled") && strings.Contains(oaiMmeStatus, "inactive") {
										break
									}
								}
								retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "start"}, "."))
								time.Sleep(5 * time.Second)
								counter = 0
								break
							}

						}
						time.Sleep(1 * time.Second)
						resp, err = http.Get(urlHssJournal)
					}
					if counter >= maxCounterHssStatus {
						break
					}
				} else {
					OaiObj.Logger.Print("oai-mme is in inactive status, restarting the service")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "stop"}, "."))
					for {
						time.Sleep(1 * time.Second)
						retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "status"}, "."))
						oaiMmeStatus = strings.Join(retStatus.Stdout, " ")
						if strings.Contains(oaiMmeStatus, "disabled") && strings.Contains(oaiMmeStatus, "inactive") {
							break
						}
					}
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "start"}, "."))
					time.Sleep(5 * time.Second)
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start oai-mme failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "start"}, "."))
				time.Sleep(2 * time.Second)
				counter = 0
			}
		}
	}
	fmt.Println("END of oai-mme configuring and starting")
	OaiObj.Logger.Print("END of oai-mme configuring and starting")
	return nil
}

// RestartMmeV2 : Restart MME as a daemon
func restartMmeV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Restart oai-mme daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.mme-restart")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Restart oai-mme failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-mme is successfully restarted")
	return nil
}

// stopMmeV2 : Stop MME as a daemon
func stopMmeV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Stop oai-mme daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.mme-stop")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Stop oai-mme failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-mme is successfully stopped")
	return nil
}
