package oai

import (
	"errors"
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"strings"
	"time"
)

// initSpgwcV2 : Init SPGW
func initSpgwcV2(OaiObj Oai) error {
	return nil
}

// configSpgwcV2 : Config oai-spgw
func configSpgwcV2(OaiObj Oai) error {
	return nil
}

// StartSpgwcV2 : Start SPGW as a daemon
func startSpgwcV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring OAI-SPGWC V2")
	OaiObj.Logger.Print("Starting configuration of OAI-SPGWC V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-spgwc.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	spgwBin := strings.Join([]string{snapBinaryPath, "oai-spgwc"}, "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "conf-get"}, "."))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]

	spgwConf := strings.Join([]string{confPath, "spgwc.conf"}, "/")

	if buildSnap == false {
		// Init spgwc
		OaiObj.Logger.Print("Start Init of oai-spgwc")
		fmt.Println("Start Init of oai-spgwc")
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Init of oai-spgwc is failed")
				fmt.Println("Init of oai-spgwc is failed")
			} else {
				OaiObj.Logger.Print("Init of oai-spgwc is successful")
				fmt.Println("Init of oai-spgwc is successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Init of oai-spgwc")
			fmt.Println("Retrying to Init of oai-spgwc")
			retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "init"}, "."))
		}

		// Configure oai-spgw
		OaiObj.Logger.Print("Configure of oai-spgwc")
		fmt.Println("Configure of oai-spgwc")

		// Get interface IP and outbound interface
		// interfaceIP := util.GetOutboundIP(OaiObj.Logger)
		// outInterface, _ := util.GetInterfaceByIP(OaiObj.Logger, interfaceIP)

		sedCommand := "s:DEFAULT_DNS_IPV4_ADDRESS.*;:DEFAULT_DNS_IPV4_ADDRESS     = \"" + OaiObj.Conf.DNS + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		// if retStatus.Exit != 0 {
		// 	return errors.New("Set DEFAULT_DNS_IPV4_ADDRESS in " + spgwConf + " failed")
		// }
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + OaiObj.Conf.DNS + " in " + spgwConf + " failed")
				fmt.Println("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + OaiObj.Conf.DNS + " in " + spgwConf + " failed")
			} else {
				OaiObj.Logger.Print("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + OaiObj.Conf.DNS + " in " + spgwConf + " successful")
				fmt.Println("Set DEFAULT_DNS_IPV4_ADDRESS to the value " + OaiObj.Conf.DNS + " in " + spgwConf + " successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Set DEFAULT_DNS_IPV4_ADDRESS to the value " + OaiObj.Conf.DNS + " in " + spgwConf)
			fmt.Println("Retrying to Set DEFAULT_DNS_IPV4_ADDRESS to the value " + OaiObj.Conf.DNS + " in " + spgwConf)
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
		}

		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/oai.openair5G.eur/"+"oai.ipv4"+"/g", spgwConf)

		if CnAllInOneMode == false {
			// Get interface IP and outbound interface
			interfaceIP := util.GetOutboundIP(OaiObj.Logger)
			outInterface, _ := util.GetInterfaceByIP(OaiObj.Logger, interfaceIP)
			// INTERFACE_NAME of S11_CP
			sedCommand := "71s:\".*;:\"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set INTERFACE_NAME of S11_CP in " + spgwConf + " failed")
			}
			// IPV4_ADDRESS of S11_CP
			sedCommand = "72s:\".*;:\"" + interfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set IPV4_ADDRESS of S11_CP in " + spgwConf + " failed")
			}
			/////////////////////
			// INTERFACE_NAME of SX
			sedCommand = "167s:\".*;:\"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set INTERFACE_NAME of SX in " + spgwConf + " failed")
			}
			// IPV4_ADDRESS of SX
			sedCommand = "168s:\".*;:\"" + interfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, spgwConf)
			if retStatus.Exit != 0 {
				return errors.New("Set IPV4_ADDRESS of SX in " + spgwConf + " failed")
			}
		}

		// oai.spgwc-start
		time.Sleep(10 * time.Second)
		OaiObj.Logger.Print("start spgwc as daemon")
		fmt.Println("start spgwc as daemon")
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
					OaiObj.Logger.Print("Waiting to make sure that oai-spgwc is working properly")
					fmt.Println("Waiting to make sure that oai-spgwc is working properly")
					if counter >= maxCounter {
						break
					}
				} else {
					OaiObj.Logger.Print("oai-spgwc is in inactive status, restarting the service")
					fmt.Println("oai-spgwc is in inactive status, restarting the service")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "stop"}, "."))
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start oai-spgwc failed, try again later")
				fmt.Println("Start oai-spgwc failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{spgwBin, "start"}, "."))
				counter = 0
			}
		}
	}

	fmt.Println("END of oai-spgwc configuring and starting")
	OaiObj.Logger.Print("END of oai-spgwc configuring and starting")
	return nil
}

// RestartSpgwcV2 : Restart SPGW as a daemon
func restartSpgwcV2(OaiObj Oai) error {
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

// stopSpgwcV2 : Stop SPGW as a daemon
func stopSpgwcV2(OaiObj Oai) error {
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
