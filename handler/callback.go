package handler

import (
	"context"
	"github.com/gorilla/sessions"
	"golang.ysitd.cloud/log"
	"net/http"
)

type CallbackHandler struct {
	ConfigLoader *ConfigLoader  `inject:""`
	Session      sessions.Store `inject:"sessions"`
	Logger       log.Logger     `inject:"callback logger"`
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, deadline := ctx.Deadline(); !deadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, requestTimeout)
		defer cancel()
	}

	state := r.FormValue("state")

	if state == "" {
		http.Error(w, "State is required", http.StatusBadRequest)
		return
	}

	session, err := h.Session.Get(r, "auth."+r.Host)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during loading session", http.StatusInternalServerError)
		return
	}

	expectedState, hasState := session.Values["state"]
	if !hasState {
		h.Logger.Errorln(err)
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	if state != expectedState.(string) {
		h.Logger.Errorln("invalid oauth state")
		http.Error(w, "invalid oauth state", http.StatusConflict)
		return
	}

	code := r.FormValue("code")
	config, err := h.ConfigLoader.Get(ctx, r)
	if err != nil {
		http.Error(w, "error occur when fetching oauth data", http.StatusInternalServerError)
	} else if config == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	_, err = config.Exchange(ctx, code)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "error occur when fetching oauth token", http.StatusInternalServerError)
	}

	var redirect string
	redirectVal, found := session.Values["next"]
	if !found {
		redirect = "/"
	} else {
		redirect = redirectVal.(string)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}
