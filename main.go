package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"config"
	"firewall"
)

type server struct {
	db *sql.DB
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(config.Server)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

func (s *server) requestHandler(w http.ResponseWriter, r *http.Request) {
	if firewall.CheckConnection(s.db, w, r) {
		proxyHandler(w, r)
	}
}

func main() {
	db, _ := sql.Open("sqlite3", "./database/sqlite.db")

	defer db.Close()
	s := server{db: db}

	http.HandleFunc("/", s.requestHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ProxyPort), nil)
}
