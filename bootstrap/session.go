package bootstrap

import (
	"github.com/facebookgo/inject"
	"github.com/gorilla/sessions"
	"os"
)

func initSession() (store *sessions.CookieStore) {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options.HttpOnly = true
	return
}

func InjectSession(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Name: "sessions", Value: initSession()},
	)
}
