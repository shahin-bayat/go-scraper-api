package utils

import "net/http"

// TODO: CHECK FRO PRODUCTION and use samesite, http only and secure
func SetSession(w http.ResponseWriter, r *http.Request, key, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		HttpOnly: false,
		// SameSite: http.SameSiteStrictMode,
	})
}

func GetSession(r *http.Request, key string) string {
	session, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	return session.Value
}

func ClearSession(w http.ResponseWriter, r *http.Request, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	})
}
