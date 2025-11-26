package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
}

type TokenClaims struct {
	UserID uint
	Expiry time.Time
}

type TokenManager interface {
	GenerateAccessToken(user *User) (string, error)
	GenerateRefreshToken(user *User) (string, error)
	ValidateToken(token string, isRefresh bool) (*TokenClaims, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) error
}

type AuthUsecase interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, accessToken string, refreshToken string) error
	GetMe(ctx context.Context, userID uint) (*User, error)
}
