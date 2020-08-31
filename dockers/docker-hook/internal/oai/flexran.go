package oai

import (
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"time"
)

func startFlexRAN(OaiObj Oai, buildSnap bool) error {
	if buildSnap == false {
		OaiObj.Logger.Print("Deployment stage; Starting the snap of FlexRAN")
		fmt.Println("Deployment stage; Starting the snap of FlexRAN")
		OaiObj.Logger.Print("Start flexran daemon")
		for {
			retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/flexran.start")
			if len(retStatus.Stderr) == 0 {
				break
			}
			OaiObj.Logger.Print("Start flexran failed, try again later")
			time.Sleep(1 * time.Second)
		}
	} else {
		OaiObj.Logger.Print("Building stage; Skiping start the snap of FlexRAN")
		fmt.Println("Building stage; Skiping start the snap of FlexRAN")
		for {
			retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/flexran.start")
			if len(retStatus.Stderr) == 0 {
				break
			}
			OaiObj.Logger.Print("Start flexran failed, try again later")
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}
