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
# file          oai.go
# brief 		initialize the docker before starting installing the required snaps
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package oai

import (
	"fmt"
	"log"
	"mosaic5g/docker-hook/internal/pkg/common"
	"os"
	"time"
)

// Oai stores the log and conf
type Oai struct {
	logFile *os.File          // File for log to write something
	Logger  *log.Logger       // Collect log
	Conf    *common.CfgGlobal // config files

}

// Init the Oai with log and conf
func (me *Oai) Init(logPath string, confPath string) error {
	newFile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	me.logFile = newFile
	me.Logger = log.New(me.logFile, "[Debug]"+time.Now().Format("2006-01-02 15:04:05")+" ", log.Lshortfile)
	me.Conf = new(common.CfgGlobal)
	err = me.Conf.GetConf(me.Logger, confPath)
	if err != nil {
		return err
	}
	me.Logger.Print("Configs: ", me.Conf)
	fmt.Println("Configs: ", *me.Conf)
	// os.Exit(0)

	return nil
}

// Clean will Close the logFile and clean up Obj
func (me *Oai) Clean() {
	me.logFile.Close()
}
