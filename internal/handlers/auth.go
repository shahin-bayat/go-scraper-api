package handlers

import (
	"context"
	"fmt"
	"net/http"

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
	url := h.services.AuthService.GetAuthCodeUrl(state, verifier)

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

	token, err := h.services.AuthService.ExchangeToken(r.Context(), code, verifier)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("failed to exchange token: %w", err))
		return
	}

	appRedirectURL := generateAppRedirectURL(h.appConfig.AppUniversalURL, token.AccessToken, token.RefreshToken)

	userInfo, err := h.services.AuthService.GetUserInfo(r.Context(), token)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// check if the user exists in database
	existingUser, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
	if err != nil {
		// user doesn't exist
		newUser := models.User{
			Email:         userInfo.Email,
			GivenName:     userInfo.GivenName,
			FamilyName:    userInfo.FamilyName,
			Name:          userInfo.Name,
			Locale:        userInfo.Locale,
			AvatarURL:     userInfo.AvatarURL,
			VerifiedEmail: userInfo.VerifiedEmail,
		}
		// create the user in the database
		userId, err := h.store.UserRepository().CreateUser(&newUser)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user: %w",
				err))
			return
		}
		// create the session in Redis
		err = h.store.UserRepository().CreateUserSession(userId, token)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user session: %w", err))
			return
		}
		http.Redirect(w, r, appRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	// user exists
	// check if the session exists
	_, err = h.store.UserRepository().GetUserSession(existingUser.ID)
	if err != nil {
		// session doesn't exist, create it
		err = h.store.UserRepository().CreateUserSession(existingUser.ID, token)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user session: %w", err))
			return
		}
	} else {
		// session exists, update it
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

	token, err := h.services.AuthService.TokenSource(r.Context(), &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}).Token()
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusUnauthorized, fmt.Errorf("failed to validate token: %w", err))
		return
	}

	// check if the token is valid
	if !h.services.AuthService.TokenValid(token) {
		token, err := h.services.AuthService.TokenSource(r.Context(), &oauth2.Token{
			RefreshToken: refreshToken,
		}).Token()
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to get new access token: %w", err))
			return
		}

		userInfo, err := h.services.AuthService.GetUserInfo(r.Context(), token)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		// get user from database
		user, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to get user by email: %w", err))
			return
		}
		// get user session
		_, err = h.store.UserRepository().GetUserSession(user.ID)
		if err != nil {
			// user session doesn't exist, create it
			err = h.store.UserRepository().CreateUserSession(user.ID, token)
			if err != nil {
				utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to create user session: %w", err))
			}
		} else {
			// user session exists, update it
			err = h.store.UserRepository().UpdateUserSession(user.ID, token)
			if err != nil {
				utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to update user session: %w", err))
				return
			}
		}
		headers := map[string]string{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		}
		utils.WriteJSON(w, http.StatusOK, user, headers)
		return
	}

	userInfo, err := h.services.AuthService.GetUserInfo(r.Context(), token)
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
	var user *models.User
	accessToken := r.Header.Get("access_token")
	refreshToken := r.Header.Get("refresh_token")

	if accessToken == "" && refreshToken == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("missing token"))
		return
	}

	if refreshToken != "" {
		userInfo, err := h.services.AuthService.GetUserInfo(r.Context(), &oauth2.Token{RefreshToken: refreshToken})
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err = h.store.UserRepository().GetUserByEmail(userInfo.Email)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to get user info: %w", err))
			return
		}
		if err := h.services.AuthService.RevokeToken(refreshToken); err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to revoke token: %w", err))
			return
		}
	}

	if accessToken != "" {
		userInfo, err := h.services.AuthService.GetUserInfo(r.Context(), &oauth2.Token{AccessToken: accessToken})
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err = h.store.UserRepository().GetUserByEmail(userInfo.Email)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to get user info: %w", err))
			return
		}
		if err := h.services.AuthService.RevokeToken(accessToken); err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to revoke token: %w", err))
			return
		}
	}

	if err := h.store.UserRepository().DeleteUserSession(user.ID); err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("failed to delete user session: %w", err))
		return
	}
	utils.ClearSession(w, r, "verifier")
	utils.ClearSession(w, r, "state")
	utils.WriteJSON(w, http.StatusNoContent, nil, nil)
}

func generateAppRedirectURL(appURL string, accessToken, refreshToken string) string {
	if refreshToken == "" {
		return fmt.Sprintf("%s?access_token=%s", appURL, accessToken)
	}
	return fmt.Sprintf("%s?access_token=%s&refresh_token=%s", appURL, accessToken, refreshToken)
}
