package firewall

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type ViewData struct {
	Url string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func createToken(ln int) string {
	b := make([]rune, ln)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func firewallTemplate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	token := createToken(24)

	_, err := db.Exec("INSERT INTO firewall_tokens (token, address, created_at) VALUES (?, ?, CURRENT_TIMESTAMP);", token, getUserIpAddress(r))

	if err != nil {
		log.Print(err)
	}

	newUrl, err := url.Parse(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	values := newUrl.Query()
	values.Set("firewall_token", token)
	newUrl.RawQuery = values.Encode()

	data := ViewData{Url: newUrl.String()}
	tmpl, _ := template.ParseFiles("views/index.html")
	tmpl.Execute(w, data)
}

func getUserIpAddress(r *http.Request) string { //todo delete port
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}

	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	return strings.SplitN(IPAddress, ":", 2)[0]
}

func checkUserCookieToken(db *sql.DB, token string, ip string) bool {
	if token == "" {
		return false
	}

	var count int
	row := db.QueryRow("SELECT count(*) as count  FROM access_tokens WHERE token = ? AND address = ?", token, ip)
	err := row.Scan(&count)

	if err != nil {
		log.Print(err)
	}

	if count == 1 {
		return true
	}

	return false
}

func checkUserQuery(db *sql.DB, w http.ResponseWriter, r *http.Request) bool {
	token := r.FormValue("firewall_token")
	ip := getUserIpAddress(r)

	if token == "" {
		return false
	}

	var id int
	row := db.QueryRow("SELECT id FROM firewall_tokens WHERE token = ? AND address = ?", token, ip)
	err := row.Scan(&id)

	if err != nil {
		log.Print(err)
	}

	if id > 0 {
		token := createToken(24)
		db.Exec("DELETE FROM firewall_tokens WHERE id = ?", id)
		_, err := db.Exec("INSERT INTO access_tokens (token, address, created_at) VALUES (?, ?, CURRENT_TIMESTAMP);", token, ip)

		if err != nil {
			log.Print(err)
		}

		cookie := http.Cookie{Name: "token_access", Path: "/", Value: token, MaxAge: 9000}
		http.SetCookie(w, &cookie)

		return true
	}

	return false
}

func CheckConnection(db *sql.DB, w http.ResponseWriter, r *http.Request) bool {
	var tokenAccess string
	tokenAccessCookie, err := r.Cookie("token_access")

	if err != nil {
		tokenAccess = ""
	} else {
		tokenAccess = tokenAccessCookie.Value
	}

	// Проверяем валидность токена доступа
	if checkUserCookieToken(db, tokenAccess, getUserIpAddress(r)) == true {
		return true
	} else if tokenAccess != "" {
		cookie := http.Cookie{Name: "token_access", Path: "/", MaxAge: 0}
		http.SetCookie(w, &cookie)
	}

	// Проверяем валидность токена проверки если он есть
	if checkUserQuery(db, w, r) == true {
		return true
	}

	firewallTemplate(db, w, r)
	return false
}
