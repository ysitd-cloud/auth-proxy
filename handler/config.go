package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"

	"code.ysitd.cloud/proxy/modals/vhost"
	"golang.ysitd.cloud/http/timing"
)

var tokenUrl string
var authUrl string

func init() {
	tokenUrl = fmt.Sprintf("https://%s/oauth/token", os.Getenv("OAUTH_HOST"))
	authUrl = fmt.Sprintf("https://%s/oauth/authorize", os.Getenv("OAUTH_HOST"))
}

func copyUrl(r *http.Request, path string) string {
	u := new(url.URL)
	u.Scheme = "https"
	u.Host = r.Host
	u.Path = path
	return u.String()
}

func newOauthConfig(r *http.Request, host *vhost.VirtualHost) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     host.OauthID,
		ClientSecret: host.OAuthSecret,
		RedirectURL:  copyUrl(r, "/auth/ycloud/callback"),
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenUrl,
			AuthURL:  authUrl,
		},
	}
}

type ConfigLoader struct {
	Vhost     *vhost.Store `inject:""`
	ConfigMap map[string]*oauth2.Config
}

func (cl *ConfigLoader) Get(ctx context.Context, r *http.Request) (config *oauth2.Config, err error) {
	collector := ctx.Value(timingKey).(*timing.Collector)
	timer := collector.New("fetch_oauth", "Fetch OAuth Config")
	timer.Start()

	if cl.ConfigMap == nil {
		cl.ConfigMap = make(map[string]*oauth2.Config)
	}

	config, exists := cl.ConfigMap[r.Host]
	timer.Stop()

	if !exists {
		host, err := cl.Vhost.GetVHost(ctx, r.Host)
		if err != nil {
			return nil, err
		} else if host == nil {
			return nil, nil
		}
		config = newOauthConfig(r, host)
		cl.ConfigMap[r.Host] = config
	}

	return
}
