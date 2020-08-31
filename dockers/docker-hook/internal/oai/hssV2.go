package oai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mosaic5g/docker-hook/internal/pkg/util"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// startHssV2 : Start HSS as a daemon
func startHssV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting configuring HSS V2")
	OaiObj.Logger.Print("Starting configuration of OAI-HSS V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-hss.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-hss.conf-get"}, "/"))
	s = strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]

	///////////////////
	// Get working path, Hostname
	// the config files of oai-hss are: hss_rel14_fd.conf, hss_rel14.conf, hss_rel14.json and oss.json
	// hssConf := strings.Join([]string{confPath, "hss_rel14.conf"}, "/")
	hssConfJSON := strings.Join([]string{confPath, "hss_rel14.json"}, "/")
	// hssConfOssJson := strings.Join([]string{confPath, "oss.json"}, "/")
	hssFdConf := strings.Join([]string{confPath, "hss_rel14_fd.conf"}, "/")
	hssBin := strings.Join([]string{snapBinaryPath, "oai-hss"}, "/")

	// // Init hss
	fmt.Println("Start Init oai-hss: ", hssBin+".init")
	OaiObj.Logger.Print("Start Init oai-hss: ", hssBin+".init")
	retStatus = util.RunCmd(OaiObj.Logger, hssBin+".init")

	// create new log files
	JSONFile, err := ioutil.ReadFile(hssConfJSON)
	if err != nil {
		fmt.Println("Error while reading the json file config og oai-hss; ", err)
		OaiObj.Logger.Print("Error while reading the json file config og oai-hss; ", err)
	} else {
		fmt.Println("JSONFile:", JSONFile)
		OaiObj.Logger.Print("JSONFile:", JSONFile)

		var DataJSON StructHssRel14

		err = json.Unmarshal(JSONFile, &DataJSON)
		if err != nil {
			fmt.Println("Error while parsing the json file of oai-hss; ", err)
			OaiObj.Logger.Print("Error while parsing the json file of oai-hss; ", err)
		} else {
			DataJSON.Hss.Logname = "/var/log/hss.log"
			DataJSON.Hss.Statlogname = "/var/log/hss_stat.log"
			DataJSON.Hss.Auditlogname = "/var/log/hss_audit.log"

			// create log files
			retStatus = util.RunCmd(OaiObj.Logger, "touch", DataJSON.Hss.Logname)
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Error creating the log file " + DataJSON.Hss.Logname)
				return errors.New("Error creating the log file " + DataJSON.Hss.Logname)
			}
			retStatus = util.RunCmd(OaiObj.Logger, "touch", DataJSON.Hss.Statlogname)
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Error creating the log file " + DataJSON.Hss.Statlogname)
				return errors.New("Error creating the log file " + DataJSON.Hss.Statlogname)
			}
			retStatus = util.RunCmd(OaiObj.Logger, "touch", DataJSON.Hss.Auditlogname)
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Error creating the log file " + DataJSON.Hss.Auditlogname)
				return errors.New("Error creating the log file " + DataJSON.Hss.Auditlogname)
			}
			result, e := json.MarshalIndent(DataJSON, "", " ")
			if e != nil {
				OaiObj.Logger.Print("error", err)
			} else {
				fmt.Println("SUCCESS")
				OaiObj.Logger.Print("Success change log files of oai-hss config")
			}
			_ = ioutil.WriteFile(hssConfJSON, result, 0644)
		}

	}
	// hostname, _ := os.Hostname()
	// Strat configuring oai-hss
	fmt.Print("Configure hss.conf")
	OaiObj.Logger.Print("Configure hss.conf")
	fmt.Print("buildSnap ", buildSnap)
	OaiObj.Logger.Print("buildSnap ", buildSnap)
	// Replace Cassandra address
	if buildSnap != true {
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

		// sudo oai-hss.add-mme -i ubuntu.openair5G.eur -C 172.18.0.2
		OaiObj.Logger.Print("Adding oai-mme to Cassanra DB ")
		retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Adding the mme to hss database failed")
				fmt.Println("Adding the mme to hss database failed")
				// return errors.New("Adding the mme to hss database failed")
			} else {
				OaiObj.Logger.Print("Adding the mme to hss database was successful")
				fmt.Println("Adding the mme to hss database was successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to add oai-mme to hss database")
			fmt.Println("Retrying to add oai-mme to hss database")
			retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
		}

		// oai-hss.add-users -I208950000000001-208950000000010 -a oai.ipv4 -C 172.18.0.2
		OaiObj.Logger.Print("Adding users to Cassanra DB ")
		retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-users", "-I", "208950000000001-208950000000010", "-a", "oai.ipv4", "-C", cassandraIP)
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
			retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-users", "-I", "208950000000001-208950000000010", "-a", "oai.ipv4", "-C", cassandraIP)
		}

		if CnAllInOneMode == false {
			outInterfaceIP := util.GetOutboundIP(OaiObj.Logger)
			// outInterface, _ := util.GetInterfaceByIP(OaiObj.Logger, outInterfaceIP)

			sedCommand := "s:ListenOn.*;:ListenOn = \"" + outInterfaceIP + "\";:g"
			retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", sedCommand, hssFdConf)
			if retStatus.Exit != 0 {
				return errors.New("Set ListenOn in " + hssFdConf + " failed")
			}
		}

		///////////////////////////////////////////////////////////
		// fmt.Println("start waiting for 20 seconds")
		// OaiObj.Logger.Print("start waiting for 20 seconds")
		// time.Sleep(60 * time.Second)
		// fmt.Println("finish waiting for 20 seconds")
		// OaiObj.Logger.Print("finish waiting for 20 seconds")
		// // Run oai-hss with specific cassandra DB
		// // oai-hss -s 172.20.0.2
		// util.RunCmdNonBlocking(OaiObj.Logger, hssBin, "-s", cassandraIP)
		///////////////////////////////////////////////////////////
		// Start oai-hss with default cassandra DB
		// retStatus = util.RunCmd(OaiObj.Logger, hssBin, "-s", cassandraIP)
		// OaiObj.Logger.Print("hssBin = ", hssBin)
		// OaiObj.Logger.Print("cassandraIP = ", cassandraIP)
		// fmt.Println("hssBin = ", hssBin)
		// fmt.Println("cassandraIP = ", cassandraIP)
		// fmt.Println("OaiObj.Conf.CassandraDomainName = ", OaiObj.Conf.CassandraDomainName)
		// fmt.Println("OaiObj.Conf.CassandraDomainName = ", OaiObj.Conf.CassandraDomainName)
		// for {
		// 	if retStatus.Exit != 0 {
		// 		OaiObj.Logger.Print("starting oai-hss with specific cassandra failed")
		// 		fmt.Println("starting oai-hss with specific cassandra failed")
		// 		// return errors.New("starting oai-hss with specific cassandra failed")
		// 	} else {
		// 		OaiObj.Logger.Print("starting oai-hss with specific cassandra successful")
		// 		fmt.Println("starting oai-hss with specific cassandra was successful")
		// 		break
		// 	}
		// 	time.Sleep(1 * time.Second)
		// 	OaiObj.Logger.Print("Retrying to strat oai-hss with specific cassandra")
		// 	fmt.Println("Retrying to start oai-hss with specific cassandra")
		// 	retStatus = util.RunCmd(OaiObj.Logger, hssBin, "-s", cassandraIP)
		// }
		fmt.Println("Finish configuring oai-hss, oai-hss is working with cassandra ", cassandraIP)
		OaiObj.Logger.Print("Finish configuring oai-hss, oai-hss is working with cassandra ", cassandraIP)

	}
	fmt.Println("END of oai-hss configuring and starting")
	OaiObj.Logger.Print("END of oai-hss configuring and starting")
	return nil
}

