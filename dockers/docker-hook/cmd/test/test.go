package main

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

// This hook is made for installing and configuring snaps inside docker
// Author: Osama Arouk, Kevin Hsi-Ping Hsu
*/
import (
	"errors"
	"fmt"
	"log"
	"mosaic5g/docker-hook/internal/oai"

	"github.com/go-cmd/cmd"
)

const (
	logPath  = "/home/cigarier/go/src/mosaic5g/docker-hook/cmd/test/hook.log"
	confPath = "/home/cigarier/go/src/mosaic5g/docker-hook/cmd/test/conf.yaml"
)

func main() {
	// Initialize oai struct
	OaiObj := oai.Oai{}
	err := OaiObj.Init(logPath, confPath)
	if err != nil {
		panic(err)
	}

	OaiObj.Logger.Print("Init of OAI is successful")
	fmt.Println("Init of OAI is successful")

	//Install snap core
	OaiObj.Logger.Print("Installing snap")
	fmt.Println("Installing snap")
	oai.InstallSnap(OaiObj)
	interfaceName := "wlp2s0"
	ret := RunCmd(OaiObj.Logger, "ifconfig", interfaceName)
	if ret.Exit != 0 {
		fmt.Println("127.0.1.10", errors.New("Fail to run ifconfig"))
	}
	if len(ret.Stdout) <= 0 {
		fmt.Println("127.0.1.10", errors.New("Fail to get result"))
	}
	OaiObj.Logger.Print("ret=", ret)
	fmt.Println("ret=", ret)
	log.Fatal(errors.New("Fail to get result log.Fatal"))
	PrintFunc(OaiObj.Logger, "finalStatus.Cmd", errors.New("Fail to get result"))
}

// RunCmd will run external commands in sync. Return stdout[0].
func RunCmd(logger *log.Logger, cmdName string, args ...string) cmd.Status {
	PrintFunc(logger, "cmdName= "+cmdName)
	for i := 0; i < len(args); i++ {
		PrintFunc(logger, "args["+string(i)+"]="+args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := <-installSnap.Start() // block and wait
	PrintFunc(logger, finalStatus)
	PrintFunc(logger, finalStatus.Cmd)

	return finalStatus
}

// //PrintFunc will return if this package is already exist or not
// func PrintFunc(logger *log.Logger, massage interface{}) {
// 	logger.Print(massage)
// 	fmt.Println(massage)
// }

//PrintFunc PrintFunc
func PrintFunc(logger *log.Logger, args ...interface{}) {
	switch len(args) {
	case 1:
		logger.Print(args[0])
		fmt.Println(args[0])
	case 2:
		logger.Print(args[0], args[1])
		fmt.Println(args[0], args[1])
	default:
		logger.Print("Unexpected number of variables")
		panic("Unexpected number of variables")
	}
}
