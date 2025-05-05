package ocookie

import (
	"net/http"
	"time"
)

func Get(w http.ResponseWriter, r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func Set(w http.ResponseWriter, r *http.Request, name, value string) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: name, Value: value, Path: "/", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func Remove(w http.ResponseWriter, r *http.Request, name string) {
	expiration := time.Now().Add(-24 * time.Hour)
	cookie := http.Cookie{Name: name, Value: "", Path: "/", Expires: expiration}
	http.SetCookie(w, &cookie)
}
