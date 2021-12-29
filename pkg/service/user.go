package service

import (
	"context"
	"mailer-auth/pkg/models"
	"mailer-auth/pkg/repository"
)

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(username string) (*models.User, error) {
	return s.repo.Read(context.TODO(), username)
}

func (s *UserService) CreateUser(username, hash string) (interface{}, error) {
	return s.repo.Create(context.TODO(), username, hash)
}
