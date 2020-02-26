package log

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var rootLogger *logrus.Logger

func init() {
	rootLogger = logrus.New()

	rootLogger.SetOutput(os.Stdout)
	rootLogger.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(rootLogger.Writer())
}

func Logger() logrus.FieldLogger {
	return rootLogger
}
