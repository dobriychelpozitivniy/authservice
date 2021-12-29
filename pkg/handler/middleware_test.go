package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mailer-auth/pkg/config"
	"mailer-auth/pkg/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_isAuthorized(t *testing.T) {
	testCases := []struct{
		Name string
		ExpectedCode int
		isNeedValidCookies bool
		isExistContextValue bool
	}{
		{
			Name:          "valid cookies",
			ExpectedCode:  200,
			isNeedValidCookies: true,
			isExistContextValue: true,
		},
		{
			Name:          "invalid cookies",
			ExpectedCode:  401,
			isNeedValidCookies: false,
			isExistContextValue: false,
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
		TokenManager:    tm,
	}

	handler := NewHandler(&ConfigHandler{
		Host:            cfg.Host,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}, services, false)

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			r := gin.New()

			r.POST("/isAuthorized", handler.isAuthorized(), handler.info)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var aT string = ""
			var rT string = ""
			if tC.isNeedValidCookies {
				aT, _ = handler.services.GenerateAccessToken("testuser")
				rT, _ = handler.services.GenerateRefreshToken("testuser")
			}


			c.Request, _ = http.NewRequest("POST", "/isAuthorized", nil)
			c.Request.AddCookie(&http.Cookie{
				Name:   "access-token",
				Value:  aT,
				Path:   "/",
				Domain: "127.0.0.1",
				MaxAge: 3600,
				Secure: false,
			})
			c.Request.AddCookie(&http.Cookie{
				Name:   "refresh-token",
				Value:  rT,
				Path:   "/",
				Domain: "127.0.0.1",
				MaxAge: 3600,
				Secure: false,
			})

			// Make Request
			r.HandleContext(c)

			_, ok := c.Get("User")

			assert.Equal(t, tC.isExistContextValue, ok)
			assert.Equal(t, tC.ExpectedCode, w.Code)
		})
	}


}
