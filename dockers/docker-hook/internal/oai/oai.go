package oai

import (
	"fmt"
	"log"
	"mosaic5g/docker-hook/internal/pkg/common"
	"mosaic5g/docker-hook/internal/pkg/util"
	"os"
)

const (
	// logPath  = "/root/hook.log"
	// confPath = "/root/config/conf.yaml"

	// Config path and log file for HSS entity V1
	oaiHssLogPathV1  = "/root/hook-oaihss-v1.log"
	oaiHssConfPathV1 = "/root/config/conf.yaml"

	// Config path and log file for MME entity V1
	oaiMmeLogPathV1  = "/root/hook-oaimme-v1.log"
	oaiMmeConfPathV1 = "/root/config/conf.yaml"

	// Config path and log file for SPGW entity V1
	oaiSpgwLogPathV1  = "/root/hook-oaispgw-v1.log"
	oaiSpgwConfPathV1 = "/root/config/conf.yaml"

	// Config path and log file for HSS entity V2
	oaiHssLogPathV2  = "/root/hook-oaihss-v2.log"
	oaiHssConfPathV2 = "/root/config/conf.yaml"

	// Config path and log file for MME entity V1
	oaiMmeLogPathV2  = "/root/hook-oaimme-v2.log"
	oaiMmeConfPathV2 = "/root/config/conf.yaml"

	// Config path and log file for SPGW entity V1
	oaiSpgwcLogPathV2  = "/root/hook-oaispgwc-v2.log"
	oaiSpgwcConfPathV2 = "/root/config/conf.yaml"

	// Config path and log file for SPGW entity V1
	oaiSpgwuLogPathV2  = "/root/hook-oaispgwu-v2.log"
	oaiSpgwuConfPathV2 = "/root/config/conf.yaml"

	// Config path and log file for RAN entities
	oaiEnbLogPathV2  = "/root/hook-oaienb.log"
	oaiEnbConfPathV2 = "/root/config/conf.yaml"
	// oaiEnbLogPathV2  = "/home/cigarier/go/src/mosaic5g/docker-hook/cmd/test/hook-oaienb.log"
	// oaiEnbConfPathV2 = "/home/cigarier/go/src/mosaic5g/docker-hook/cmd/test/oai-conf.yml"

	oaiCuLogPathV2  = "/root/hook-oaicu.log"
	oaiCuConfPathV2 = "/root/config/conf.yaml"

	oaiDuLogPathV2  = "/root/hook-oaidu.log"
	oaiDuConfPathV2 = "/root/config/conf.yaml"

	oaiRccLogPathV2  = "/root/hook-oaircc.log"
	oaiRccConfPathV2 = "/root/config/conf.yaml"

	oaiRruLogPathV2  = "/root/hook-oairru.log"
	oaiRruConfPathV2 = "/root/config/conf.yaml"
)

// const (
// 	// logPath  = "/root/hook.log"
// 	// confPath = "/root/config/conf.yaml"

// 	// Config path and log file for HSS entity V1
// 	oaiHssLogPathV1  = "/root/hook-oaihss-v1.log"
// 	oaiHssConfPathV1 = "/root/config/conf-oaihss-v1.yaml"

// 	// Config path and log file for MME entity V1
// 	oaiMmeLogPathV1  = "/root/hook-oaimme-v1.log"
// 	oaiMmeConfPathV1 = "/root/config/conf-oaimme-v1.yaml"

// 	// Config path and log file for SPGW entity V1
// 	oaiSpgwLogPathV1  = "/root/hook-oaispgw-v1.log"
// 	oaiSpgwConfPathV1 = "/root/config/conf-oaispgw-v1.yaml"

// 	// Config path and log file for HSS entity V2
// 	oaiHssLogPathV2  = "/root/hook-oaihss-v2.log"
// 	oaiHssConfPathV2 = "/root/config/conf-oaihss-v2.yaml"

// 	// Config path and log file for MME entity V1
// 	oaiMmeLogPathV2  = "/root/hook-oaimme-v2.log"
// 	oaiMmeConfPathV2 = "/root/config/conf-oaimme-v2.yaml"

