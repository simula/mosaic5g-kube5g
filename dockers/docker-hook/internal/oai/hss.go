package oai

import (
	"errors"
	"fmt"
	"docker-hook/internal/pkg/util"
	"os"
	"strings"
)

// StartHss : Start HSS as a daemon
func startHss(OaiObj Oai) error {
	fmt.Println("hss.go Starting configuring HSS")
	// Get working path, Hostname
	hssConf := OaiObj.Conf.ConfigurationPathofCN + "hss.conf"
	hssFdConf := OaiObj.Conf.ConfigurationPathofCN + "hss_fd.conf"
	hssBin := OaiObj.Conf.SnapBinaryPath + "oai-cn.hss"
	hostname, _ := os.Hostname()
	fmt.Println("hssConf=", hssConf)
	fmt.Println("hssFdConf=", hssFdConf)
	fmt.Println("hssBin=", hssBin)
	fmt.Println("hostname=", hostname)
	// Strat configuring oai-hss
	OaiObj.Logger.Print("Configure hss.conf")
	//Replace MySQL address
	retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.1/"+OaiObj.Conf.MysqlDomainName+"/g", hssConf)
	fmt.Println("retStatus.Exit=", retStatus.Exit)
	OaiObj.Logger.Print("retStatus.Exit=", retStatus.Exit)
	if retStatus.Exit != 0 {
		OaiObj.Logger.Print("Set mysql IP in " + hssConf + " failed")
		fmt.Println("Set mysql IP in " + hssConf + " failed")
		return errors.New("Set mysql IP in " + hssConf + " failed")
	}
	// Replace Identity
	OaiObj.Logger.Print("Configure hss_fd.conf")
	fmt.Println("Configure hss_fd.conf")
	identity := hostname + ".openair4G.eur" // use the Hostname we got before
	syntax := "s/ubuntu.openair4G.eur/" + identity + "/g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", syntax, hssFdConf)
	OaiObj.Logger.Print("identity=", identity)
	fmt.Println("identity=", identity)
	OaiObj.Logger.Print("syntax=", syntax)
	fmt.Println("syntax=", syntax)
	OaiObj.Logger.Print("retStatus=", retStatus)
	fmt.Println("retStatus=", retStatus)

	if retStatus.Exit != 0 {
		OaiObj.Logger.Print("Set realm in " + hssFdConf + " failed")
		fmt.Println("Set realm in " + hssFdConf + " failed")
		return errors.New("Set realm in " + hssFdConf + " failed")
	}
	// Init hss
	if OaiObj.Conf.Test == false {
		fmt.Println("Init hss")
		OaiObj.Logger.Print("Init hss")
		fmt.Println(OaiObj.Logger, hssBin+"-init")
		retStatus = util.RunCmd(OaiObj.Logger, hssBin+"-init")
		fmt.Println("retStatus", retStatus)
		fmt.Println("retStatus.Stderr", retStatus.Stderr)
		for {
			fail := false
			for i := 0; i < len(retStatus.Stderr); i++ {
				if strings.Contains(retStatus.Stderr[i], "ERROR") {
					fmt.Println("Init error, re-run again")
					OaiObj.Logger.Println("Init error, re-run again")
					fail = true
				}
			}
			if fail {
				retStatus = util.RunCmd(OaiObj.Logger, hssBin+"-init")
			} else {
				break
			}
		}
		// oai-cn.hss-start
		fmt.Println("start hss as daemon")
		OaiObj.Logger.Print("start hss as daemon")
		util.RunCmd(OaiObj.Logger, hssBin+"-start")
	} else {
		fmt.Println("OaiObj.Conf.Test=True")
	}
	fmt.Println("END")
	return nil
}
