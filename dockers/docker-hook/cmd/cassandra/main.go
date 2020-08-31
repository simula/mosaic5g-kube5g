/*
# Copyright (c) 2020 Eurecom
################################################################################
# Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The OpenAirInterface Software Alliance licenses this file to You under
# the Apache License, Version 2.0  (the "License"); you may not use this file
# except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#-------------------------------------------------------------------------------
# For more information about the OpenAirInterface (OAI) Software Alliance:
#      contact@openairinterface.org
################################################################################

// This APP is made for configuring the database cassandra for oai-hss inside docker
// Author: Osama Arouk
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
	// Conf    *common.Cfg // config files
}

const (
	logPath = "/root/hook.log"
)

func main() {
	OaiObj := Oai{}
	newFile, _ := os.Create(logPath)
	OaiObj.logFile = newFile
	// OaiObj.Logger = log.New(os.Stdout, "[Mosaic5g] ", log.Ldate|log.Ltime|log.Llongfile)
	OaiObj.Logger = log.New(OaiObj.logFile, "[Mosaic5G-Cassandra-DB] ", log.Ldate|log.Ltime|log.Llongfile)

	OaiObj.Logger.Print("Star Initializing the database Cassandra")
	fmt.Println("Star Initializing database Cassandra")

	// Getting the ip address
	hostname, err := os.Hostname()
	if err != nil {
		OaiObj.Logger.Fatalln("The following error is arised while getting the hostname: ", err)
		fmt.Println("The following error is arised while getting the hostname: ", err)
	}

	// Getting the ip address of the docker container
	addr, err := net.LookupHost(hostname)
	if err != nil {
		OaiObj.Logger.Fatalln("The following error is arised while getting the ip address of", hostname, ": ", err)
		fmt.Println("The following error is arised while getting the ip address of", hostname, ": ", err)
	}
	cassandraIP := addr[0]

	OaiObj.Logger.Print("Database cassandra ip address is ", cassandraIP)
	fmt.Println("Database cassandra ip address is ", cassandraIP)

	//nodetool status
	retStatus := RunCmd(OaiObj.Logger, "nodetool", "status")
	for {
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error of getting the status of database cassandra: ", retStatus.Error)
			fmt.Println("Error of getting the status of Database cassandra: ", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Database cassandra status: ", retStatus.Stdout)
			fmt.Println("Database cassandra status: ", retStatus.Stdout)
			break
		}
		time.Sleep(time.Second * 1)
		retStatus = RunCmd(OaiObj.Logger, "nodetool", "status")
	}

	//cqlsh --file /oai_db.cql cassandraIP
	retStatus = RunCmd(OaiObj.Logger, "cqlsh", "--file", "/oai_db.cql", cassandraIP)
	for {
		if retStatus.Exit != 0 {
			OaiObj.Logger.Print("Error of executing cqlsh for Database cassandra: ", retStatus.Error)
			fmt.Println("Error of executing cqlsh for Database cassandra: ", retStatus.Error)
		} else {
			OaiObj.Logger.Print("Database cassandra cqlsh: ", retStatus.Stdout)
			fmt.Println("Database cassandra cqlsh: ", retStatus.Stdout)
			break
		}
		time.Sleep(1 * time.Second)
		retStatus = RunCmd(OaiObj.Logger, "cqlsh", "--file", "/oai_db.cql", cassandraIP)
	}
	OaiObj.logFile.Close()
}

// RunCmd will run external commands in sync. Return stdout[0].
func RunCmd(logger *log.Logger, cmdName string, args ...string) cmd.Status {
	logger.Print("cmdName=", cmdName)
	for i := 0; i < len(args); i++ {
		logger.Print("args[", i, "]=", args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := <-installSnap.Start() // block and wait
	return finalStatus
}
