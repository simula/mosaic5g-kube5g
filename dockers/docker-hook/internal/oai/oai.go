package oai

import (
	"docker-hook/internal/pkg/common"
	"log"
	"os"
	"time"
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
		return err
	}
	me.logFile = newFile
	me.Logger = log.New(me.logFile, "[Debug]"+time.Now().Format("2006-01-02 15:04:05")+" ", log.Lshortfile)
	me.Conf = new(common.Cfg)
	err = me.Conf.GetConf(me.Logger, confPath)
	if err != nil {
		return err
	}
	me.Logger.Print("Configs:")
	me.Logger.Print(me.Conf)
	return nil
}

// Clean will Close the logFile and clean up Obj
func (me *Oai) Clean() {
	me.logFile.Close()
}
