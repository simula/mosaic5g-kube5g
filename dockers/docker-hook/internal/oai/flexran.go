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
# file          flexran.go
# brief 		configure the snap of flexran, and start it
# authors:
	- Osama Arouk (arouk@eurecom.fr)
	- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
*-------------------------------------------------------------------------------
*/
package oai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ValueFlexran
type ValueFlexran struct {
	BsID int64 `json:"bs_id"`
}

type FlexranStats struct {
	DateTime  string         `json:"date_time"`
	EnbConfig []ValueFlexran `json:"eNB_config"`
}

var flexranStats string = `
{
	"reports": [
	  {
		"reportFrequency": "FLSRF_PERIODICAL",
		"sf": 100,
		"cellReports": [
		  "FLCST_NOISE_INTERFERENCE"
		],
		"ueReports": [
		  "FLUST_PHR",
		  "FLUST_DL_CQI",
		  "FLUST_BSR",
		  "FLUST_RLC_BS",
		  "FLUST_MAC_CE_BS",
		  "FLUST_UL_CQI",
		  "FLUST_RRC_MEASUREMENTS",
		  "FLUST_PDCP_STATS",
		  "FLUST_MAC_STATS",
		  "FLUST_GTP_STATS",
		  "FLUST_S1AP_STATS"
		]
	  }
	]
  }
`

// func startFlexRAN(OaiObj Oai) error {
func startFlexRAN(OaiObj Oai) {
	var msg string = ""
	OaiObj.Logger.Print("Start flexran daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/flexran.start")
		oairanStatus := strings.Join(retStatus.Stdout, " ")
		OaiObj.Logger.Print(oairanStatus)
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Start flexran failed, try again later")
		time.Sleep(1 * time.Second)
	}

	OaiObj.Logger.Print("Start flexran daemon")
	var flexranGetStats, currentEnbID string = "http://localhost:9999/stats", ""
	var valListEnbID FlexranStats

	var eNBIdList = make([]string, 1000)
	var counter int = 0
	for {
		time.Sleep(time.Duration(5) * time.Second)
		resp, err := http.Get(flexranGetStats)
		// for {
		if err != nil {
			OaiObj.Logger.Print(err)
		} else {
			defer resp.Body.Close()
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			err := json.Unmarshal([]byte(bodyString), &valListEnbID)
			if err != nil {
				msg = "Error parsing the stats\n Error:" + err.Error() + "\n"
				fmt.Println(msg)
				OaiObj.Logger.Println(msg)
			} else {
				// fmt.Println(valListEnbID.EnbConfig)
				// OaiObj.Logger.Println(valListEnbID.EnbConfig)
				for i := 0; i < len(valListEnbID.EnbConfig); i++ {
					// fmt.Println(valListEnbID.EnbConfig[i].BsID)
					// OaiObj.Logger.Println(valListEnbID.EnbConfig[i].BsID)
					currentEnbID = strconv.FormatInt(valListEnbID.EnbConfig[i].BsID, 10)
					found, _ := util.FindValInList(eNBIdList, currentEnbID)
					// OaiObj.Logger.Println(len(eNBIdList))
					if !found || len(eNBIdList) == 0 {
						eNBIdList[counter] = currentEnbID
						counter++
						msg = "new eNB is attached to flexran with the id " + strconv.FormatInt(valListEnbID.EnbConfig[i].BsID, 10)
						fmt.Println(msg)
						OaiObj.Logger.Println(msg)
						SetConfStats(OaiObj, currentEnbID)
					}

				}
			}
		}
		// }
	}
	// return nil
}

// SetConfStats set the configuration of the stats for sepcific eNB
func SetConfStats(OaiObj Oai, enbID string) {
	var msg string = ""
	// add the stats here
	if OaiObj.Conf.Flexran[0].Stats != "" {
		flexranStats = OaiObj.Conf.Flexran[0].Stats
	}
	// dump the users in json file
	err := ioutil.WriteFile(OaiObj.FlexranStatsPath, []byte(flexranStats), 0644)
	if err != nil {
		msg = "Error while trying to dump flexran stats to json file \n Users=" + flexranStats + "\n Error=" + err.Error()
		OaiObj.Logger.Print(msg)
		fmt.Println(msg)
	} else {
		// curl -XPOST http://127.0.0.1:9999/stats/conf/enb/ --data-binary @stats.json
		urlStats := "http://127.0.0.1:9999/stats/conf/enb/" + enbID

		file, err := os.Open(OaiObj.FlexranStatsPath)
		if err != nil {
			msg = "error while openning the file " + OaiObj.FlexranStatsPath + "\nError: " + err.Error()
			OaiObj.Logger.Print(msg)
			fmt.Println(msg)
		} else {
			_, err := http.Post(urlStats, "application/json", file)
			if err != nil {
				msg = "error while setting Set statistics configuration. urlStats=" + urlStats
			} else {
				msg = "The following config stats is applied to the enb " + enbID + "\n config-stats:" + flexranStats
			}
			OaiObj.Logger.Print(msg)
			fmt.Println(msg)
		}

	}

}
