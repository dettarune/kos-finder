package util

import (
	"fmt"
	"net/http"

	"github.com/dettarune/kos-finder/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type CustomError struct {
	StatusCode int
	Message    string
	Suggestion string
}

func HashPassword(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}
func VerifyPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return err
	}
	
	return nil  
}



func (e *CustomError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func CheckAuthConflict(existingUser, reqUser *model.RegisterRequest) error {
	switch {
	case existingUser.Username == reqUser.Username && existingUser.Email == reqUser.Email:
		return &CustomError{http.StatusConflict, "Username and Email are already registered", "Choose different ones"}
	case existingUser.Username == reqUser.Username:
		return &CustomError{http.StatusConflict, "Username is already registered", "Choose a different username"}
	case existingUser.Email == reqUser.Email:
		return &CustomError{http.StatusConflict, "Email is already registered", "Choose a different email"}
	default:
		return nil
	}
}
