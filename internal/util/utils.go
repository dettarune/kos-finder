package util

import (
	"fmt"
	"net/http"

	"github.com/dettarune/kos-finder/internal/entity"
)

type CustomError struct {
	StatusCode int
	Message    string
	Suggestion string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func CheckAuthConflict(existingUser, reqUser *entity.User) error {
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
