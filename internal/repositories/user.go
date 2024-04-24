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
	ErrorGetUser           = errors.New("failed to get user info")
	ErrorUpdateUser        = errors.New("failed to update user info")
)

type UserRepository interface {
	GetUserSession(userId uint) (*oauth2.Token, error)
	CreateUserSession(userID uint, token *oauth2.Token) error
	UpdateUserSession(userId uint, token *oauth2.Token) error
	DeleteUserSession(userId uint) error
	CreateUser(user *models.User) (uint, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(userId uint) (*models.User, error)
	UpdateUser(userId uint, user *models.UpdateUserRequest) error
}

type userRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewUserRepository(db *sqlx.DB, redis *redis.Client) UserRepository {
	return &userRepository{db: db, redis: redis}
}

func (ur *userRepository) GetUserSession(userId uint) (*oauth2.Token, error) {
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

func (ur *userRepository) CreateUserSession(userID uint, token *oauth2.Token) error {
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

func (ur *userRepository) UpdateUserSession(userId uint, token *oauth2.Token) error {
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

func (ur *userRepository) DeleteUserSession(userId uint) error {
	redisKey := fmt.Sprintf("user:%d", userId)
	err := ur.redis.Del(context.Background(), redisKey).Err()
	if err != nil {
		return ErrorDeleteUserSession
	}
	return nil
}

func (ur *userRepository) CreateUser(user *models.User) (uint, error) {
	var newUserId uint
	err := ur.db.QueryRow(
		"INSERT INTO users (email, given_name, family_name, name, locale, avatar_url, verified_email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		user.Email, user.GivenName, user.FamilyName, user.Name, user.Locale, user.AvatarURL, user.VerifiedEmail,
	).Scan(&newUserId)
	if err != nil {
		return 0, ErrorCreateUser
	}
	return newUserId, nil
}

func (ur *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return &models.User{}, ErrorGetUser
	}
	return &user, nil
}

func (ur *userRepository) GetUserById(userId uint) (*models.User, error) {
	var user models.User
	err := ur.db.Get(&user, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		return &models.User{}, ErrorGetUser
	}
	return &user, nil
}

func (ur *userRepository) UpdateUser(userId uint, user *models.UpdateUserRequest) error {
	_, err := ur.db.Exec(
		"UPDATE users SET stripe_customer_id = $1 WHERE id = $2",
		user.StripeCustomerID, userId,
	)
	if err != nil {
		return ErrorUpdateUser
	}
	return nil
}
