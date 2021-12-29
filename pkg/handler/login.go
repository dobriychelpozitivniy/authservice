package handler

import (
	"github.com/gin-gonic/gin"
	"mailer-auth/pkg/models"
	"net/http"
)

// @Summary Login
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 400,404 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure default {object} models.Response
// @Router /login [post]

// Checks the credentials and sets the access-token and refresh-token in the cookie if the credentials are valid.
func (h *Handler) login(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Wrong authorization header or value",
		})

		return
	}

	user, err := h.services.User.GetUser(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "Cant get user on login",
		})

		return
	}

	err = h.services.PasswordManager.Check(password, user.PasswordHash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &models.Response{
			Code:    http.StatusUnauthorized,
			Message: "Wrong password",
		})

		return
	}

	refreshToken, err := h.services.TokenManager.GenerateRefreshToken(user.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Cant generate refresh-token",
		})

		return
	}

	accessToken, err := h.services.TokenManager.GenerateAccessToken(user.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Cant generate access-token",
		})

		return
	}

	c.SetCookie(
		"refresh-token",
		refreshToken,
		h.config.RefreshTokenTTL,
		"/",
		h.config.Host,
		false,
		true,
	)
	c.SetCookie(
		"access-token",
		accessToken,
		h.config.AccessTokenTTL,
		"/",
		h.config.Host,
		false,
		true,
	)

	redirect := c.Query("redirect_uri")
	if redirect != "" {
		c.Redirect(http.StatusPermanentRedirect, redirect)

		return
	}

	c.JSON(http.StatusOK, &models.Response{
		Code:    http.StatusOK,
		Message: "Successful login",
	})
}
