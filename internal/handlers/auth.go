package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/models"
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

	token, err := h.services.AuthService.Google.Exchange(r.Context(), code, oauth2.SetAuthURLParam("code_verifier", verifier), oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to exchange token: %w", err))
		return
	}

	appRedirectURL := generateAppRedirectURL(h.appConfig.AppUniversalURL, token.AccessToken, token.RefreshToken)

	userInfo, err := getUserInfo(r, h.services.AuthService.Google, token, h.appConfig.GoogleUserInfoURL)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to get user info: %w", err))
		return
	}

	// check if the user exists in database
	existingUser, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
	if err != nil {
		// User doesn't exist
		newUser := models.User{
			Email:         userInfo.Email,
			GivenName:     userInfo.GivenName,
			FamilyName:    userInfo.FamilyName,
			Name:          userInfo.Name,
			Locale:        userInfo.Locale,
			AvatarURL:     userInfo.AvatarURL,
			VerifiedEmail: userInfo.VerifiedEmail,
		}
		// Create the user in the database
		userId, err := h.store.UserRepository().CreateUser(&newUser)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w",
				err))
			return
		}
		// Create the session in Redis
		err = h.store.UserRepository().CreateUserSession(userId, token)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user session: %w", err))
			return
		}
		http.Redirect(w, r, appRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	// User exists
	// Check if the user session exists
	_, err = h.store.UserRepository().GetUserSession(existingUser.ID)
	if err != nil {
		// User session doesn't exist, create it
		err = h.store.UserRepository().CreateUserSession(existingUser.ID, token)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user session: %w", err))
			return
		}
	} else {
		// User session exists, update it
		err = h.store.UserRepository().UpdateUserSession(existingUser.ID, token)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to update user session: %w", err))
			return
		}
	}

	http.Redirect(w, r, appRedirectURL, http.StatusTemporaryRedirect)
}

func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var token *oauth2.Token
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))

	accessToken := r.Header.Get("access_token")
	if accessToken == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("missing access token"))
		return
	}
	refreshToken := r.Header.Get("refresh_token")
	if refreshToken == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("missing refresh token"))
		return
	}

	token, err := h.services.AuthService.Google.TokenSource(r.Context(), &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}).Token()
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusUnauthorized, fmt.Errorf("failed to validate token: %w", err))
		return
	}

	// check if the token is valid
	if !tokenValid(token) {
		token, err := h.services.AuthService.Google.TokenSource(r.Context(), &oauth2.Token{
			RefreshToken: refreshToken,
		}).Token()
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to get new access token: %w", err))
			return
		}
		// get user info
		userInfo, err := getUserInfo(r, h.services.AuthService.Google, token, h.appConfig.GoogleUserInfoURL)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		// get user from db
		user, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to get user by email: %w", err))
			return
		}
		// get the user session
		_, err = h.store.UserRepository().GetUserSession(user.ID)
		if err != nil {
			fmt.Println("User session doesn't exist")
			// user session doesn't exist, create it
			err = h.store.UserRepository().CreateUserSession(user.ID, token)
			if err != nil {
				utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user session: %w", err))
			}
			return
		}
		// user session exists, update it
		err = h.store.UserRepository().UpdateUserSession(user.ID, token)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to update user session: %w", err))
			return
		}
		headers := map[string]string{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		}
		utils.WriteJSON(w, http.StatusOK, user, headers)
		return
	}

	// get user info
	userInfo, err := getUserInfo(r, h.services.AuthService.Google, token, h.appConfig.GoogleUserInfoURL)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// get user from db
	user, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to get user by email: %w", err))
		return
	}

	headers := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": refreshToken,
	}
	utils.WriteJSON(w, http.StatusOK, user, headers)
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

func getUserInfo(r *http.Request, oAuth2config *oauth2.Config, token *oauth2.Token, googleUserInfoUrl string) (models.GoogleUserInfo, error) {
	client := oAuth2config.Client(r.Context(), token)
	resp, err := client.Get(googleUserInfoUrl)
	if err != nil {
		return models.GoogleUserInfo{}, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()
	userInfo := models.GoogleUserInfo{}
	err = utils.DecodeResponseBody(resp.Body, &userInfo)
	if err != nil {
		return models.GoogleUserInfo{}, fmt.Errorf("failed to decode user info: %w", err)
	}
	return userInfo, nil
}

func generateAppRedirectURL(appURL string, accessToken, refreshToken string) string {
	if refreshToken == "" {
		return fmt.Sprintf("%s?access_token=%s", appURL, accessToken)
	}
	return fmt.Sprintf("%s?access_token=%s&refresh_token=%s", appURL, accessToken, refreshToken)
}

func tokenValid(token *oauth2.Token) bool {
	if token == nil {
		return false
	}
	return token.Valid() && token.Expiry.After(time.Now())
}

// TODO: : delete user session on user delete
// TODO: : delete access token of the user session on logout
