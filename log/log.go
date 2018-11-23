package log

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

type Logger struct {
	loggers [] *logrus.Logger
}

// Creates a new logger
func NewLogger() *Logger{
	return &Logger{}
}

// Enable console output
func (logger *Logger) EnableConsole() *Logger{
	var logrus = logrus.New()
	logrus.Out = os.Stdout
	logger.loggers = append(logger.loggers, logrus)
	return logger
}

// Enable log file
func (logger *Logger) SetLogFile(logfile string) *Logger{
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logrus := logrus.New()
	logrus.Out = file
	logger.loggers = append(logger.loggers, logrus)
	return logger
}

// Trace logs a message at level Trace on the standard logger.
func (logger *Logger) Trace(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Trace(args...)
	}
}

// Debug logs a message at level Debug on the standard logger.
func (logger *Logger) Debug(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Debug(args...)
	}
}

// Print logs a message at level Info on the standard logger.
func (logger *Logger) Print(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Print(args...)
	}
}

// Info logs a message at level Info on the standard logger.
func (logger *Logger) Info(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Info(args...)
	}
}

// Warn logs a message at level Warn on the standard logger.
func (logger *Logger) Warn(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Warn(args...)
	}
}

// Warning logs a message at level Warn on the standard logger.
func (logger *Logger) Warning(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Warning(args...)
	}
}

// Error logs a message at level Error on the standard logger.
func (logger *Logger) Error(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Error(args...)
	}
}

// Panic logs a message at level Panic on the standard logger.
func (logger *Logger) Panic(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Panic(args...)
	}
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logger *Logger) Fatal(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Fatal(args...)
	}
}

// Tracef logs a message at level Trace on the standard logger.
func (logger *Logger) Tracef(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Tracef(format, args...)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func (logger *Logger) Debugf(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Debugf(format, args...)
	}
}

// Printf logs a message at level Info on the standard logger.
func (logger *Logger) Printf(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Printf(format, args...)
	}
}

// Infof logs a message at level Info on the standard logger.
func (logger *Logger) Infof(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Infof(format, args...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func (logger *Logger) Warnf(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Warnf(format, args...)
	}
}

// Warningf logs a message at level Warn on the standard logger.
func (logger *Logger) Warningf(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Warningf(format, args...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func (logger *Logger) Errorf(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Errorf(format, args...)
	}
}

// Panicf logs a message at level Panic on the standard logger.
func (logger *Logger) Panicf(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Panicf(format, args...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Fatalf(format, args...)
	}
}

// Traceln logs a message at level Trace on the standard logger.
func (logger *Logger) Traceln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Traceln(args...)
	}
}

// Debugln logs a message at level Debug on the standard logger.
func (logger *Logger) Debugln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Debugln(args...)
	}
}

// Println logs a message at level Info on the standard logger.
func (logger *Logger) Println(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Println(args...)
	}
}

// Infoln logs a message at level Info on the standard logger.
func (logger *Logger) Infoln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Infoln(args...)
	}
}

// Warnln logs a message at level Warn on the standard logger.
func (logger *Logger) Warnln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Warnln(args...)
	}
}

// Warningln logs a message at level Warn on the standard logger.
func (logger *Logger) Warningln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Warningln(args...)
	}
}

// Errorln logs a message at level Error on the standard logger.
func (logger *Logger) Errorln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Errorln(args...)
	}
}

// Panicln logs a message at level Panic on the standard logger.
func (logger *Logger) Panicln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Panicln(args...)
	}
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logger *Logger) Fatalln(args ...interface{}) {
	for  _, logrus := range logger.loggers {
		logrus.Fatalln(args...)
	}
}
