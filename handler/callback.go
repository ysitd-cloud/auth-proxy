package handler

import (
	"net/http"

	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.ysitd.cloud/log"
	"io/ioutil"
	"os"
)

type CallbackHandler struct {
	Client       *http.Client   `inject:""`
	ConfigLoader *ConfigLoader  `inject:""`
	Session      sessions.Store `inject:"sessions"`
	Logger       log.Logger     `inject:"callback logger"`
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := newContext(r)
	defer cancel()

	state := r.FormValue("state")

	if state == "" {
		http.Error(w, "State is required", http.StatusBadRequest)
		return
	}

	session, err := h.Session.Get(r, sessionName(r))
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

	token, err := config.Exchange(ctx, code)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "error occur when fetching oauth token", http.StatusInternalServerError)
		return
	}

	var redirect string
	redirectVal, found := session.Values["next"]
	if !found {
		redirect = "/"
	} else {
		redirect = redirectVal.(string)
	}

	delete(session.Values, "next")
	delete(session.Values, "state")

	infoUrl := fmt.Sprintf("https://%s/api/v1/user/info", os.Getenv("OAUTH_HOST"))
	req, err := http.NewRequest("GET", infoUrl, nil)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "error occur when fetching user detail", http.StatusInternalServerError)
		return
	}

	authHeader := fmt.Sprintf("Bearer %s", token.AccessToken)
	req.Header.Set("Authorization", authHeader)

	resp, err := h.Client.Do(req)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "error occur when fetching user detail", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "error occur when fetching user detail", http.StatusInternalServerError)
		return
	}

	var user User

	if err := json.Unmarshal(body, &user); err != nil {
		h.Logger.Debugln(string(body))
		h.Logger.Errorln(err)
		http.Error(w, "error occur when fetching user detail", http.StatusInternalServerError)
		return
	}

	session.Values["user"] = user.Username
	if err := session.Save(r, w); err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during store session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}
