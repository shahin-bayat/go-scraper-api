package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"golang.org/x/oauth2"
)

type contextKey string

const providerKey contextKey = "provider"

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))

	verifier := oauth2.GenerateVerifier()
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate random string", http.StatusInternalServerError)
		return
	}
	utils.SetSession(w, r, "verifier", verifier)
	utils.SetSession(w, r, "state", state)
	url := h.services.AuthService.Google.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	verifier := utils.GetSession(r, "verifier")
	sessionState := utils.GetSession(r, "state")

	if state != sessionState {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("state does not match"))
		return
	}

	token, err := h.services.AuthService.Google.Exchange(r.Context(), code, oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to exchange token: %w", err))
		return
	}

	// TODO: use user email or id and check if the user is already in the database
	// TODO:if the user is not in the database, add the user with access token and refresh token to the database
	// TODO: if the user is in the database, update the user info, access token and refresh token in the database
	googleAuthClientUrl := os.Getenv("GOOGLE_AUTH_CLIENT_URL")
	if googleAuthClientUrl == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to get google auth client url"))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%saccess_token=%s", googleAuthClientUrl, token.AccessToken), http.StatusTemporaryRedirect)
}

func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var token *oauth2.Token
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))

	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("missing access token"))
		return
	}

	// validate the access token
	token, err := h.services.AuthService.Google.TokenSource(r.Context(), &oauth2.Token{
		AccessToken: accessToken,
	}).Token()
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to validate token: %w", err))
		return
	}

	if !token.Valid() {
		// TODO: fetch the user using provided access token
		// TODO: use the refresh token in the database to get a new access token
		token, err := h.services.AuthService.Google.TokenSource(r.Context(), &oauth2.Token{
			// this should be the refresh token saved in the database
			RefreshToken: token.RefreshToken,
		}).Token()

		if err != nil {
			utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to get new access token: %w", err))
			return
		}

		// TODO: save the new refresh token and access token in the database
		_ = token.RefreshToken
		w.Header().Set("Authorization", token.AccessToken)

	}

	client := h.services.AuthService.Google.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to get user info: %w", err))
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to decode user info: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token.AccessToken)
	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, http.StatusOK, userData)

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
