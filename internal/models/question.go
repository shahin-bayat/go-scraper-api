package models

import (
	"time"
)

type Category struct {
	ID          uint       `json:"id" db:"id"`
	Text        string     `json:"text" db:"text"`
	CategoryKey string     `json:"-" db:"category_key"`
	CreatedAt   time.Time  `json:"-"  db:"created_at"`
	UpdatedAt   time.Time  `json:"-" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
}

type Question struct {
	ID             uint   `json:"id" db:"id"`
	QuestionNumber string `json:"question_number" db:"question_number"`
	QuestionKey    string `json:"-" db:"question_key"`
	CreatedAt      string `json:"-" db:"created_at"`
	UpdatedAt      string `json:"-" db:"updated_at"`
	DeletedAt      string `json:"-" db:"deleted_at"`
}

type Answer struct {
	ID         uint       `json:"id" db:"id"`
	QuestionID uint       `json:"-" db:"question_id"`
	Text       string     `json:"answer" db:"text"`
	IsCorrect  bool       `json:"is_correct" db:"is_correct"`
	CreatedAt  time.Time  `json:"-" db:"created_at"`
	UpdatedAt  time.Time  `json:"-" db:"updated_at"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`
}

type CategoryDetailResponse struct {
	QuestionNumber string `json:"question_number" db:"question_number"`
	QuestionId     uint   `json:"question_id" db:"id"`
}

type QuestionDetailResponse struct {
	QuestionNumber string   `json:"question_number" db:"question_number"`
	Question       string   `json:"question" db:"extracted_text"`
	HasImage       bool     `json:"has_image" db:"has_image"`
	Filename       string   `json:"-" db:"file_name"`
	FileURL        string   `json:"file_url"`
	Answers        []Answer `json:"answers"`
	IsBookmarked   bool     `json:"is_bookmarked" db:"is_bookmarked"`
}

type Bookmark struct {
	ID         uint      `db:"id"`
	UserID     uint      `db:"user_id"`
	QuestionID uint      `db:"question_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type BookmarkResponse struct {
	QuestionNumber string `json:"question_number" db:"question_number"`
	QuestionID     uint   `json:"question_id" db:"question_id"`
}

type BookmarkRequest struct {
	QuestionId uint `json:"question_id"`
}

func (r *BookmarkRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if r.QuestionId == 0 {
		errors["question_id"] = "question_id is required"
	}
	return errors
}
