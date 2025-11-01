package entity

import "time"

type User struct {
	Id                   int        `json:"id" db:"id"`
	Email                string     `json:"email" db:"email" validate:"required,email,max=128"`
	Full_name            string     `json:"full_name" db:"full_name" validate:"required,min=2,max=50"`
	Username             string     `json:"username" db:"username" validate:"required,alphanum,min=3,max=20"`
	Password             string     `json:"password" db:"password" validate:"required,min=8"`
	Phone                string     `json:"phone" db:"phone" validate:"required,numeric,min=8,max=15"`
	Role                 string     `json:"role" db:"role" validate:"omitempty,oneof=customer admin"`
	Last_password_change *time.Time `json:"last_password_change" db:"last_password_change" validate:"omitempty"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at" validate:"omitempty"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at" validate:"omitempty"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" db:"deleted_at" validate:"omitempty"`
}

type Address struct {
	Id            int        `json:"id" db:"id"`
	UserId        int        `json:"user_id" db:"user_id" validate:"required"`
	RecipientName string     `json:"recipient_name" db:"recipient_name" validate:"omitempty,max=200"`
	Phone         string     `json:"phone" db:"phone" validate:"omitempty,max=32"`
	Street        string     `json:"street" db:"street" validate:"omitempty,max=32"`
	City          string     `json:"city" db:"city" validate:"omitempty,max=120"`
	Province      string     `json:"province" db:"province" validate:"omitempty,max=120"`
	PostalCode    string     `json:"postal_code" db:"postal_code" validate:"omitempty,max=20"`
	IsDefault     bool       `json:"is_default" db:"is_default" validate:"omitempty"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at" validate:"omitempty"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at" validate:"omitempty"`
}
