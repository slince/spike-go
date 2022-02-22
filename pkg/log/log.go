package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	loggers [] *logrus.Logger
}

func NewLogger() *Logger{
	return &Logger{}
}

func (logger *Logger) EnableConsole() *Logger{
	var log = logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.TraceLevel
	logger.loggers = append(logger.loggers, log)
	return logger
}

func (logger *Logger) SetLogFile(logfile string) *Logger{
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log := logrus.New()
	log.Out = file
	logger.loggers = append(logger.loggers, log)
	return logger
}

func (logger *Logger) Trace(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Trace(args...)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Debug(args...)
	}
}

func (logger *Logger) Print(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Print(args...)
	}
}

func (logger *Logger) Info(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Info(args...)
	}
}

func (logger *Logger) Warn(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Warn(args...)
	}
}

func (logger *Logger) Warning(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Warning(args...)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Error(args...)
	}
}

func (logger *Logger) Panic(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Panic(args...)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	for  _, log := range logger.loggers {
		log.Fatal(args...)
	}
}