package handler

import "net/http"

type MainHandler struct {
	mux             *http.ServeMux
	CallbackHandler *CallbackHandler `inject:""`
	LoginHandler    *LoginHandler    `inject:""`
	Proxy           *Proxy           `inject:""`
}

func (h *MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.mux == nil {
		h.setupMux()
	}
	h.mux.ServeHTTP(w, r)
}

func (h *MainHandler) setupMux() {
	h.mux = http.NewServeMux()
	h.mux.Handle("/", h.Proxy)
	h.mux.Handle("/auth/ysitd", h.LoginHandler)
	h.mux.Handle("/auth/ysitd/callback", h.CallbackHandler)
}
