package exceptions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dettarune/kos-finder/internal/model"
	"github.com/go-playground/validator/v10"
)

type GlobalError interface {
	Error() string
	GetCode() int
}

type HttpError struct {
	Msg  string
	Code int
}

func (e HttpError) Error() string {
	return e.Msg
}

func (e HttpError) GetCode() int {
	return e.Code
}

func 	NewBadRequestError(msg string) HttpError {
	return HttpError{Msg: msg, Code: http.StatusBadRequest}
}

func NewUnauthorizedError(msg string) HttpError {
	return HttpError{Msg: msg, Code: http.StatusUnauthorized}
}

func NewForbiddenError(msg string) HttpError {
	return HttpError{Msg: msg, Code: http.StatusForbidden}
}

func NewConflictError(msg string) HttpError {
	return HttpError{Msg: msg, Code: http.StatusConflict}
}

func NewNotFoundError(msg string) HttpError {
	return HttpError{Msg: msg, Code: http.StatusNotFound}
}

func NewInternalServerError() HttpError {
	return HttpError{Msg: "Internal Server Error", Code: http.StatusInternalServerError}
}


func NewFailedValidationError(err *validator.ValidationErrors) *model.ValidationError {
	message := make(map[string]string)

	for _, e := range *err {
		field := e.Field()
		message[field] = handleValidationErrorMessage(e.Tag(), e.Param(), field)
	}

	return &model.ValidationError{
		Status: false,
		Message: "Validation Error",
		StatusCode: http.StatusUnprocessableEntity,
		Errors:  message,
	}
}

func handleValidationErrorMessage(tag, param, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s field can't be null", field)
	case "email":
		return fmt.Sprint("Email must be a valid email")
	case "min":
		return fmt.Sprintf("%s field must be at least %s characters", strings.ToLower(field), param)
	case "max":
		return fmt.Sprintf("%s field must be at most %s characters", strings.ToLower(field), param)
	}
	return "SABARRR, Error Belum dihandle ges. intinya validasi lu error dah"
}
