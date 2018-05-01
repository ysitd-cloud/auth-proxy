package vhost

type VirtualHost struct {
	Hostname    string
	OauthID     string
	OAuthSecret string
	BackendPath string
}
