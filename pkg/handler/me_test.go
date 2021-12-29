package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mailer-auth/pkg/config"
	"mailer-auth/pkg/models"
	"mailer-auth/pkg/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_me(t *testing.T) {
	testCases := []struct {
		Name          string
		ExpectedCode  int
		ExpectedValue interface{}
		ContextValue  models.UserClaims
		ContextKey    string
		BodyToken     string
	}{
		{
			Name:          "valid key and value",
			ExpectedCode:  200,
			ExpectedValue: "testuser",
			ContextValue:  models.UserClaims{Username: "testuser"},
			ContextKey:    "User",
			BodyToken:     "",
		},
		{
			Name:          "invalid key",
			ExpectedCode:  500,
			ExpectedValue: "",
			ContextValue:  models.UserClaims{Username: "testuser"},
			ContextKey:    "invalid",
			BodyToken:     "none",
		},
		{
			Name:          "invalid value",
			ExpectedCode:  500,
			ExpectedValue: "",
			ContextValue:  models.UserClaims{},
			ContextKey:    "User",
			BodyToken:     "none",
		},
		{
			Name:          "valid access token",
			ExpectedCode:  200,
			ExpectedValue: "testuser",
			ContextValue:  models.UserClaims{Username: "testuser"},
			ContextKey:    "User",
			BodyToken:     "",
		},
		{
			Name:          "invalid access token",
			ExpectedCode:  400,
			ExpectedValue: "",
			ContextValue:  models.UserClaims{Username: "testuser"},
			ContextKey:    "User",
			BodyToken:     "wqjdoqwhdiquwfhaosldaiwodj",
		},
	}

	cfg, err := config.Load("../../configs/local")
	if err != nil {
		t.Errorf("Error init config: %s", err)
	}

	tm := service.NewTokenManagerService(&service.ConfigService{
		SigningKey:      cfg.Key,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	})

	services := &service.Service{
		TokenManager: tm,
	}

	handler := NewHandler(&ConfigHandler{
		Host:            cfg.Host,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}, services, false)

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			if tC.BodyToken == "" {
				tC.BodyToken, _ = handler.services.GenerateAccessToken("testuser")
			}

			r := gin.New()

			r.Use(func(c *gin.Context) {
				c.Set(tC.ContextKey, tC.ContextValue)
			})

			r.POST("/me", handler.me)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(models.BodyRequest{AccessToken: tC.BodyToken})

			c.Request, _ = http.NewRequest("POST", "/me", strings.NewReader(string(body)))

			r.ServeHTTP(w, c.Request)

			var user models.UserClaims

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
