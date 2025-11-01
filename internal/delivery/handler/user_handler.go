package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dettarune/kos-finder/internal/entity"
	"github.com/dettarune/kos-finder/internal/model"
	"github.com/dettarune/kos-finder/internal/usecase"
	"github.com/dettarune/kos-finder/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	usecase *usecase.UserUseCase
	log     *logrus.Logger
}

func NewUserHandler(usecase *usecase.UserUseCase, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		usecase: usecase,
		log:     log,
	}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.usecase.Register(r.Context(), &user); err != nil {
		var ce *util.CustomError
		if ok := errors.As(err, &ce); ok {
			writeJSONError(w, ce.StatusCode, ce.Message, nil)
			return
		}

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errs := make(map[string]string)
			for _, fe := range ve {
				errs[fe.Field()] = fe.Tag()
			}
			writeJSONError(w, http.StatusBadRequest, "Validation failed", errs)
			return
		}

		writeJSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully. Please check your email for verification.",
	})
}

func writeJSONError(w http.ResponseWriter, status int, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := map[string]any{"message": message}
	if details != nil {
		resp["errors"] = details
	}
	json.NewEncoder(w).Encode(resp)
}


// func (h *UserHandler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
// 	email := r.URL.Query().Get("email")
// 	if email != "email" && email== "" {
// 		http.Error(w, "Parameter Must be email", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.usecase.VerifyWithEmail(r.Context(), email); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
	

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": fmt.Sprintf("Success Sending Email to %s", email),
// 	})
// }


func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
    var reqUser model.UserLogin
    if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    token, err := h.usecase.Login(r.Context(), &reqUser)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

	http.SetCookie(w, &http.Cookie{
        Name:     "Authorization",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteNoneMode,
        Expires:  time.Now().Add(1 * time.Hour),
    })
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Authorization", token)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "User Successfully Logged In",
	"status": http.StatusText(http.StatusOK),
    })
}