// startHssV2 : Start HSS as a daemon
func startAndBlockHssV2(OaiObj Oai, CnAllInOneMode bool, buildSnap bool) error {
	fmt.Println("Starting and block configuring HSS V2")
	OaiObj.Logger.Print("Starting and block configuration of OAI-HSS V2")

	retStatus := util.RunCmd(OaiObj.Logger, "which", "oai-hss.status")
	s := strings.Split(retStatus.Stdout[0], "/")
	snapBinaryPath := strings.Join(s[0:len(s)-1], "/")

	retStatus = util.RunCmd(OaiObj.Logger, strings.Join([]string{snapBinaryPath, "oai-hss.conf-get"}, "/"))
	s = strings.Split(retStatus.Stdout[0], "/")
	// confPath := strings.Join(s[0:len(s)-1], "/")
	// confFileName := s[len(s)-1]

	///////////////////
	// Get working path, Hostname
	// the config files of oai-hss are: hss_rel14_fd.conf, hss_rel14.conf, hss_rel14.json and oss.json
	// hssConf := strings.Join([]string{confPath, "hss_rel14.conf"}, "/")
	// hssConfJSON := strings.Join([]string{confPath, "hss_rel14.json"}, "/")
	// hssConfOssJson := strings.Join([]string{confPath, "oss.json"}, "/")
	// hssFdConf := strings.Join([]string{confPath, "hss_rel14_fd.conf"}, "/")
	hssBin := strings.Join([]string{snapBinaryPath, "oai-hss"}, "/")

	// hostname, _ := os.Hostname()
	// Strat configuring oai-hss
	fmt.Print("Configure hss.conf")
	OaiObj.Logger.Print("Configure hss.conf")
	// Replace Cassandra address
	if buildSnap != true {
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

		// sudo oai-hss.add-mme -i ubuntu.openair5G.eur -C 172.18.0.2
		OaiObj.Logger.Print("Adding oai-mme to Cassanra DB ")
		retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
		for {
			if retStatus.Exit != 0 {
				OaiObj.Logger.Print("Adding the mme to hss database failed")
				fmt.Println("Adding the mme to hss database failed")
				// return errors.New("Adding the mme to hss database failed")
			} else {
				OaiObj.Logger.Print("Adding the mme to hss database was successful")
				fmt.Println("Adding the mme to hss database was successful")
				break
			}
			time.Sleep(1 * time.Second)
			OaiObj.Logger.Print("Retrying to add oai-mme to hss database")
			fmt.Println("Retrying to add oai-mme to hss database")
			retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-mme", "-i", "ubuntu.openair5G.eur", "-C", cassandraIP)
		}

		// oai-hss.add-users -I208950000000001-208950000000010 -a oai.ipv4 -C 172.18.0.2
		OaiObj.Logger.Print("Adding users to Cassanra DB ")
		retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-users", "-I", "208950000000001-208950000000010", "-a", "oai.ipv4", "-C", cassandraIP)
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
			retStatus = util.RunCmd(OaiObj.Logger, hssBin+".add-users", "-I", "208950000000001-208950000000010", "-a", "oai.ipv4", "-C", cassandraIP)
		}

		fmt.Println("start waiting for 20 seconds")
		OaiObj.Logger.Print("start waiting for 20 seconds")
		// time.Sleep(60 * time.Second)
		fmt.Println("finish waiting for 20 seconds")
		OaiObj.Logger.Print("finish waiting for 20 seconds")
		// Run oai-hss with specific cassandra DB
		// oai-hss -s 172.20.0.2
		fmt.Println("Start and block: ", hssBin, " -s ", cassandraIP)
		OaiObj.Logger.Print("Start and block: ", hssBin, " -s ", cassandraIP)
		//////////////
		// restart oai-mme
		// /snap/bin/oai-mme.restart
		if CnAllInOneMode == true {
			app := "/snap/bin/oai-mme.restart"

			cmd := exec.Command(app)
			stdout, err := cmd.Output()

			if err != nil {
				fmt.Println(err.Error())
				OaiObj.Logger.Print(err.Error())
			} else {
				fmt.Println(string(stdout))
				OaiObj.Logger.Print(string(stdout))
			}
		}
		// time.Sleep(10 * time.Second)
		//////////////
		app := hssBin
		arg0 := "-s"
		arg1 := cassandraIP

		cmd := exec.Command(app, arg0, arg1)
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			OaiObj.Logger.Print(err.Error())
		} else {
			fmt.Println(string(stdout))
			OaiObj.Logger.Print(string(stdout))
		}
		/////////////////

		// util.RunCmd(OaiObj.Logger, hssBin, "-s", cassandraIP)

	}
	fmt.Println("END of oai-hss configuring and starting")
	OaiObj.Logger.Print("END of oai-hss configuring and starting")
	return nil
}

