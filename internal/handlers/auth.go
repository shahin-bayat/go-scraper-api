package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shahin-bayat/scraper-api/internal/utils"
	"golang.org/x/oauth2"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	verifier := oauth2.GenerateVerifier()
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate random string", http.StatusInternalServerError)
		return
	}
	utils.SetSession(w, r, "verifier", verifier)
	utils.SetSession(w, r, "state", state)
	url := h.config.OAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state != utils.GetSession(r, "state") {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	verifier := utils.GetSession(r, "verifier")
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

	// TODO: use user email or id and check if the user is already in the database
	// TODO:if the user is not in the database, add the user with access token and refresh token to the database
	// TODO: if the user is in the database, update the user info, access token and refresh token in the database

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token.AccessToken)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}

func (h *Handler) HandleAuthStatus(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// validate the access token
	token, err := h.config.OAuth2Config.TokenSource(r.Context(), &oauth2.Token{
		AccessToken: accessToken,
	}).Token()
	if err != nil {
		http.Error(w, "Failed to validate token", http.StatusInternalServerError)
		return
	}

	if !token.Valid() {
		// TODO: fetch the user using provided access token
		// TODO: use the refresh token in the database to get a new access token
		token, err := h.config.OAuth2Config.TokenSource(r.Context(), &oauth2.Token{
			// this should be the refresh token saved in the database
			RefreshToken: token.RefreshToken,
		}).Token()

		if err != nil {
			http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
			return
		}

		// TODO: save the new refresh token and access token in the database
		_ = token.RefreshToken
		w.Header().Set("Authorization", token.AccessToken)

	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token.AccessToken)
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	utils.ClearSession(w, r, "verifier")
	utils.ClearSession(w, r, "state")
	w.WriteHeader(http.StatusOK)
	// remove the refresh token and access token from the database
	// redirect to login
	w.Header().Set("Location", "/auth/login")
	w.WriteHeader(http.StatusTemporaryRedirect)

}
