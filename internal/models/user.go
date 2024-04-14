package models

import "time"

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
	Locale        string `json:"locale"`
	AvatarURL     string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type User struct {
	ID            uint       `json:"-" db:"id"`
	Email         string     `json:"email" db:"email"`
	GivenName     string     `json:"given_name" db:"given_name"`
	FamilyName    string     `json:"family_name" db:"family_name"`
	Name          string     `json:"name" db:"name"`
	Locale        string     `json:"locale" db:"locale"`
	AvatarURL     string     `json:"avatar_url" db:"avatar_url"`
	VerifiedEmail bool       `json:"verified_email" db:"verified_email"`
	CreatedAt     time.Time  `json:"-" db:"created_at"`
	UpdatedAt     time.Time  `json:"-" db:"updated_at"`
	DeletedAt     *time.Time `json:"-" db:"deleted_at"`
}
