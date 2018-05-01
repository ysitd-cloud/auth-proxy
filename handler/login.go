package handler

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"golang.ysitd.cloud/log"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var size = len(letterRunes)

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(size)]
	}
	return string(b)
}

type LoginHandler struct {
	ConfigLoader *ConfigLoader  `inject:""`
	Session      sessions.Store `inject:"sessions"`
	Logger       log.Logger     `inject:"login logger"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, deadline := ctx.Deadline(); !deadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, requestTimeout)
		defer cancel()
	}

	config, err := h.ConfigLoader.Get(ctx, r)
	if err != nil {
		http.Error(w, "error occur when fetching oauth data", http.StatusInternalServerError)
	} else if config == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	session, err := h.Session.Get(r, "auth."+r.Host)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during loading session", http.StatusInternalServerError)
		return
	}

	state := randString(8)

	session.Values["state"] = state
	if err := session.Save(r, w); err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during store session", http.StatusInternalServerError)
		return
	}

	redirectUrl := config.AuthCodeURL(state)

	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
