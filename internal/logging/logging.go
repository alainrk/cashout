package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func GetLogger(level string) *logrus.Logger {
	log := logrus.New()

	switch strings.ToLower(level) {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn", "warning":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		// Default to info level if the provided level is invalid
		log.SetLevel(logrus.InfoLevel)
		log.Warnf("Invalid log level '%s', defaulting to 'info'", level)
	}

	return log
}
