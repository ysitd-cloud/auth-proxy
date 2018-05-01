package vhost

type VirtualHost struct {
	Hostname    string
	OauthID     string
	OAuthSecret string
	BackendHost string
	BackendPort int
}
