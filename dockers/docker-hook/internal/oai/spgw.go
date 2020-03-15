package oai

import (
	"docker-hook/internal/pkg/util"
)

// StartSpgw : Start SPGW as a daemon
func startSpgw(OaiObj Oai) {
	spgwConf := OaiObj.Conf.ConfigurationPathofCN + "spgw.conf"
	spgwBin := OaiObj.Conf.SnapBinaryPath + "oai-cn.spgw"
	// Init spgw
	OaiObj.Logger.Print("Init spgw")
	if OaiObj.Conf.Test == false {
		util.RunCmd(OaiObj.Logger, spgwBin+"-init")
	}
	// Configure oai-spgw
	OaiObj.Logger.Print("Configure spgw.conf")
	// Get interface IP and configure the spgw.conf
	interfaceIP := util.GetOutboundIP()
	sedCommand := "31s:\".*;:\"" + interfaceIP + "/24\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	// Get outbound interface
	outInterface,_ := util.GetInterfaceByIP(interfaceIP)
	sedCommand = "30s/lo/" + outInterface + "/g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	// Get the nameserver from conf
	sedCommand = "101s:\".*;:\"" + OaiObj.Conf.DNS + "\";:g"
	util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
	// oai-cn.spgw-start
	if OaiObj.Conf.Test == false {
		OaiObj.Logger.Print("start spgw as daemon")
		util.RunCmd(OaiObj.Logger, spgwBin+"-start")
	}
}
