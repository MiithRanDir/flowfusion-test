package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-auth-service/internal/domain"
	"go-auth-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockTokenManager
type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateAccessToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) GenerateRefreshToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) ValidateToken(token string, isRefresh bool) (*domain.TokenClaims, error) {
	args := m.Called(token, isRefresh)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

// MockPasswordHasher
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) CheckPassword(hash, password string) error {
	args := m.Called(hash, password)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenManager := new(MockTokenManager)
	mockPasswordHasher := new(MockPasswordHasher)

	authUsecase := usecase.NewAuthUsecase(mockUserRepo, mockTokenManager, mockPasswordHasher, nil)

	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			Email:    "test@example.com",
			Password: "password",
			Name:     "Test User",
		}

		mockUserRepo.On("GetByEmail", mock.Anything, user.Email).Return(nil, errors.New("not found"))
		mockPasswordHasher.On("HashPassword", user.Password).Return("hashed_password", nil)
		mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == user.Email && u.Password == "hashed_password"
		})).Return(nil)

		err := authUsecase.Register(context.Background(), user)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockPasswordHasher.AssertExpectations(t)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		user := &domain.User{
			Email:    "existing@example.com",
			Password: "password",
		}

		existingUser := &domain.User{ID: 1, Email: "existing@example.com"}

		mockUserRepo.On("GetByEmail", mock.Anything, user.Email).Return(existingUser, nil)

		err := authUsecase.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
		mockUserRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenManager := new(MockTokenManager)
	mockPasswordHasher := new(MockPasswordHasher)

	authUsecase := usecase.NewAuthUsecase(mockUserRepo, mockTokenManager, mockPasswordHasher, nil)

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		password := "password"
		hashedPassword := "hashed_password"
		user := &domain.User{ID: 1, Email: email, Password: hashedPassword}

		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
		mockPasswordHasher.On("CheckPassword", hashedPassword, password).Return(nil)
		mockTokenManager.On("GenerateAccessToken", user).Return("access_token", nil)
		mockTokenManager.On("GenerateRefreshToken", user).Return("refresh_token", nil)

		accessToken, refreshToken, err := authUsecase.Login(context.Background(), email, password)

		assert.NoError(t, err)
		assert.Equal(t, "access_token", accessToken)
		assert.Equal(t, "refresh_token", refreshToken)
		mockUserRepo.AssertExpectations(t)
		mockPasswordHasher.AssertExpectations(t)
		mockTokenManager.AssertExpectations(t)
	})

	t.Run("InvalidCredentials_UserNotFound", func(t *testing.T) {
		email := "nonexistent@example.com"
		password := "password"

		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(nil, errors.New("not found"))

		_, _, err := authUsecase.Login(context.Background(), email, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("InvalidCredentials_WrongPassword", func(t *testing.T) {
		email := "test@example.com"
		password := "wrong_password"
		hashedPassword := "hashed_password"
		user := &domain.User{ID: 1, Email: email, Password: hashedPassword}

		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
		mockPasswordHasher.On("CheckPassword", hashedPassword, password).Return(errors.New("password mismatch"))

		_, _, err := authUsecase.Login(context.Background(), email, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		mockUserRepo.AssertExpectations(t)
		mockPasswordHasher.AssertExpectations(t)
	})
}

func TestRefreshToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenManager := new(MockTokenManager)
	mockPasswordHasher := new(MockPasswordHasher)

	// Note: Redis client is nil here, so blacklist check is skipped
	authUsecase := usecase.NewAuthUsecase(mockUserRepo, mockTokenManager, mockPasswordHasher, nil)

	t.Run("Success", func(t *testing.T) {
		refreshToken := "valid_refresh_token"
		claims := &domain.TokenClaims{UserID: 1, Expiry: time.Now().Add(time.Hour)}
		user := &domain.User{ID: 1, Email: "test@example.com"}

		mockTokenManager.On("ValidateToken", refreshToken, true).Return(claims, nil)
		mockUserRepo.On("GetByID", mock.Anything, claims.UserID).Return(user, nil)
		mockTokenManager.On("GenerateAccessToken", user).Return("new_access_token", nil)
		mockTokenManager.On("GenerateRefreshToken", user).Return("new_refresh_token", nil)

		newAccess, newRefresh, err := authUsecase.RefreshToken(context.Background(), refreshToken)

		assert.NoError(t, err)
		assert.Equal(t, "new_access_token", newAccess)
		assert.Equal(t, "new_refresh_token", newRefresh)
		mockTokenManager.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
