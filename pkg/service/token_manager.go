package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"mailer-auth/pkg/models"
	"time"
)

type TokenManagerService struct {
	signingKey string
	accessTokenTTL int
	refreshTokenTTL int
}

func NewTokenManagerService(cfg *ConfigService) *TokenManagerService {
	return &TokenManagerService{
		signingKey: cfg.SigningKey,
		accessTokenTTL: cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}



func (s *TokenManagerService) GenerateRefreshToken(userID string) (string, error) {
	claims := &models.RefreshTokenClaims{
		Username: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(s.refreshTokenTTL)).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString([]byte(s.signingKey))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}


func (s *TokenManagerService) GenerateAccessToken(userID string) (string, error) {
	claims := &models.AccessTokenClaims{
		User: models.UserClaims{
			Username: userID,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(s.accessTokenTTL)).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(s.signingKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// ParseAccessToken Check the validity of the access token and return information from it
func (s *TokenManagerService) ParseAccessToken(AccessTokenClaims string) (*models.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(AccessTokenClaims, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.AccessTokenClaims)
	if !ok {
		return nil, errors.New("cant get user claims from token")
	}

	return claims, nil
}

// ParseRefreshToken Check the validity of the refresh token and return information from it
func (s *TokenManagerService) ParseRefreshToken(RefreshToken string) (*models.RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(RefreshToken, &models.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.RefreshTokenClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}


