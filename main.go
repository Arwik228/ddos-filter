package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(config.server)
	if err != nil {
		log.Println(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

	for _, cookie := range r.Cookies() {
		log.Println("Found a cookie named:", cookie.Name)
	}
	proxyHandler(w, r)
}

func main() {
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", config.proxyPort), nil)
}