// configHssV2 : Config oai-hss
func configHssV2(OaiObj Oai) error {
	fmt.Println("hss.go Starting initializing OAI-HSS")
	///////////////////
	//c := OaiObj.Conf
	retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-conf-get")
	s := strings.Split(retStatus.Stdout[0], "/")
	confPath := strings.Join(s[0:len(s)-1], "/")
	snapBinaryPath := "/snap/bin/"
	///////////////////
	// Get working path, Hostname
	hssConf := confPath + "hss.conf"
	hssFdConf := confPath + "hss_fd.conf"
	hssBin := snapBinaryPath + "oai-cn.hss"
	hostname, _ := os.Hostname()
	fmt.Println("hssConf=", hssConf)
	fmt.Println("hssFdConf=", hssFdConf)
	fmt.Println("hssBin=", hssBin)
	fmt.Println("hostname=", hostname)
	// Strat configuring oai-hss
	OaiObj.Logger.Print("Configure hss.conf")
	//Replace MySQL address
	retStatus = util.RunCmd(OaiObj.Logger, "sed", "-i", "s/127.0.0.1/"+OaiObj.Conf.CassandraDomainName+"/g", hssConf)
	fmt.Println("retStatus.Exit=", retStatus.Exit)
	OaiObj.Logger.Print("retStatus.Exit=", retStatus.Exit)
	if retStatus.Exit != 0 {
		OaiObj.Logger.Print("Set mysql IP in " + hssConf + " failed")
		fmt.Println("Set mysql IP in " + hssConf + " failed")
		return errors.New("Set mysql IP in " + hssConf + " failed")
	}

	// oai-cn.hss-start
	fmt.Println("start hss as daemon")
	OaiObj.Logger.Print("start hss as daemon")
	util.RunCmd(OaiObj.Logger, hssBin+".start")
	return nil
}

