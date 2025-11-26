package service

import (
	"errors"
	"fmt"
	"time"

	"go-auth-service/config"
	"go-auth-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	cfg config.Config
}

func NewTokenService(cfg config.Config) *TokenService {
	return &TokenService{cfg: cfg}
}

func (t *TokenService) GenerateAccessToken(user *domain.User) (string, error) {
	expiry, err := time.ParseDuration(t.cfg.JWTAccessExpiry)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(expiry).Unix(),
		"type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.cfg.JWTSecret))
}

func (t *TokenService) GenerateRefreshToken(user *domain.User) (string, error) {
	expiry, err := time.ParseDuration(t.cfg.JWTRefreshExpiry)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(expiry).Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.cfg.JWTRefreshSecret))
}

func (t *TokenService) ValidateToken(tokenString string, isRefresh bool) (*domain.TokenClaims, error) {
	secret := t.cfg.JWTSecret
	if isRefresh {
		secret = t.cfg.JWTRefreshSecret
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			return nil, errors.New("invalid subject in token")
		}

		expFloat, ok := claims["exp"].(float64)
		if !ok {
			return nil, errors.New("invalid expiry in token")
		}

		return &domain.TokenClaims{
			UserID: uint(userIDFloat),
			Expiry: time.Unix(int64(expFloat), 0),
		}, nil
	}

	return nil, errors.New("invalid token")
}
