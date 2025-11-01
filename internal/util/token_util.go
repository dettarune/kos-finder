package util

import (
	"fmt"
	"time"

	"github.com/dettarune/kos-finder/internal/entity"
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

func (t *TokenUtil) CreateToken(payload *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": payload.Username,
		"email": payload.Email,
		"exp":      time.Now().Add(time.Hour).UnixMilli(),
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (t *TokenUtil) ParseToken(jwtToken string) (*entity.User, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token")
	}

	claims := token.Claims.(jwt.MapClaims)

	exp := claims["exp"].(float64)
	if exp < float64(time.Now().UnixMilli()) {
		return nil, fmt.Errorf("expired token")
	}

	username := claims["username"].(string)
	user := &entity.User{
		Username: username,
	}

	return user, nil
}
