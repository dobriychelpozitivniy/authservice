package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mailer-auth/pkg/models"
	"net/http"
)

// @Summary Me
// @Tags user
// @Description get info about logged in user or info about user in body
// @ID me
// @Produce  json
// @Param input body models.BodyRequest false "user token"
// @Security ApiKeyAuth
// @Success 200 {object} models.UserClaims
// @Failure 400,401,500 {object} models.Response
// @Failure default {object} models.Response
// @Router /me [post]

// Gets a user from middleware or from body if request from profile service, converts to models.UserClaims and returns it.
func (h *Handler) me(c *gin.Context) {
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

	if c.GetHeader("X-SERVICE-NAME") == "profile/info" {
		c.JSON(200, user)

		return
	}

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Some problems in server",
		})

		return
	}

	var bodyRequest models.BodyRequest

	if err = json.Unmarshal(jsonData, &bodyRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		})

		return
	}

	userFromBody, err := h.services.ParseAccessToken(bodyRequest.AccessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid access-token in body",
		})

		return
	}

	c.JSON(200, userFromBody.User)
}