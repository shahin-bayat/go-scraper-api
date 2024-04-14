package config

import (
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var AppConf *AppConfig
var PostgresConf *PostgresConfig
var RedisConf *RedisConfig

type PostgresConfig struct {
	PgDatabase    string
	PgPassword    string
	PgUser        string
	PgPort        string
	PgHost        string
	PgInternalUrl string
}

type RedisConfig struct {
	RedisUser        string
	RedisPassword    string
	RedisHost        string
	RedisPort        string
	RedisDB          string
	RedisInternalUrl string
}

type AppConfig struct {
	Port                 int
	APIBaseURL           string
	AppUniversalURL      string
	StripePublishableKey string
	StripeWebhookSecret  string
	StripeSecretKey      string
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURL    string
	GoogleUserInfoURL    string
	GoogleRevokeURL      string
}

func init() {
	PostgresConf = &PostgresConfig{
		PgDatabase:    os.Getenv("PG_DATABASE"),
		PgPassword:    os.Getenv("PG_PASSWORD"),
		PgUser:        os.Getenv("PG_USER"),
		PgPort:        os.Getenv("PG_PORT"),
		PgHost:        os.Getenv("PG_HOST"),
		PgInternalUrl: os.Getenv("PG_INTERNAL_URL"),
	}
	RedisConf = &RedisConfig{
		RedisUser:        os.Getenv("REDIS_USER"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		RedisDB:          os.Getenv("REDIS_DB"),
		RedisInternalUrl: os.Getenv("REDIS_INTERNAL_URL"),
	}

	AppConf = &AppConfig{
		Port:                 getIntEnv("PORT"),
		APIBaseURL:           getStringEnv("API_BASE_URL"),
		AppUniversalURL:      getStringEnv("APP_UNIVERSAL_URL"),
		StripePublishableKey: getStringEnv("STRIPE_PUBLISHABLE_KEY"),
		StripeWebhookSecret:  getStringEnv("STRIPE_WEBHOOK_SECRET"),
		StripeSecretKey:      getStringEnv("STRIPE_SECRET_KEY"),
		GoogleClientID:       getStringEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   getStringEnv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:    getStringEnv("GOOGLE_REDIRECT_URL"),
		GoogleUserInfoURL:    getStringEnv("GOOGLE_USER_INFO_URL"),
		GoogleRevokeURL:      getStringEnv("GOOGLE_REVOKE_URL"),
	}

}

func getStringEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func getIntEnv(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Environment variable %s is not an integer", key)
	}
	return intValue
}
