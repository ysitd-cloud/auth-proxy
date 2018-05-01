package bootstrap

import (
	"net/http"

	"code.ysitd.cloud/proxy/handler"
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
)

var h handler.MainHandler

func init() {
	var graph inject.Graph
	graph.Logger = initLogger()

	fns := []func(*inject.Graph){
		InjectDB,
		InjectLogger,
		InjectSession,
	}

	for _, fn := range fns {
		fn(&graph)
	}

	graph.Provide(&inject.Object{Value: &h})

	graph.Populate()
}

func GetHandler() http.Handler {
	return &h
}

func GetMainLogger() logrus.FieldLogger {
	return logger.WithField("source", "main")
}