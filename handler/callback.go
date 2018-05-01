package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.ysitd.cloud/log"

	"code.ysitd.cloud/proxy/timing"
)

type CallbackHandler struct {
	Client       *http.Client   `inject:""`
	ConfigLoader *ConfigLoader  `inject:""`
	Session      sessions.Store `inject:"sessions"`
	Logger       log.Logger     `inject:"callback logger"`
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	collector := r.Context().Value("timing").(*timing.Collector)

	timer := collector.New("state_check", "Check State")
	timer.Start()
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
	timer.Stop()
	delete(session.Values, "state")

	if state != expectedState.(string) {
		h.Logger.Errorln("invalid oauth state")
		http.Error(w, "invalid oauth state", http.StatusConflict)
		return
	}

	config, err := h.ConfigLoader.Get(r.Context(), r)
	if err != nil {
		http.Error(w, "error occur when fetching oauth data", http.StatusInternalServerError)
	} else if config == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	timer = collector.New("exchange_token", "Exchange Token")
	timer.Start()
	code := r.FormValue("code")
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "error occur when fetching oauth token", http.StatusInternalServerError)
		return
	}
	timer.Stop()

	timer = collector.New("fetch_user", "Fetch User")
	timer.Start()
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

	timer.Stop()

	var redirect string
	redirectVal, found := session.Values["next"]
	delete(session.Values, "next")
	if !found {
		redirect = "/"
	} else {
		redirect = redirectVal.(string)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}
