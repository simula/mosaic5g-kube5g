package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mosaic5g/docker-hook/internal/pkg/util"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type conf struct {
	//////////////////////
	logFile *os.File    // File for log to write something
	Logger  *log.Logger // Collect log
	Snap    struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"snap"`
	Node_function struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"node_function"`
	Target_hardware struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"target_hardware"`
	Mme_ip_addr struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"mme_ip_addr"`
	Eutra_band struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"eutra_band"`
	Downlink_frequency struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"downlink_frequency"`
	Uplink_frequency_offset struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"uplink_frequency_offset"`
	N_RB_DL struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"N_RB_DL"`
	Nb_antennas_tx struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"nb_antennas_tx"`
	Nb_antennas_rx struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"nb_antennas_rx"`
	Tx_gain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"tx_gain"`
	Rx_gain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"rx_gain"`
	Enb_name struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"enb_name"`
	Enb_id struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"enb_id"`
	Parallel_config struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"parallel_config"`
	Max_rxgain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"max_rxgain"`
	////////////////////////
	// Configurations of ENB
	MCC                       string `yaml:"mcc"`
	MNC                       string `yaml:"mnc"`
	EutraBand_old             string `yaml:"eutraBand"`
	DownlinkFrequency_old     string `yaml:"downlinkFrequency"`
	UplinkFrequencyOffset_old string `yaml:"uplinkFrequencyOffset"`
	FlexRAN                   bool   `yaml:"flexRAN"`
	// Global setting
	ConfigurationPathofCN  string `yaml:"configurationPathofCN"`
	ConfigurationPathofRAN string `yaml:"configurationPathofRAN"`
	SnapBinaryPath         string `yaml:"snapBinaryPath"`
	DNS                    string `yaml:"dns"`
	HssDomainName          string `yaml:"hssDomainName"`
	MmeDomainName          string `yaml:"mmeDomainName"`
	SpgwDomainName         string `yaml:"spgwDomainName"`
	MysqlDomainName        string `yaml:"mysqlDomainName"`
	FlexRANDomainName      string `yaml:"flexRANDomainName"`
	Test                   bool   `yaml:"test"` //test configuring without changing any file; No snap is installed
	Bar                    struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"bar"`
}

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("test_conf.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	//fmt.Println("yamlFile=", c)
	//var logg *loggerTest
	newFile, err := os.Create("/home/borer/go/src/dockerhook-test/hook.log")
	if err != nil {
		fmt.Printf("error occured= \t ")
	}
	//fmt.Printf(newFile.Name())
	c.logFile = newFile

	//////////////
	c.Logger = log.New(c.logFile, "[Debug]"+time.Now().Format("2006-01-02 15:04:05")+" ", log.Lshortfile)

	enbConf := "/home/gatto/go/src/mosaic5g/docker-hook/cmd/test/enb.config"
	//enbConf := c.ConfigurationPathofRAN + "enb.band7.tm1.50PRB.usrpb210.conf"
	sedCommand := ""
	// mmeIP := "sdvjnsd"
	// sedCommand = "s:mme_ip_address *= *( *{ *ipv4 *= *\".*\" *;:mme_ip_address      = ( { ipv4       = \"" + mmeIP + "\"" + ";:g"
	// // sedCommand = "175s:\".*;:\"" + mmeIP + "\";:g"
	// util.RunCmd(c.Logger, "sed", "-i", sedCommand, enbConf)

	// sed -n "/NETWORK_INTERFACES/="  mante.eucnc.orig.95.conf
	sedCommand = "/NETWORK_INTERFACES/="
	retStatus := util.RunCmd(c.Logger, "sed", "-n", sedCommand, enbConf)
	fmt.Println("retStatus.Stdout=", retStatus.Stdout[0])

	// lineNumber, err := strconv.ParseUint(retStatus.Stdout[0], 10, 32)

	number, err := strconv.ParseUint(retStatus.Stdout[0], 10, 32)
	lineNumber := strconv.Itoa(int(number - 1))
	sedCommand1 := lineNumber + " a \t\t\t  target_enb_x2_ip_address      = ( {"
	lineNumber = strconv.Itoa(int(number))
	sedCommand2 := lineNumber + " a 	 ipv4       = \"192.168.12.161\";"
	lineNumber = strconv.Itoa(int(number + 1))
	sedCommand3 := lineNumber + " a 	  ipv6       = \"192:168:30::17\";"
	lineNumber = strconv.Itoa(int(number + 2))
	sedCommand4 := lineNumber + " a 	  preference = \"ipv4\";"
	lineNumber = strconv.Itoa(int(number + 3))
	sedCommand5 := lineNumber + " a 	  }"
	lineNumber = strconv.Itoa(int(number + 4))
	sedCommand6 := lineNumber + " a 	  );"

	// target_enb_x2_ip_address      = (
	// 	{ ipv4       = "192.168.12.161";
	// 	  ipv6       = "192:168:30::17";
	// 	  preference = "ipv4";
	// 	}
	//    );

	retStatus = util.RunCmd(c.Logger, "sed", "-i", sedCommand1, enbConf)
	retStatus = util.RunCmd(c.Logger, "sed", "-i", sedCommand2, enbConf)
	retStatus = util.RunCmd(c.Logger, "sed", "-i", sedCommand3, enbConf)
	retStatus = util.RunCmd(c.Logger, "sed", "-i", sedCommand4, enbConf)
	retStatus = util.RunCmd(c.Logger, "sed", "-i", sedCommand5, enbConf)
	retStatus = util.RunCmd(c.Logger, "sed", "-i", sedCommand6, enbConf)
	// sed -i '192 a line\n2.5skdvds' mante.eucnc.orig.95.conf
	fmt.Println("retStatus.Stdout=", retStatus.Stdout)
	fmt.Println("retStatus.Stderr=", retStatus.Stderr)

	// host, err := net.LookupAddr("192.168.12.85")
	// fmt.Println("err=", err)
	// fmt.Println("HOST=", host)

	return c
}

type loggerTest struct {
	logFile *os.File    // File for log to write something
	Logger  *log.Logger // Collect log

}

func main() {
	var c conf
	c.getConf()

	//fmt.Println("Hellow\n sdfs")
	//fmt.Println(c.Bar)
	//fmt.Println("NodeFunction=", c.Node_function.Default)
	//fmt.Println("Description=", c.Bar.Description)
	//fmt.Println("MNC=", c.MCC)

	//v := reflect.ValueOf(c)

	//values := make([]interface{}, v.NumField())

	//for i := 0; i < v.NumField(); i++ {
	//	values[i] = v.Field(i).Interface()
	//fmt.Println("values[", i, "]=", values[i])
	//if i == 16 {
	//fmt.Println("HELLO WORLD")
	//fmt.Printf("Value: %#v \n", c.Bar)
	//}
	//}

}
