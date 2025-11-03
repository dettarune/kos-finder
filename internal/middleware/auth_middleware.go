package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dettarune/kos-finder/internal/model"
	"github.com/dettarune/kos-finder/internal/util"
)

type AuthMiddleware struct {
	tokenUtil *util.TokenUtil
}

func NewAuthMiddleware(tokenUtil *util.TokenUtil) *AuthMiddleware {
	return &AuthMiddleware{tokenUtil: tokenUtil}
}

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			model.ErrorResponse(w, http.StatusUnauthorized, "Missing access token. Please login first", nil)
			return
		}

		claims, err := m.tokenUtil.ParseToken(token)
		if err != nil {
			model.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserClaims(ctx context.Context) (*model.TokenClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*model.TokenClaims)
	return claims, ok
}

func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r.Context())
			if !ok {
				model.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}

			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				model.ErrorResponse(w, http.StatusForbidden, "Role Unauthorized for this service", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}