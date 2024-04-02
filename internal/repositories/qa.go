package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/models"
)

type QARepository struct {
	db *sqlx.DB
}

func NewQARepository(db *sqlx.DB) *QARepository {
	return &QARepository{
		db: db,
	}
}

func (qar *QARepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := qar.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (qar *QARepository) GetCategoryDetail(categoryId int) ([]models.CategoryDetailResponse, error) {
	var categoryDetailResponse = make([]models.CategoryDetailResponse, 0)
	rows, err := qar.db.Queryx(`
			SELECT q.question_number, q.id FROM category_questions AS cq 
			JOIN questions AS q ON cq.question_id = q.id 
			WHERE category_id = $1 
			ORDER BY q.id ASC
	`, categoryId)
	if err != nil {
		return nil, fmt.Errorf("error getting category detail: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cdr models.CategoryDetailResponse
		err := rows.StructScan(&cdr)
		if err != nil {
			return nil, fmt.Errorf("error getting category detail: %w", err)
		}
		categoryDetailResponse = append(categoryDetailResponse, cdr)
	}
	return categoryDetailResponse, nil
}

func (qar *QARepository) GetQuestionDetail(categoryId int, questionId int) (models.QuestionDetailResponse, error) {
	var questionDetailResponse models.QuestionDetailResponse
	err := qar.db.Get(&questionDetailResponse, `
			SELECT q.question_number, i.extracted_text, i.has_image, i.file_name 
			FROM category_questions as cq 
			JOIN questions AS q ON cq.question_id = q.id 
			JOIN images AS i ON i.question_id = q.id 
			WHERE cq.category_id = $1 AND cq.question_id = $2
	`, categoryId, questionId)

	if err != nil {
		return models.QuestionDetailResponse{}, fmt.Errorf("error getting question detail: %w", err)
	}

	var answers []models.Answer
	err = qar.db.Select(&answers, `
        SELECT id, question_id, text, is_correct, created_at, updated_at, deleted_at
        FROM answers
        WHERE question_id = $1
    `, questionId)
	if err != nil {
		return models.QuestionDetailResponse{}, fmt.Errorf("error getting answers: %w", err)
	}

	questionDetailResponse.Answers = answers

	return questionDetailResponse, nil
}
