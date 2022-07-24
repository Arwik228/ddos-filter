package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"config"
	"firewall"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(config.Server)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if firewall.CheckConnection(w, r) {
		proxyHandler(w, r)
	}
}

func main() {
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ProxyPort), nil)
}
