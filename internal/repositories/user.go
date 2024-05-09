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
	ErrorMissingToken = errors.New("token is missing")
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
		return nil, err
	}
	rt, err := ur.redis.HGet(context.Background(), fmt.Sprintf("user:%d", userId), "refresh_token").Result()
	if err != nil {
		return nil, err
	}
	eStr, err := ur.redis.HGet(context.Background(), fmt.Sprintf("user:%d", userId), "expiry").Result()
	if err != nil {
		return nil, err
	}
	e, err := time.Parse(time.RFC3339, eStr)
	if err != nil {
		return nil, err
	}

	tt, err := ur.redis.HGet(context.Background(), fmt.Sprintf("user:%d", userId), "token_type").Result()
	if err != nil {
		return nil, err
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

	if err := ur.redis.HSet(context.Background(), redisKey, "access_token", token.AccessToken).Err(); err != nil {
		return err
	}
	if err := ur.redis.HSet(context.Background(), redisKey, "refresh_token", token.RefreshToken).Err(); err != nil {
		return err
	}
	if err := ur.redis.HSet(context.Background(), redisKey, "expiry", token.Expiry).Err(); err != nil {
		return err
	}
	if err := ur.redis.HSet(context.Background(), redisKey, "token_type", token.TokenType).Err(); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUserSession(userId uint, token *oauth2.Token) error {
	if token == nil {
		return ErrorMissingToken
	}
	redisKey := fmt.Sprintf("user:%d", userId)
	if err := ur.redis.HSet(context.Background(), redisKey, "access_token", token.AccessToken).Err(); err != nil {
		return err
	}
	if err := ur.redis.HSet(context.Background(), redisKey, "expiry", token.Expiry).Err(); err != nil {
		return err
	}
	if err := ur.redis.HSet(context.Background(), redisKey, "token_type", token.TokenType).Err(); err != nil {
		return err
	}
	if token.RefreshToken != "" {
		if err := ur.redis.HSet(context.Background(), redisKey, "refresh_token", token.RefreshToken).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (ur *userRepository) DeleteUserSession(userId uint) error {
	redisKey := fmt.Sprintf("user:%d", userId)
	if err := ur.redis.Del(context.Background(), redisKey).Err(); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *models.User) (uint, error) {
	var newUserId uint
	if err := ur.db.QueryRow(
		"INSERT INTO users (email, given_name, family_name, name, locale, avatar_url, verified_email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		user.Email, user.GivenName, user.FamilyName, user.Name, user.Locale, user.AvatarURL, user.VerifiedEmail,
	).Scan(&newUserId); err != nil {
		return 0, err
	}
	return newUserId, nil
}

func (ur *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Get(&user, "SELECT * FROM users WHERE email = $1", email); err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (ur *userRepository) GetUserById(userId uint) (*models.User, error) {
	var user models.User
	if err := ur.db.Get(&user, "SELECT * FROM users WHERE id = $1", userId); err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (ur *userRepository) UpdateUser(userId uint, user *models.UpdateUserRequest) error {
	if _, err := ur.db.Exec(
		"UPDATE users SET stripe_customer_id = $1 WHERE id = $2",
		user.StripeCustomerID, userId,
	); err != nil {
		return err
	}
	return nil
}
