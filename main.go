package main

import (
	"code.ysitd.cloud/proxy/bootstrap"
	"net/http"
	"os"
)

func main() {
	handler := bootstrap.GetHandler()
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		bootstrap.GetMainLogger().Errorln(err)
	}
}