// // StartHss : Start HSS as a daemon
// func startHss(OaiObj Oai) error {
// 	///////////////
// 	OaiObj.Logger.Print("Start oai-hss daemon")
// 	for {
// 		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-start")
// 		if len(retStatus.Stderr) == 0 {
// 			break
// 		}
// 		OaiObj.Logger.Print("Start oai-hss failed, try again later")
// 		time.Sleep(1 * time.Second)
// 	}
// 	fmt.Println("oai-hss is successfully started")
// 	return nil
// }

// // RestartHss : Restart HSS as a daemon
// func restartHss(OaiObj Oai) error {
// 	///////////////
// 	OaiObj.Logger.Print("Retart oai-hss daemon")
// 	for {
// 		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-restart")
// 		if len(retStatus.Stderr) == 0 {
// 			break
// 		}
// 		OaiObj.Logger.Print("Restart oai-hss failed, try again later")
// 		time.Sleep(1 * time.Second)
// 	}
// 	fmt.Println("oai-hss is successfully restarted")
// 	return nil
// }

// // stopHss : Stop HSS as a daemon
// func stopHss(OaiObj Oai) error {
// 	///////////////
// 	OaiObj.Logger.Print("Stop oai-hss daemon")
// 	for {
// 		retStatus := util.RunCmd(OaiObj.Logger, "/snap/bin/oai-cn.hss-stop")
// 		if len(retStatus.Stderr) == 0 {
// 			break
// 		}
// 		OaiObj.Logger.Print("Stop oai-hss failed, try again later")
// 		time.Sleep(1 * time.Second)
// 	}
// 	fmt.Println("oai-hss is successfully stopped")
// 	return nil
// }

//StructHssRel14 :
type StructHssRel14 struct {
	Common struct {
		Fdcfg       string `json:"fdcfg"`
		Originhost  string `json:"originhost"`
		Originrealm string `json:"originrealm"`
	} `json:"common"`
	Hss struct {
		Gtwhost             string `json:"gtwhost"`
		Gtwport             int    `json:"gtwport"`
		Restport            int    `json:"restport"`
		Ossport             int    `json:"ossport"`
		Casssrv             string `json:"casssrv"`
		Cassusr             string `json:"cassusr"`
		Casspwd             string `json:"casspwd"`
		Cassdb              string `json:"cassdb"`
		Casscoreconnections int    `json:"casscoreconnections"`
		Cassmaxconnections  int    `json:"cassmaxconnections"`
		Cassioqueuesize     int    `json:"cassioqueuesize"`
		Cassiothreads       int    `json:"cassiothreads"`
		Randv               bool   `json:"randv"`
		Optkey              string `json:"optkey"`
		Reloadkey           bool   `json:"reloadkey"`
		Roamallow           bool   `json:"roamallow"`
		Logsize             int    `json:"logsize"`
		Lognumber           int    `json:"lognumber"`
		Logname             string `json:"logname"`
		Logqsize            int    `json:"logqsize"`
		Statlogsize         int    `json:"statlogsize"`
		Statlognumber       int    `json:"statlognumber"`
		Statlogname         string `json:"statlogname"`
		Auditlogsize        int    `json:"auditlogsize"`
		Auditlognumber      int    `json:"auditlognumber"`
		Auditlogname        string `json:"auditlogname"`
		Statfreq            int    `json:"statfreq"`
		Numworkers          int    `json:"numworkers"`
		Concurrent          int    `json:"concurrent"`
		Ossfile             string `json:"ossfile"`
	} `json:"hss"`
}
