package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mosaic5g/docker-hook/internal/pkg/common"
	"mosaic5g/docker-hook/internal/pkg/util"
	"os"

	"gopkg.in/yaml.v2"
)

type conf struct {
	// Configurations of ENB
	MCC                   string `yaml:"mcc"`
	MNC                   string `yaml:"mnc"`
	EutraBand             string `yaml:"eutraBand"`
	DownlinkFrequency     string `yaml:"downlinkFrequency"`
	UplinkFrequencyOffset string `yaml:"uplinkFrequencyOffset"`
	FlexRAN               bool   `yaml:"flexRAN"`
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
	logFile *os.File    // File for log to write something
	Logger  *log.Logger // Collect log
	Conf    *common.Cfg // config files
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Println("yamlFile=", c)

	return c
}

func main() {
	var c conf
	c.getConf()

	fmt.Println(c.Bar)
	fmt.Println(c.Bar.Default)

	// v := reflect.ValueOf(c)

	// values := make([]interface{}, v.NumField())

	// for i := 0; i < v.NumField(); i++ {
	// 	values[i] = v.Field(i).Interface()
	// 	fmt.Println("values[", i, "]=", values[i])
	// 	if i == 16 {
	// 		fmt.Println("HELLO WORLD")
	// 		fmt.Printf("Value: %#v \n", c.Bar)
	// 	}
	// }

	retStatus := util.RunCmd(c.Logger, "which", "oai-ran.enb-status")
	if len(retStatus.Stderr) == 0 {
		fmt.Println("retStatus.Stderr=", retStatus.Stderr)
		fmt.Println("retStatus.Stdout[0]=", retStatus.Stdout[0])
	} else {
		fmt.Println("retStatus.Stdout[0]=", retStatus.Stdout[0])
	}

}
