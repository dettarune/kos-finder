package util

import (
	"fmt"

	"github.com/dettarune/kos-finder/internal/exceptions"
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
		return exceptions.NewConflictError("Username And Email are already registed")
	case existingUser.Username == reqUser.Username:
		return exceptions.NewConflictError("Username are already registed")
	case existingUser.Email == reqUser.Email:
		return exceptions.NewConflictError("Email are already registed")
	default:
		return nil
	}
}
