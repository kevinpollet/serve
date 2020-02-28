package log

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var root *logrus.Logger

func init() {
	root = logrus.New()

	root.SetOutput(os.Stdout)
	root.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(root.Writer())
}

func Logger() logrus.FieldLogger {
	return root
}
