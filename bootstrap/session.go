package bootstrap

import (
	"github.com/facebookgo/inject"
	"github.com/gorilla/sessions"
	"os"
)

func initSession() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}

func InjectSession(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Name: "sessions", Value: initSession()},
	)
}
