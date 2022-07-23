package firewall

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

type ViewData struct {
	Token string
}

var server = "http://localhost:8080/"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func createToken(ln int) string {
	b := make([]rune, ln)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func firewallTemplate(w http.ResponseWriter) {
	token := createToken(18)
	hash := sha256.Sum256([]byte(token))

	log.Printf(hex.EncodeToString(hash))
	expire := time.Now().Add(1440 * time.Minute)
	cookie := http.Cookie{Name: "token", Value: string(hash[:]), Path: "/", Expires: expire, MaxAge: 90000}
	http.SetCookie(w, &cookie)

	data := ViewData{Token: token}
	tmpl, _ := template.ParseFiles("assets/html/index.html")
	tmpl.Execute(w, data)
}

func checkConnection(w http.ResponseWriter, r *http.Request) {
	//	firewallTemplate(w)
}
