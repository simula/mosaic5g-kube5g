package oai

import (
	"errors"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"strings"
	"time"
)

func startENBV2(OaiObj Oai, buildSnap bool) error {
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
	sedCommand := "s:node_function.*;:node_function             = \"" + c.NodeFunction.Default + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Replace MCC
	sedCommand = "s/mcc =.[^;]*/mcc = " + c.MCC + "/g"
	OaiObj.Logger.Print("Replace MCC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MCC in " + enbConf + " failed")
	}

	//Replace MNC
	sedCommand = "s/mnc =.[^;]*/mnc = " + c.MNC + "/g"
	OaiObj.Logger.Print("Replace MNC")
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MNC in " + enbConf + " failed")
	}

	//eutra_band
	sedCommand = "s:eutra_band.*;:eutra_band                                      = " + c.EutraBand.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// downlink_frequency
	sedCommand = "s:downlink_frequency.*;:downlink_frequency                              = " + c.DownlinkFrequency.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// uplink_frequency_offset
	sedCommand = "s:uplink_frequency_offset.*;:uplink_frequency_offset                         = " + c.UplinkFrequencyOffset.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// N_RB_DL
	sedCommand = "s:N_RB_DL.*;:N_RB_DL                                         = " + c.NumberRbDl.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Get Outbound IP and Interface name
	outIP := util.GetOutboundIP(OaiObj.Logger)
	outInterface, err := util.GetInterfaceByIP(OaiObj.Logger, outIP)
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

	// sedCommand = "192s:\".*;:\"" + outIP + "/23\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// sedCommand = "194s:\".*;:\"" + outIP + "/23\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// sedCommand = "197s:\".*;:\"" + outIP + "/24\";:g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Set up FlexRAN
	if (OaiObj.Conf.FlexRAN == true) && (buildSnap == false) {
		// Get flexRAN ip
		var flexranIP string
		OaiObj.Logger.Print("Configure FlexRAN Parameters")
		flexranIP, err = util.GetIPFromDomain(OaiObj.Logger, c.FlexRANDomainName)
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
	sedCommand = "s:parallel_config.*;:parallel_config    = \"" + c.ParallelConfig.Default + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// max_rxgain
	sedCommand = "s:max_rxgain.*;:max_rxgain     = " + c.MaxRxGain.Default + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Get the IP address of oai-mme
	if buildSnap == false {
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
			OaiObj.Logger.Print("Valid ip address for oai-hss not get retreived")
			time.Sleep(1 * time.Second)
			mmeIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.MmeDomainName)
		}
		sedCommand = "175s:\".*;:\"" + mmeIP + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

		OaiObj.Logger.Print("Start waiting for 25 seconds before running oai-enb")
		time.Sleep(170 * time.Second) // 170
		// time.Sleep(150 * time.Second)
		OaiObj.Logger.Print("Finish waiting for 25 seconds before running oai-enb")

		OaiObj.Logger.Print("Start enb daemon")
		// util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-ran.enb-start"}, "/"))
		// util.RunCmd(OaiObj.Logger, "/snap/bin/oai-ran.enb-start")

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
	OaiObj.Logger.Print("enb daemon Started")
	return nil
}
