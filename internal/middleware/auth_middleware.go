package middleware

import (
	"net/http"
	"strings"

	"github.com/dettarune/kos-finder/internal/exceptions"
	"github.com/dettarune/kos-finder/internal/util"
)

type UserMiddleware struct {
	tokenUtil *util.TokenUtil
}

func NewUserMiddleware(token string) *util.TokenUtil {
	return &UserMiddleware{tokenUtil: token}

}

func AuthMiddleware(r *http.Request) exceptions.HttpError{
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return exceptions.NewBadRequestError("Missing Access Token, Please login first")
	}

	util.VerifyJwt()
	
	return exceptions.NewInternalServerError()
	

}