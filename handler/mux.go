package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/handlers"

	"golang.ysitd.cloud/log"

	"code.ysitd.cloud/proxy/timing"
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

	ctx, cancel := newContext(r)
	defer cancel()

	done := make(chan bool)

	go func(done chan bool) {
		h.handleHTTP(ctx, w, r)
		done <- true
	}(done)

	select {
	case <-done:
		break
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			h.Logger.Errorln(err)
		}
		break
	}
}

func (h *MainHandler) handleHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	collector := timing.NewCollector()
	ctx = context.WithValue(ctx, timingKey, collector)

	h.Handler.ServeHTTP(collector.Response(w), r.WithContext(ctx))
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
