package repositories

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/models"
)

var (
	ErrorGetCategories          = errors.New("error getting categories")
	ErrorGetCategoryDetail      = errors.New("error getting category detail")
	ErrorGetQuestionDetail      = errors.New("error getting question detail")
	ErrorGetQuestionAnswers     = errors.New("error getting question answers")
	ErrorGetAnswersTranslations = errors.New("error getting answers translations")
	ErrorMissingCategoryId      = errors.New("category id is required")
	ErrorUnsupportedLanguage    = errors.New("language not supported")
	ErrorMissingQuestionId      = errors.New("question id is required")
	ErrorMissingFilename        = errors.New("filename is required")
	ErrorFileNotFound           = errors.New("file not found")
	ErrorBookmarkQuestion       = errors.New("error bookmarking question")
)

type QuestionRepository interface {
	GetCategories() ([]models.Category, error)
	GetCategoryDetail(categoryId uint) ([]models.CategoryDetailResponse, error)
	GetFreeCategoryDetail(categoryId uint, freeQuestionIds [3]uint) ([]models.CategoryDetailResponse, error)
	GetQuestionDetail(questionId uint, lang string, apiBaseUrl string) (models.QuestionDetailResponse, error)
	BookmarkQuestion(questionId uint, userId uint) (uint, error)
	ErrorMissingCategoryId() error
	ErrorUnsupportedLanguage() error
	ErrorMissingQuestionId() error
	ErrorMissingFilename() error
	ErrorFileNotFound() error
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
	err := qr.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return nil, ErrorGetCategories
	}
	return categories, nil
}

func (qr *questionRepository) GetCategoryDetail(categoryId uint) ([]models.CategoryDetailResponse, error) {
	var categoryDetailResponse = make([]models.CategoryDetailResponse, 0)
	rows, err := qr.db.Queryx(
		`
			SELECT q.question_number, q.id 
			FROM category_questions AS cq 
			JOIN questions AS q ON cq.question_id = q.id 
			WHERE category_id = $1 
			ORDER BY q.id ASC
	`, categoryId,
	)
	if err != nil {
		return nil, ErrorGetCategoryDetail
	}
	defer rows.Close()

	for rows.Next() {
		var cdr models.CategoryDetailResponse
		err := rows.StructScan(&cdr)
		if err != nil {
			return nil, ErrorGetCategoryDetail
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
			ORDER BY q.id ASC
	`, categoryId, freeQuestionIds[0], freeQuestionIds[1], freeQuestionIds[2],
	)
	if err != nil {
		return nil, ErrorGetCategoryDetail
	}
	defer rows.Close()

	for rows.Next() {
		var cdr models.CategoryDetailResponse
		err := rows.StructScan(&cdr)
		if err != nil {
			return nil, ErrorGetCategoryDetail
		}
		categoryDetailResponse = append(categoryDetailResponse, cdr)
	}
	return categoryDetailResponse, nil
}

func (qr *questionRepository) GetQuestionDetail(questionId uint, lang string, apiBaseUrl string) (models.QuestionDetailResponse, error) {

	var questionTranslation models.Translation
	var answersTranslation []models.Translation
	var response models.QuestionDetailResponse

	err := qr.db.Get(
		&response, `
			SELECT q.question_number, i.extracted_text, i.has_image, i.file_name 
			FROM questions AS q
			JOIN images AS i ON i.question_id = q.id
			WHERE q.id = $1
	`, questionId,
	)

	if err != nil {
		return models.QuestionDetailResponse{}, ErrorGetQuestionDetail
	}

	response.FileURL = fmt.Sprintf("%s/image/%s", apiBaseUrl, response.Filename)

	var answers []models.Answer
	err = qr.db.Select(
		&answers, `
				SELECT id, question_id, text, is_correct, created_at, updated_at, deleted_at
				FROM answers
				WHERE question_id = $1
		`, questionId,
	)
	if err != nil {
		return models.QuestionDetailResponse{}, ErrorGetQuestionAnswers
	}

	response.Answers = answers

	if lang != "" {
		err = qr.db.Get(
			&questionTranslation, `
			SELECT * from translations
			WHERE refer_id = $1 AND type = $2 AND lang = $3
		`, questionId, models.QuestionType, lang,
		)
		if err != nil {
			return models.QuestionDetailResponse{}, ErrorGetAnswersTranslations
		}
		err = qr.db.Select(
			&answersTranslation, `
			SELECT * from translations
			WHERE refer_id IN ($1, $2, $3, $4) AND type = $5 AND lang = $6
		`, answers[0].ID, answers[1].ID, answers[2].ID, answers[3].ID, models.AnswerType, lang,
		)
		if err != nil {
			return models.QuestionDetailResponse{}, ErrorGetAnswersTranslations
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
		err := qr.db.QueryRow(
			`
			INSERT INTO bookmarks (user_id, question_id) VALUES ($1, $2) RETURNING id
	`, userId, questionId,
		).Scan(&bookmarkId)
		if err != nil {
			return 0, ErrorBookmarkQuestion
		}
		return bookmarkId, nil
	} else {
		_, err := qr.db.Exec("DELETE FROM bookmarks WHERE user_id = $1 AND question_id = $2", userId, questionId)
		if err != nil {
			return 0, ErrorBookmarkQuestion
		}
	}
	return 0, nil
}

func (qr *questionRepository) ErrorMissingCategoryId() error {
	return ErrorMissingCategoryId
}

func (qr *questionRepository) ErrorUnsupportedLanguage() error {
	return ErrorUnsupportedLanguage
}
func (qr *questionRepository) ErrorMissingQuestionId() error {
	return ErrorMissingQuestionId
}

func (qr *questionRepository) ErrorMissingFilename() error {
	return ErrorMissingFilename
}

func (qr *questionRepository) ErrorFileNotFound() error {
	return ErrorFileNotFound
}
