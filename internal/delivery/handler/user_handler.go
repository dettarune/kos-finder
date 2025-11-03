package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dettarune/kos-finder/internal/model"
	"github.com/dettarune/kos-finder/internal/usecase"
	"github.com/dettarune/kos-finder/internal/exceptions"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	usecase *usecase.UserUseCase
	log     *logrus.Logger
}

func NewUserHandler(usecase *usecase.UserUseCase, log *logrus.Logger) *UserHandler {
	return &UserHandler{usecase: usecase, log: log}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		model.BadRequestResponse(w, "Format JSON tidak valid")
		return
	}

	if err := h.usecase.Register(r.Context(), &user); err != nil {
		h.handleError(w, err)
		return
	}

	model.SuccessResponse(w, http.StatusCreated, "Registration successful. Please check your email for verification token", nil)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.BadRequestResponse(w, "Format JSON tidak valid")
		return
	}

	token, err := h.usecase.Login(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
		// Secure:   true, //(HTTPS only)
	})

	model.SuccessResponse(w, http.StatusOK, "Login successful", map[string]interface{}{
		"token":      token,
		"expires_in": 86400,
	})
}

func (h *UserHandler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		var req struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			token = req.Token
		}
	}

	if token == "" {
		model.BadRequestResponse(w, "Verification token is required")
		return
	}

	if err := h.usecase.VerifyPassword(r.Context(), token); err != nil {
		h.handleError(w, err)
		return
	}

	model.SuccessResponse(w, http.StatusOK, "Email verified successfully. You can now login", nil)
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-1 * time.Hour), 
		MaxAge:   -1,
	})

	model.SuccessResponse(w, http.StatusOK, "Logout successful", nil)
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *model.ValidationError:
		model.ErrorResponse(w, e.StatusCode, e.Message, e.Errors)

	case exceptions.GlobalError:
		model.ErrorResponse(w, e.GetCode(), e.Error(), nil)

	default:
		h.log.WithError(err).Error("Internal error")
		model.InternalServerErrorResponse(w)
	}
}