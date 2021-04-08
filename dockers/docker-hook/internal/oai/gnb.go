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
# file          gnb.go
# brief 		configure the snap of oai-ran, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
	- Alireza Mohammadi
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

// replaceExistingPuschProcThreads changes the parameter pusch_proc_threads to the value defined inside yaml only if this parameter
// exists in the currrent conf file. Confirmed this function tested with the examples.
// TODO: what if this parameter is not in conf file but in the yaml file, how do we proceed?
func replaceExistingPuschProcThreads(c *common.CfgGlobal, OaiObj Oai, gnbConf string) int {
	// check if text exists in the file
	// naming to self-explain variable
	notFoundParameter := util.RunCmd(OaiObj.Logger, "grep", "-iq", "pusch_proc_threads", gnbConf)
	fmt.Printf("Here you go %v", notFoundParameter.Exit)
	if notFoundParameter.Exit != 0 {
		// we only need to handle the case when we don't find the parameter in config file
		// insert the parameter to the "right" place.
		// normally,it is better to have a better interface that allows to insert the parameter
		// under the struct L1s of conf file. But since we don't have it at this moment,
		// we will do this by assume that the order of parameters inside the L1s struct doesn't matter.
		sedCommand := "N;/L1s.*{/a \\ \\ \\ \\ \\ \\ \\ \\ pusch_proc_threads     = " + c.OaiGnb[0].PuschProcThreads + ";"
		retStatusNotFound := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
		if retStatusNotFound.Exit != 0 {
			fmt.Println(retStatusNotFound.Complete)
			fmt.Println(retStatusNotFound.Stderr)
		}
		return retStatusNotFound.Exit
	} else {
		sedCommand := "s:pusch_proc_threads.*;:pusch_proc_threads     = " + c.OaiGnb[0].PuschProcThreads + ";:g"
		retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
		return retStatus.Exit
	}
}

