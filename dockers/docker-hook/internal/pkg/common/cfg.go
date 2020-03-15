package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Cfg stores available settings
type Cfg struct {
	// Configurations of ENB
	MCC                   string `yaml:"mcc"`
	MNC                   string `yaml:"mnc"`
	EutraBand             string `yaml:"eutraBand"`
	DownlinkFrequency     string `yaml:"downlinkFrequency"`
	UplinkFrequencyOffset string `yaml:"uplinkFrequencyOffset"`
	NumberRbDl            string `yaml:"NumberRbDl"`
	ParallelConfig        string `yaml:"ParallelConfig"`
	MaxRxGain             string `yaml:"MaxRxGain"`
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
}

// GetConf : read yaml into struct
func (c *Cfg) GetConf(logger *log.Logger, path string) error {
	//Read yaml here
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Println(err.Error())
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		logger.Println(err.Error())
		return err
	}

	return nil
}

// ToMap converts config to map[string]string in GO
func (c *Cfg) ToMap(logger *log.Logger) error {
	datas := make(map[string]string)
	vn := reflect.ValueOf(c).Elem()
	for i := 0; i < vn.NumField(); i++ {
		if vn.Field(i).Kind().String() == "bool" {
			datas[vn.Type().Field(i).Name] = strconv.FormatBool(vn.Field(i).Interface().(bool))
		} else if vn.Field(i).Kind().String() == "string" {
			datas[vn.Type().Field(i).Name] = vn.Field(i).Interface().(string)
		} else {
			logger.Println("No matched kind for element ", i)
		}
	}
	for k, v := range datas {
		fmt.Println(k, " is ", v)
	}

	return nil
}
