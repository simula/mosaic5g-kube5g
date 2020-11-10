package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// CfgOaiEnb the configuration of oai-enb
type CfgOaiEnb struct {
	// Configurations of ENB
	K8sDeploymentName     string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName        string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace    string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources       k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector      []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter       []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiEnbSize            int32                         `json:"oaiEnbSize" yaml:"oaiEnbSize"`
	OaiEnbImage           string                        `json:"oaiEnbImage" yaml:"oaiEnbImage"`
	MCC                   string                        `json:"mcc" yaml:"mcc"`
	MNC                   string                        `json:"mnc" yaml:"mnc"`
	MmeService            MmeServiceDescription         `json:"mmeService" yaml:"mmeService"`
	FlexRAN               bool                          `json:"flexRAN" yaml:"flexRAN"`
	FlexranServiceName    string                        `json:"flexranServiceName" yaml:"flexranServiceName"`
	Snap                  SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	EutraBand             GeneralDescription            `json:"eutra_band" yaml:"eutra_band"`
	DownlinkFrequency     GeneralDescription            `json:"downlink_frequency" yaml:"downlink_frequency"`
	UplinkFrequencyOffset GeneralDescription            `json:"uplink_frequency_offset" yaml:"uplink_frequency_offset"`
	NumberRbDl            GeneralDescription            `json:"N_RB_DL" yaml:"N_RB_DL"`
	ParallelConfig        GeneralDescription            `json:"parallel_config" yaml:"parallel_config"`
	MaxRxGain             GeneralDescription            `json:"max_rxgain" yaml:"max_rxgain"`
}

// CfgFlexran the configuration of flexran
type CfgFlexran struct {
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	FlexranSize        int32                         `json:"flexranSize" yaml:"flexranSize"`
	FlexranImage       string                        `json:"flexranImage" yaml:"flexranImage"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
}

// CfgDatabase the configuration of database
type CfgDatabase struct {
	DatabaseType       string                        `json:"databaseType" yaml:"databaseType"`
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`

	DatabaseSize  int32  `json:"databaseSize" yaml:"databaseSize"`
	DatabaseImage string `json:"databaseImage" yaml:"databaseImage"`
}

// CfgLlMec the configuration of flexran
type CfgLlMec struct {
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	LlmecSize          int32                         `json:"llmecSize" yaml:"llmecSize"`
	LlmecImage         string                        `json:"llmecImage" yaml:"llmecImage"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
}

// CfgHssV1 the configuration of flexran
type CfgHssV1 struct {
	K8sDeploymentName   string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName      string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace  string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources     k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector    []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter     []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiHssSize          int32                         `json:"oaiHssSize" yaml:"oaiHssSize"`
	OaiHssImage         string                        `json:"oaiHssImage" yaml:"oaiHssImage"`
	Realm               GeneralDescription            `json:"realm" yaml:"realm"`
	Snap                SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	DatabaseServiceName string                        `json:"databaseServiceName" yaml:"databaseServiceName"`
	MmeServiceName      string                        `json:"mmeServiceName" yaml:"mmeServiceName"`
}

// CfgHssV2 the configuration of flexran
type CfgHssV2 struct {
	K8sDeploymentName   string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName      string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace  string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources     k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector    []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter     []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiHssSize          int32                         `json:"oaiHssSize" yaml:"oaiHssSize"`
	OaiHssImage         string                        `json:"oaiHssImage" yaml:"oaiHssImage"`
	Realm               GeneralDescription            `json:"realm" yaml:"realm"`
	Snap                SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	DatabaseServiceName string                        `json:"databaseServiceName" yaml:"databaseServiceName"`
	HssServiceName      string                        `json:"hssServiceName" yaml:"hssServiceName"`
	MmeServiceName      string                        `json:"mmeServiceName" yaml:"mmeServiceName"`
	SpgwcServiceName    string                        `json:"spgwcServiceName" yaml:"spgwcServiceName"`
}

// CfgMmeV1 the configuration of flexran
type CfgMmeV1 struct {
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiMmeSize         int32                         `json:"oaiMmeSize" yaml:"oaiMmeSize"`
	OaiMmeImage        string                        `json:"oaiMmeImage" yaml:"oaiMmeImage"`
	Realm              GeneralDescription            `json:"realm" yaml:"realm"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	MCC                string                        `json:"mcc" yaml:"mcc"`
	MNC                string                        `json:"mnc" yaml:"mnc"`
	HssServiceName     string                        `json:"hssServiceName" yaml:"hssServiceName"`
	SpgwServiceName    string                        `json:"spgwServiceName" yaml:"spgwServiceName"`
}

