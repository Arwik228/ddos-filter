package firewall

import (
	"crypto/sha256"
	//	"encoding/hex"
	//"log"
	"math/rand"
	"net/http"
	"text/template"
	//	"time"
)

type ViewData struct {
	Token string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func createToken(ln int) string {
	b := make([]rune, ln)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func firewallTemplate(w http.ResponseWriter) {
	/*token := createToken(18)
	hash := sha256.Sum256([]byte(token))

	log.Printf(hex.EncodeToString(hash))
	expire := time.Now().Add(1440 * time.Minute)
	cookie := http.Cookie{Name: "token", Value: string(hash[:]), Path: "/", Expires: expire, MaxAge: 90000}
	http.SetCookie(w, &cookie)
	*/
	data := ViewData{Token: "str"}
	tmpl, _ := template.ParseFiles("views/index.html")
	tmpl.Execute(w, data)
}

func getUserIpAddress(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func checkUserCookieToken(token string, ip string) bool {
	if token == "" {
		return false
	}

	if token == sha256.Sum256([]byte(ip)) {
		return true
	}

	return false
}

func CheckConnection(w http.ResponseWriter, r *http.Request) bool {

	// Отрабатывает для пользователей прошедших авторизацию
	tokenAccess, _ := r.Cookie("token_access")
	if checkUserCookieToken(tokenAccess.Value, getUserIpAddress(r)) {
		return true
	}

	//queryToken := r.URL.Query().Get("fire_token")

	firewallTemplate(w)
	return false
}
