package handler

import (
	"net/http"
	"strings"
)

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func sessionName(r *http.Request) string {
	return "auth." + stripPort(r.Host)
}
