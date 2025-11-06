package util

import (
	"fmt"
	"time"

	"github.com/dettarune/kos-finder/internal/exceptions"
	"github.com/dettarune/kos-finder/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type TokenUtil struct {
	SecretKey string
}

func NewTokenUtils(v *viper.Viper) *TokenUtil {
	secretKey := v.GetString("jwt.secretKey")
	if secretKey == "" {
		panic("JWT secretKey not found in config")
	}

	return &TokenUtil{
		SecretKey: secretKey,
	}
}

func (t *TokenUtil) CreateToken(payload *model.CreateToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": payload.Username,
		"role":     payload.Role,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (t *TokenUtil) ParseToken(jwtToken string) (*model.TokenClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration claim")
	}
	
	if expFloat < float64(time.Now().Unix()) {
		return nil, fmt.Errorf("token expired")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username claim")
	}
	
	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role claim")
	}

	return &model.TokenClaims{
		Username: username,
		Role:     role,
	}, nil
}

func (t *TokenUtil) VerifyJwt(jwtToken string) error {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, exceptions.NewBadRequestError("Invalid signing method")
		}
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return exceptions.NewUnauthorizedError("Failed to parse token")
	}

	if !token.Valid {
		return exceptions.NewUnauthorizedError("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return exceptions.NewUnauthorizedError("Invalid token claims")
	}

	if expFloat, ok := claims["exp"].(float64); ok {
		if expFloat < float64(time.Now().Unix()) {
			return exceptions.NewUnauthorizedError("Token expired")
		}
	}

	return nil
}