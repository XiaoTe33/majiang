package model

import (
	"github.com/sirupsen/logrus"
	"majiang/log"
)

var myLog = log.Log

func init() {
	myLog.SetLevel(logrus.TraceLevel)
}
