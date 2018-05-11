package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	goLog "golang.ysitd.cloud/log"

	"app.ysitd/proxy/modals/vhost"
	"golang.ysitd.cloud/http/timing"
)

type Proxy struct {
	VHost   *vhost.Store   `inject:""`
	Logger  goLog.Logger   `inject:"proxy logger"`
	Session sessions.Store `inject:"sessions"`
	Proxies map[string]*httputil.ReverseProxy
}

func (h *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debugf("Fetch host %s", r.Host)
	host, err := h.VHost.GetVHost(r.Context(), r.Host)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during routing", http.StatusInternalServerError)
		return
	} else if host == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	session, err := h.Session.Get(r, sessionName(r))
	_, isLogin := session.Values["user"]
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during loading session", http.StatusInternalServerError)
		return
	} else if session.IsNew || !isLogin {
		session.Values["next"] = r.URL.RequestURI()
		if err := session.Save(r, w); err != nil {
			h.Logger.Errorln(err)
			http.Error(w, "Error during store session", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/auth/ycloud", http.StatusFound)
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
		proxy = httputil.NewSingleHostReverseProxy(backendUrl)
		entry := h.Logger.WithFields(logrus.Fields{
			"source": "reverse_proxy",
			"target": host.BackendPath,
		})
		proxy.ErrorLog = log.New(entry.WriterLevel(logrus.ErrorLevel), "", 0)
		h.Proxies[host.BackendPath] = proxy
	}

	collector := r.Context().Value(timingKey).(*timing.Collector)
	timer := collector.New("backend_proxy", "Backend Proxy")
	timer.Start()
	proxy.ServeHTTP(timer.Response(w), r)
}
