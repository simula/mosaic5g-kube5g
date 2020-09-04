package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

// CfgRan Configuration of oai-ran
type CfgRan struct {
	OaianConf struct {
		Snap struct {
			Name     string `yaml:"name"`
			Channel  string `yaml:"channel"`
			Devmode  bool   `yaml:"devmode"`
			Jailmode bool   `yaml:"jailmode"`
			Refresh  bool   `yaml:"refresh"`
		} `yaml:"snap"`
		EnbID             string `yaml:"eNB_ID"`
		EnbName           string `yaml:"eNB_name"`
		Realm             string `yaml:"realm"`
		Mcc               []uint `yaml:"mcc"`
		Mnc               []uint `yaml:"mnc"`
		ComponentCarriers struct {
			NodeFunction          string `yaml:"node_function"`
			EutraBand             string `yaml:"eutra_band"`
			DownlinkFrequency     string `yaml:"downlink_frequency"`
			UplinkFrequencyOffset string `yaml:"uplink_frequency_offset"`
			NidCell               int    `yaml:"Nid_cell"`
			NRbDl                 int    `yaml:"N_RB_DL"`
		} `yaml:"component_carriers"`
		MmeIPAddress struct {
			MmeDomainName string `yaml:"mmeDomainName"`
			Ipv4          string `yaml:"ipv4"`
			Ipv6          string `yaml:"ipv6"`
			Active        string `yaml:"active"`
			Preference    string `yaml:"preference"`
		} `yaml:"mme_ip_address"`
		EnableMeasurementReports string `yaml:"enable_measurement_reports"`
		X2Ho                     struct {
			EnableX2             string           `yaml:"enable_x2"`
			MasterNode           bool             `yaml:"master_node"`
			TargetEnbX2IPAddress []listMasterEnbs `yaml:"target_enb_x2_ip_address"`
		} `yaml:"x2_ho"`
		NetworkInterfaces struct {
			EnbInterfaceNameForS1Mme string `yaml:"ENB_INTERFACE_NAME_FOR_S1_MME"`
			EnbIPv4AddressForS1Mme   string `yaml:"ENB_IPV4_ADDRESS_FOR_S1_MME"`
			EnbInterfaceNameforS1U   string `yaml:"ENB_INTERFACE_NAME_FOR_S1U"`
			EnbIPv4AddressForS1U     string `yaml:"ENB_IPV4_ADDRESS_FOR_S1U"`
			EnbPortForS1U            uint   `yaml:"ENB_PORT_FOR_S1U"`
			EnbIPv4AddressForS1X2C   string `yaml:"ENB_IPV4_ADDRESS_FOR_X2C"`
			EnbPortForS1UX2C         uint   `yaml:"ENB_PORT_FOR_X2C"`
		} `yaml:"NETWORK_INTERFACES"`
		Rus struct {
			MaxRxGain                    uint `yaml:"max_rxgain"`
			MaxPdschReferenceSignalPower int  `yaml:"max_pdschReferenceSignalPower"`
			// Bands                        []uint `yaml:"bands"`
		} `yaml:"RUs"`
		NetworkController struct {
			FlexranEnabled       string `yaml:"FLEXRAN_ENABLED"`
			FlexRANDomainName    string `yaml:"flexRANDomainName"`
			FlexRANInterfaceName string `yaml:"FLEXRAN_INTERFACE_NAME"`
			FlexRANIPv4Address   string `yaml:"FLEXRAN_IPV4_ADDRESS"`
			FlexRANPort          uint   `yaml:"FLEXRAN_PORT"`
			FlexRANCache         string `yaml:"FLEXRAN_CACHE"`
			FlexRANAwaitReconf   string `yaml:"FLEXRAN_AWAIT_RECONF"`
		} `yaml:"NETWORK_CONTROLLER"`
		ThreadStruct struct {
			ParallelConfig string `yaml:"parallel_config"`
			WorkerConfig   string `yaml:"worker_config"`
		} `yaml:"THREAD_STRUCT"`
	} `yaml:"oai-ran-conf"`
}

// listMasterEnbs List of all enbs that will be connected (via X2) to the current enb
type listMasterEnbs struct {
	RanDomainName string `yaml:"ranDomainName"`
	Ipv4          string `yaml:"ipv4"`
	Ipv6          string `yaml:"ipv6"`
	Preference    string `yaml:"preference"`
}

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
	EnbID struct {
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
func (c *CfgRan) GetConf(logger *log.Logger, path string) error {
	//Read yaml here
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Panicln(err.Error())
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		fmt.Println("ERROR=", err)
		logger.Panicln(err.Error())
		return err
	}

	return nil
}

// func getConfTemp(logger *log.Logger, path string, cfg interface{}) error {
// 	//Read yaml here
// 	yamlFile, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		logger.Panicln(err.Error())
// 		return err
// 	}

// 	err = yaml.Unmarshal(yamlFile, cfg)

// 	if err != nil {
// 		logger.Panicln(err.Error())
// 		return err
// 	}
// 	if true {
// 		fmt.Println("path", path)
// 		fmt.Println("me.ConfOaiRan", cfg)
// 		panic("test panic")
// 	}

// 	return nil
// }

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
