package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"golang.org/x/oauth2"
)

var (
	ErrorCreateUser        = errors.New("failed to create user")
	ErrorCreateUserSession = errors.New("failed to create user session")
	ErrorMissingToken      = errors.New("token is missing")
	ErrorGetUserSession    = errors.New("failed to get user session")
	ErrorUpdateUserSession = errors.New("failed to update user session")
	ErrorDeleteUserSession = errors.New("failed to delete user session")
	ErrorParseTokenExpiry  = errors.New("failed to parse token expiry time")
	ErrorGetUserByEmail    = errors.New("failed to get user by email")
)

type UserRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewUserRepository(db *sqlx.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{db: db, redis: redis}
}

func (ur *UserRepository) GetUserSession(userId uint) (*oauth2.Token, error) {
	at, err := ur.redis.HGet(context.Background(), fmt.Sprintf("user:%d", userId), "access_token").Result()
	if err != nil {
		return nil, ErrorGetUserSession
	}
	rt, err := ur.redis.HGet(context.Background(), fmt.Sprintf("user:%d", userId), "refresh_token").Result()
	if err != nil {
		return nil, ErrorGetUserSession
	}
	eStr, err := ur.redis.HGet(context.Background(), fmt.Sprintf("user:%d", userId), "expiry").Result()
	if err != nil {
		return nil, ErrorGetUserSession
	}
	e, err := time.Parse(time.RFC3339, eStr)
	if err != nil {
		return nil, ErrorParseTokenExpiry
	}

	tt, err := ur.redis.HGet(context.Background(), fmt.Sprintf("user:%d", userId), "token_type").Result()
	if err != nil {
		return nil, ErrorGetUserSession
	}

	token := &oauth2.Token{
		AccessToken:  at,
		RefreshToken: rt,
		Expiry:       e,
		TokenType:    tt,
	}
	return token, nil
}

func (ur *UserRepository) CreateUserSession(userID uint, token *oauth2.Token) error {
	if token == nil {
		return ErrorMissingToken
	}
	redisKey := fmt.Sprintf("user:%d", userID)

	err := ur.redis.HSet(context.Background(), redisKey, "access_token", token.AccessToken).Err()
	if err != nil {
		return ErrorCreateUserSession
	}
	err = ur.redis.HSet(context.Background(), redisKey, "refresh_token", token.RefreshToken).Err()
	if err != nil {
		return ErrorCreateUserSession
	}
	err = ur.redis.HSet(context.Background(), redisKey, "expiry", token.Expiry).Err()
	if err != nil {
		return ErrorCreateUserSession
	}
	err = ur.redis.HSet(context.Background(), redisKey, "token_type", token.TokenType).Err()
	if err != nil {
		return ErrorCreateUserSession
	}

	return nil
}

func (ur *UserRepository) UpdateUserSession(userId uint, token *oauth2.Token) error {
	if token == nil {
		return ErrorMissingToken
	}
	redisKey := fmt.Sprintf("user:%d", userId)
	err := ur.redis.HSet(context.Background(), redisKey, "access_token", token.AccessToken).Err()
	if err != nil {
		return ErrorUpdateUserSession
	}
	err = ur.redis.HSet(context.Background(), redisKey, "expiry", token.Expiry).Err()
	if err != nil {
		return ErrorUpdateUserSession
	}
	err = ur.redis.HSet(context.Background(), redisKey, "token_type", token.TokenType).Err()
	if err != nil {
		return ErrorUpdateUserSession
	}

	if token.RefreshToken != "" {
		err = ur.redis.HSet(context.Background(), redisKey, "refresh_token", token.RefreshToken).Err()
		if err != nil {
			return ErrorUpdateUserSession
		}
	}
	return nil
}

func (ur *UserRepository) DeleteUserSession(userId uint) error {
	redisKey := fmt.Sprintf("user:%d", userId)
	err := ur.redis.Del(context.Background(), redisKey).Err()
	if err != nil {
		return ErrorDeleteUserSession
	}
	return nil
}

func (ur *UserRepository) CreateUser(user *models.User) (uint, error) {
	var newUserId uint
	err := ur.db.QueryRow("INSERT INTO users (email, given_name, family_name, name, locale, avatar_url, verified_email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", user.Email, user.GivenName, user.FamilyName, user.Name, user.Locale, user.AvatarURL, user.VerifiedEmail).Scan(&newUserId)
	if err != nil {
		return 0, ErrorCreateUser
	}
	return newUserId, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return &models.User{}, ErrorGetUserByEmail
	}
	return &user, nil
}
