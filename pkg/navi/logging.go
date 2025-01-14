package navi

import (
	"log"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var accessLogger *log.Logger

func InitializeLogging() {
	accessLogger = log.New(&lumberjack.Logger{
		Filename:   viper.GetString("logging.access"),
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}, "", 0)
}
