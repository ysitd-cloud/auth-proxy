package bootstrap

import (
	"net/http"
	"time"

	"github.com/facebookgo/inject"
)

func InjectHTTPClient(graph *inject.Graph) {
	graph.Provide(&inject.Object{Value: &http.Client{
		Timeout: 30 * time.Second,
	}})
}
