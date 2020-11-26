/*
#!/usr/local/go/bin/go
################################################################################
* Copyright 2016-2019 Eurecom and Mosaic5G Platforms Authors
* Licensed to the Mosaic5G under one or more contributor license
* agreements. See the NOTICE file distributed with this
* work for additional information regarding copyright ownership.
* The Mosaic5G licenses this file to You under the
* Apache License, Version 2.0  (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
################################################################################
#-------------------------------------------------------------------------------
# For more information about Mosaic5G:
#                                   admin@mosaic-5g.io
# file          hss.go
# brief 		configure the snap of oai-hss v1, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/
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
func startHssV1(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("hss.go Starting configuring HSS")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-cn.hss-status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-cn.hss-conf-get"}, "/"))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")

	hssConf := strings.Join([]string{confPath, "hss.conf"}, "/")
	hssFdConf := strings.Join([]string{confPath, "hss_fd.conf"}, "/")
	hssBin := strings.Join([]string{snapBinaryPath, "oai-cn.hss"}, "/")

	OaiObj.Logger.Print("hssConf=", hssConf)
	fmt.Println("hssConf=", hssConf)

	OaiObj.Logger.Print("hssFdConf=", hssFdConf)
	fmt.Println("hssFdConf=", hssFdConf)

	OaiObj.Logger.Print("hssBin=", hssBin)
	fmt.Println("hssBin=", hssBin)

	hostname, _ := os.Hostname()

	// Strat configuring oai-hss
	OaiObj.Logger.Print("Configure hss.conf")
	//
	var databaseMysqlName, realm string
	if CnAllInOneMode == true {
		databaseMysqlName = OaiObj.Conf.OaiCn.V1[0].OaiHss.DatabaseServiceName
		realm = OaiObj.Conf.OaiCn.V1[0].Realm.Default
	} else {
		databaseMysqlName = OaiObj.Conf.OaiHss.V1[0].DatabaseServiceName
		realm = OaiObj.Conf.OaiHss.V1[0].Realm.Default
	}

	//Replace MySQL address
	mysqlIP, err := util.GetIPFromDomain(OaiObj.Logger, databaseMysqlName)
	if buildSnap == true {
		mysqlIP = databaseMysqlName
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
			OaiObj.Logger.Print("Valid ip address for " + databaseMysqlName + " not yet retreived")
			time.Sleep(1 * time.Second)
			mysqlIP, err = util.GetIPFromDomain(OaiObj.Logger, databaseMysqlName)
		}
	}
	// MYSQL_server
	sedCommand := "s:MYSQL_server.*;:MYSQL_server = \"" + mysqlIP + "\";:g"
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, hssConf)
	if retStatus.Exit != 0 {
		OaiObj.Logger.Print("Set " + databaseMysqlName + " IP in " + hssConf + " failed")
		fmt.Println("Set " + databaseMysqlName + " IP in " + hssConf + " failed")
		return errors.New("Set MYSQL_server in " + hssConf + " failed")
	}

	// Identity
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

	if buildSnap != true {
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

	}
	OaiObj.Logger.Print("END")
	fmt.Println("END")
	return nil
}
