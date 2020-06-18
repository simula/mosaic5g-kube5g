package oai

import (
	"mosaic5g/docker-hook/internal/pkg/util"
	"time"
)

func startFlexRAN(OaiObj Oai) error {
	OaiObj.Logger.Print("Start flexran daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/flexran.start")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Start flexran failed, try again later")
		time.Sleep(1 * time.Second)
	}
	return nil
}
