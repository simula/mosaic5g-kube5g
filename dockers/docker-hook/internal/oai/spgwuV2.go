package oai

import (
	"errors"
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"strings"
	"time"
)

// initSpgwuV2 : Init SPGW
func initSpgwuV2(OaiObj Oai) error {
	return nil
}

// configSpgwuV2 : Config oai-spgw
func configSpgwuV2(OaiObj Oai) error {
	return nil
}

// StartSpgwuV2 : Start SPGW as a daemon
func startSpgwuV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring OAI-SPGWu V2")
	OaiObj.Logger.Print("Starting configuration of OAI-SPGWu V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-spgwu.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	spgwBin := strings.Join([]string{snapBinaryPath, "oai-spgwu"}, "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "conf-get"}, "."))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]

	spgwConf := strings.Join([]string{confPath, "spgwu.conf"}, "/")
	// spgwBin := strings.Join([]string{confPath, "oai-spgwu"}, "/")

	if buildSnap == false {
		// Init spgwu
		OaiObj.Logger.Print("Start Init of oai-spgwu")
		fmt.Println("Start Init of oai-spgwu")
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Init of oai-spgwu is failed")
				fmt.Println("Init of oai-spgwu is failed")
			} else {
				OaiObj.Logger.Print("Init of oai-spgwu is successful")
				fmt.Println("Init of oai-spgwu is successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Init of oai-spgwu")
			fmt.Println("Retrying to Init of oai-spgwu")
			retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		}

		// Configure oai-spgw
		OaiObj.Logger.Print("Configure of oai-spgwu")
		fmt.Println("Configure of oai-spgwu")

		// Get interface IP and outbound interface
		interfaceIP := util.GetOutboundIP(OaiObj.Logger)
		outInterface, _ := util.GetInterfaceByIP(OaiObj.Logger, interfaceIP)
		// INTERFACE_NAME of S1U_S12_S4_UP
		sedCommand := "59s:\".*;:\"" + outInterface + "\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		// IPV4_ADDRESS of S1U_S12_S4_UP
		sedCommand = "60s:\".*;:\"" + interfaceIP + "/24\";:g"
		util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)

		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/wlp2s0/"+outInterface+"/g", spgwConf)
		if CnAllInOneMode == false {
			// Get interface IP and outbound interface
			interfaceIP := util.GetOutboundIP(OaiObj.Logger)
			outInterface, _ := util.GetInterfaceByIP(OaiObj.Logger, interfaceIP)
			// INTERFACE_NAME of SX
			sedCommand := "72s:\".*;:\"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set INTERFACE_NAME of SX in " + spgwConf + " failed")
			}
			// IPV4_ADDRESS of SX
			sedCommand = "73s:\".*;:\"" + interfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set IPV4_ADDRESS of SX in " + spgwConf + " failed")
			}
			//////////////////////////////////////////////////////////////
			spgwcIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.SpgwcDomainName)
			if buildSnap == true {
				spgwcIP = "127.0.12.1"
			} else {
				for {
					if err != nil {
						OaiObj.Logger.Print(err)
					} else {
						hostNameSpgwc, err := net.LookupHost(spgwcIP)
						if len(hostNameSpgwc) > 0 {
							break
						} else {
							OaiObj.Logger.Print(err)
						}
					}
					OaiObj.Logger.Print("Valid ip address for oai-spgwc not yet retreived")
					time.Sleep(1 * time.Second)
					spgwcIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.SpgwcDomainName)
				}
			}
			sedCommand = "s:IPV4_ADDRESS=\"127.0.12.1.*;:IPV4_ADDRESS=\"" + spgwcIP + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				fmt.Println("Set IPV4_ADDRESS in " + spgwConf + " failed")
			}
			//////////////////////////////////////////////////////////////
		}
		// oai.spgwu-start
		time.Sleep(10 * time.Second)
		OaiObj.Logger.Print("start spgwu as daemon")
		fmt.Println("start spgwu as daemon")

		retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
		counter := 0
		maxCounter := 2
		for {
			if len(retStatus.Stderr) == 0 {
				time.Sleep(5 * time.Second)
				counter = counter + 1
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "status"}, "."))
				oairanStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oairanStatus, "inactive")
				if checkInactive != true {
					OaiObj.Logger.Print("Waiting to make sure that oai-spgwu is working properly")
					fmt.Println("Waiting to make sure that oai-spgwu is working properly")
					if counter >= maxCounter {
						break
					}
				} else {
					OaiObj.Logger.Print("oai-spgwu is in inactive status, restarting the service")
					fmt.Println("oai-spgwu is in inactive status, restarting the service")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "stop"}, "."))
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start oai-spgwu failed, try again later")
				fmt.Println("Start oai-spgwu failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
				counter = 0
			}
		}

	}
	fmt.Println("END of oai-spgwu configuring and starting")
	OaiObj.Logger.Print("END of oai-spgwu configuring and starting")
	return nil
}

// RestartSpgwuV2 : Restart SPGW as a daemon
func restartSpgwuV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Restart oai-spgw daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.spgw-restart")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Restart oai-spgw failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-spgw is successfully restarted")
	return nil
}

// stopSpgwuV2 : Stop SPGW as a daemon
func stopSpgwuV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Stop oai-spgw daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.spgw-stop")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Stop oai-spgw failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-spgw is successfully stopped")
	return nil
}
