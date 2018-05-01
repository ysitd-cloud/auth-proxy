package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"code.ysitd.cloud/proxy/modals/vhost"
)

func copyUrl(r *http.Request, path string) string {
	u := new(url.URL)
	u.Scheme = r.URL.Scheme
	u.Host = r.URL.Host
	u.Path = path
	return u.String()
}

func newOauthConfig(r *http.Request, host *vhost.VirtualHost) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     host.OauthID,
		ClientSecret: host.OAuthSecret,
		RedirectURL:  fmt.Sprintf("https://%s/oauth/authorize", r.Host),
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			TokenURL: copyUrl(r, "/oauth/token"),
			AuthURL:  copyUrl(r, "/oauth/authorize"),
		},
	}
}

type ConfigLoader struct {
	Vhost     *vhost.Store `inject:""`
	ConfigMap map[string]*oauth2.Config
}

func (cl *ConfigLoader) Get(ctx context.Context, r *http.Request) (config *oauth2.Config, err error) {
	if cl.ConfigMap == nil {
		cl.ConfigMap = make(map[string]*oauth2.Config)
	}

	config, exists := cl.ConfigMap[r.Host]

	if !exists {
		host, err := cl.Vhost.GetVHost(ctx, r.Host)
		if err != nil {
			return
		} else if host == nil {
			return nil, nil
		}
		config = newOauthConfig(r, host)
		cl.ConfigMap[r.Host] = config
	}

	return
}