// CfgMmeV2 the configuration of flexran
type CfgMmeV2 struct {
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiMmeSize         int32                         `json:"oaiMmeSize" yaml:"oaiMmeSize"`
	OaiMmeImage        string                        `json:"oaiMmeImage" yaml:"oaiMmeImage"`
	Realm              GeneralDescription            `json:"realm" yaml:"realm"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	MCC                string                        `json:"mcc" yaml:"mcc"`
	MNC                string                        `json:"mnc" yaml:"mnc"`
	HssServiceName     string                        `json:"hssServiceName" yaml:"hssServiceName"`
	SpgwcServiceName   string                        `json:"spgwcServiceName" yaml:"spgwcServiceName"`
}

// CfgSpgwV1 the configuration of flexran
type CfgSpgwV1 struct {
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiSpgwSize        int32                         `json:"oaiSpgwSize" yaml:"oaiSpgwSize"`
	OaiSpgwImage       string                        `json:"oaiSpgwImage" yaml:"oaiSpgwImage"`
	Realm              GeneralDescription            `json:"realm" yaml:"realm"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	DNS                string                        `json:"dns" yaml:"dns"`
	HssServiceName     string                        `json:"hssServiceName" yaml:"hssServiceName"`
	MmeServiceName     string                        `json:"mmeServiceName" yaml:"mmeServiceName"`
}

// CfgSpgwcV2 the configuration of flexran
type CfgSpgwcV2 struct {
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiSpgwcSize       int32                         `json:"oaiSpgwcSize" yaml:"oaiSpgwcSize"`
	OaiSpgwcImage      string                        `json:"oaiSpgwcImage" yaml:"oaiSpgwcImage"`
	Realm              GeneralDescription            `json:"realm" yaml:"realm"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	DNS                string                        `json:"dns" yaml:"dns"`
}

// CfgSpgwuV2 the configuration of flexran
type CfgSpgwuV2 struct {
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	OaiSpgwuSize       int32                         `json:"oaiSpgwuSize" yaml:"oaiSpgwuSize"`
	OaiSpgwuImage      string                        `json:"oaiSpgwuImage" yaml:"oaiSpgwuImage"`
	Realm              GeneralDescription            `json:"realm" yaml:"realm"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	SpgwcServiceName   string                        `json:"spgwcServiceName" yaml:"spgwcServiceName"`
}

// CfgHssGlobal CfgHssGlobal
type CfgHssGlobal struct {
	V1 []CfgHssV1 `json:"v1" yaml:"v1"`
	V2 []CfgHssV2 `json:"v2" yaml:"v2"`
}

// CfgMmeGlobal CfgMmeGlobal
type CfgMmeGlobal struct {
	V1 []CfgMmeV1 `json:"v1" yaml:"v1"`
	V2 []CfgMmeV2 `json:"v2" yaml:"v2"`
}

// CfgSpgwGlobal CfgSpgwGlobal
type CfgSpgwGlobal struct {
	V1 []CfgSpgwV1 `json:"v1" yaml:"v1"`
}

// CfgSpgwcGlobal CfgSpgwcGlobal
type CfgSpgwcGlobal struct {
	V2 []CfgSpgwcV2 `json:"v2" yaml:"v2"`
}

