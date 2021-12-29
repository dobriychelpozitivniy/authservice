package service

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordManagerService struct {
}

func NewPasswordManagerService() *PasswordManagerService {
	return &PasswordManagerService{}
}

// Hash Generate password hash
func (m *PasswordManagerService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Check Checks the validity of a hash to a password
func (m *PasswordManagerService) Check(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
