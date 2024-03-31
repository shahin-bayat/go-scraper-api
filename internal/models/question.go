package models

type Question struct {
	ID             uint   `json:"id" db:"id"`
	QuestionNumber string `json:"question_number" db:"question_number"`
	QuestionKey    string `json:"-" db:"question_key"`
	CreatedAt      string `json:"-" db:"created_at"`
	UpdatedAt      string `json:"-" db:"updated_at"`
	DeletedAt      string `json:"-" db:"deleted_at"`
}
