package middlewares

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	ErrorUserNotAuthorized = errors.New("user is not authorized")
)

type contextKey string

const (
	User         contextKey = "user"
	Subscription contextKey = "subscription"
)

func (am *Middlewares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accessToken := r.Header.Get("access_token")
		refreshToken := r.Header.Get("refresh_token")

		if accessToken == "" && refreshToken == "" {
			ctx = context.WithValue(ctx, User, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		token, err := am.authService.Token(r.Context(), &oauth2.Token{AccessToken: accessToken})
		if err != nil {
			ctx = context.WithValue(ctx, User, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userInfo, err := am.authService.ValidateToken(r.Context(), token)
		// access token is invalid
		if err != nil {
			if err == am.authService.ErrorDecodeUserInfo() {
				ctx = context.WithValue(ctx, User, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// refresh the token
			if refreshToken == "" {
				ctx = context.WithValue(ctx, User, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			token, err = am.authService.Token(r.Context(), &oauth2.Token{
				RefreshToken: refreshToken,
			})
			if err != nil {
				ctx = context.WithValue(ctx, User, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			userInfo, err = am.authService.ValidateToken(r.Context(), token)
			// refresh token is invalid
			if err != nil {
				ctx = context.WithValue(ctx, User, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			w.Header().Add("access_token", token.AccessToken)
			w.Header().Add("refresh_token", token.RefreshToken)

			user, err := am.store.UserRepository().GetUserByEmail(userInfo.Email)
			if err != nil {
				ctx = context.WithValue(ctx, User, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			_, err = am.store.UserRepository().GetUserSession(user.ID)
			if err != nil {
				err = am.store.UserRepository().CreateUserSession(user.ID, token)
				if err != nil {
					ctx = context.WithValue(ctx, User, nil)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				ctx = context.WithValue(ctx, User, user.ID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			err = am.store.UserRepository().UpdateUserSession(user.ID, token)
			if err != nil {
				ctx = context.WithValue(ctx, User, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			ctx = context.WithValue(ctx, User, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		user, err := am.store.UserRepository().GetUserByEmail(userInfo.Email)
		if err != nil {
			ctx = context.WithValue(ctx, User, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, User, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIdFromContext(ctx context.Context) (uint, error) {
	if ctx.Value(User) == nil {
		return 0, ErrorUserNotAuthorized
	}
	userId := ctx.Value(User).(uint)
	return userId, nil
}
