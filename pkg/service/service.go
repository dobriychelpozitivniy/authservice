package service

import (
	"mailer-auth/pkg/models"
	"mailer-auth/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type ConfigService struct {
	SigningKey string
	AccessTokenTTL int
	RefreshTokenTTL int
}

type PasswordManager interface {
	Hash(password string) (string, error)
	Check(password string, hash string) error
}

type User interface {
	GetUser(username string) (*models.User, error)
	CreateUser(username string, hash string) (interface{}, error)
}

type TokenManager interface {
	GenerateRefreshToken(userID string) (string, error)
	GenerateAccessToken(userID string) (string, error)
	ParseAccessToken(AccessTokenClaims string) (*models.AccessTokenClaims, error)
	ParseRefreshToken(RefreshToken string) (*models.RefreshTokenClaims, error)
}

type Service struct {
	PasswordManager
	TokenManager
	User
}

func NewService(cfg *ConfigService, repos *repository.Repository) *Service {
	return &Service{
		PasswordManager: NewPasswordManagerService(),
		TokenManager: NewTokenManagerService(cfg),
		User: NewUserService(repos.Users),
	}
}