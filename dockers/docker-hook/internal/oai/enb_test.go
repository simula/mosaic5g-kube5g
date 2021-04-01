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
# file          enb.go
# brief 		configure the snap of oai-ran, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/
package oai

import (
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/common"
	"testing"
)

func Test_naive_changeParamTxGain(t *testing.T) {

	logPath := "./hook.log"
	confPath := "./conf.yaml"
	usersPath := "./users.json"
	flexranStatsPath := "./flexran_stats.json"

	instance_c := common.CfgGlobal{}
	c := &instance_c
	c.OaiEnb = []common.CfgOaiEnb{
		common.CfgOaiEnb{
			NidCellMbsfn: struct {
				Default     string "yaml:\"default\""
				Description string "yaml:\"description\""
			}{
				Default:     "5",
				Description: "nothing",
			},
		},
	}

	// fmt.Printf("Len : %#v \n", len(c.OaiEnb))
	// fmt.Printf("Value : %#v \n", c.OaiEnb[0].TxGain.Default)
	oai := Oai{}

	err := oai.Init(logPath, confPath, usersPath, flexranStatsPath)
	if err != nil {
		t.Error("oai created failed")
		return
	}
	enb_path := "enb_sample.conf"
	status := changeParamNidCellMbsfn(c, oai, enb_path)
	if status != 0 {
		t.Error(status)
		t.Error("Test failed")
	} else {
		fmt.Printf("test passed")
	}
}

func Test_naive_changeParamPuschProcThreads(t *testing.T) {

	logPath := "./hook.log"
	confPath := "./conf.yaml"
	usersPath := "./users.json"
	flexranStatsPath := "./flexran_stats.json"

	instance_c := common.CfgGlobal{}
	c := &instance_c
	c.OaiGnb = []common.CfgOaiGnb{
		common.CfgOaiGnb{
			PuschProcThreads: "10",
		},
	}

	// fmt.Printf("Len : %#v \n", len(c.OaiEnb))
	// fmt.Printf("Value : %#v \n", c.OaiEnb[0].TxGain.Default)
	oai := Oai{}

	err := oai.Init(logPath, confPath, usersPath, flexranStatsPath)
	if err != nil {
		t.Error("oai created failed")
		return
	}
	gnb_path := "gnb.band261.tm1.32PRB.usrpn300.conf.txt"
	status := replaceExistingPuschProcThreads(c, oai, gnb_path)
	if status != 0 {
		t.Error("Test failed")
	} else {
		fmt.Printf("test passed")
	}
}
