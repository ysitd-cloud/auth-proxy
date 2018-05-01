package handler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"code.ysitd.cloud/proxy/modals/vhost"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const requestTimeout = 10 * time.Second

type Proxy struct {
	VHost   *vhost.Store       `inject:""`
	Logger  logrus.FieldLogger `inject:"proxy handler"`
	Client  *http.Client       `inject:""`
	Session sessions.Store     `inject:"sessions"`
}

func (h *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, deadline := ctx.Deadline(); !deadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, requestTimeout)
		defer cancel()
	}

	host, err := h.VHost.GetVHost(ctx, r.Host)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during routing", http.StatusInternalServerError)
		return
	} else if host == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	session, err := h.Session.Get(r, "auth."+r.Host)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during loading session", http.StatusInternalServerError)
		return
	} else if session.IsNew {
		http.Redirect(w, r, "/auth/ysitd", http.StatusFound)
		return
	}

	u := new(url.URL)
	*u = *r.URL
	u.Host = host.BackendHost + ":" + strconv.Itoa(host.BackendPort)

	defer r.Body.Close()

	req, err := http.NewRequest(r.Method, u.String(), r.Body)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during connecting backend", http.StatusBadGateway)
		return
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during connecting backend", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	io.Copy(w, resp.Body)
}
