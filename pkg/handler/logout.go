package handler

import (
	"github.com/gin-gonic/gin"
	"mailer-auth/pkg/models"
	"net/http"
)

// @Summary Logout
// @Tags user
// @Description delete cookies with tokens
// @ID logout
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure default {object} models.Response
// @Router /logout [get]

// Remove refresh-token and access-token from cookies.
func(h *Handler) logout(c *gin.Context) {
	c.SetCookie(
		"refresh-token",
		"",
		-1,
		"/",
		h.config.Host,
		false,
		true,
	)
	c.SetCookie(
		"access-token",
		"",
		-1,
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

	c.AbortWithStatusJSON(http.StatusOK, &models.Response{
		Code:    http.StatusOK,
		Message: "Successful logout",
	})
}