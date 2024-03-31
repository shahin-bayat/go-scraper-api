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

func (cr *CategoryRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := cr.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (cr *CategoryRepository) GetCategoryDetail(categoryId int) (models.CategoryDetailResponse, error) {
	var CategoryResponse models.CategoryDetailResponse
	err := cr.db.Get(&CategoryResponse.QuestionsCount, "SELECT Count(*) FROM category_questions WHERE category_id = $1", categoryId)
	if err != nil {
		return models.CategoryDetailResponse{}, err
	}
	return CategoryResponse, nil

}
