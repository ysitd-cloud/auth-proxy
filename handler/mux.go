package handler

import (
	"net/http"

	"github.com/gorilla/handlers"

	"golang.ysitd.cloud/log"
)

type MainHandler struct {
	http.Handler
	Logger          log.Logger       `inject:"main logger"`
	CallbackHandler *CallbackHandler `inject:""`
	LoginHandler    *LoginHandler    `inject:""`
	Proxy           *Proxy           `inject:""`
}

func (h *MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Handler == nil {
		h.setupMux()
	}

	h.Handler.ServeHTTP(w, r)
}

func (h *MainHandler) setupMux() {
	mux := http.NewServeMux()
	mux.Handle("/", h.Proxy)
	mux.Handle("/auth/ycloud", h.LoginHandler)
	mux.Handle("/auth/ycloud/callback", h.CallbackHandler)

	httpLogger := h.Logger.WithField("source", "http")
	handler := handlers.CombinedLoggingHandler(httpLogger.Writer(), mux)
	logOpt := handlers.RecoveryLogger(h.Logger.WithField("source", "recover"))
	printOpt := handlers.PrintRecoveryStack(true)
	handler = handlers.RecoveryHandler(logOpt, printOpt)(handler)

	h.Handler = handler
}
