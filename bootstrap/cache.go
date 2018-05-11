package bootstrap

import (
	"github.com/facebookgo/inject"

	"app.ysitd/proxy/modals/vhost"
)

var cache *vhost.Cache

func InjectCache(graph *inject.Graph) {
	cache = vhost.NewCache()
	graph.Provide(&inject.Object{Value: cache})
}
