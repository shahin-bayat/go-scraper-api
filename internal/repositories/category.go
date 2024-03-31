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

func (cr *CategoryRepository) GetCategoryDetail(categoryId int) ([]models.CategoryDetailResponse, error) {
	var categoryDetailResponse = make([]models.CategoryDetailResponse, 0)
	rows, err := cr.db.Queryx("SELECT q.question_number, q.id FROM category_questions AS cq JOIN questions AS q ON cq.question_id = q.id WHERE category_id = $1 ORDER BY q.id ASC", categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cdr models.CategoryDetailResponse
		err := rows.StructScan(&cdr)
		if err != nil {
			return nil, err
		}
		categoryDetailResponse = append(categoryDetailResponse, cdr)
	}
	return categoryDetailResponse, nil
}
