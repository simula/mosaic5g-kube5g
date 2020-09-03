package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// type testType struct {
// 	Hello   string `yaml:"Hello"`
// 	EnbID   string `yaml:"eNB_ID"`
// 	EnbName string `yaml:"eNB_name"`
// 	Realm   string `yaml:"realm"`
// }

// CfgRan Configuration of oai-ran
type CfgRan struct {
	// OaianConf testType `yaml:"oai-ran-conf"`
	OaianConf struct {
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

type Salary struct {
	Basic, HRA, TA float64
}

type Employee struct {
	FirstName, LastName, Email string
	Age                        int
	MonthlySalary              []Salary
}

//GetConf d
func (c *CfgRan) GetConf(path string) error {
	//Read yaml here
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		fmt.Println("ERROR=", err)
		return err
	}

	return nil
}

func main() {
	e := Employee{
		FirstName: "Mark",
		LastName:  "Jones",
		Email:     "mark@gmail.com",
		Age:       25,
		MonthlySalary: []Salary{
			Salary{
				Basic: 15000.00,
				HRA:   5000.00,
				TA:    2000.00,
			},
			Salary{
				Basic: 16000.00,
				HRA:   5000.00,
				TA:    2100.00,
			},
			Salary{
				Basic: 17000.00,
				HRA:   5000.00,
				TA:    2200.00,
			},
		},
	}
	fmt.Println(e.FirstName, e.LastName)
	fmt.Println(e.Age)
	fmt.Println(e.Email)
	for i := 0; i < len(e.MonthlySalary); i++ {
		fmt.Println((e.MonthlySalary[i]).Basic)

	}

	// value := "123"
	// number, err := strconv.ParseUint(value, 10, 32)
	// number = number - 1
	// lineNumber := strconv.Itoa(int(number - 1))
	// fmt.Println(number)
	// fmt.Println(lineNumber)
	// fmt.Println(err)

	cnfRan := CfgRan{}
	confPath := "/home/gatto/go/src/mosaic5g/docker-hook/cmd/test/oai-conf-2.yml"
	err := cnfRan.GetConf(confPath)
	fmt.Println(err)
	fmt.Println("cnfRan.ComponentCarriers=", cnfRan.OaianConf.EnbID)
	fmt.Println("cnfRan.Hello=", cnfRan.OaianConf.ComponentCarriers)
	fmt.Println("cnfRan.Hello=", cnfRan.OaianConf.ThreadStruct)
	fmt.Println("cnfRan.Hello=", cnfRan.OaianConf)

}
