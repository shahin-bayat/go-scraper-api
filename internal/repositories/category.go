package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/models"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (qr *CategoryRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := qr.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	return categories, nil
}
