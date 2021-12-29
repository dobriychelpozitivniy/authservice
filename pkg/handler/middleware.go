package handler

import (
	"github.com/gin-gonic/gin"
	"mailer-auth/pkg/models"
	"net/http"
)

// Check the validity of the access-token and get the user from it.
// If the access token is invalid or does not exist, then we do the same with the refresh token
// If one of them is valid,then their values in the cookie are updated and user info add in context
// If all cookies not valid or not exist response 401.
func (h *Handler) isAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessClaims, response := checkAccessToken(c, h)
		if response != nil {
			c.AbortWithStatusJSON(response.Code, response)

			return
		}

		if accessClaims == nil {
			accessClaims, response = checkRefreshToken(c, h)
			if response != nil {
				c.AbortWithStatusJSON(response.Code, response)

				return
			}

			c.Set("User", accessClaims.User)
			c.Next()

			return
		}

		c.Set("User", accessClaims.User)
		c.Next()
	}
}

// Check is valid and is exist refresh-token in cookie
// If error - return not nil models.Response with code and error message
func checkRefreshToken(c *gin.Context, h *Handler) (*models.AccessTokenClaims, *models.Response) {
	refreshToken, ok := getRefreshCookieToken(c)
	if !ok {
		return nil, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "empty cookies",
		}
	}

	refreshClaims, err := h.services.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "invalid refresh-token",
		}
	}

	accessToken, response := setNewAccessToken(c, h, refreshClaims.Username)
	if response != nil {
		return nil, response
	}

	accessClaims, err := h.services.ParseAccessToken(accessToken)
	if err != nil {
		return nil, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Cant parse access-token",
		}
	}

	return accessClaims, nil
}

func setNewRefreshToken(c *gin.Context, h *Handler, username string) *models.Response {
	newRefreshToken, err := h.services.GenerateRefreshToken(username)
	if err != nil {
		return &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "cant generate refresh token",
		}
	}

	c.SetCookie(
		"refresh-token",
		newRefreshToken,
		h.config.RefreshTokenTTL,
		"/",
		h.config.Host,
		false,
		true,
	)

	return nil
}

func setNewAccessToken(c *gin.Context, h *Handler, username string) (string, *models.Response) {
	newAccessToken, err := h.services.GenerateAccessToken(username)
	if err != nil {
		return "", &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "cant generate access token",
		}
	}

	c.SetCookie(
		"access-token",
		newAccessToken,
		h.config.AccessTokenTTL,
		"/",
		h.config.Host,
		false,
		true,
	)

	return newAccessToken, nil
}

// Check is valid and is exist access-token in cookie
// If error - return not nil models.Response with code and error message
func checkAccessToken(c *gin.Context, h *Handler) (*models.AccessTokenClaims, *models.Response) {
	accessToken, ok := getAccessCookieToken(c)
	if !ok {
		return nil, nil
	}

	accessClaims, err := h.services.ParseAccessToken(accessToken)
	if err != nil {
		return nil, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "invalid access-token",
		}
	}

	setNewRefreshToken(c, h, accessClaims.User.Username)

	return accessClaims, nil
}

func getAccessCookieToken(c *gin.Context) (string, bool) {
	token, err := c.Cookie("access-token")
	if err != nil {
		return "", false
	}

	return token, true
}

func getRefreshCookieToken(c *gin.Context) (string, bool) {
	token, err := c.Cookie("refresh-token")
	if err != nil {
		return "", false
	}

	return token, true
}
