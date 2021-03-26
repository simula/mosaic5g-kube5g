/*
#!/usr/local/go/bin/go
################################################################################
* Copyright 2016-2019 Eurecom and Mosaic5G Platforms Authors
* Licensed to the Mosaic5G under one or more contributor license
* agreements. See the NOTICE file distributed with this
* work for additional information regarding copyright ownership.
* The Mosaic5G licenses this file to You under the
* Apache License, Version 2.0  (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
################################################################################
#-------------------------------------------------------------------------------
# For more information about Mosaic5G:
#                                   admin@mosaic-5g.io
# file          cfg.go
# brief 		define the configuration of the snaps, check the file cmd/test/conf.yaml to see an example of such configuration
# authors:
		- Osama Arouk (arouk@eurecom.fr)
		- Kevin Hsi-Ping Hsu (hsuh@eurecom.fr)
		- Alireza Mohammadi (alireza.mohamamdi@eurecom.fr)
*-------------------------------------------------------------------------------
*/

package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

// MmeStatus parse the status of mme when interacting with it using url
type MmeStatus struct {
	Service string `json:"service" yaml:"service"`
	Startup string `json:"startup" yaml:"startup"`
	Current string `json:"current" yaml:"current"`
	Notes   string `json:"notes" yaml:"notes"`
}

// CnEntityV2Status parse the status of oai-cn entity (i.e., oai-hss, oai-mme, oai-spgwc, and oai-spgwu) when interacting with it using url
type CnEntityV2Status struct {
	Service string `json:"service" yaml:"service"`
	Startup string `json:"startup" yaml:"startup"`
	Current string `json:"current" yaml:"current"`
	Notes   string `json:"notes" yaml:"notes"`
}

// CfgOaiEnb the configuration of oai-enb
type CfgOaiEnb struct {
	// Configurations of ENB
	NodeFunction string `json:"nodeFunction" yaml:"nodeFunction"`
	MCC          string `yaml:"mcc"`
	MNC          string `yaml:"mnc"`
	MmeService   struct {
		Description string `yaml:"description"`
		Name        string `yaml:"name"`
		SnapVersion string `yaml:"snapVersion"`
		IPV4        string `yaml:"ipv4"`
	} `yaml:"mmeService"`
	FlexranService struct {
		Description   string `yaml:"description"`
		Enabled       bool   `yaml:"enabled"`
		Name          string `yaml:"name"`
		InterfaceName string `yaml:"interface_name"`
		IPV4          string `yaml:"ipv4"`
		Port          string `yaml:"port"`
	} `yaml:"flexranService"`
	FlexRAN            bool                 `yaml:"flexRAN"`
	FlexRANServiceName string               `yaml:"flexRANServiceName"`
	Snap               SnapDescriptionFinal `json:"snap" yaml:"snap"`
	EutraBand          struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"eutra_band"`

	NidCellMbsfn struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"Nid_cell_mbsfn"`

	DownlinkFrequency struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"downlink_frequency"`
	UplinkFrequencyOffset struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"uplink_frequency_offset"`
	NumberRbDl struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"N_RB_DL"`
	TxGain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"tx_gain"`
	RxGain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"rx_gain"`
	PuschP0Nominal struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"pusch_p0_Nominal"`
	PucchP0Nominal struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"pucch_p0_Nominal"`
	PdschReferenceSignalPower struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"pdsch_referenceSignalPower"`
	PuSch10xSnr struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"puSch10xSnr"`
	PuCch10xSnr struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"puCch10xSnr"`
	ParallelConfig struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"parallel_config"`
	MaxRxGain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"max_rxgain"`
	Usrp struct {
		N3xx N3xxDescription `yaml:"n3xx"`
	} `json:"usrp" yaml:"usrp"`
	X2Config struct {
		Description     string `yaml:"description"`
		EnableX2        string `yaml:"enable_x2"`
		TRelocPrep      string `yaml:"t_reloc_prep"`
		TX2RelocOverall string `yaml:"tx2_reloc_overall"`
		TDCPrep         string `yaml:"t_dc_prep"`
		TDCOverall      string `yaml:"t_dc_overall"`
	} `yaml:"x2_config"`
}

