package repositories

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/models"
)

type QuestionRepository struct {
	db *sqlx.DB
}

func NewQuestionRepository(db *sqlx.DB) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (qar *QuestionRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := qar.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (qar *QuestionRepository) GetCategoryDetail(categoryId int) ([]models.CategoryDetailResponse, error) {
	var categoryDetailResponse = make([]models.CategoryDetailResponse, 0)
	rows, err := qar.db.Queryx(`
			SELECT q.question_number, q.id 
			FROM category_questions AS cq 
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

func (qar *QuestionRepository) GetQuestionDetail(questionId int, lang string) (models.QuestionDetailResponse, error) {
	var apiBaseUrl = os.Getenv("API_BASE_URL")
	var questionTranslation models.Translation
	var answersTranslation []models.Translation
	var response models.QuestionDetailResponse

	err := qar.db.Get(&response, `
			SELECT q.question_number, i.extracted_text, i.has_image, i.file_name 
			FROM questions AS q
			JOIN images AS i ON i.question_id = q.id
			WHERE q.id = $1
	`, questionId)

	if err != nil {
		return models.QuestionDetailResponse{}, fmt.Errorf("error getting question detail: %w", err)
	}

	response.FileURL = fmt.Sprintf("%s/image/%s", apiBaseUrl, response.Filename)

	var answers []models.Answer
	err = qar.db.Select(&answers, `
				SELECT id, question_id, text, is_correct, created_at, updated_at, deleted_at
				FROM answers
				WHERE question_id = $1
		`, questionId)
	if err != nil {
		return models.QuestionDetailResponse{}, fmt.Errorf("error getting answers: %w", err)
	}

	response.Answers = answers

	if lang != "" {
		err = qar.db.Get(&questionTranslation, `
			SELECT * from translations
			WHERE refer_id = $1 AND type = $2 AND lang = $3
		`, questionId, models.QuestionType, lang)
		if err != nil {
			return models.QuestionDetailResponse{}, fmt.Errorf("error getting question translation: %w", err)
		}
		err = qar.db.Select(&answersTranslation, `
			SELECT * from translations
			WHERE refer_id IN ($1, $2, $3, $4) AND type = $5 AND lang = $6
		`, answers[0].ID, answers[1].ID, answers[2].ID, answers[3].ID, models.AnswerType, lang)
		if err != nil {
			return models.QuestionDetailResponse{}, fmt.Errorf("error getting answers translation: %w", err)
		}

		response.Question = questionTranslation.Translation

		for i, answer := range response.Answers {
			for _, translation := range answersTranslation {
				if uint(answer.ID) == uint(translation.ReferID) {
					fmt.Printf("answer id: %d, translation refer id: %d\n", answer.ID, translation.ReferID)
					fmt.Printf("answer text: %s, translation text: %s\n", answer.Text, translation.Translation)
					response.Answers[i].Text = translation.Translation
				}
			}
		}
	}

	return response, nil
}
