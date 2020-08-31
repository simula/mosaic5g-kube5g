package oai

import (
	"errors"
	"fmt"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"strings"
	"time"
)

// StartMmeV2 : Start MME as a daemon
func startMmeV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring OAI-MME V2")
	OaiObj.Logger.Print("Starting configuration of OAI-MME V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-mme.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")
	mmeBin := strings.Join([]string{snapBinaryPath, "oai-mme"}, "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "conf-get"}, "."))
	fmt.Println(retStatus)
	OaiObj.Logger.Print(retStatus)
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]
	mmeConf := strings.Join([]string{confPath, "mme.conf"}, "/")
	mmeFdConf := strings.Join([]string{confPath, "mme_fd.conf"}, "/")

	// hostname, _ := os.Hostname()

	if buildSnap == false {
		retStatus := util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "init"}, "."))
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Init of oai-mme is failed")
				fmt.Println("Init of oai-mme is failed")
			} else {
				OaiObj.Logger.Print("Init of oai-mme is successful")
				fmt.Println("Init of oai-mme is successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Init of oai-mme")
			fmt.Println("Retrying to Init of oai-mme")
			retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "init"}, "."))
		}

		// Get interface ip and replace the default one
		outInterfaceIP := util.GetOutboundIP(OaiObj.Logger)
		outInterface, _ := util.GetInterfaceByIP(OaiObj.Logger, outInterfaceIP)

		// MME_INTERFACE_NAME_FOR_S1_MME
		sedCommand := "s:MME_INTERFACE_NAME_FOR_S1_MME.*;:MME_INTERFACE_NAME_FOR_S1_MME               = \"" + outInterface + "\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " failed")
				fmt.Println("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " failed")
			} else {
				OaiObj.Logger.Print("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " successful")
				fmt.Println("Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf + " successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf)
			fmt.Println("Retrying to Set MME_INTERFACE_NAME_FOR_S1_MME to the value " + outInterface + " in " + mmeConf)
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		}

		// MME_IPV4_ADDRESS_FOR_S1_MME
		sedCommand = "s:MME_IPV4_ADDRESS_FOR_S1_MME.*;:MME_IPV4_ADDRESS_FOR_S1_MME          = \"" + outInterfaceIP + "/24\";:g"
		retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)

		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf + " failed")
				fmt.Println("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf + " failed")
			} else {
				OaiObj.Logger.Print("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + " in " + mmeConf + " successful")
				fmt.Println("Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf + " successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf)
			fmt.Println("Retrying to Set MME_IPV4_ADDRESS_FOR_S1_MME to the value " + outInterfaceIP + "/24 in " + mmeConf)
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
		}
		if CnAllInOneMode == false {
			outInterfaceIP := util.GetOutboundIP(OaiObj.Logger)
			outInterface, _ := util.GetInterfaceByIP(OaiObj.Logger, outInterfaceIP)

			sedCommand := "s:MME_INTERFACE_NAME_FOR_S11.*;:MME_INTERFACE_NAME_FOR_S11 = \"" + outInterface + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
			if retStatus.Exit != 0 {
				return errors.New("Set MME_INTERFACE_NAME_FOR_S11 in " + mmeConf + " failed")
			}

			sedCommand = "s:MME_IPV4_ADDRESS_FOR_S11.*;:MME_IPV4_ADDRESS_FOR_S11          = \"" + outInterfaceIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
			if retStatus.Exit != 0 {
				return errors.New("Set MME_IPV4_ADDRESS_FOR_S11 in " + mmeConf + " failed")
			}

			sedCommand = "s:ListenOn.*;:ListenOn = \"" + outInterfaceIP + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeFdConf)
			if retStatus.Exit != 0 {
				return errors.New("Set ListenOn in " + mmeFdConf + " failed")
			}
			/////////////////////////////////////////////////////////
			hssIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.HssDomainName)
			if buildSnap == true {
				hssIP = "127.0.0.10"
			} else {
				for {
					if err != nil {
						OaiObj.Logger.Print(err)
					} else {
						hostNameHss, err := net.LookupHost(hssIP)
						if len(hostNameHss) > 0 {
							break
						} else {
							OaiObj.Logger.Print(err)
						}
					}
					OaiObj.Logger.Print("Valid ip address for oai-hss not yet retreived")
					time.Sleep(1 * time.Second)
					hssIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.HssDomainName)
				}
			}
			//////////////////////////////////////////////////////////////
			spgwcIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.SpgwcDomainName)
			if buildSnap == true {
				spgwcIP = "127.0.11.2"
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
			//////////////////////////////////////////////////////////////

			// replace the ip address of hss
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.10/"+hssIP+"/g", mmeFdConf)
			if retStatus.Exit != 0 {
				return errors.New("Set the ip address of oai-hss in " + mmeFdConf + " failed")
			}

			// SGW_IPV4_ADDRESS_FOR_S11
			// replace SGW_IPV4_ADDRESS_FOR_S11
			sedCommand = "s:SGW_IPV4_ADDRESS_FOR_S11=\"127.0.11.2.*;:SGW_IPV4_ADDRESS_FOR_S11=\"" + spgwcIP + "/24\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, mmeConf)
			if retStatus.Exit != 0 {
				return errors.New("Set the ip address of oai-spgwc SGW_IPV4_ADDRESS_FOR_S11 in " + mmeConf + " failed")
			}

		}

		/////////////////////////////////////////////////////////////////////////////////
		if CnAllInOneMode == true {
			// deuplicated from oai-hss
			cassandraIP, err := util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.CassandraDomainName)
			OaiObj.Logger.Print("cassandraIP " + cassandraIP)
			for {
				if err != nil {
					OaiObj.Logger.Print(err)
				} else {
					hostNameCassandra, err := net.LookupHost(cassandraIP)

					if len(hostNameCassandra) > 0 {
						break
					} else {
						OaiObj.Logger.Print(err)
					}
				}
				OaiObj.Logger.Print("Valid ip address for mysql not yet retreived")
				time.Sleep(1 * time.Second)
				cassandraIP, err = util.GetIPFromDomain(OaiObj.Logger, OaiObj.Conf.CassandraDomainName)
			}
			// deuplicated from oai-hss
			// sudo oai-hss.add-mme -i ubuntu.openair5G.eur -C 172.18.0.2
			OaiObj.Logger.Print("Adding oai-mme to Cassanra DB ")
			retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
			for {
				if retStatus.Exit != 0 {
					OaiObj.Logger.Print("Adding the mme to hss database failed")
					fmt.Println("Adding the mme to hss database failed")
				} else {
					OaiObj.Logger.Print("Adding the mme to hss database was successful")
					fmt.Println("Adding the mme to hss database was successful")
					break
				}
				time.Sleep(1 * time.Second)
				OaiObj.Logger.Print("Retrying to add oai-mme to hss database")
				fmt.Println("Retrying to add oai-mme to hss database")
				retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
			}
			// oai-hss.add-users -I208950000000001-208950000000010 -a oai.ipv4 -C 172.18.0.2
			OaiObj.Logger.Print("Adding users to Cassanra DB ")
			retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-users", "-I", "208950000000001-208950000000010", "-a", "oai.ipv4", "-C", cassandraIP)
			for {
				if retStatus.Exit != 0 {
					OaiObj.Logger.Print("Adding users to hss database failed")
					fmt.Println("Adding users to hss database failed")
					// return errors.New("Adding users to hss database failed")
				} else {
					OaiObj.Logger.Print("Adding users to hss database was successful")
					fmt.Println("Adding users to hss database was successful")
					break
				}
				time.Sleep(1 * time.Second)
				OaiObj.Logger.Print("Retrying to add users to hss database")
				fmt.Println("Retrying to add users to hss database")
				retStatus = util.RunCmd(OaiObj.Logger, "/snap/bin/oai-hss.add-users", "-I", "208950000000001-208950000000010", "-a", "oai.ipv4", "-C", cassandraIP)
			}
		}
		/////////////////////////////
		// oai-mme.start
		time.Sleep(10 * time.Second)
		retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "start"}, "."))
		counter := 0
		maxCounter := 2
		for {
			if len(retStatus.Stderr) == 0 {
				time.Sleep(5 * time.Second)
				counter = counter + 1
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "status"}, "."))
				oairanStatus := strings.Join(retStatus.Stdout, " ")
				checkInactive := strings.Contains(oairanStatus, "inactive")
				if checkInactive != true {
					OaiObj.Logger.Print("Waiting to make sure that oai-mme is working properly")
					fmt.Println("Waiting to make sure that oai-mme is working properly")
					if counter >= maxCounter {
						break
					}
				} else {
					OaiObj.Logger.Print("oai-mme is in inactive status, restarting the service")
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "stop"}, "."))
					retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "start"}, "."))
					counter = 0
				}
			} else {
				OaiObj.Logger.Print("Start oai-mme failed, try again later")
				retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{mmeBin, "start"}, "."))
				counter = 0
			}
		}
	}
	fmt.Println("END of oai-mme configuring and starting")
	OaiObj.Logger.Print("END of oai-mme configuring and starting")
	return nil
}

// RestartMmeV2 : Restart MME as a daemon
func restartMmeV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Restart oai-mme daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.mme-restart")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Restart oai-mme failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-mme is successfully restarted")
	return nil
}

// stopMmeV2 : Stop MME as a daemon
func stopMmeV2(OaiObj Oai) error {
	OaiObj.Logger.Print("Stop oai-mme daemon")
	for {
		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.mme-stop")
		if len(retStatus.Stderr) == 0 {
			break
		}
		OaiObj.Logger.Print("Stop oai-mme failed, try again later")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("oai-mme is successfully stopped")
	return nil
}
