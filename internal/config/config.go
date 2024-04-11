package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type EnvError struct {
	Key     string
	Message string
}

func (e *EnvError) Error() string {
	return fmt.Sprintf("environment variable '%s': %s", e.Key, e.Message)
}

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

type config struct {
	Postgres *PostgresConfig
	Redis    *RedisConfig
}

var Config *config

func init() {
	postgresConfig := &PostgresConfig{
		PgDatabase:    os.Getenv("PG_DATABASE"),
		PgPassword:    os.Getenv("PG_PASSWORD"),
		PgUser:        os.Getenv("PG_USER"),
		PgPort:        os.Getenv("PG_PORT"),
		PgHost:        os.Getenv("PG_HOST"),
		PgInternalUrl: os.Getenv("PG_INTERNAL_URL"),
	}
	redisConfig := &RedisConfig{
		RedisUser:        os.Getenv("REDIS_USER"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		RedisDB:          os.Getenv("REDIS_DB"),
		RedisInternalUrl: os.Getenv("REDIS_INTERNAL_URL"),
	}

	Config = &config{
		Postgres: postgresConfig,
		Redis:    redisConfig,
	}

}

// func getStringEnv(key string) (string, error) {
// 	value := os.Getenv(key)
// 	if value == "" {
// 		return "", &EnvError{Key: key, Message: "environment variable is not set"}
// 	}
// 	return value, nil
// }
