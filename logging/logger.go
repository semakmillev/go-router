package logging

import (
	"github.com/chasex/glog"
	"fmt"
)

func GetLogger(fileName string) *glog.Logger {

	options := glog.LogOptions{
		File:  "./logs/" + fileName,
		Flag:  glog.LstdFlags | glog.Lmicroseconds | glog.Lshortfile,
		Level: glog.Ldebug,
		Mode:  glog.R_Day,
		//Maxsize: 1024 * 1024 * 16,
	}
	logger, err := glog.New(options)
	if err != nil {
		panic(err)
	}
	return logger
}


func Warn(logger *glog.Logger, message string){
	fmt.Println("warn: "+message)
	logger.Warn(message)
}

func Info(logger *glog.Logger, message string){
	fmt.Println("info: "+message)
	logger.Info(message)
}

func Error(logger *glog.Logger, message string){
	fmt.Println("err: "+message)
	logger.Error(message)
}