// CfgSpgwuGlobal CfgSpgwuGlobal
type CfgSpgwuGlobal struct {
	V2 []CfgSpgwuV2 `json:"v2" yaml:"v2"`
}

// CfgCnV1 the configuration of flexran
type CfgCnV1 struct {
	OaiCnSize          int32                         `json:"oaiCnSize" yaml:"oaiCnSize"`
	OaiCnImage         string                        `json:"oaiCnImage" yaml:"oaiCnImage"`
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	Realm              GeneralDescription            `json:"realm" yaml:"realm"`
	Snap               SnapDescriptionFinal          `json:"snap" yaml:"snap"`
	OaiHss             CnV1OaiHssDescription         `json:"oaiHss" yaml:"oaiHss"`
	OaiMme             CnV1OaiMmeDescription         `json:"oaiMme" yaml:"oaiMme"`
	OaiSpgw            CnV1OaiSpgwDescription        `json:"oaiSpgw" yaml:"oaiSpgw"`
}

// CfgCnV2 the configuration of flexran
type CfgCnV2 struct {
	OaiCnSize          int32                         `json:"oaiCnSize" yaml:"oaiCnSize"`
	OaiCnImage         string                        `json:"oaiCnImage" yaml:"oaiCnImage"`
	K8sDeploymentName  string                        `json:"k8sDeploymentName" yaml:"k8sDeploymentName"`
	K8sServiceName     string                        `json:"k8sServiceName" yaml:"k8sServiceName"`
	K8sEntityNamespace string                        `json:"k8sEntityNamespace" yaml:"k8sEntityNamespace"`
	K8sPodResources    k8sPodResourcesDescription    `json:"k8sPodResources" yaml:"k8sPodResources"`
	K8sLabelSelector   []K8sLabelSelectorDescription `json:"k8sLabelSelector" yaml:"k8sLabelSelector"`
	K8sNodeSelecter    []K8sNodeSelecterDescription  `json:"k8sNodeSelecter" yaml:"k8sNodeSelecter"`
	Realm              GeneralDescription            `json:"realm" yaml:"realm"`
	OaiHss             CnV2OaiHssDescription         `json:"oaiHss" yaml:"oaiHss"`
	OaiMme             CnV2OaiMmeDescription         `json:"oaiMme" yaml:"oaiMme"`
	OaiSpgwc           CnV2OaiSpgwcDescription       `json:"oaiSpgwc" yaml:"oaiSpgwc"`
	OaiSpgwu           CnV2OaiSpgwuDescription       `json:"oaiSpgwu" yaml:"oaiSpgwu"`
}

// CfgCnGlobal CfgCnGlobal
type CfgCnGlobal struct {
	V1 []CfgCnV1 `json:"v1" yaml:"v1"`
	V2 []CfgCnV2 `json:"v2" yaml:"v2"`
}

// K8sNodeSelecterDescription list of lables accepted for K8S
type K8sNodeSelecterDescription struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
	// Usrp  bool   `json:"usrp" yaml:"usrp"`
}

// K8sLabelSelectorDescription list of lables accepted for K8S
type K8sLabelSelectorDescription struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
	// App   string `json:"app" yaml:"app"`
}

// k8sPodResourcesDescription request and limits of the resources for the pods
type k8sPodResourcesDescription struct {
	Limits   ResourcesDescription `json:"limits" yaml:"limits"`
	Requests ResourcesDescription `json:"requests" yaml:"requests"`
}

// ResourcesDescription Define the limits of the resources
type ResourcesDescription struct {
	ResourceCPU    string `json:"resourceCPU" yaml:"resourceCPU"`
	ResourceMemory string `json:"resourceMemory" yaml:"resourceMemory"`
}

