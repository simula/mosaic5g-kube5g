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
# file          main.go
# brief         create docker-hook for oai-cassandra so that the database of oai-cn v2 will be added to cassandra when starting docker
# authors:
		- Osama Arouk (arouk@eurecom.fr)
*-------------------------------------------------------------------------------
*/
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-cmd/cmd"
)

// Oai stores the log and conf
type Oai struct {
	logFile *os.File    // File for log to write something
	Logger  *log.Logger // Collect log
}

const (
	logPath = "/root/hook.log"
)

func main() {
	OaiObj := Oai{}
	newFile, _ := os.Create(logPath)
	OaiObj.logFile = newFile
	OaiObj.Logger = log.New(OaiObj.logFile, "[Debug]"+time.Now().Format("2006-01-02 15:04:05")+" ", log.Lshortfile)

	OaiObj.Logger.Print("Star Initializing Cassandra DB")
	// Getting the ip address
	hostname, _ := os.Hostname()
	cassandraIP, err := GetIPFromDomain(OaiObj.Logger, hostname)
	for {
		if err != nil {
			OaiObj.Logger.Print(err)
		} else {
			hostNameMysql, err := net.LookupHost(cassandraIP)
			if len(hostNameMysql) > 0 {
				break
			} else {
				OaiObj.Logger.Print(err)
			}
		}
		OaiObj.Logger.Print("Valid ip address for mysql not yet retreived")
		time.Sleep(1 * time.Second)
		cassandraIP, err = GetIPFromDomain(OaiObj.Logger, cassandraIP)
	}
	OaiObj.Logger.Print("assandra DB Ip address, \t", cassandraIP)

	//nodetool status
	retStatus := RunCmd(OaiObj.Logger, "nodetool", "status")
	for {
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error of getting the status of assandra DB: \n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("assandra DB status: \n", retStatus.Stdout)
			break
		}
		time.Sleep(1 * time.Second)
		retStatus = RunCmd(OaiObj.Logger, "nodetool", "status")
	}

	//cqlsh --file /oai_db.cql cassandraIP
	retStatus = RunCmd(OaiObj.Logger, "cqlsh", "--file", "/oai_db.cql", cassandraIP)

	for {
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error of executing cqlsh for assandra DB: \n", retStatus.Error)
		} else {
			OaiObj.Logger.Print("assandra DB cqlsh: \n", retStatus.Stdout)
			break
		}
		time.Sleep(1 * time.Second)
		retStatus = RunCmd(OaiObj.Logger, "cqlsh", "--file", "/oai_db.cql", cassandraIP)
	}

	OaiObj.logFile.Close()
}

// RunCmd will run external commands in sync. Return stdout[0].
func RunCmd(logger *log.Logger, cmdName string, args ...string) cmd.Status {
	fmt.Println("cmdName=", cmdName)
	for i := 0; i < len(args); i++ {
		fmt.Println("args[", i, "]=", args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := <-installSnap.Start() // block and wait
	return finalStatus
}

// GetIPFromDomain GetIPFromDomain
func GetIPFromDomain(logger *log.Logger, domain string) (string, error) {
	addr, err := net.LookupHost(domain)
	if err != nil {
		logger.Print("Failed to get IP from domain,err: ", err)
		return "", err
	}
	return addr[0], nil
}
