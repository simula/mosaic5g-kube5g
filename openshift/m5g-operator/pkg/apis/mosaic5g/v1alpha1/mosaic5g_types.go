package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Mosaic5gSpec defines the desired state of Mosaic5g
// +k8s:openapi-gen=true
type Mosaic5gSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	// Size         int32  `json:"size" yaml:"size"`
	CoreNetworkAllInOne bool   `json:"coreNetworkAllInOne" yaml:"coreNetworkAllInOne"`
	MysqlSize           int32  `json:"mysqlSize" yaml:"mysqlSize"`
	OaiCnSize           int32  `json:"oaicnsize" yaml:"oaicnsize"`
	OaiHssSize          int32  `json:"oaihsssize" yaml:"oaihsssize"`
	OaiMmeSize          int32  `json:"oaimmesize" yaml:"oaimmesize"`
	OaiSpgwSize         int32  `json:"oaispgwsize" yaml:"oaispgwsize"`
	OaiRanSize          int32  `json:"oairansize" yaml:"oairansize"`
	MysqlImage          string `json:"mysqlImage" yaml:"mysqlImage"`
	CNImage             string `json:"cnImage" yaml:"cnImage"`
	OaiHssImage         string `json:"oaiHssImage" yaml:"oaiHssImage"`
	OaiMmeImage         string `json:"oaiMmeImage" yaml:"oaiMmeImage"`
	OaiSpgwImage        string `json:"oaiSpgwImage" yaml:"oaiSpgwImage"`
	RANImage            string `json:"ranImage" yaml:"ranImage"`

	MCC string `json:"mcc" yaml:"mcc"`
	MNC string `json:"mnc" yaml:"mnc"`
	// EutraBand              string `json:"eutraBand" yaml:"eutraBand"`
	// DownlinkFrequency      string `json:"downlinkFrequency" yaml:"downlinkFrequency"`
	// UplinkFrequencyOffset  string `json:"uplinkFrequencyOffset" yaml:"uplinkFrequencyOffset"`
	FlexRAN                bool   `json:"flexRAN" yaml:"flexRAN"`
	ConfigurationPathofCN  string `json:"configurationPathofCN" yaml:"configurationPathofCN"`
	ConfigurationPathofRAN string `json:"configurationPathofRAN" yaml:"configurationPathofRAN"`
	SnapBinaryPath         string `json:"snapBinaryPath" yaml:"snapBinaryPath"`
	DNS                    string `json:"dns" yaml:"dns"`
	HssDomainName          string `json:"hssDomainName" yaml:"hssDomainName"`
	MmeDomainName          string `json:"mmeDomainName" yaml:"mmeDomainName"`
	SpgwDomainName         string `json:"spgwDomainName" yaml:"spgwDomainName"`
	MysqlDomainName        string `json:"mysqlDomainName" yaml:"mysqlDomainName"`
	FlexRANDomainName      string `json:"flexRANDomainName" yaml:"flexRANDomainName"`
	// New Configurations of ENB
	Snap         SnapDescription    `json:"snap" yaml:"snap"`
	NodeFunction GeneralDescription `json:"node_function" yaml:"node_function"`
	MmeIpAddr    GeneralDescription `json:"mme_ip_addr" yaml:"mme_ip_addr"`

	EutraBand             GeneralDescription `json:"eutra_band" yaml:"eutra_band"`
	DownlinkFrequency     GeneralDescription `json:"downlink_frequency" yaml:"downlink_frequency"`
	UplinkFrequencyOffset GeneralDescription `json:"uplink_frequency_offset" yaml:"uplink_frequency_offset"`
	NumberRbDl            GeneralDescription `json:"N_RB_DL" yaml:"N_RB_DL"`
	NbAntennasTx          GeneralDescription `json:"nb_antennas_tx" yaml:"nb_antennas_tx"`
	NbAntennasRx          GeneralDescription `json:"nb_antennas_rx" yaml:"nb_antennas_rx"`
	TxGain                GeneralDescription `json:"tx_gain" yaml:"tx_gain"`
	RxGain                GeneralDescription `json:"rx_gain" yaml:"rx_gain"`
	EnbName               GeneralDescription `json:"enb_name" yaml:"enb_name"`
	EnbId                 GeneralDescription `json:"enb_id" yaml:"enb_id"`
	ParallelConfig        GeneralDescription `json:"parallel_config" yaml:"parallel_config"`

	MaxRxGain GeneralDescription `json:"max_rxgain" yaml:"max_rxgain"`

	CuPortc            GeneralDescription `json:"cu_portc" yaml:"cu_portc"`
	DuPortc            GeneralDescription `json:"du_portc" yaml:"du_portc"`
	CuPortd            GeneralDescription `json:"cu_portd" yaml:"cu_portd"`
	DuPortd            GeneralDescription `json:"du_portd" yaml:"du_portd"`
	RruPortc           GeneralDescription `json:"rru_portc" yaml:"rru_portc"`
	RruPortd           GeneralDescription `json:"rru_portd" yaml:"rru_portd"`
	RccPortc           GeneralDescription `json:"rcc_portc" yaml:"rcc_portc"`
	RccPortd           GeneralDescription `json:"rcc_portd" yaml:"rcc_portd"`
	RccRruTrPreference GeneralDescription `json:"rcc_rru_tr_preference" yaml:"rcc_rru_tr_preference"`
	CuDomainName       string             `json:"cuDomainName" yaml:"cuDomainName"`
	DuDomainName       string             `json:"duDomainName" yaml:"duDomainName"`
	RccDomainName      string             `json:"rccDomainName" yaml:"rccDomainName"`
	RruDomainName      string             `json:"rruDomainName" yaml:"rruDomainName"`

	// HssDomainName          string `yaml:"hssDomainName"`
	// MmeDomainName          string `yaml:"mmeDomainName"`
	// SpgwDomainName         string `yaml:"spgwDomainName"`
	// MysqlDomainName        string `yaml:"mysqlDomainName"`
	//FlexRANDomainName      string `yaml:"flexRANDomainName"`
}

// GeneralDescription This is general description for every parameters defined above
type GeneralDescription struct {
	Default     string `json:"default" yaml:"default"`
	Description string `json:"description" yaml:"description"`
	// Type        string `json:"type" yaml:"type"`
}

// SnapDescription this type is to descripe the details of the snap to be used for the current application
type SnapDescription struct {
	Description string `json:"description" yaml:"description"`
	Name        string `json:"name" yaml:"name"`
	Channel     string `json:"channel" yaml:"channel"`
	Devmode     bool   `json:"devmode" yaml:"devmode"`
	Jailmode    bool   `json:"jailmode" yaml:"jailmode"`
	// Type        string `json:"type" yaml:"type"`
}

// Mosaic5gStatus defines the observed state of Mosaic5g
// +k8s:openapi-gen=true
type Mosaic5gStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mosaic5g is the Schema for the mosaic5gs API
// +k8s:openapi-gen=true
type Mosaic5g struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Mosaic5gSpec   `json:"spec,omitempty"`
	Status Mosaic5gStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mosaic5gList contains a list of Mosaic5g
type Mosaic5gList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mosaic5g `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mosaic5g{}, &Mosaic5gList{})
}
