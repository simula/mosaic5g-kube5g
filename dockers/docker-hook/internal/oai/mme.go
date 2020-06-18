package oai

import (
	"errors"
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"os"
	"time"
)

// StartMme : Start MME as a daemon
func startMme(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	c := OaiObj.Conf
	mmeConf := c.ConfigurationPathofCN + "mme.conf"
	mmeFdConf := c.ConfigurationPathofCN + "mme_fd.conf"
	mmeBin := c.SnapBinaryPath + "oai-cn.mme"
	hostname, _ := os.Hostname()

	// Init mme
	if OaiObj.Conf.Test == false {
		OaiObj.Logger.Print("Init mme")
		retStatus := util.RunCmd(OaiObj.Logger, mmeBin+"-init")
		if retStatus.Exit != 0 {
			return errors.New("mme init failed ")
		}
	}

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

	spgwIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.SpgwDomainName)
	///////////////////////////////
	if buildSnap == true {
		spgwIP = "127.0.11.2"
	} else {
		for {
			if err != nil {
				OaiObj.Logger.Print(err)
			} else {
				hostNameSpgw, err := net.LookupHost(spgwIP)
				if len(hostNameSpgw) > 0 {
					break
				} else {
					OaiObj.Logger.Print(err)
				}
			}
			OaiObj.Logger.Print("Valid ip address for spgw not yet retreived")
			time.Sleep(1 * time.Second)
			spgwIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.SpgwDomainName)
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

	// Identity
	realm := "openair4G.eur"           // define the realm
	identity := hostname + "." + realm // use the Hostname we got before
	sedCommand = "s:Identity.*;:Identity = \"" + identity + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set Identity in " + mmeFdConf + " failed")
	}
	// Realm
	sedCommand = "s:Realm.*;:Realm = \"" + realm + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set Realm in " + mmeFdConf + " failed")
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
		if buildSnap == true {
			hssIP = "127.0.0.1"
		} else {
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
				hssIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.HssDomainName)
			}
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

// RestartMme : Restart MME as a daemon
func restartMme(OaiObj Oai) error {
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

// stopMme : Stop MME as a daemon
func stopMme(OaiObj Oai) error {
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
