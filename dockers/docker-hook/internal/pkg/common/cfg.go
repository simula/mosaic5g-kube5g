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
// type Cfg struct {
// 	// Configurations of ENB
// 	MCC                   string `yaml:"mcc"`
// 	MNC                   string `yaml:"mnc"`
// 	EutraBand             string `yaml:"eutraBand"`
// 	DownlinkFrequency     string `yaml:"downlinkFrequency"`
// 	UplinkFrequencyOffset string `yaml:"uplinkFrequencyOffset"`
// 	NumberRbDl            string `yaml:"NumberRbDl"`
// 	ParallelConfig        string `yaml:"ParallelConfig"`
// 	MaxRxGain             string `yaml:"MaxRxGain"`
// 	FlexRAN               bool   `yaml:"flexRAN"`
// 	// Global setting
// 	ConfigurationPathofCN  string `yaml:"configurationPathofCN"`
// 	ConfigurationPathofRAN string `yaml:"configurationPathofRAN"`
// 	SnapBinaryPath         string `yaml:"snapBinaryPath"`
// 	DNS                    string `yaml:"dns"`
// 	HssDomainName          string `yaml:"hssDomainName"`
// 	MmeDomainName          string `yaml:"mmeDomainName"`
// 	SpgwDomainName         string `yaml:"spgwDomainName"`
// 	MysqlDomainName        string `yaml:"mysqlDomainName"`
// 	FlexRANDomainName      string `yaml:"flexRANDomainName"`
// 	Test                   bool   `yaml:"test"` //test configuring without changing any file; No snap is installed
// 	// New Configurations of ENB
// 	Snap struct {
// 		Description string `yaml:"description"`
// 		Name        string `yaml:"name"`
// 		Channel     string `yaml:"channel"`
// 		Devmode     bool   `yaml:"devmode"`
// 		Jailmode    bool   `yaml:"jailmode"`
// 		Type        string `yaml:"type"`
// 	} `yaml:"snap"`
// 	NodeFunction struct {
// 		Default     string `yaml:"default"`
// 		Description string `yaml:"description"`
// 		Type        string `yaml:"type"`
// 	} `yaml:"node_function"`
// }

// // CommonCfg: common configurations for oai, like the realm
// type CommonCfg struct {
// 	Realm struct {
// 		Default     string `yaml:"default"`
// 		Description string `yaml:"description"`
// 		Type        string `yaml:"type"`
// 	} `yaml:"realm"`
// 	DNS string `yaml:"dns"`
// }

// Cfg stores available settings
type Cfg struct {
	// New Configurations of ENB

	Realm struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"realm"`

	Snap struct {
		Description string `yaml:"description"`
		Name        string `yaml:"name"`
		Channel     string `yaml:"channel"`
		Devmode     bool   `yaml:"devmode"`
		Jailmode    bool   `yaml:"jailmode"`
		Type        string `yaml:"type"`
	} `yaml:"snap"`
	NodeFunction struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"node_function"`
	MmeIPAddr struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"mme_ip_addr"`
	EutraBand struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"eutra_band"`
	DownlinkFrequency struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"downlink_frequency"`
	UplinkFrequencyOffset struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"uplink_frequency_offset"`
	NumberRbDl struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"N_RB_DL"`
	NbAntennasTx struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"nb_antennas_tx"`
	NbAntennasRx struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"nb_antennas_rx"`
	TxGain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"tx_gain"`
	RxGain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"rx_gain"`
	EnbName struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"enb_name"`
	EnbId struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"enb_id"`
	ParallelConfig struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"parallel_config"`
	MaxRxGain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"max_rxgain"`
	CuPortc struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"cu_portc"`
	DuPortc struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"du_portc"`
	CuPortd struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"cu_portd"`
	DuPortd struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"du_portd"`
	///
	RruPortc struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"rru_portc"`
	RruPortd struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"rru_portd"`
	RccPortc struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"rcc_portc"`
	RccPortd struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"rcc_portd"`
	RccRruTrPreference struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
		Type        string `yaml:"type"`
	} `yaml:"rcc_rru_tr_preference"`
	///
	// Configurations of ENB
	MCC string `yaml:"mcc"`
	MNC string `yaml:"mnc"`
	//EutraBand             string `yaml:"eutraBand"`
	//DownlinkFrequency     string `yaml:"downlinkFrequency"`
	//UplinkFrequencyOffset string `yaml:"uplinkFrequencyOffset"`
	FlexRAN bool `yaml:"flexRAN"`
	// Global setting
	ConfigurationPathofCN  string `yaml:"configurationPathofCN"`
	ConfigurationPathofRAN string `yaml:"configurationPathofRAN"`
	SnapBinaryPath         string `yaml:"snapBinaryPath"`
	DNS                    string `yaml:"dns"`
	CuDomainName           string `yaml:"cuDomainName"`
	DuDomainName           string `yaml:"duDomainName"`
	RccDomainName          string `yaml:"rccDomainName"`
	RruDomainName          string `yaml:"rruDomainName"`
	HssDomainName          string `yaml:"hssDomainName"`
	MmeDomainName          string `yaml:"mmeDomainName"`
	SpgwDomainName         string `yaml:"spgwDomainName"`
	SpgwcDomainName        string `yaml:"spgwcDomainName"`
	SpgwuDomainName        string `yaml:"spgwuDomainName"`
	MysqlDomainName        string `yaml:"mysqlDomainName"`
	CassandraDomainName    string `yaml:"cassandraDomainName"`
	FlexRANDomainName      string `yaml:"flexRANDomainName"`
	Test                   bool   `yaml:"test"` //test configuring without changing any file; No snap is installed
}

// GetConf : read yaml into struct
func (c *Cfg) GetConf(logger *log.Logger, path string) error {
	//Read yaml here
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Panicln(err.Error())
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		logger.Panicln(err.Error())
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
