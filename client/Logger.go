package client

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func initLogger(level string) *logrus.Logger {
	logLevel,  err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}
	return &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			ForceColors: true,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  "02/01/2006 03:04 PM",
		},
		Hooks: make(logrus.LevelHooks),
		Level: logLevel,
	}
}