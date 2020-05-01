package oai

import (
	"docker-hook/internal/pkg/util"
	"errors"
)

// StartSpgw : Start SPGW as a daemon
func startSpgw(OaiObj Oai, CnAllInOneMode bool) error {
	spgwConf := OaiObj.Conf.ConfigurationPathofCN + "spgw.conf"
	spgwBin := OaiObj.Conf.SnapBinaryPath + "oai-cn.spgw"
	// Init spgw
	OaiObj.Logger.Print("Init spgw")
	if OaiObj.Conf.Test == false {
		util.RunCmd(OaiObj.Logger, spgwBin+"-init")
	}
	// Configure oai-spgw
	OaiObj.Logger.Print("Configure spgw.conf")

	// Get interface IP and outbound interface
	interfaceIP := util.GetOutboundIP()
	outInterface, _ := util.GetInterfaceByIP(interfaceIP)

	if CnAllInOneMode == true {
		// S-GW binded interface for S11 communication (GTPV2-C): interface name
		sedCommand := "s:SGW_INTERFACE_NAME_FOR_S11.*;:SGW_INTERFACE_NAME_FOR_S11              = \"" + "lo" + "\";:g"
		retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_INTERFACE_NAME_FOR_S11 in " + spgwConf + " failed")
		}
		// S-GW binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11.*;:SGW_IPV4_ADDRESS_FOR_S11                = \"" + "127.0.11.2" + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_IPV4_ADDRESS_FOR_S11 in " + spgwConf + " failed")
		}
	} else {
		// S-GW binded interface for S11 communication (GTPV2-C): interface name
		sedCommand := "s:SGW_INTERFACE_NAME_FOR_S11.*;:SGW_INTERFACE_NAME_FOR_S11              = \"" + outInterface + "\";:g"
		retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_INTERFACE_NAME_FOR_S11 in " + spgwConf + " failed")
		}
		// S-GW binded interface for S11 communication (GTPV2-C): ip address
		sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11.*;:SGW_IPV4_ADDRESS_FOR_S11                = \"" + interfaceIP + "/8\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		if retStatus.Exit != 0 {
			return errors.New("Set SGW_IPV4_ADDRESS_FOR_S11 in " + spgwConf + " failed")
		}
	}

	// S-GW binded interface for S1-U communication (GTPV1-U): interface name
	sedCommand := "s:SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP.*;:SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP    = \"" + outInterface + "\";:g"
	retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set SGW_INTERFACE_NAME_FOR_S1U_S12_S4_UP in " + spgwConf + " failed")
	}
	// S-GW binded interface for S1-U communication (GTPV1-U): ip address
	sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP.*;:SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP      = \"" + interfaceIP + "/24\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set SGW_IPV4_ADDRESS_FOR_S1U_S12_S4_UP in " + spgwConf + " failed")
	}

	// # P-GW binded interface for SGI (egress/ingress internet traffic): interface name
	sedCommand = "s:PGW_INTERFACE_NAME_FOR_SGI.*;:PGW_INTERFACE_NAME_FOR_SGI            = \"" + outInterface + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set PGW_INTERFACE_NAME_FOR_SGI in " + spgwConf + " failed")
	}

	// # DNS address communicated to UEs
	sedCommand = "s:DEFAULT_DNS_IPV4_ADDRESS.*;:DEFAULT_DNS_IPV4_ADDRESS     = \"" + OaiObj.Conf.DNS + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set DEFAULT_DNS_IPV4_ADDRESS in " + spgwConf + " failed")
	}

	secondaryDNS := "8.8.4.4"
	sedCommand = "s:DEFAULT_DNS_SEC_IPV4_ADDRESS.*;:DEFAULT_DNS_SEC_IPV4_ADDRESS = \"" + secondaryDNS + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	if retStatus.Exit != 0 {
		return errors.New("Set DEFAULT_DNS_SEC_IPV4_ADDRESS in " + spgwConf + " failed")
	}

	// oai-cn.spgw-start
	if OaiObj.Conf.Test == false {
		OaiObj.Logger.Print("start spgw as daemon")
		util.RunCmd(OaiObj.Logger, spgwBin+"-start")
	}
	return nil
}
