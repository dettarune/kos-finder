package model

type TokenClaims struct {
	Username string
	Role     string
	ExpiredAt int64
}

type CreateToken struct {
	Username string
}