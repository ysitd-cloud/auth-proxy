package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"

	"golang.ysitd.cloud/log"
)

var logger *logrus.Logger

func initLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}

	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		var err error
		logger, err = log.NewForContainer(logFile)
		if err != nil {
			panic(err)
		}
	} else {
		logger = logrus.New()
	}

	if os.Getenv("VERBOSE") != "" {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}

func InjectLogger(graph *inject.Graph) {
	logger := initLogger()
	graph.Provide(
		&inject.Object{Value: logger},
		&inject.Object{Name: "main logger", Value: logger.WithField("source", "main")},
		&inject.Object{Name: "callback logger", Value: logger.WithField("source", "callback")},
		&inject.Object{Name: "login logger", Value: logger.WithField("source", "login")},
		&inject.Object{Name: "proxy logger", Value: logger.WithField("source", "proxy")},
	)
}
