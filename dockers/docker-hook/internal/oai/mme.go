package oai

import (
	"errors"
	"docker-hook/internal/pkg/util"
	"os"
)

// StartMme : Start MME as a daemon
func startMme(OaiObj Oai) error {
	c := OaiObj.Conf
	mmeConf := c.ConfigurationPathofCN + "mme.conf"
	mmeFdConf := c.ConfigurationPathofCN + "mme_fd.conf"
	mmeBin := c.SnapBinaryPath + "oai-cn.mme"
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
	sedCommand := "56s/ubuntu/" + hostname + "/g"
	retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set hss domain name in " + mmeConf + " failed")
	}
	// Get interface ip and replace the default one
	outInterfaceIP := util.GetOutboundIP()
	sedCommand = "154s:\".*;:\"" + outInterfaceIP + "/24\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set interface IP in " + mmeConf + " failed")
	}
	// Configure interface name
	outInterface,_ := util.GetInterfaceByIP(outInterfaceIP)
	sedCommand = "153s/lo/" + outInterface + "/g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
	if retStatus.Exit != 0 {
		return errors.New("Set interface name in " + mmeConf + " failed")
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
	// oai-cn.mme-start
	if OaiObj.Conf.Test == false {
		OaiObj.Logger.Print("start mme as daemon")
		util.RunCmd(OaiObj.Logger, mmeBin+"-start")
	}
	return nil
}
