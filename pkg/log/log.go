package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	loggers [] *logrus.Logger
}

type Config struct {
	Console bool
	File  string
	Level string
}

func NewLogger(config Config) (logger *Logger, err error){
	var level logrus.Level
	level, err = logrus.ParseLevel(config.Level)
	if err != nil {
		return
	}
	logger = new(Logger)
	if config.Console {
		logger.EnableConsole(level)
	}
	if len(config.File) > 0 {
		err = logger.SetLogFile(config.File, level)
	}
	return
}

func (l *Logger) EnableConsole(level logrus.Level){
	var log = logrus.New()
	log.Out = os.Stdout
	log.Level = level
	l.loggers = append(l.loggers, log)
}

func (l *Logger) SetLogFile(logfile string, level logrus.Level) error {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	var logger = logrus.New()
	logger.Out = file
	l.loggers = append(l.loggers, logger)
	return nil
}

func (l *Logger) Trace(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Trace(args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Debug(args...)
	}
}

func (l *Logger) Print(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Print(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Info(args...)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Warn(args...)
	}
}

func (l *Logger) Warning(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Warning(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Error(args...)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Panic(args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	for  _, log := range l.loggers {
		log.Fatal(args...)
	}
}