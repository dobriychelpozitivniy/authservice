package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"mailer-auth/pkg/config"
	"mailer-auth/pkg/models"
	"mailer-auth/pkg/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Info(t *testing.T) {
	testCases := []struct{
		Name string
		ExpectedCode int
		ExpectedValue interface{}
		ContextValue models.UserClaims
		ContextKey string
		ExpectedMessage string
	}{
		{
			Name:          "valid key and value",
			ExpectedCode:  200,
			ExpectedValue: "testuser",
			ContextValue:  models.UserClaims{Username: "testuser"},
			ContextKey:    "User",
		},
		{
			Name:          "invalid key",
			ExpectedCode:  500,
			ExpectedValue: "",
			ContextValue:  models.UserClaims{Username: "testuser"},
			ContextKey:    "invalid",
		},
		{
			Name:          "invalid value",
			ExpectedCode:  500,
			ExpectedValue: "",
			ContextValue:  models.UserClaims{},
			ContextKey:    "User",
		},

	}

	cfg, err := config.Load("../../configs/local")
	if err != nil {
		t.Errorf("Error init config: %s", err)
	}

	services := &service.Service{}
	handler := NewHandler(&ConfigHandler{
		Host:            cfg.Host,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}, services, false)

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			r := gin.New()

			r.Use(func(c *gin.Context) {
				c.Set(tC.ContextKey, tC.ContextValue)
			})

			r.POST("/i", handler.info)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest("POST", "/i", nil)

			// Make Request
			r.ServeHTTP(w, c.Request)

			var user models.UserClaims

			zap.S().Info(w.Body)
			zap.S().Info(w.Code)
			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Error(err.Error())
			}

			err = json.Unmarshal(body, &user)
			if err != nil {
				t.Error(err.Error())
			}

			assert.Equal(t, tC.ExpectedValue, user.Username)
			assert.Equal(t, tC.ExpectedCode, w.Code)
		})
	}
}