// 	// Config path and log file for SPGW entity V1
// 	oaiSpgwcLogPathV2  = "/root/hook-oaispgwc-v2.log"
// 	oaiSpgwcConfPathV2 = "/root/config/conf-oaispgwc-v2.yaml"

// 	// Config path and log file for SPGW entity V1
// 	oaiSpgwuLogPathV2  = "/root/hook-oaispgwu-v2.log"
// 	oaiSpgwuConfPathV2 = "/root/config/conf-oaispgwu-v2.yaml"

// 	// Config path and log file for RAN entities
// 	oaiEnbLogPathV2  = "/root/hook-oaienb.log"
// 	oaiEnbConfPathV2 = "/root/config/conf-oaienb.yaml"

// 	oaiCuLogPathV2  = "/root/hook-oaicu.log"
// 	oaiCuConfPathV2 = "/root/config/conf-oaicu.yaml"

// 	oaiDuLogPathV2  = "/root/hook-oaidu.log"
// 	oaiDuConfPathV2 = "/root/config/conf-oaidu.yaml"

// 	oaiRccLogPathV2  = "/root/hook-oaircc.log"
// 	oaiRccConfPathV2 = "/root/config/conf-oaircc.yaml"

// 	oaiRruLogPathV2  = "/root/hook-oairru.log"
// 	oaiRruConfPathV2 = "/root/config/conf-oairru.yaml"
// )

//OaiEntity define custom type
type OaiEntity int

const (
	enb OaiEntity = iota
	cu
	du
	rcc
	rru
	hssV1
	mmeV1
	spgwV1
	hssV2
	mmeV2
	spgwcV2
	spgwuV2
)

// Oai stores the log and conf
type Oai struct {
	logFile *os.File    // File for log to write something
	Logger  *log.Logger // Collect log
	// There is different log files for hss, mme, spgwc, and spgwu in case of all-in-one deployment
	logFileHss   *os.File       // File for log to write something
	LoggerHss    *log.Logger    // Collect log
	logFileMme   *os.File       // File for log to write something
	LoggerMme    *log.Logger    // Collect log
	logFileSpgw  *os.File       // File for log to write something
	LoggerSpgw   *log.Logger    // Collect log
	logFileSpgwc *os.File       // File for log to write something
	LoggerSpgwc  *log.Logger    // Collect log
	logFileSpgwu *os.File       // File for log to write something
	LoggerSpgwu  *log.Logger    // Collect log
	Conf         *common.Cfg    // config files
	ConfOaiRan   *common.CfgRan // config files
	OaiEntity    string         // enb
}

// Init the Oai with log and conf
func (me *Oai) Init(entity string) {
	var newFile *os.File
	var err error
	var logPath, confPath string

	switch entity {
	case "ran":
		logPath = oaiEnbLogPathV2
		confPath = oaiEnbConfPathV2
		me.ConfOaiRan = new(common.CfgRan)
		err = me.ConfOaiRan.GetConf(me.Logger, confPath)

		fmt.Println("me.ConfOaiRan", me.ConfOaiRan)

		ranEntity := me.ConfOaiRan.OaianConf.ComponentCarriers.NodeFunction
		if err != nil {
			panic(err)
		}
		newFile, err = os.Create(logPath)
		me.logFile = newFile
		me.Logger = log.New(me.logFile, "[Mosaic5G-"+ranEntity+"-] ", log.Ldate|log.Ltime|log.Llongfile)
		util.PrintFunc(me.Logger, "Configuration is successfully retreived")
		util.PrintFunc(me.Logger, "Configs:", me.ConfOaiRan)
	default:
		me.Conf = new(common.Cfg)
		err = me.Conf.GetConf(me.Logger, confPath)
		if err != nil {
			panic(err)
		}
		util.PrintFunc(me.Logger, "Configuration is successfully retreived")
		util.PrintFunc(me.Logger, "Configs:", me.Conf)
	}

}

// Clean will Close the logFile and clean up Obj
func (me *Oai) Clean() {
	util.PrintFunc(me.Logger, "Closing the config file: ", me.logFile)
	me.logFile.Close()
}

// // Print will Close the logFile and clean up Obj
// func (me *Oai) Print(message string) {
// 	me.Logger.Print(message)
// 	fmt.Println(message)
// }
