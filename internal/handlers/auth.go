package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shahin-bayat/scraper-api/internal/utils"
	"golang.org/x/oauth2"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	verifier := oauth2.GenerateVerifier()
	// generate a random string for state
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate random string", http.StatusInternalServerError)
		return
	}

	setSession(w, r, "verifier", verifier)
	setSession(w, r, "state", state)

	url := h.config.OAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state != getSession(r, "state") {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	verifier := getSession(r, "verifier")
	token, err := h.config.OAuth2Config.Exchange(r.Context(), code, oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := h.config.OAuth2Config.Client(r.Context(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to read user info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}

func setSession(w http.ResponseWriter, r *http.Request, key, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func getSession(r *http.Request, key string) string {
	session, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	return session.Value
}
