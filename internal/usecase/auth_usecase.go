package usecase

import (
	"context"
	"errors"
	"time"

	constant "go-auth-service/internal/constants"
	"go-auth-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type authUsecase struct {
	userRepo       domain.UserRepository
	tokenManager   domain.TokenManager
	passwordHasher domain.PasswordHasher
	redisClient    *redis.Client
}

func NewAuthUsecase(userRepo domain.UserRepository, tokenManager domain.TokenManager, passwordHasher domain.PasswordHasher, redisClient *redis.Client) domain.AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		tokenManager:   tokenManager,
		passwordHasher: passwordHasher,
		redisClient:    redisClient,
	}
}

func (u *authUsecase) Register(ctx context.Context, user *domain.User) error {
	existingUser, _ := u.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := u.passwordHasher.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return u.userRepo.Create(ctx, user)
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := u.passwordHasher.CheckPassword(user.Password, password); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := u.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.tokenManager.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Check if token is blacklisted
	if u.redisClient != nil {
		val, _ := u.redisClient.Get(ctx, constant.STR_BLACKLIST+refreshToken).Result()
		if val != "" {
			return "", "", errors.New("token is blacklisted")
		}
	}

	claims, err := u.tokenManager.ValidateToken(refreshToken, true)
	if err != nil {
		return "", "", err
	}

	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := u.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	// Optionally rotate refresh token
	newRefreshToken, err := u.tokenManager.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	// Invalidate old refresh token if rotation is enabled
	if u.redisClient != nil {
		u.redisClient.Set(ctx, constant.STR_BLACKLIST+refreshToken, "true", time.Until(claims.Expiry))
	}

	return newAccessToken, newRefreshToken, nil
}

func (u *authUsecase) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	if u.redisClient == nil {
		return nil
	}

	// Blacklist access token
	accessClaims, err := u.tokenManager.ValidateToken(accessToken, false)
	if err == nil {
		u.redisClient.Set(ctx, constant.STR_BLACKLIST+accessToken, "true", time.Until(accessClaims.Expiry))
	}

	// Blacklist refresh token
	refreshClaims, err := u.tokenManager.ValidateToken(refreshToken, true)
	if err == nil {
		u.redisClient.Set(ctx, constant.STR_BLACKLIST+refreshToken, "true", time.Until(refreshClaims.Expiry))
	}

	return nil
}

func (u *authUsecase) GetMe(ctx context.Context, userID uint) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, userID)
}
