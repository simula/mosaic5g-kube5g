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
# file          cfg.go
# brief 		define the configuration of the snaps, check the file cmd/test/conf.yaml to see an example of such configuration
# authors:
		- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
		- Osama Arouk (arouk@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package util

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/go-cmd/cmd"
)

// RunCmd will run external commands in sync. Return stdout[0].
func RunCmd(logger *log.Logger, cmdName string, args ...string) cmd.Status {
	fmt.Println("cmdName=", cmdName)
	logger.Print("cmdName=", cmdName)
	for i := 0; i < len(args); i++ {
		fmt.Println("args[", i, "]=", args[i])
		logger.Print("args[", i, "]=", args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := <-installSnap.Start() // block and wait
	// logger.Print(finalStatus.Cmd)
	logger.Print(finalStatus)
	fmt.Println("19 installSnap=", installSnap)
	fmt.Println("20 finalStatus=", finalStatus)
	return finalStatus
}

//CheckSnapPackageExist will return if this package is already exist or not
func CheckSnapPackageExist(logger *log.Logger, packageName string) (bool, error) {
	if len(packageName) <= 0 {
		return false, errors.New("Input package name is empty")
	}
	retStatus := RunCmd(logger, "snap", "list")
	if retStatus.Exit != 0 {
		return false, errors.New("snap list return non-zero")
	}
	for i := 0; i < len(retStatus.Stdout); i++ {
		if strings.Contains(retStatus.Stdout[i], packageName) {
			logger.Println("Package: ", packageName, " Exist")
			return true, nil
		}

	}
	logger.Println("Package: ", packageName, " does not Exist")
	return false, nil
}

//GetInterfaceIP will get the ip of the interface. If failed, it'll return a default (127.0.1.10) value
func GetInterfaceIP(logger *log.Logger, interfaceName string) (string, error) {
	ret := RunCmd(logger, "ifconfig", interfaceName)
	if ret.Exit != 0 {
		return "127.0.1.10", errors.New("Fail to run ifconfig")
	}
	if len(ret.Stdout) <= 0 {
		return "127.0.1.10", errors.New("Fail to get result")
	}
	i := 0
	space := " "
	for {
		if ret.Stdout[1][27+i+1] == space[0] {
			break
		}
		i++
	}
	return ret.Stdout[1][20 : 27+i+1], nil
}

//GetIPFromDomain will get the IP of the domain
func GetIPFromDomain(logger *log.Logger, domain string) (string, error) {
	addr, err := net.LookupHost(domain)
	if err != nil {
		logger.Print("Failed to get IP from domain,err: ", err)
		return "", err
	}
	return addr[0], nil
}

// GetOutboundIP gets preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// GetInterfaceByIP can get interface name from IP
func GetInterfaceByIP(targetIP string) (string, error) {
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.String() == targetIP {

				return i.Name, nil
			}
			// process IP address
		}
	}
	return "", err
}
