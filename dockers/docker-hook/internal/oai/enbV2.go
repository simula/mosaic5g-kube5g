package oai

import (
	"errors"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"strconv"
	"strings"
	"time"
)

func startENBV2(OaiObj Oai, buildSnap bool) error {
	// get the configuration
	c := OaiObj.Conf
	cnf := OaiObj.ConfOaiRan.OaiRanConf

	// config filename of the snap
	confFileName := "enb.band7.tm1.50PRB.usrpb210.conf"
	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-ran.enb-status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")

	// Stop oai-enb
	OaiObj.Logger.Print("Stop enb daemon")
	for {
		// "/snap/bin/oai-ran.enb-stop"
		// retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-stop")
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-stop"}, "/"))
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Stop oai-enb failed, try again later")
		time.Sleep(1 * time.Second)
	}
	///////////////////

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-conf-get"}, "/"))
	// retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-conf-get")
	s = strings.Split(retStatus.Stdout[0], "/")
	enbConf := strings.Join(s[0:len(s)-1], "/")
	enbConf = strings.Join([]string{enbConf, confFileName}, "/")
	OaiObj.Logger.Print("enbConf=", enbConf)
	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-conf-set"}, "/"), enbConf)
	// retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-conf-set", enbConf)

	// //Active_eNBs
	// Active_eNBs := "eNB-Eurecom-LTEBox"
	// sedCommand := "s:Active_eNBs.*;:Active_eNBs = ( \"" + Active_eNBs + "\");:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// //eNB_ID
	// sedCommand = "s:eNB_ID.*;:eNB_ID    =  " + c.EnbId.Default + ";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// //eNB_name
	// sedCommand := "s:eNB_name.*;:eNB_name  =  \"" + c.EnbName.Default + "\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// node function
	sedCommand := "s:node_function.*;:node_function             = \"" + cnf.ComponentCarriers.NodeFunction + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Replace MCC
	sedCommand = "s/mcc =.[^;]*/mcc = " + string(cnf.Mcc[0]) + "/g"
	OaiObj.Logger.Print("Replace MCC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MCC in " + enbConf + " failed")
	}

	//Replace MNC
	sedCommand = "s/mnc =.[^;]*/mnc = " + string(cnf.Mnc[0]) + "/g"
	OaiObj.Logger.Print("Replace MNC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MNC in " + enbConf + " failed")
	}

	//eutra_band
	sedCommand = "s:eutra_band.*;:eutra_band                                      = " + cnf.ComponentCarriers.EutraBand + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// downlink_frequency
	sedCommand = "s:downlink_frequency.*;:downlink_frequency                              = " + cnf.ComponentCarriers.DownlinkFrequency + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// uplink_frequency_offset
	sedCommand = "s:uplink_frequency_offset.*;:uplink_frequency_offset                         = " + cnf.ComponentCarriers.UplinkFrequencyOffset + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Nid_cell
	sedCommand = "s:Nid_cell.*;:Nid_cell                                         = " + string(cnf.ComponentCarriers.NidCell) + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// N_RB_DL
	sedCommand = "s:N_RB_DL.*;:N_RB_DL                                         = " + string(cnf.ComponentCarriers.NRbDl) + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Get Outbound IP and Interface name
	outIP := util.GetOutboundIP(OaiObj.Logger)
	outInterface, err := util.GetInterfaceByIP(OaiObj.Logger, outIP)
	if err != nil {
		util.PrintFuncFatal(OaiObj.Logger, err)
	}
	util.PrintFunc(OaiObj.Logger, "Outbound Interface and IP is "+outInterface+" "+outIP)
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
	if (cnf.NetworkController.FlexranEnabled == "yes") && (buildSnap == false) {
		var flexranIP, flexranIface string
		if cnf.NetworkController.FlexRANDomainName == "" {
			flexranIP = cnf.NetworkController.FlexRANIPv4Address
			flexranIface = cnf.NetworkController.FlexRANInterfaceName
		} else {
			// Get flexRAN ip
			flexranIface = "eth0"
			OaiObj.Logger.Print("Configure FlexRAN Parameters")
			flexranIP, err = util.GetIPFromDomain(OaiObj.Logger, cnf.NetworkController.FlexRANDomainName)
			if err != nil {
				OaiObj.Logger.Print(err)
				OaiObj.Logger.Print("Getting IP of FlexRAN failed, try again later")
			}
		}
		sedCommand = "s:FLEXRAN_ENABLED.*;:FLEXRAN_ENABLED=        \"yes\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
		sedCommand = "s:FLEXRAN_INTERFACE_NAME.*;:FLEXRAN_INTERFACE_NAME= \"" + flexranIface + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
		sedCommand = "s:FLEXRAN_IPV4_ADDRESS.*;:FLEXRAN_IPV4_ADDRESS   = \"" + flexranIP + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	} else {
		OaiObj.Logger.Print("Disable FlexRAN Feature")
		sedCommand = "s:FLEXRAN_ENABLED.*;:FLEXRAN_ENABLED=        \"no\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	}

	// parallel_config
	sedCommand = "s:parallel_config.*;:parallel_config    = \"" + cnf.ThreadStruct.ParallelConfig + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// max_rxgain
	sedCommand = "s:max_rxgain.*;:max_rxgain     = " + string(cnf.Rus.MaxRxGain) + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// max_pdschReferenceSignalPower
	sedCommand = "s:max_pdschReferenceSignalPower.*;:max_pdschReferenceSignalPower     = " + string(cnf.Rus.MaxPdschReferenceSignalPower) + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// X2-HO
	if cnf.X2Ho.EnableX2 == "yes" {
		// enable_x2
		sedCommand = "s:enable_x2.*;:enable_x2     = " + "yes" + ";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
		if cnf.X2Ho.MasterNode {
			// Master eNB

		} else {
			// Slave eNB
			for i := 0; i < len(cnf.X2Ho.TargetEnbX2IPAddress); i++ {
				var remoteEnbIp string
				if (cnf.X2Ho.TargetEnbX2IPAddress[i]).RanDomainName != "" {
					// Master eNB defined by domain; getting the ip address
					remoteEnbIp, err = util.GetIPFromDomain(OaiObj.Logger, (cnf.X2Ho.TargetEnbX2IPAddress[i]).RanDomainName)
					for {
						if err != nil {
							util.PrintFunc(OaiObj.Logger, err)
						} else {
							hostNameRemoteEnb, err := net.LookupHost(remoteEnbIp)
							if len(hostNameRemoteEnb) > 0 {
								// time.Sleep(3 * time.Second)
								break
							} else {
								util.PrintFunc(OaiObj.Logger, err)
							}
						}
						util.PrintFunc(OaiObj.Logger, "Valid ip address for master eNB not yet retreived")
						time.Sleep(1 * time.Second)
						remoteEnbIp, err = util.GetIPFromDomain(OaiObj.Logger, (cnf.X2Ho.TargetEnbX2IPAddress[i]).RanDomainName)
					}
				} else {
					// The ip address of master eNB is defined
					remoteEnbIp = (cnf.X2Ho.TargetEnbX2IPAddress[i]).Ipv4
				}
				// sed -n "/NETWORK_INTERFACES/="  mante.eucnc.orig.95.conf
				sedCommand = "/NETWORK_INTERFACES/="
				retStatus := util.RunCmd(OaiObj.Logger, "sed", "-n", sedCommand, enbConf)

				number, _ := strconv.ParseUint(retStatus.Stdout[0], 10, 32)
				lineNumber := strconv.Itoa(int(number - 1))
				sedCommand1 := lineNumber + " a target_enb_x2_ip_address      = ( {"
				lineNumber = strconv.Itoa(int(number))
				sedCommand2 := lineNumber + " a 	 ipv4       = \"" + remoteEnbIp + "\";"
				lineNumber = strconv.Itoa(int(number + 1))
				sedCommand3 := lineNumber + " a 	  ipv6       = \"192:168:30::17\";"
				lineNumber = strconv.Itoa(int(number + 2))
				sedCommand4 := lineNumber + " a 	  preference = \"ipv4\";"
				lineNumber = strconv.Itoa(int(number + 3))
				sedCommand5 := lineNumber + " a 	  }"
				lineNumber = strconv.Itoa(int(number + 4))
				sedCommand6 := lineNumber + " a 	  );"

				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand1, enbConf)
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand2, enbConf)
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand3, enbConf)
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand4, enbConf)
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand5, enbConf)
				retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand6, enbConf)
			}
		}
	}
	// Get the IP address of oai-mme
	if (buildSnap == false) && (cnf.MmeIPAddress.Ipv4 != "") {
		var mmeIP string = cnf.MmeIPAddress.Ipv4
		if mmeIP == "" {
			mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.MmeDomainName)
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
				OaiObj.Logger.Print("Valid ip address for oai-hss not yet retreived")
				time.Sleep(1 * time.Second)
				mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.MmeDomainName)
			}
		}
		sedCommand = "s:mme_ip_address *= *( *{ *ipv4 *= *\".*\" *;:mme_ip_address      = ( { ipv4       = \"" + mmeIP + "\"" + ";:g"
		// sedCommand = "175s:\".*;:\"" + mmeIP + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

		util.PrintFunc(OaiObj.Logger, "Start waiting for 170 seconds before running oai-enb")
		time.Sleep(170 * time.Second) // 170
		util.PrintFunc(OaiObj.Logger, "Finish waiting for 170 seconds before running oai-enb")

		util.PrintFunc(OaiObj.Logger, "Start enb daemon")
		OaiObj.Logger.Print("Start enb daemon")

		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
		counter := 0
		maxCounter := 5 //30
		for {
			if len(retStatus.Stderr) == 0 {
				time.Sleep(5 * time.Second)
				counter = counter + 1
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-status"}, "/"))
				oairanStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oairanStatus, "inactive")
				if checkInactive != true {
					if counter >= maxCounter {
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
	util.PrintFunc(OaiObj.Logger, "enb daemon Started")
	return nil
}