// Mosaic5gSpec defines the desired state of Mosaic5g
// +k8s:openapi-gen=true
//Mosaic5gSpec Mosaic5gSpec
type Mosaic5gSpec struct {
	K8sGlobalNamespace string         `json:"k8sGlobalNamespace" yaml:"k8sGlobalNamespace"`
	OaiEnb             []CfgOaiEnb    `json:"oaiEnb" yaml:"oaiEnb"`
	Flexran            []CfgFlexran   `json:"flexran" yaml:"flexran"`
	LlMec              []CfgLlMec     `json:"llmec" yaml:"llmec"`
	Database           []CfgDatabase  `json:"database" yaml:"database"`
	OaiCn              CfgCnGlobal    `json:"oaiCn" yaml:"oaiCn"`
	OaiHss             CfgHssGlobal   `json:"oaiHss" yaml:"oaiHss"`
	OaiMme             CfgMmeGlobal   `json:"oaiMme" yaml:"oaiMme"`
	OaiSpgw            CfgSpgwGlobal  `json:"oaiSpgw" yaml:"oaiSpgw"`
	OaiSpgwc           CfgSpgwcGlobal `json:"oaiSpgwc" yaml:"oaiSpgwc"`
	OaiSpgwu           CfgSpgwuGlobal `json:"oaiSpgwu" yaml:"oaiSpgwu"`
}

// GlobalConf global configuration
type GlobalConf struct {
	ConfYaml Mosaic5gSpec `json:"conf.yaml" yaml:"conf.yaml"`
}

// CfgGlobal defines the desired state of Mosaic5g
// +k8s:openapi-gen=true
//CfgGlobal CfgGlobal
type CfgGlobal struct {
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

	// HssDomainName          string `json:"mcc" yaml:"hssDomainName"`
	// MmeDomainName          string `json:"mcc" yaml:"mmeDomainName"`
	// SpgwDomainName         string `json:"mcc" yaml:"spgwDomainName"`
	// MysqlDomainName        string `json:"mcc" yaml:"mysqlDomainName"`
	//FlexRANDomainName      string `json:"mcc" yaml:"flexRANDomainName"`
}

// GeneralDescription This is general description for every parameters defined above
type GeneralDescription struct {
	Default     string `json:"default" yaml:"default"`
	Description string `json:"description" yaml:"description"`
	// Type        string `json:"type" yaml:"type"`
}

// SnapDescriptionFinal this type is to descripe the details of the snap to be used for the current application
type SnapDescriptionFinal struct {
	Description string `json:"description" yaml:"description"`
	Name        string `json:"name" yaml:"name"`
	Channel     string `json:"channel" yaml:"channel"`
	Devmode     bool   `json:"devmode" yaml:"devmode"`
	Refresh     bool   `json:"refresh" yaml:"refresh"`
}

// MmeServiceDescription is to descripe the details of the mme service for oai-ran
type MmeServiceDescription struct {
	Description string `json:"description" yaml:"description"`
	Name        string `json:"name" yaml:"name"`
	SnapVersion string `json:"snapVersion" yaml:"snapVersion"`
	IPV4        string `json:"ipv4" yaml:"ipv4"`
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

type CnV1OaiHssDescription struct {
	DatabaseServiceName string `json:"databaseServiceName" yaml:"databaseServiceName"`
}

type CnV2OaiHssDescription struct {
	Snap                SnapDescriptionFinal `json:"snap" yaml:"snap"`
	DatabaseServiceName string               `json:"databaseServiceName" yaml:"databaseServiceName"`
}

type CnV1OaiMmeDescription struct {
	MCC string `json:"mcc" yaml:"mcc"`
	MNC string `json:"mnc" yaml:"mnc"`
}

type CnV2OaiMmeDescription struct {
	Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
	MCC  string               `json:"mcc" yaml:"mcc"`
	MNC  string               `json:"mnc" yaml:"mnc"`
}

type CnV1OaiSpgwDescription struct {
	DNS string `json:"dns" yaml:"dns"`
}

type CnV2OaiSpgwcDescription struct {
	Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
	DNS  string               `json:"dns" yaml:"dns"`
}

type CnV2OaiSpgwuDescription struct {
	Snap SnapDescriptionFinal `json:"snap" yaml:"snap"`
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
