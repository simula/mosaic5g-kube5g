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
	"mosaic5g/docker-hook/internal/pkg/util"
	"strings"
	"time"
)

func startFlexRAN(OaiObj Oai) error {
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
	return nil
}
