package model

type RegisterRequest struct {
	Id        int    `json:"id" db:"id"`
	Email     string `json:"email" db:"email" validate:"required,email,max=128"`
	Full_name string `json:"full_name" db:"full_name" validate:"required,min=2,max=50"`
	Role      string `json:"role" db:"role" validate:"required,oneof=owner customer"`
	Username  string `json:"username" db:"username" validate:"required,alphanum,min=3,max=20"`
	Password  string `json:"password" db:"password" validate:"required,min=8"`
	Phone     string `json:"phone" db:"phone" validate:"required,numeric,min=8,max=15"`
}

type LoginRequest struct {
	Username string `json:"username" db:"username" validate:"required,alphanum,min=3,max=20"`
	Password string `json:"password" db:"password" validate:"required,min=8"`
}
