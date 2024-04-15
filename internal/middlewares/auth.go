package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
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

		token, err := am.authService.TokenSource(r.Context(), &oauth2.Token{AccessToken: accessToken, RefreshToken: refreshToken}).Token()
		if err != nil {
			ctx = context.WithValue(ctx, User, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if !am.authService.TokenValid(token) {
			token, err = am.authService.TokenSource(r.Context(), &oauth2.Token{
				RefreshToken: refreshToken,
			}).Token()
			if err != nil {
				ctx = context.WithValue(ctx, User, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		userInfo, err := am.authService.GetUserInfo(r.Context(), token)
		if err != nil {
			ctx = context.WithValue(ctx, User, nil)
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
		return 0, fmt.Errorf("user is not authorized")
	}
	userId := ctx.Value(User).(uint)
	return userId, nil
}
