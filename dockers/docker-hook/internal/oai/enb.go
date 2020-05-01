package oai

import (
	"docker-hook/internal/pkg/util"
	"errors"
	"net"
	"strings"
	"time"
)

func startENB(OaiObj Oai) error {
	c := OaiObj.Conf
	enbConf := c.ConfigurationPathofRAN + "enb.band7.tm1.50PRB.usrpb210.conf"
	// Replace MCC
	sedCommand := "s/mcc =.[^;]*/mcc = " + c.MCC + "/g"
	OaiObj.Logger.Print(sedCommand)
	retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
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

	sedCommand = "s:eutra_band.*;:eutra_band              			      = " + c.EutraBand + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	sedCommand = "s:downlink_frequency.*;:downlink_frequency                              = " + c.DownlinkFrequency + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	sedCommand = "s:uplink_frequency_offset.*;:uplink_frequency_offset                         = " + c.UplinkFrequencyOffset + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// max_rxgain
	sedCommand = "s:max_rxgain.*;:max_rxgain     = " + c.MaxRxGain + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// N_RB_DL
	sedCommand = "s:N_RB_DL.*;:N_RB_DL              			      = " + c.NumberRbDl + ";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// parallel_config
	sedCommand = "s:parallel_config.*;:parallel_config    = \"" + c.ParallelConfig + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// configure the oai-ran with mme
	// Get the IP address of oai-hss
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

	time.Sleep(20 * time.Second)
	sedCommand = "175s:\".*;:\"" + mmeIP + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// Get Outbound IP
	outIP := util.GetOutboundIP()
	outInterface, err := util.GetInterfaceByIP(outIP)
	if err != nil {
		OaiObj.Logger.Print(err)
	}
	OaiObj.Logger.Print("Outbound Interface and IP is ", outInterface, " ", outIP)
	// Replace interface
	// sedCommand = "s/eno1/" + outInterface + "/g"
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	sedCommand = "s:ENB_INTERFACE_NAME_FOR_S1_MME.*;:ENB_INTERFACE_NAME_FOR_S1_MME    = \"" + outInterface + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	sedCommand = "s:ENB_INTERFACE_NAME_FOR_S1U.*;:ENB_INTERFACE_NAME_FOR_S1U    = \"" + outInterface + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)

	// Replace enb IP
	sedCommand = "192s:\".*;:\"" + outIP + "/23\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	sedCommand = "194s:\".*;:\"" + outIP + "/23\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	sedCommand = "197s:\".*;:\"" + outIP + "/24\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, enbConf)
	// Set up FlexRAN
	if OaiObj.Conf.FlexRAN == true {
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
	// Start enb
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
