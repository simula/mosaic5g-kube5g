package oai

import (
	"docker-hook/internal/pkg/util"
	"errors"
	"net"
	"os"
	"time"
)

// StartMme : Start MME as a daemon
func startMme(OaiObj Oai, CnAllInOneMode bool) error {
	c := OaiObj.Conf
	mmeConf := c.ConfigurationPathofCN + "mme.conf"
	mmeFdConf := c.ConfigurationPathofCN + "mme_fd.conf"
	mmeBin := c.SnapBinaryPath + "oai-cn.mme"

	spgwIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.SpgwDomainName)
	if err != nil {
		OaiObj.Logger.Print(err)
		spgwIP = OaiObj.Conf.SpgwDomainName
	}

	// Init mme
	if OaiObj.Conf.Test == false {
		OaiObj.Logger.Print("Init mme")
		retStatus := util.RunCmd(OaiObj.Logger, mmeBin+"-init")
		if retStatus.Exit != 0 {
			return errors.New("mme init failed ")
		}
	}
	hostname, _ := os.Hostname()
	// Configure oai-mme
	OaiObj.Logger.Print("Configure mme.conf")

	// hostname
	sedCommand := "s:HSS_HOSTNAME.*;:HSS_HOSTNAME               = \"" + hostname + "\";:g"
	retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set hss domain name in " + mmeConf + " failed")
	}
	// Replace GUMMEI
	OaiObj.Logger.Print("Replace MNC")
	sedCommand = "s/MNC=\"93\"/MNC=\"" + c.MNC + "\\\"/g"
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set GUMMEI in " + mmeConf + " failed")
	}
	OaiObj.Logger.Print("Replace MCC")
	//Replace MCC
	sedCommand = "s:{MCC=\"208\":{MCC=\"" + c.MCC + "\":g"
	OaiObj.Logger.Print(sedCommand)
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set TAI in " + mmeConf + " failed")
	}

	// Get interface ip and replace the default one
	outInterfaceIP := util.GetOutboundIP()
	outInterface, _ := util.GetInterfaceByIP(outInterfaceIP)

	// MME binded interface for S1-C or S1-MME  communication (S1AP): interface name
	sedCommand = "s:MME_INTERFACE_NAME_FOR_S1_MME.*;:MME_INTERFACE_NAME_FOR_S1_MME         = \"" + outInterface + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MME_INTERFACE_NAME_FOR_S1_MME in " + mmeConf + " failed")
	}
	// MME binded interface for S1-C or S1-MME  communication (S1AP): ip address
	sedCommand = "s:MME_IPV4_ADDRESS_FOR_S1_MME.*;:MME_IPV4_ADDRESS_FOR_S1_MME           = \"" + outInterfaceIP + "/24\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set MME_IPV4_ADDRESS_FOR_S1_MME in " + mmeConf + " failed")
	}
	if CnAllInOneMode == true {
		// MME binded interface for S11 communication (GTPV2-C): interface name
		sedCommand = "s:MME_INTERFACE_NAME_FOR_S11_MME.*;:MME_INTERFACE_NAME_FOR_S11_MME        = \"" + "lo" + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_INTERFACE_NAME_FOR_S11_MME in " + mmeConf + " failed")
		}
		// MME binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:MME_IPV4_ADDRESS_FOR_S11_MME.*;:MME_IPV4_ADDRESS_FOR_S11_MME          = \"" + "127.0.11.1" + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_IPV4_ADDRESS_FOR_S11_MME in " + mmeConf + " failed")
		}
	} else {
		// MME binded interface for S11 communication (GTPV2-C): interface name
		sedCommand = "s:MME_INTERFACE_NAME_FOR_S11_MME.*;:MME_INTERFACE_NAME_FOR_S11_MME        = \"" + outInterface + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_INTERFACE_NAME_FOR_S11_MME in " + mmeConf + " failed")
		}
		// MME binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:MME_IPV4_ADDRESS_FOR_S11_MME.*;:MME_IPV4_ADDRESS_FOR_S11_MME          = \"" + outInterfaceIP + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set MME_IPV4_ADDRESS_FOR_S11_MME in " + mmeConf + " failed")
		}
	}

	if CnAllInOneMode == true {
		//S-GW binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11.*;:SGW_IPV4_ADDRESS_FOR_S11          = \"" + "127.0.11.2" + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_IPV4_ADDRESS_FOR_S11 in " + mmeConf + " failed")
		}
	} else {
		//S-GW binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11.*;:SGW_IPV4_ADDRESS_FOR_S11          = \"" + spgwIP + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_IPV4_ADDRESS_FOR_S11 in " + mmeConf + " failed")
		}
	}

	// Replace Identity
	OaiObj.Logger.Print("Configure mme_fd.conf")
	sedCommand = "4s/ubuntu/" + hostname + "/g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set Identity in " + mmeFdConf + " failed")
	}
	// Replace the hostname of Peer conectivity address
	sedCommand = "103s/ubuntu/" + hostname + "/g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set hostname in " + mmeFdConf + " failed")
	}
	// Get the IP address of oai-hss

	if CnAllInOneMode == false {
		hssIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.HssDomainName)
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
			OaiObj.Logger.Print("Valid ip address for oai-hss not get retreived")
			time.Sleep(1 * time.Second)
			hssIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.HssDomainName)
		}

		// replace the ip address of hss
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.1/"+hssIP+"/g", mmeFdConf)
		if retStatus.Exit != 0 {
			return errors.New("Set the ip address of oai-hss in " + mmeFdConf + " failed")
		}
	}

	// oai-cn.mme-start
	if OaiObj.Conf.Test == false {
		OaiObj.Logger.Print("start mme as daemon")
		util.RunCmd(OaiObj.Logger, mmeBin+"-start")
	}
	return nil
}
