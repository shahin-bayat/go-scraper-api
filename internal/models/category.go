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

type CategoryDetailResponse struct {
	QuestionsCount int `json:"questions_count"`
}
