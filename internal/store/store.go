package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/repositories"
)

type Store interface {
	UserRepository() repositories.UserRepository
	HealthRepository() repositories.HealthRepository
	QuestionRepository() repositories.QuestionRepository
	SubscriptionRepository() repositories.SubscriptionRepository
}

type store struct {
	db                     *sqlx.DB
	redis                  *redis.Client
	userRepository         repositories.UserRepository
	healthRepository       repositories.HealthRepository
	questionRepository     repositories.QuestionRepository
	subscriptionRepository repositories.SubscriptionRepository
}

func New(db *sqlx.DB, redis *redis.Client) Store {
	return &store{
		db:    db,
		redis: redis,
	}

}

func (s *store) UserRepository() repositories.UserRepository {
	if s.userRepository == nil {
		s.userRepository = repositories.NewUserRepository(s.db, s.redis)
	}
	return s.userRepository
}

func (s *store) HealthRepository() repositories.HealthRepository {
	if s.healthRepository == nil {
		s.healthRepository = repositories.NewHealthRepository(s.db, s.redis)
	}
	return s.healthRepository
}

func (s *store) QuestionRepository() repositories.QuestionRepository {
	if s.questionRepository == nil {
		s.questionRepository = repositories.NewQuestionRepository(s.db)
	}
	return s.questionRepository
}

func (s *store) SubscriptionRepository() repositories.SubscriptionRepository {
	if s.subscriptionRepository == nil {
		s.subscriptionRepository = repositories.NewSubscriptionRepository(s.db)
	}
	return s.subscriptionRepository
}