// CfgOaiGnb  the configuration of oai-gnb
type CfgOaiGnb struct {
	// Configurations of GNB
	NodeFunction            string `json:"nodeFunction" yaml:"nodeFunction"`
	MCC                     string `json:"mcc" yaml:"mcc"`
	MNC                     string `json:"mnc" yaml:"mnc"`
	SsbSubcarrierOffset     string `yaml:"ssb_SubcarrierOffset"`
	PdschAntennaPorts       string `yaml:"pdsch_AntennaPorts"`
	PuschTargetSNRx10       string `yaml:"pusch_TargetSNRx10"`
	PucchTargetSNRx10       string `yaml:"pucch_TargetSNRx10"`
	ServingCellConfigCommon []struct {
		AbsoluteFrequencySSB             string `yaml:"absoluteFrequencySSB"`
		DLFrequencyBand                  string `yaml:"dl_frequencyBand"`
		DLAbsoluteFrequencyPointA        string `yaml:"dl_absoluteFrequencyPointA"`
		DLOffstToCarrier                 string `yaml:"dl_offstToCarrier"`
		DLSubcarrierSpacing              string `yaml:"dl_subcarrierSpacing"`
		DLCarrierBandwidth               string `yaml:"dl_carrierBandwidth"`
		InitialDLBWPlocationAndBandwidth string `yaml:"initialDLBWPlocationAndBandwidth"`
	} `yaml:"servingCellConfigCommon"`
	MmeService struct {
		Description string `yaml:"description"`
		Name        string `yaml:"name"`
		SnapVersion string `yaml:"snapVersion"`
		IPV4        string `yaml:"ipv4"`
		IPV6        string `yaml:"ipv6"`
		Preference  string `yaml:"preference"`
	} `yaml:"mmeService"`
	Snap           SnapDescriptionFinal `json:"snap" yaml:"snap"`
	ParallelConfig struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"parallel_config"`
	MaxRxGain struct {
		Default     string `yaml:"default"`
		Description string `yaml:"description"`
	} `yaml:"max_rxgain"`
	Usrp struct {
		N3xx N3xxDescription `yaml:"n3xx"`
	} `json:"usrp" yaml:"usrp"`
	X2Config struct {
		Description     string `yaml:"description"`
		EnableX2        string `yaml:"enable_x2"`
		TRelocPrep      string `yaml:"t_reloc_prep"`
		TX2RelocOverall string `yaml:"tx2_reloc_overall"`
		TDCPrep         string `yaml:"t_dc_prep"`
		TDCOverall      string `yaml:"t_dc_overall"`
		TargetENBX2IP   struct {
			Name       string `yaml:"name"`
			IPV4       string `yaml:"ipv4"`
			IPV6       string `yaml:"ipv6"`
			Preference string `yaml:"preference"`
		} `yaml:"target_enb_x2_ip_address"`
	} `yaml:"x2_config"`
}

// CfgFlexran the configuration of flexran
type CfgFlexran struct {
	Snap  SnapDescriptionFinal `json:"snap" yaml:"snap"`
	Stats string               `json:"stats" yaml:"stats"`
}

// CfgLlMec the configuration of flexran
type CfgLlMec struct {
	Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
}

// CfgHssV1 the configuration of flexran
type CfgHssV1 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	Snap                SnapDescriptionFinal `json:"snap" yaml:"snap"`
	DatabaseServiceName string               `yaml:"databaseServiceName"`
	MmeServiceName      string               `yaml:"mmeServiceName"`
}

// CfgHssV2 the configuration of flexran
type CfgHssV2 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	ApnNi               GeneralDescription   `json:"APN_NI" yaml:"APN_NI"`
	Imsi                GeneralDescription   `json:"imsi" yaml:"imsi"`
	Users               string               `json:"users" yaml:"users"`
	Snap                SnapDescriptionFinal `json:"snap" yaml:"snap"`
	DatabaseServiceName string               `yaml:"databaseServiceName"`
	HssServiceName      string               `yaml:"hssServiceName"`
	MmeServiceName      string               `yaml:"mmeServiceName"`
	SpgwcServiceName    string               `yaml:"spgwcServiceName"`
}

// CfgMmeV1 the configuration of flexran
type CfgMmeV1 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	Snap            SnapDescriptionFinal `json:"snap" yaml:"snap"`
	MCC             string               `yaml:"mcc"`
	MNC             string               `yaml:"mnc"`
	HssServiceName  string               `yaml:"hssServiceName"`
	SpgwServiceName string               `yaml:"spgwServiceName"`
}

// CfgMmeV2 the configuration of flexran
type CfgMmeV2 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	Snap             SnapDescriptionFinal `json:"snap" yaml:"snap"`
	MCC              string               `yaml:"mcc"`
	MNC              string               `yaml:"mnc"`
	HssServiceName   string               `yaml:"hssServiceName"`
	SpgwcServiceName string               `yaml:"spgwcServiceName"`
}

// CfgSpgwV1 the configuration of flexran
type CfgSpgwV1 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	Snap           SnapDescriptionFinal `json:"snap" yaml:"snap"`
	DNS            string               `yaml:"dns"`
	HssServiceName string               `yaml:"hssServiceName"`
	MmeServiceName string               `yaml:"mmeServiceName"`
}

// CfgSpgwcV2 the configuration of flexran
type CfgSpgwcV2 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	ApnNi GeneralDescription   `json:"APN_NI" yaml:"APN_NI"`
	Snap  SnapDescriptionFinal `json:"snap" yaml:"snap"`
	DNS   string               `yaml:"dns"`
}

// CfgSpgwuV2 the configuration of flexran
type CfgSpgwuV2 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	Snap             SnapDescriptionFinal `json:"snap" yaml:"snap"`
	SpgwcServiceName string               `yaml:"spgwcServiceName"`
}

// CfgHssGlobal CfgGlobal
type CfgHssGlobal struct {
	V1 []CfgHssV1 `json:"v1" yaml:"v1"`
	V2 []CfgHssV2 `json:"v2" yaml:"v2"`
}

// CfgMmeGlobal CfgGlobal
type CfgMmeGlobal struct {
	V1 []CfgMmeV1 `json:"v1" yaml:"v1"`
	V2 []CfgMmeV2 `json:"v2" yaml:"v2"`
}

// CfgSpgwGlobal CfgGlobal
type CfgSpgwGlobal struct {
	V1 []CfgSpgwV1 `json:"v1" yaml:"v1"`
}

// CfgSpgwcGlobal CfgGlobal
type CfgSpgwcGlobal struct {
	V2 []CfgSpgwcV2 `json:"v2" yaml:"v2"`
}

// CfgSpgwuGlobal CfgGlobal
type CfgSpgwuGlobal struct {
	V2 []CfgSpgwuV2 `json:"v2" yaml:"v2"`
}

// CfgCnV1 the configuration of flexran
type CfgCnV1 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
	// OaiCnServiceName string `yaml:"oaiCnServiceName"`
	OaiHss struct {
		DatabaseServiceName string `yaml:"databaseServiceName"`
	} `yaml:"oaiHss"`
	OaiMme struct {
		MCC string `yaml:"mcc"`
		MNC string `yaml:"mnc"`
	} `yaml:"oaiMme"`
	OaiSpgw struct {
		DNS string `yaml:"dns"`
	} `yaml:"oaiSpgw"`
}

// CfgCnV2 the configuration of flexran
type CfgCnV2 struct {
	Realm struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"realm"`
	ApnNi  GeneralDescription `json:"APN_NI" yaml:"APN_NI"`
	Imsi   GeneralDescription `json:"imsi" yaml:"imsi"`
	Users  string             `json:"users" yaml:"users"`
	OaiHss struct {
		Snap                SnapDescriptionFinal `json:"snap" yaml:"snap"`
		DatabaseServiceName string               `yaml:"databaseServiceName"`
	} `yaml:"oaiHss"`
	OaiMme struct {
		Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
		MCC  string               `yaml:"mcc"`
		MNC  string               `yaml:"mnc"`
	} `yaml:"oaiMme"`
	OaiSpgwc struct {
		Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
		DNS  string               `yaml:"dns"`
	} `yaml:"oaiSpgwc"`
	OaiSpgwu struct {
		Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
	} `yaml:"oaiSpgwu"`
}

