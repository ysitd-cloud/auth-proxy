package main

import (
	"net/http"
	"os"

	"app.ysitd/proxy/bootstrap"
)

func main() {
	go bootstrap.GetCache().Run()
	handler := bootstrap.GetHandler()
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		bootstrap.GetMainLogger().Errorln(err)
	}
}
