package handler

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"code.ysitd.cloud/proxy/modals/vhost"
)

const requestTimeout = 10 * time.Second

type Proxy struct {
	VHost   *vhost.Store       `inject:""`
	Logger  logrus.FieldLogger `inject:"proxy logger"`
	Session sessions.Store     `inject:"sessions"`
	Proxies map[string]*httputil.ReverseProxy
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

	if h.Proxies == nil {
		h.Proxies = make(map[string]*httputil.ReverseProxy)
	}

	proxy, exists := h.Proxies[host.BackendPath]

	if !exists {
		backendUrl, err := url.Parse(host.BackendPath)
		if err != nil {
			h.Logger.Errorln(err)
			http.Error(w, "Error during parse backend", http.StatusInternalServerError)
			return
		}
		h.Proxies[host.BackendPath] = httputil.NewSingleHostReverseProxy(backendUrl)
		proxy = h.Proxies[host.BackendPath]
	}

	proxy.ServeHTTP(w, r)
}
