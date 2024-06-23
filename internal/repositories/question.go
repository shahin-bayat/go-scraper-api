package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/models"
)

type QuestionRepository interface {
	GetCategories() ([]models.Category, error)
	GetCategoryDetail(categoryId uint, questionType string) ([]models.CategoryDetailResponse, error)
	GetFreeCategoryDetail(categoryId uint, freeQuestionIds [3]uint) ([]models.CategoryDetailResponse, error)
	GetQuestionDetail(questionId, userId uint, lang string, apiBaseUrl string) (models.QuestionDetailResponse, error)
	BookmarkQuestion(questionId uint, userId uint) (uint, error)
	GetBookmarks(userId uint) ([]models.BookmarkResponse, error)
}
type questionRepository struct {
	db *sqlx.DB
}

func NewQuestionRepository(db *sqlx.DB) QuestionRepository {
	return &questionRepository{
		db: db,
	}
}

func (qr *questionRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := qr.db.Select(&categories, "SELECT * FROM categories"); err != nil {
		return nil, err
	}
	return categories, nil
}

func (qr *questionRepository) GetCategoryDetail(categoryId uint, questionType string) ([]models.CategoryDetailResponse, error) {
	var categoryDetailResponse = make([]models.CategoryDetailResponse, 0)
	var query string
	if questionType == "image" {
		query = `SELECT q.question_number, q.id 
			FROM category_questions AS cq 
			JOIN questions AS q ON cq.question_id = q.id
			JOIN images AS i ON i.question_id = q.id
			WHERE category_id = $1 AND i.has_image = true
			ORDER BY q.id`
	} else {
		query = `SELECT q.question_number, q.id 
			FROM category_questions AS cq 
			JOIN questions AS q ON cq.question_id = q.id 
			WHERE category_id = $1 
			ORDER BY q.id`
	}

	rows, err := qr.db.Queryx(query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cdr models.CategoryDetailResponse
		if err := rows.StructScan(&cdr); err != nil {
			return nil, err
		}
		categoryDetailResponse = append(categoryDetailResponse, cdr)
	}
	return categoryDetailResponse, nil
}

func (qr *questionRepository) GetFreeCategoryDetail(categoryId uint, freeQuestionIds [3]uint) ([]models.CategoryDetailResponse, error) {
	var categoryDetailResponse = make([]models.CategoryDetailResponse, 0)
	rows, err := qr.db.Queryx(
		`
			SELECT q.question_number, q.id 
			FROM category_questions AS cq 
			JOIN questions AS q ON cq.question_id = q.id 
			WHERE category_id = $1 AND q.id IN ($2, $3, $4)
			ORDER BY q.id
			`, categoryId, freeQuestionIds[0], freeQuestionIds[1], freeQuestionIds[2],
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cdr models.CategoryDetailResponse
		if err := rows.StructScan(&cdr); err != nil {
			return nil, err
		}
		categoryDetailResponse = append(categoryDetailResponse, cdr)
	}
	return categoryDetailResponse, nil
}

func (qr *questionRepository) GetQuestionDetail(questionId, userId uint, lang string, apiBaseUrl string) (models.QuestionDetailResponse, error) {
	var questionTranslation models.Translation
	var answersTranslation []models.Translation
	var response models.QuestionDetailResponse

	if err := qr.db.Get(
		&response, `
			SELECT q.question_number, i.extracted_text, i.has_image, i.file_name, EXISTS(SELECT 1 FROM bookmarks WHERE question_id = $1 AND user_id = $2) AS is_bookmarked  
			FROM questions AS q
			JOIN images AS i ON i.question_id = q.id
			WHERE q.id = $1
			`, questionId, userId,
	); err != nil {
		return models.QuestionDetailResponse{}, err
	}
	response.FileURL = fmt.Sprintf("%s/image/%s", apiBaseUrl, response.Filename)

	var answers []models.Answer
	if err := qr.db.Select(
		&answers, `
				SELECT id, question_id, text, is_correct, created_at, updated_at, deleted_at
				FROM answers
				WHERE question_id = $1
				`, questionId,
	); err != nil {
		return models.QuestionDetailResponse{}, err
	}
	response.Answers = answers

	if lang != "" {
		if err := qr.db.Get(
			&questionTranslation, `
			SELECT * from translations
			WHERE refer_id = $1 AND type = $2 AND lang = $3
		`, questionId, models.QuestionType, lang,
		); err != nil {
			return models.QuestionDetailResponse{}, err
		}
		if err := qr.db.Select(
			&answersTranslation, `
			SELECT * from translations
			WHERE refer_id IN ($1, $2, $3, $4) AND type = $5 AND lang = $6
		`, answers[0].ID, answers[1].ID, answers[2].ID, answers[3].ID, models.AnswerType, lang,
		); err != nil {
			return models.QuestionDetailResponse{}, err
		}
		response.Question = questionTranslation.Translation
		for i, answer := range response.Answers {
			for _, translation := range answersTranslation {
				if answer.ID == translation.ReferID {
					response.Answers[i].Text = translation.Translation
				}
			}
		}
	}
	return response, nil
}

func (qr *questionRepository) BookmarkQuestion(questionId uint, userId uint) (uint, error) {
	var bookmark models.Bookmark
	err := qr.db.Get(
		&bookmark, `SELECT * from bookmarks WHERE user_id = $1 AND question_id = $2`, userId, questionId,
	)

	if err != nil && err.Error() == "sql: no rows in result set" {
		var bookmarkId uint
		if err := qr.db.QueryRow(
			`
			INSERT INTO bookmarks (user_id, question_id) VALUES ($1, $2) RETURNING id
			`, userId, questionId,
		).Scan(&bookmarkId); err != nil {
			return 0, err
		}
		return bookmarkId, nil
	} else {
		if _, err := qr.db.Exec(
			"DELETE FROM bookmarks WHERE user_id = $1 AND question_id = $2", userId, questionId,
		); err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func (qr *questionRepository) GetBookmarks(userId uint) ([]models.BookmarkResponse, error) {
	var bookmarks []models.BookmarkResponse
	if err := qr.db.Select(
		&bookmarks, `
				SELECT q.question_number, q.id AS question_id FROM bookmarks AS b
				JOIN questions AS q ON b.question_id = q.id
				WHERE b.user_id = $1
				ORDER BY q.id
				`, userId,
	); err != nil {
		return nil, err
	}
	return bookmarks, nil
}
