package handler

import (
	"github.com/gin-gonic/gin"
	"mailer-auth/pkg/models"
	"net/http"
)

// @Summary Info
// @Tags user
// @Description get info about logged in user
// @ID info
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserClaims
// @Failure 401,500 {object} models.Response
// @Failure default {object} models.Response
// @Router /i [post]

// Gets a user from middleware, converts to models.UserClaims and returns it.
func (h *Handler) info(c *gin.Context) {
	ctxUser, ok := c.Get("User")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Some problems in server",
		})

		return
	}

	user := ctxUser.(models.UserClaims)
	if user.Username == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Some problems in server",
		})

		return
	}

	c.JSON(200, user)
}
