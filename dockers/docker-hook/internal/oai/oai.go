package oai

import (
	"log"
	"mosaic5g/docker-hook/internal/pkg/common"
	"mosaic5g/docker-hook/internal/pkg/util"
	"os"
)

// Oai stores the log and conf
type Oai struct {
	logFile *os.File    // File for log to write something
	Logger  *log.Logger // Collect log
	Conf    *common.Cfg // config files
}

// Init the Oai with log and conf
func (me *Oai) Init(logPath string, confPath string) error {
	newFile, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}
	me.logFile = newFile
	me.Logger = log.New(me.logFile, "[Mosaic5G-] ", log.Ldate|log.Ltime|log.Llongfile)

	me.Conf = new(common.Cfg)
	err = me.Conf.GetConf(me.Logger, confPath)
	if err != nil {
		panic(err)
	}

	util.PrintFunc(me.Logger, "Configuration is successfully retreived")
	util.PrintFunc(me.Logger, "Configs:", me.Conf)
	return nil
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
