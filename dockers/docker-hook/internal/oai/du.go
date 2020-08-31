package oai

import (
	"errors"
	"mosaic5g/docker-hook/internal/pkg/util"
	"strings"
	"time"
)

func startDu(OaiObj Oai) error {
	///////////////
	OaiObj.Logger.Print("Stop du daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-stop")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Start du failed, try again later")
		time.Sleep(1 * time.Second)
	}
	///////////////////
	c := OaiObj.Conf
	confFileName := "du.lte.band7.10MHz.if4p5.conf"
	retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-conf-get")
	s := strings.Split(retStatus.Stdout[0], "/")
	enbConf := strings.Join(s[0:len(s)-1], "/")
	enbConf = strings.Join([]string{enbConf, confFileName}, "/")
	OaiObj.Logger.Print("enbConf=", enbConf)
	retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-conf-set", enbConf)
	//enbConf := c.ConfigurationPathofRAN + "cu.lte.conf"

	// mmeDomain := c.MmeDomainName
	cuDomain := c.CuDomainName

	// Get Outbound IP OaiObj.Logger
	outIP := util.GetOutboundIP(OaiObj.Logger)
	outInterface, err := util.GetInterfaceByIP(OaiObj.Logger, outIP)
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	OaiObj.Logger.Print("Outbound Interfacea and IP is ", outInterface, " ", outIP)

	// Get du ip
	OaiObj.Logger.Print("cuDomain=", cuDomain)
	cuIP, err := util.GetIPFromDomain(OaiObj.Logger, cuDomain)
	if err != nil {
		OaiObj.Logger.Print(err)
		cuIP = "192.168.12.45"
	}

	// Replace MCC
	sedCommand := "s/mcc =.[^;]*/mcc = " + c.MCC + "/g"
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MCC in " + enbConf + " failed")
	}
	OaiObj.Logger.Print("Replace MNC")
	//Replace MNC
	sedCommand = "s/mnc =.[^;]*/mnc = " + c.MNC + "/g"
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MNC in " + enbConf + " failed")
	}
	// Get mme ip
	// mmeIP, err := util.GetIPFromDomain(OaiObj.Logger, mmeDomain)
	// if err != nil {
	// 	OaiObj.Logger.Print(err)
	// 	mmeIP = "10.10.10.10"
	// }
	//eNB_name
	sedCommand = "s:eNB_name.*;:eNB_name              			      = \"" + c.EnbName.Default + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	//Active_eNBs
	sedCommand = "s:Active_eNBs.*;:Active_eNBs              			      = ( \"" + c.EnbName.Default + "\");:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	//eNB_ID
	sedCommand = "s:eNB_ID.*;:eNB_ID              			      = " + c.EnbId.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// max_rxgain
	sedCommand = "s:max_rxgain.*;:max_rxgain              			      = " + c.MaxRxGain.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// node function
	sedCommand = "s:node_function.*;:node_function              			      = \"" + c.NodeFunction.Default + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// eutra_band
	sedCommand = "s:eutra_band.*;:      eutra_band              			      = " + c.EutraBand.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// downlink_frequency
	sedCommand = "s:downlink_frequency.*;:      downlink_frequency      			      = " + c.DownlinkFrequency.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	///uplink_frequency_offset
	sedCommand = "s:uplink_frequency_offset.*;:      uplink_frequency_offset 			      = " + c.UplinkFrequencyOffset.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	//NumberRbDl
	sedCommand = "s:N_RB_DL.*;:      N_RB_DL 			      = " + c.NumberRbDl.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	//config oai-cu
	// local_n_if_name
	sedCommand = "s:local_n_if_name.*;:      local_n_if_name 			      = \"" + outInterface + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// local_n_address
	sedCommand = "s:local_n_address.*;:      local_n_address 			      = \"" + outIP + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// remote_n_address
	sedCommand = "s:remote_n_address.*;:      remote_n_address 			      = \"" + cuIP + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// local_n_portc
	sedCommand = "s:local_n_portc.*;:      local_n_portc 			      = " + c.DuPortc.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// remote_n_portc
	sedCommand = "s:remote_n_portc.*;:      remote_n_portc 			      = " + c.CuPortc.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// local_n_portd
	sedCommand = "s:local_n_portd.*;:      local_n_portd 			      = " + c.DuPortd.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// remote_n_portd
	sedCommand = "s:remote_n_portd.*;:      remote_n_portd 			      = " + c.CuPortd.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// tr_n_preference
	// sedCommand = "s:tr_n_preference.*;:      tr_n_preference 			      = \"" + "local_L1" + "\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// sedCommand = "s:tr_n_preference.*= \"f1\";:tr_n_preference  = \"" + "local_L1" + "\";:g"
	// util.RunCmd(c.Logger, "sed", "-i", sedCommand, enbConf)

	//TO BE REVISED
	// sedCommand = "182s:\".*;:\"" + mmeIP + "\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	//NETWORK_INTERFACES
	// // ENB_INTERFACE_NAME_FOR_S1_MME
	// sedCommand = "s:ENB_INTERFACE_NAME_FOR_S1_MME.*;:      ENB_INTERFACE_NAME_FOR_S1_MME 			      = " + outInterface + ";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// //ENB_IPV4_ADDRESS_FOR_S1_MME
	// sedCommand = "s:ENB_IPV4_ADDRESS_FOR_S1_MME.*;:ENB_IPV4_ADDRESS_FOR_S1_MME 							 = \"" + outIP + "/23\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// // ENB_INTERFACE_NAME_FOR_S1U
	// sedCommand = "s:ENB_INTERFACE_NAME_FOR_S1U.*;:      ENB_INTERFACE_NAME_FOR_S1U 			      = " + outInterface + ";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// //ENB_IPV4_ADDRESS_FOR_S1U
	// sedCommand = "s:ENB_IPV4_ADDRESS_FOR_S1U.*;:ENB_IPV4_ADDRESS_FOR_S1U 							 = \"" + outIP + "/23\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// //ENB_IPV4_ADDRESS_FOR_X2C
	// sedCommand = "s:ENB_IPV4_ADDRESS_FOR_X2C.*;:ENB_IPV4_ADDRESS_FOR_X2C 							 = \"" + outIP + "/24\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// sedCommand = "s/eno1/" + outInterface + "/g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// // Replace enb IP
	// sedCommand = "192s:\".*;:\"" + outIP + "/23\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// sedCommand = "194s:\".*;:\"" + outIP + "/23\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// sedCommand = "197s:\".*;:\"" + outIP + "/24\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Set up FlexRAN
	if OaiObj.Conf.FlexRAN == true {
		// Get flexRAN ip
		var flexranIP string
		var flexranIface string
		OaiObj.Logger.Print("Configure FlexRAN Parameters")
		flexranIP, err = util.GetIPFromDomain(OaiObj.Logger, c.FlexRANDomainName)
		if err != nil {
			OaiObj.Logger.Print(err)
			OaiObj.Logger.Print("Getting IP of FlexRAN failed, try again later")
		} else {
			flexranIface, err = util.GetInterfaceByIP(OaiObj.Logger, outIP)
			if err != nil {
				OaiObj.Logger.Print(err)
				OaiObj.Logger.Print("Getting Interface of FlexRAN failed, try again later")
			} else {
				OaiObj.Logger.Print("FlexRAN Interfacea and IP is ", flexranIface, " ", flexranIP)
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
	// // Start enb
	// if OaiObj.Conf.Test == false {
	// 	OaiObj.Logger.Print("Start enb daemon")
	// 	statusEnabledActive := false
	// 	// wait until the oai-cu start
	// 	for {
	// 		time.Sleep(20 * time.Second)
	// 		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-start")
	// 		if len(retStatus.Stderr) == 0 {
	// 			time.Sleep(15 * time.Second)
	// 			runStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-status")
	// 			for i := 0; i < len(runStatus.Stdout); i++ {
	// 				checkVal := strings.Contains(runStatus.Stdout[i], " enabled ")
	// 				if checkVal == true {
	// 					checkVal = strings.Contains(runStatus.Stdout[i], " active ")
	// 					statusEnabledActive = true
	// 					break
	// 				}
	// 			}
	// 			if statusEnabledActive == true {
	// 				OaiObj.Logger.Print("Service enb is enabled and active yet, exit!!!")
	// 				break
	// 			} else {
	// 				OaiObj.Logger.Print("Service enb is not enabled and active yet, try again later")
	// 				retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-stop")
	// 				time.Sleep(3 * time.Second)
	// 			}
	// 		} else {
	// 			OaiObj.Logger.Print("Start enb failed, try again later")
	// 			retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-stop")
	// 			time.Sleep(3 * time.Second)
	// 		}

	// 	}
	// }
	////////////////
	// Start enb
	time.Sleep(30 * time.Second)
	if OaiObj.Conf.Test == false {
		counter := 0
		OaiObj.Logger.Print("Start enb daemon")
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-start")
		for {
			if len(retStatus.Stderr) == 0 {
				time.Sleep(1 * time.Second)
				counter = counter + 1
				retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-status")
				oairanStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oairanStatus, "inactive")
				if checkInactive != true {
					if counter >= 30 {
						break
					}
				} else {
					OaiObj.Logger.Print("enb is in inactive status, restarting the service")
					util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-stop")
					retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-start")
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start enb failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-start")
				counter = 0
			}
		}
	}
	return nil
}
