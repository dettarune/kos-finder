package model

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func SuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	SendJSON(w, status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(w http.ResponseWriter, status int, message string, detail interface{}) {
	SendJSON(w, status, APIResponse{
		Success: false,
		Message: message,
		Error:   detail,
	})
}

func InternalServerErrorResponse(w http.ResponseWriter) {
	ErrorResponse(w, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
}

func BadRequestResponse(w http.ResponseWriter, msg string) {
	ErrorResponse(w, http.StatusBadRequest, msg, nil)
}