// CfgCnGlobal CfgGlobal
type CfgCnGlobal struct {
	V1 []CfgCnV1 `json:"v1" yaml:"v1"`
	V2 []CfgCnV2 `json:"v2" yaml:"v2"`
}

// CfgGlobal CfgGlobal
type CfgGlobal struct {
	OaiEnb   []CfgOaiEnb    `json:"oaiEnb" yaml:"oaiEnb"`
	OaiGnb   []CfgOaiGnb    `json:"oaiGnb" yaml:"oaiGnb"`
	Flexran  []CfgFlexran   `json:"flexran" yaml:"flexran"`
	LlMec    []CfgLlMec     `json:"llmec" yaml:"llmec"`
	OaiCn    CfgCnGlobal    `json:"oaiCn" yaml:"oaiCn"`
	OaiHss   CfgHssGlobal   `json:"oaiHss" yaml:"oaiHss"`
	OaiMme   CfgMmeGlobal   `json:"oaiMme" yaml:"oaiMme"`
	OaiSpgw  CfgSpgwGlobal  `json:"oaiSpgw" yaml:"oaiSpgw"`
	OaiSpgwc CfgSpgwcGlobal `json:"oaiSpgwc" yaml:"oaiSpgwc"`
	OaiSpgwu CfgSpgwuGlobal `json:"oaiSpgwu" yaml:"oaiSpgwu"`
}

// SnapDescriptionFinal this type is to descripe the details of the snap to be used for the current application
type SnapDescriptionFinal struct {
	Description string   `json:"description" yaml:"description"`
	Name        string   `json:"name" yaml:"name"`
	Channel     string   `json:"channel" yaml:"channel"`
	Devmode     bool     `json:"devmode" yaml:"devmode"`
	Refresh     bool     `json:"refresh" yaml:"refresh"`
	Plugs       []string `json:"plugs" yaml:"plugs"`
}

// GeneralDescription This is general description for every parameters defined above
type GeneralDescription struct {
	Default     string `json:"default" yaml:"default"`
	Description string `json:"description" yaml:"description"`
}

// N3xxDescription define the parameters of N3XX family of USRP
type N3xxDescription struct {
	Enabled  bool `json:"enabled" yaml:"enabled"`
	SdrAddrs struct {
		Addr       string `yaml:"addr"`
		SecondAddr string `yaml:"second_addr"`
	} `yaml:"sdr_addrs"`
	ClockSrc string `yaml:"clock_src"`
}

// GetConf : read yaml into struct
func (c *CfgGlobal) GetConf(logger *log.Logger, path string) error {
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
func (c *CfgGlobal) ToMap(logger *log.Logger) error {
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
