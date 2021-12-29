package models

type BodyRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}
