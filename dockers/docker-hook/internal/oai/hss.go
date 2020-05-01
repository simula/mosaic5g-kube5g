package oai

import (
	"docker-hook/internal/pkg/util"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// StartHss : Start HSS as a daemon
func startHss(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
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
	mysqlIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.MysqlDomainName)
	if buildSnap == true {
		mysqlIP = OaiObj.Conf.MysqlDomainName
	} else {
		for {
			if err != nil {
				OaiObj.Logger.Print(err)
			} else {
				hostNameMysql, err := net.LookupHost(mysqlIP)
				if len(hostNameMysql) > 0 {
					break
				} else {
					OaiObj.Logger.Print(err)
				}
			}
			OaiObj.Logger.Print("Valid ip address for mysql not yet retreived")
			time.Sleep(1 * time.Second)
			mysqlIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.MysqlDomainName)
		}
	}
	// MYSQL_server
	sedCommand := "s:MYSQL_server.*;:MYSQL_server = \"" + mysqlIP + "\";:g"
	retStatus := util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, hssConf)
	if retStatus.Exit != 0 {
		OaiObj.Logger.Print("Set mysql IP in " + hssConf + " failed")
		fmt.Println("Set mysql IP in " + hssConf + " failed")
		return errors.New("Set MYSQL_server in " + hssConf + " failed")
	}

	// Identity
	realm := "openair4G.eur"           // define the realm
	identity := hostname + "." + realm // use the Hostname we got before
	sedCommand = "s:Identity.*;:Identity = \"" + identity + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, hssFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set Identity in " + hssFdConf + " failed")
	}
	// Realm
	sedCommand = "s:Realm.*;:Realm = \"" + realm + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, hssFdConf)
	if retStatus.Exit != 0 {
		return errors.New("Set Realm in " + hssFdConf + " failed")
	}

	if retStatus.Exit != 0 {
		OaiObj.Logger.Print("Set realm in " + hssFdConf + " failed")
		fmt.Println("Set realm in " + hssFdConf + " failed")
		return errors.New("Set realm in " + hssFdConf + " failed")
	}
	// Init hss
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
	fmt.Println("END")
	return nil
}
