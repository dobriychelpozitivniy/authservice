package models

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	Username string `json:"username" binding:"required"`
}

type AccessTokenClaims struct {
	User UserClaims
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	Username string `json:"username" binding:"required"`
	jwt.StandardClaims
}