func startGNB(OaiObj Oai, buildSnap bool) error {
	var msg string = ""
	// get the configuration
	c := OaiObj.Conf
	// config filename of the snap
	// confFileName := "enb.band7.tm1.50PRB.usrpb210.conf"
	nodeFunction := OaiObj.Conf.OaiGnb[0].NodeFunction
	if nodeFunction == "" {
		nodeFunction = "gnb"
	}
	cmdNodeFunction := "oai-ran." + nodeFunction
	msg = "getting the config file of " + nodeFunction
	OaiObj.Logger.Print(msg)
	fmt.Println(msg)
	retStatus := util.RunCmd(OaiObj.Logger, cmdNodeFunction+"-conf-get")
	confFileName := ""
	if retStatus.Exit == 0 {
		s := strings.Split(retStatus.Stdout[0], "/")
		confFileName = s[len(s)-1]
		OaiObj.Logger.Print("the config file of " + nodeFunction + " is " + confFileName)
	} else {
		var outError string
		for i := 0; i < len(retStatus.Stderr); i++ {
			outError += retStatus.Stderr[i] + "\n"
		}
		return errors.New("Error while getting the config file of " + nodeFunction + "\n" + outError)
	}

	ranConfLinesStr := util.RunCmd(OaiObj.Logger, cmdNodeFunction+"-conf-show")
	var ranConfStr string
	for i := 0; i < len(ranConfLinesStr.Stdout); i++ {
		ranConfStr += ranConfLinesStr.Stdout[i] + "\n"
	}

	retStatus = util.RunCmd(OaiObj.Logger, "which", cmdNodeFunction+"-status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")

	// Stop oai-gnb
	OaiObj.Logger.Print("Stop gnb daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-stop"}, "/"))

		if len(retStatus.Stderr) == 0 {
			oaiRanDisabledInactive := 0
			for {
				time.Sleep(1 * time.Second)
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-status"}, "/"))
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
		OaiObj.Logger.Print("Stop oai-gnb failed, try again later")
		time.Sleep(1 * time.Second)
	}

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-conf-get"}, "/"))

	OaiObj.Logger.Print("confFileName=" + confFileName)
	s = strings.Split(retStatus.Stdout[0], "/")
	gnbConf := strings.Join(s[0:len(s)-1], "/")
	gnbConf = strings.Join([]string{gnbConf, confFileName}, "/")
	OaiObj.Logger.Print("confFileName=" + confFileName)
	OaiObj.Logger.Print("gnbConf=", gnbConf)
	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-conf-set"}, "/"), gnbConf)

	// Replace MCC
	sedCommand := "s/mcc =.[^;]*/mcc = " + c.OaiGnb[0].MCC + "/g"
	OaiObj.Logger.Print("Replace MCC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MCC in " + gnbConf + " failed")
	}

	//Replace MNC
	sedCommand = "s/mnc =.[^;]*/mnc = " + c.OaiGnb[0].MNC + "/g"
	OaiObj.Logger.Print("Replace MNC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MNC in " + gnbConf + " failed")
	}

	// Get Outbound IP and Interface name
	outIP := util.GetOutboundIP()
	outInterface, err := util.GetInterfaceByIP(outIP)
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	OaiObj.Logger.Print("Outbound Interface and IP is ", outInterface, " ", outIP)
	// Replace interface
	sedCommand = "s:GNB_INTERFACE_NAME_FOR_S1_MME.*;:GNB_INTERFACE_NAME_FOR_S1_MME            = \"" + outInterface + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
	sedCommand = "s:GNB_INTERFACE_NAME_FOR_S1U.*;:GNB_INTERFACE_NAME_FOR_S1U               = \"" + outInterface + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

	// GNB_IPV4_ADDRESS_FOR_S1_MME
	sedCommand = "s:GNB_IPV4_ADDRESS_FOR_S1_MME.*;:GNB_IPV4_ADDRESS_FOR_S1_MME              = \"" + outIP + "/23\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

	sedCommand = "s:GNB_IPV4_ADDRESS_FOR_S1U.*;:GNB_IPV4_ADDRESS_FOR_S1U                 = \"" + outIP + "/23\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

	// GNB_IPV4_ADDRESS_FOR_X2C
	sedCommand = "s:GNB_IPV4_ADDRESS_FOR_X2C.*;:GNB_IPV4_ADDRESS_FOR_X2C                 = \"" + outIP + "/24\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

	// parallel_config
	sedCommand = "s:parallel_config.*;:parallel_config    = \"" + c.OaiGnb[0].ParallelConfig.Default + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

	// max_rxgain
	sedCommand = "s:max_rxgain.*;:max_rxgain     = " + c.OaiGnb[0].MaxRxGain.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

	// pusch_proc_threads
	replaceExistingPuschProcThreads(c, OaiObj, gnbConf)

	// Get the IP address of oai-mme
	mmeServiceName := OaiObj.Conf.OaiGnb[0].MmeService.Name
	mmeIPV4Customized := OaiObj.Conf.OaiGnb[0].MmeService.IPV4
	mmeSnapVersion := OaiObj.Conf.OaiGnb[0].MmeService.SnapVersion
	if buildSnap == false {

		if c.OaiGnb[0].Usrp.N3xx.Enabled {
			firstAdd := c.OaiGnb[0].Usrp.N3xx.SdrAddrs.Addr
			secondAdd := c.OaiGnb[0].Usrp.N3xx.SdrAddrs.SecondAddr
			clockSrcVal := c.OaiGnb[0].Usrp.N3xx.ClockSrc
			clockSrcEnbStr := `clock_src\s*=\s*"[a-z]*[A-Z]*[0-9]*";`
			cdrAddrsEnbStr := `sdr_addrs\s*=\s*"`
			OaiObj.Logger.Print("clockSrcEnbStr=", clockSrcEnbStr)
			OaiObj.Logger.Print("cdrAddrsEnbStr=", cdrAddrsEnbStr)
			reclockSrcEnb := regexp.MustCompile(clockSrcEnbStr)
			submatchallClockSrc := reclockSrcEnb.FindAllString(ranConfStr, -1)

			cdrAddrsSrcEnb := regexp.MustCompile(cdrAddrsEnbStr)
			submatchallCdrAddrs := cdrAddrsSrcEnb.FindAllString(ranConfStr, -1)

			if (len(submatchallClockSrc) == 0) && (len(submatchallCdrAddrs) == 0) {
				// sdr_addrs and clock_src do not exist in the file, add it
				OaiObj.Logger.Print("Adding sdr_addrs and clock_src to the config file of " + nodeFunction)
				lineNumberToAddSdrAddrsClockSrc := 0
				eNBInstancesStr := `eNB_instances\s*=\s*`
				regeNBInstances := regexp.MustCompile(eNBInstancesStr)
				for i := 0; i < len(ranConfLinesStr.Stdout); i++ {
					submatchall := regeNBInstances.FindAllString(ranConfLinesStr.Stdout[i], -1)
					if len(submatchall) >= 1 {
						lineNumberToAddSdrAddrsClockSrc = i + 2
						break
					}
				}
				// sdr_addrs      = "addr=192.168.20.2,second_addr=192.168.10.2";
				StrToAdd := strconv.Itoa(lineNumberToAddSdrAddrsClockSrc) + `i         sdr_addrs      = "addr=` + firstAdd + `,second_addr=` + secondAdd + `";`
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
				if retStatus.Exit != 0 {
					return errors.New(`Insert the following is failed:` + StrToAdd)
				}
				//     clock_src      = "internal";
				StrToAdd = strconv.Itoa(lineNumberToAddSdrAddrsClockSrc+1) + `i clock_src      = "` + clockSrcVal + `";`
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
				if retStatus.Exit != 0 {
					return errors.New(`Insert clock_src      = "internal"; failed`)
				}

			} else if len(submatchallClockSrc) == 0 {
				// sdr_addrs and clock_src do not exist in the file, add it
				OaiObj.Logger.Print("Adding sdr_addrs and clock_src to the config file of " + nodeFunction)
				lineNumberToAddSdrAddrsClockSrc := 0
				for i := 0; i < len(ranConfLinesStr.Stdout); i++ {
					submatchall := cdrAddrsSrcEnb.FindAllString(ranConfLinesStr.Stdout[i], -1)
					if len(submatchall) >= 1 {
						lineNumberToAddSdrAddrsClockSrc = i + 1
						break
					}
				}
				// sdr_addrs      = "addr=192.168.20.2,second_addr=192.168.10.2";
				StrToAdd := strconv.Itoa(lineNumberToAddSdrAddrsClockSrc) + `s/.*/         sdr_addrs      = "addr=` + firstAdd + `,second_addr=` + secondAdd + `";/`
				// sed -i '60s/.*/new-content-line/' file.conf
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
				if retStatus.Exit != 0 {
					return errors.New(`Insert the following is failed:` + StrToAdd)
				}
				//     clock_src      = "internal";

				StrToAdd = strconv.Itoa(lineNumberToAddSdrAddrsClockSrc+1) + `i clock_src      = "` + clockSrcVal + `";`
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
				if retStatus.Exit != 0 {
					return errors.New(`Insert clock_src      = "internal"; failed`)
				}

			} else if len(submatchallCdrAddrs) == 0 {
				// sdr_addrs and clock_src do not exist in the file, add it
				OaiObj.Logger.Print("Adding sdr_addrs and clock_src to the config file of " + nodeFunction)
				lineNumberToAddSdrAddrsClockSrc := 0
				for i := 0; i < len(ranConfLinesStr.Stdout); i++ {
					submatchall := reclockSrcEnb.FindAllString(ranConfLinesStr.Stdout[i], -1)
					if len(submatchall) >= 1 {
						lineNumberToAddSdrAddrsClockSrc = i + 1
						break
					}
				}
				// sdr_addrs      = "addr=192.168.20.2,second_addr=192.168.10.2";
				StrToAdd := strconv.Itoa(lineNumberToAddSdrAddrsClockSrc) + `i sdr_addrs      = "addr=` + firstAdd + `,second_addr=` + secondAdd + `";`
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
				if retStatus.Exit != 0 {
					return errors.New(`Insert the following is failed:` + StrToAdd)
				}

				// sed -i '60s/.*/new-content-line/' file.conf
				StrToAdd = strconv.Itoa(lineNumberToAddSdrAddrsClockSrc+1) + `s/.*/         clock_src      = "` + clockSrcVal + `";/`
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
				if retStatus.Exit != 0 {
					return errors.New(`Insert clock_src      = "internal"; failed`)
				}
				////////////////////////////////////////////////////////////////
			} else {
				OaiObj.Logger.Print("Configure sdr_addrs and clock_src in the config file of " + nodeFunction)
				cdrAddrsFound := false
				clockSrcFound := false
				for i := 0; i < len(ranConfLinesStr.Stdout); i++ {
					submatchall := cdrAddrsSrcEnb.FindAllString(ranConfLinesStr.Stdout[i], -1)
					if len(submatchall) >= 1 {
						// sdr_addrs      = "addr=192.168.20.2,second_addr=192.168.10.2";
						StrToAdd := strconv.Itoa(i+1) + `s/.*/         sdr_addrs      = "addr=` + firstAdd + `,second_addr=` + secondAdd + `";/`
						// sed -i '60s/.*/new-content-line/' file.conf
						retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
						if retStatus.Exit != 0 {
							return errors.New(`Insert the following is failed:` + StrToAdd)
						}
						cdrAddrsFound = true
					}

					submatchall = reclockSrcEnb.FindAllString(ranConfLinesStr.Stdout[i], -1)
					if len(submatchall) >= 1 {

						// sed -i '60s/.*/new-content-line/' file.conf
						StrToAdd := strconv.Itoa(i+1) + `s/.*/         clock_src      = "` + clockSrcVal + `";/`
						retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", StrToAdd, gnbConf)
						if retStatus.Exit != 0 {
							return errors.New(`Insert clock_src      = "internal"; failed`)
						}
						clockSrcFound = true
					}
					if cdrAddrsFound && clockSrcFound {
						break
					}
				}
			}
		}

		// X2 Interface:
		// enable_x2
		sedCommand = "s:enable_x2.*;:enable_x2    = \"" + c.OaiGnb[0].X2Config.EnableX2 + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
		if c.OaiGnb[0].X2Config.EnableX2 == "yes" {
			// t_reloc_prep
			sedCommand = "s:t_reloc_prep.*;:t_reloc_prep    = " + c.OaiGnb[0].X2Config.TRelocPrep + ";:g"
			util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
			// tx2_reloc_overall
			sedCommand = "s:tx2_reloc_overall.*;:tx2_reloc_overall    = " + c.OaiGnb[0].X2Config.TX2RelocOverall + ";:g"
			util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
			// t_dc_prep
			sedCommand = "s:t_dc_prep.*;:t_dc_prep    = " + c.OaiGnb[0].X2Config.TDCPrep + ";:g"
			util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

			// t_dc_overall
			sedCommand = "s:t_dc_overall.*;:t_dc_overall    = " + c.OaiGnb[0].X2Config.TDCOverall + ";:g"
			util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
			var enbIP string = ""
			var enbServiceName string = c.OaiGnb[0].X2Config.TargetENBX2IP.Name
			if enbServiceName != "" {
				// getting the ip address of the eNB
				enbIP, err = util.GetIPFromDomain(OaiObj.Logger, enbServiceName)
				OaiObj.Logger.Print("err = ", err)
				OaiObj.Logger.Print("enbIP = ", enbIP)
				for {
					if err != nil {
						OaiObj.Logger.Print(err)
					} else {
						hostNameMme, err := net.LookupHost(enbIP)
						OaiObj.Logger.Print("err = ", err)
						OaiObj.Logger.Print("hostNameMme = ", hostNameMme)
						if len(hostNameMme) > 0 {
							// time.Sleep(3 * time.Second)
							break
						} else {
							OaiObj.Logger.Print(err)
						}
					}
					OaiObj.Logger.Print("Valid ip address for eNB not get retreived")
					time.Sleep(1 * time.Second)
					enbIP, err = util.GetIPFromDomain(OaiObj.Logger, enbServiceName)
				}
			} else {
				enbIP = c.OaiGnb[0].X2Config.TargetENBX2IP.IPV4
			}
			// getting the ip address: target_enb_x2_ip_address

			// Replace enb IP
			sedCommand = "220s:\".*;:\"" + enbIP + "\";:g"
			util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)
		}

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

		sedCommand = "206s:\".*;:\"" + mmeIP + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, gnbConf)

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

		// TODO: gNB cannot be run as daemon! because it needs -E option which is not supported by the daemon
		OaiObj.Logger.Print("Start gnb not as a daemon!")
		// TODO temporary solution until supporting passing parameters when starting gnb
		mmeSnapVersion = "v1"
		if mmeSnapVersion == "v2" {
			// curl http://127.0.0.1:5552/mme/journal
			urlMmeJournal := "http://" + mmeIP + ":5552/mme/journal"
			OaiObj.Logger.Print("urlMmeJournal=", urlMmeJournal)

			retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-start"}, "/"))
			time.Sleep(5 * time.Second)

			var counterOairabActiveTime int64
			var counterGlobalMaxTime int64
			var counterGlobal int64
			counter = 0
			counterOairabActiveTime = 5
			counterGlobalMaxTime = 15
			counterGlobal = 0

			for {
				time.Sleep(1 * time.Second)
				if len(retStatus.Stderr) == 0 {
					counter++
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-status"}, "/"))
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
									// gnb is connected
									// counter = counterOairabActiveTime
									OaiObj.Logger.Print("Found gNB connected to mme, exit...")
									break
								} else {
									// gnb is NOT connected yet, restart the gNB
									OaiObj.Logger.Print("gnb is in not connected to the mme/cn, restarting the service")
									// TODO this is temporary solution, to be changed later
									retStatus := util.RunCmd(OaiObj.Logger, "pkill oai-nr")
									// util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-stop"}, "/"))
									// time.Sleep(5 * time.Second)
									for {
										time.Sleep(1 * time.Second)
										retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-status"}, "/"))
										oairanStatus := strings.Join(retStatus.Stdout, " ")
										if strings.Contains(oairanStatus, "disabled") && strings.Contains(oairanStatus, "inactive") {
											break
										}
									}
									// time.Sleep(3 * time.Second)
									retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-start"}, "/"))
									// time.Sleep(5 * time.Second)
									time.Sleep(5 * time.Second)
									counter = 0
									break
								}

							}

						}
						if counter >= counterOairabActiveTime {
							OaiObj.Logger.Print("gnb is working, exit...")
							break
						}
					} else {
						OaiObj.Logger.Print("gnb is in inactive status, restarting the service")
						util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-stop"}, "/"))
						// time.Sleep(5 * time.Second)
						for {
							time.Sleep(1 * time.Second)
							retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-status"}, "/"))
							oairanStatus := strings.Join(retStatus.Stdout, " ")
							if strings.Contains(oairanStatus, "disabled") && strings.Contains(oairanStatus, "inactive") {
								break
							}
						}
						// time.Sleep(3 * time.Second)
						retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-start"}, "/"))
						// time.Sleep(5 * time.Second)
						time.Sleep(5 * time.Second)
						counter = 0
					}
				} else {
					OaiObj.Logger.Print("Start gnb failed, try again later")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-start"}, "/"))
					time.Sleep(5 * time.Second)
					// time.Sleep(5 * time.Second)
					counter = 0
				}
				counterGlobal++
				OaiObj.Logger.Print("counterGlobal: ", counterGlobal)
				if counterGlobal >= counterGlobalMaxTime {
					OaiObj.Logger.Print("Maximum time ", counterGlobalMaxTime, " reached, exit...")
					break
				}
			}
		} else {
			retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-start"}, "/"))
			time.Sleep(5 * time.Second)
			var counterOairabActiveTime int64

			counter = 0
			counterOairabActiveTime = 30
			for {
				time.Sleep(1 * time.Second)
				if len(retStatus.Stderr) == 0 {
					counter++
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-status"}, "/"))
					oairanStatus := strings.Join(retStatus.Stdout, " ")
					checkInactive := strings.Contains(oairanStatus, "inactive")
					if checkInactive != true {
						if counter >= counterOairabActiveTime {
							OaiObj.Logger.Print("gnb is working, exit...")
							break
						}
					} else {
						OaiObj.Logger.Print("gnb is in inactive status, restarting the service")
						util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-stop"}, "/"))
						for {
							time.Sleep(1 * time.Second)
							retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-status"}, "/"))
							oairanStatus := strings.Join(retStatus.Stdout, " ")
							if strings.Contains(oairanStatus, "disabled") && strings.Contains(oairanStatus, "inactive") {
								break
							}
						}
						// time.Sleep(5 * time.Second)
						retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-start"}, "/"))
						time.Sleep(5 * time.Second)
						counter = 0
					}
				} else {
					OaiObj.Logger.Print("Start gnb failed, try again later")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, cmdNodeFunction + "-start"}, "/"))
					time.Sleep(5 * time.Second)
					counter = 0
				}
			}
		}
	}
	OaiObj.Logger.Print("gnb started as not a daemon!")
	return nil
}
