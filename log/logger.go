package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var std *logrus.Logger

func SetLogFile(filepath string) {

	logFile, err := os.OpenFile(filepath, )
}