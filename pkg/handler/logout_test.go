package handler

import (
	"fmt"
	"mailer-auth/pkg/config"
	"mailer-auth/pkg/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_logout(t *testing.T) {
	t.Run("Check delete refresh and access cookies", func(t *testing.T) {
		// Init Dependencies

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

		// Init Endpoint
		r := gin.New()
		r.GET("/logout", handler.logout)

		// Create Request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/logout", nil)
		c.Request.AddCookie(&http.Cookie{
			Name:   "access-token",
			Value:  "qwduiqwhd1289e128asdzhxcAH",
			Path:   "/",
			Domain: "127.0.0.1",
			MaxAge: 3600,
			Secure: false,
		})
		c.Request.AddCookie(&http.Cookie{
			Name:   "refresh-token",
			Value:  "jdiqwd818u11u1usnxnxn",
			Path:   "/",
			Domain: "127.0.0.1",
			MaxAge: 3600,
			Secure: false,
		})

		r.HandleContext(c)

		for _, cookie := range w.Result().Cookies() {
			if cookie.Name == "refresh-token" || cookie.Name == "access-token" {
				assert.Equal(t, "", cookie.Value, fmt.Sprintf(
					"Expected value: empty. In cookie %v is %v", cookie.Name, cookie.Value,
				))
				assert.Equal(t, -1, cookie.MaxAge, fmt.Sprintf(
					"Expected maxage: -1. In cookie %v is %v", cookie.Name, cookie.MaxAge,
				))
			}
		}
	})
}
