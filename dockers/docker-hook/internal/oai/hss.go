package oai

import (
	"errors"
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"os"
	"strings"
	"time"
)

// StartHss : Start HSS as a daemon
func startHss(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("hss.go Starting configuring HSS")
	///////////////////
	retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-conf-get")
	s := strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	snapBinaryPath := "/snap/bin/"
	///////////////////
	// Get working path, Hostname
	hssConf := confPath + "hss.conf"
	hssFdConf := confPath + "hss_fd.conf"
	hssBin := snapBinaryPath + "oai-cn.hss"
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
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, hssConf)
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

// configHss : Config oai-hss
func configHss(OaiObj Oai) error {
	fmt.Println("hss.go Starting initializing OAI-HSS")
	///////////////////
	//c := OaiObj.Conf
	retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-conf-get")
	s := strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	snapBinaryPath := "/snap/bin/"
	///////////////////
	// Get working path, Hostname
	hssConf := confPath + "hss.conf"
	hssFdConf := confPath + "hss_fd.conf"
	hssBin := snapBinaryPath + "oai-cn.hss"
	hostname, _ := os.Hostname()
	fmt.Println("hssConf=", hssConf)
	fmt.Println("hssFdConf=", hssFdConf)
	fmt.Println("hssBin=", hssBin)
	fmt.Println("hostname=", hostname)
	// Strat configuring oai-hss
	OaiObj.Logger.Print("Configure hss.conf")
	//Replace MySQL address
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.1/"+OaiObj.Conf.MysqlDomainName+"/g", hssConf)
	fmt.Println("retStatus.Exit=", retStatus.Exit)
	OaiObj.Logger.Print("retStatus.Exit=", retStatus.Exit)
	if retStatus.Exit != 0 {
		OaiObj.Logger.Print("Set mysql IP in " + hssConf + " failed")
		fmt.Println("Set mysql IP in " + hssConf + " failed")
		return errors.New("Set mysql IP in " + hssConf + " failed")
	}

	// oai-cn.hss-start
	fmt.Println("start hss as daemon")
	OaiObj.Logger.Print("start hss as daemon")
	util.RunCmd(OaiObj.Logger, hssBin+"-start")
	return nil
}

// // StartHss : Start HSS as a daemon
// func startHss(OaiObj Oai) error {
// 	///////////////
// 	OaiObj.Logger.Print("Start oai-hss daemon")
// 	for {
// 		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-start")
// 		if len(retStatus.Stderr) == 0 {
// 			break
// 		}
// 		OaiObj.Logger.Print("Start oai-hss failed, try again later")
// 		time.Sleep(1 * time.Second)
// 	}
// 	fmt.Println("oai-hss is successfully started")
// 	return nil
// }

// // RestartHss : Restart HSS as a daemon
// func restartHss(OaiObj Oai) error {
// 	///////////////
// 	OaiObj.Logger.Print("Retart oai-hss daemon")
// 	for {
// 		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-restart")
// 		if len(retStatus.Stderr) == 0 {
// 			break
// 		}
// 		OaiObj.Logger.Print("Restart oai-hss failed, try again later")
// 		time.Sleep(1 * time.Second)
// 	}
// 	fmt.Println("oai-hss is successfully restarted")
// 	return nil
// }

// // stopHss : Stop HSS as a daemon
// func stopHss(OaiObj Oai) error {
// 	///////////////
// 	OaiObj.Logger.Print("Stop oai-hss daemon")
// 	for {
// 		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-stop")
// 		if len(retStatus.Stderr) == 0 {
// 			break
// 		}
// 		OaiObj.Logger.Print("Stop oai-hss failed, try again later")
// 		time.Sleep(1 * time.Second)
// 	}
// 	fmt.Println("oai-hss is successfully stopped")
// 	return nil
// }
