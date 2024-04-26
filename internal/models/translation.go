package models

import "time"

// CREATE TYPE type AS ENUM ('question', 'answer');
// CREATE TYPE lang AS ENUM ('en');
// CREATE TABLE IF NOT EXISTS translations (
//   id SERIAL PRIMARY KEY,
//   refer_id INTEGER NOT NULL,
//   type type NOT NULL,
//   lang lang NOT NULL,
//   translation TEXT NOT NULL
// )

// how to define enum in go
// https://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-go

const (
	QuestionType = "question"
	AnswerType   = "answer"
)

type Translation struct {
	ID          uint       `db:"id" json:"-"`
	ReferID     uint       `db:"refer_id" json:"-"`
	Type        string     `db:"type" json:"-"`
	Lang        string     `db:"lang" json:"-"`
	Translation string     `db:"translation" json:"translation"`
	CreatedAt   time.Time  `db:"created_at" json:"-"`
	UpdatedAt   time.Time  `db:"updated_at" json:"-"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`
}
