package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/repositories"
)

type Store interface {
	UserRepository() *repositories.UserRepository
	HealthRepository() *repositories.HealthRepository
	CategoryRepository() *repositories.CategoryRepository
}

type store struct {
	db                 *sqlx.DB
	redis              *redis.Client
	userRepository     *repositories.UserRepository
	healthRepository   *repositories.HealthRepository
	categoryRepository *repositories.CategoryRepository
	questionRepository *repositories.QuestionRepository
}

func New(db *sqlx.DB, redis *redis.Client) Store {
	return &store{
		db:    db,
		redis: redis,
	}

}

func (s *store) UserRepository() *repositories.UserRepository {
	if s.userRepository == nil {
		s.userRepository = repositories.NewUserRepository(s.db, s.redis)
	}
	return s.userRepository
}

func (s *store) HealthRepository() *repositories.HealthRepository {
	if s.healthRepository == nil {
		s.healthRepository = repositories.NewHealthRepository(s.db, s.redis)
	}
	return s.healthRepository
}

func (s *store) CategoryRepository() *repositories.CategoryRepository {
	if s.categoryRepository == nil {
		s.categoryRepository = repositories.NewCategoryRepository(s.db)
	}
	return s.categoryRepository
}

func (s *store) QuestionRepository() *repositories.QuestionRepository {
	if s.questionRepository == nil {
		s.questionRepository = repositories.NewQuestionRepository(s.db)
	}
	return s.questionRepository
}
