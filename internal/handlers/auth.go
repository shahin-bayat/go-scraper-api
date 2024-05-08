package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"golang.org/x/oauth2"
)

type contextKey string

const providerKey contextKey = "provider"

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) error {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))

	verifier := oauth2.GenerateVerifier()
	state, err := utils.GenerateRandomString(32)
	if err != nil {
		return err
	}
	utils.SetSession(w, r, "verifier", verifier)
	utils.SetSession(w, r, "state", state)
	url := h.services.AuthService.GetAuthCodeUrl(state, verifier)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (h *Handler) HandleProviderCallback(w http.ResponseWriter, r *http.Request) error {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	verifier := utils.GetSession(r, "verifier")
	sessionState := utils.GetSession(r, "state")

	if state != sessionState {
		return h.services.AuthService.ErrorAuthStateMissmatch()
	}

	token, err := h.services.AuthService.ExchangeToken(r.Context(), code, verifier)
	if err != nil {
		return err
	}

	appRedirectURL := generateAppRedirectURL(h.appConfig.AppUniversalURL, token.AccessToken, token.RefreshToken)

	userInfo, err := h.services.AuthService.ValidateToken(r.Context(), token)
	if err != nil {
		return err
	}

	// check if the user exists in database
	existingUser, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
	if err != nil {
		// user doesn't exist
		newUser := models.NewUser(userInfo)
		userId, err := h.store.UserRepository().CreateUser(&newUser)
		if err != nil {
			return err
		}
		err = h.store.UserRepository().CreateUserSession(userId, token)
		if err != nil {
			return err
		}
		http.Redirect(w, r, appRedirectURL, http.StatusTemporaryRedirect)
		return nil
	}
	// user exists
	_, err = h.store.UserRepository().GetUserSession(existingUser.ID)
	if err != nil {
		// session doesn't exist
		err = h.store.UserRepository().CreateUserSession(existingUser.ID, token)
		if err != nil {
			return err
		}
	} else {
		// session exists
		err = h.store.UserRepository().UpdateUserSession(existingUser.ID, token)
		if err != nil {
			return err
		}
	}
	http.Redirect(w, r, appRedirectURL, http.StatusTemporaryRedirect)
	return nil
}

func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) error {
	var token *oauth2.Token
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerKey, provider))

	accessToken := r.Header.Get("access_token")
	refreshToken := r.Header.Get("refresh_token")

	if accessToken == "" && refreshToken == "" {
		return utils.NewAPIError(
			http.StatusUnprocessableEntity, h.services.AuthService.ErrorMissingAuthorizationHeader(),
		)
	}

	token, err := h.services.AuthService.Token(
		r.Context(), &oauth2.Token{
			AccessToken: accessToken,
		},
	)
	if err != nil {
		return err
	}

	userInfo, err := h.services.AuthService.ValidateToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, h.services.AuthService.ErrorDecodeUserInfo()) {
			return err
		} else {
			// access token is invalid
			if refreshToken == "" {
				return utils.NewAPIError(http.StatusUnprocessableEntity, h.services.AuthService.ErrorInvalidToken())
			}
			// refresh the token
			token, err := h.services.AuthService.Token(
				r.Context(), &oauth2.Token{
					RefreshToken: refreshToken,
				},
			)
			if err != nil {
				return err
			}
			userInfo, err := h.services.AuthService.ValidateToken(r.Context(), token)
			// refresh token is invalid
			if err != nil {
				if errors.Is(err, h.services.AuthService.ErrorDecodeUserInfo()) {
					return err
				} else {
					return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorInvalidToken())
				}
			}
			// get user from database
			user, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
			if err != nil {
				return err
			}
			// get user session
			_, err = h.store.UserRepository().GetUserSession(user.ID)
			if err != nil {
				// user session doesn't exist, create it
				err = h.store.UserRepository().CreateUserSession(user.ID, token)
				if err != nil {
					return err
				}
			} else {
				// user session exists, update it
				err = h.store.UserRepository().UpdateUserSession(user.ID, token)
				if err != nil {
					return err
				}
			}
			headers := map[string]string{
				"access_token":  token.AccessToken,
				"refresh_token": token.RefreshToken,
			}
			utils.WriteJSON(w, http.StatusOK, user, headers)
			return nil
		}
	}

	// get user from db
	user, err := h.store.UserRepository().GetUserByEmail(userInfo.Email)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": refreshToken,
	}
	utils.WriteJSON(w, http.StatusOK, user, headers)
	return nil
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	var user *models.User
	accessToken := r.Header.Get("access_token")
	refreshToken := r.Header.Get("refresh_token")

	if accessToken == "" && refreshToken == "" {
		return utils.NewAPIError(
			http.StatusUnprocessableEntity, h.services.AuthService.ErrorMissingAuthorizationHeader(),
		)
	}

	if refreshToken != "" {
		userInfo, err := h.services.AuthService.ValidateToken(r.Context(), &oauth2.Token{RefreshToken: refreshToken})
		if err != nil {
			if errors.Is(err, h.services.AuthService.ErrorDecodeUserInfo()) {
				return err
			}
			return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorInvalidToken())
		}
		user, err = h.store.UserRepository().GetUserByEmail(userInfo.Email)
		if err != nil {
			return err
		}
		if err := h.services.AuthService.RevokeToken(refreshToken); err != nil {
			return err
		}
	} else {
		userInfo, err := h.services.AuthService.ValidateToken(r.Context(), &oauth2.Token{AccessToken: accessToken})
		if err != nil {
			if errors.Is(err, h.services.AuthService.ErrorDecodeUserInfo()) {
				return err
			}
			return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorInvalidToken())
		}
		user, err = h.store.UserRepository().GetUserByEmail(userInfo.Email)
		if err != nil {
			return err
		}
		if err := h.services.AuthService.RevokeToken(accessToken); err != nil {
			return err
		}
	}

	if err := h.store.UserRepository().DeleteUserSession(user.ID); err != nil {
		return err
	}
	utils.ClearSession(w, r, "verifier")
	utils.ClearSession(w, r, "state")
	utils.WriteJSON(w, http.StatusNoContent, nil, nil)
	return nil
}

func generateAppRedirectURL(appURL string, accessToken, refreshToken string) string {
	if refreshToken == "" {
		return fmt.Sprintf("%s?access_token=%s", appURL, accessToken)
	}
	return fmt.Sprintf("%s?access_token=%s&refresh_token=%s", appURL, accessToken, refreshToken)
}
