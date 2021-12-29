package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mailer-auth/pkg/config"
	"mailer-auth/pkg/models"
	"mailer-auth/pkg/service"
	mock_service "mailer-auth/pkg/service/mocks"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Exist user login: testuser, password: wwww
func TestHandler_login(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, username string)

	testCases := []struct {
		Name             string
		HeaderKey        string
		Login            string
		Password         string
		mockBehavior mockBehavior
		ExpectedCode     int
		IsExpectedCookie bool
	}{
		{
			Name:             "not existing key Authorization",
			HeaderKey:        "",
			Login:            "",
			Password:         "",
			mockBehavior: func(s *mock_service.MockUser, username string) {},
			ExpectedCode:     http.StatusBadRequest,
			IsExpectedCookie: false,
		},
		{
			Name:             "not existing/wrong login",
			HeaderKey:        "Authorization",
			Login:            "wrong",
			mockBehavior: func(s *mock_service.MockUser, username string) {
				s.EXPECT().GetUser(username).Return(nil, errors.New("err"))
			},
			Password:         "wwww",
			ExpectedCode:     http.StatusUnauthorized,
			IsExpectedCookie: false,
		},
		{
			Name:             "not existing/wrong password",
			HeaderKey:        "Authorization",
			Login:            "testuser",
			mockBehavior: func(s *mock_service.MockUser, username string) {
				s.EXPECT().GetUser(username).Return(nil, errors.New("err"))
			},
			Password:         "wrong",
			ExpectedCode:     http.StatusUnauthorized,
			IsExpectedCookie: false,
		},
		{
			Name:             "all okey",
			HeaderKey:        "Authorization",
			Login:            "testuser",
			Password:         "wwww",
			mockBehavior: func(s *mock_service.MockUser, username string) {
				u := &models.User{
					ID:           primitive.ObjectID{},
					Username:     "testuser",
					PasswordHash: "$2a$12$tc548.1ls5q7Enkgsj5ivuP0LRU1ATyp0TaAWaSWFveZZd59TmeZm",
				}
				s.EXPECT().GetUser(username).Return(u, nil)
			},
			ExpectedCode:     http.StatusOK,
			IsExpectedCookie: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			cfg, err := config.Load("../../configs/local")
			if err != nil {
				t.Errorf("Error init config: %s", err)
			}

			userService := mock_service.NewMockUser(controller)
			tc.mockBehavior(userService, tc.Login)

			services := &service.Service{
				PasswordManager: service.NewPasswordManagerService(),
				TokenManager:    service.NewTokenManagerService(&service.ConfigService{
					SigningKey:      cfg.Key,
					AccessTokenTTL:  cfg.AccessTokenTTL,
					RefreshTokenTTL: cfg.RefreshTokenTTL,
				}),
				User:            userService,
			}

			handler := NewHandler(&ConfigHandler{
				Host:            cfg.Host,
				AccessTokenTTL:  cfg.AccessTokenTTL,
				RefreshTokenTTL: cfg.RefreshTokenTTL,
			}, services, false)

			// Init Endpoint
			r := gin.New()
			r.GET("/login", handler.login)

			// Create Request
			w := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/login", nil)
			if err != nil {
				t.Errorf("Failed to create request: %v", err)
			}

			req.Header.Set(tc.HeaderKey, fmt.Sprintf("Basic %s", toBase64(tc.Login, tc.Password)))

			// Make Request
			r.ServeHTTP(w, req)

			cookies := w.Result().Cookies()

			assert.Equal(t, tc.ExpectedCode, w.Code, fmt.Sprintf(
				"Expected code value %v. Code from response %v. Message in body %s.",
				tc.ExpectedCode, w.Code, w.Body,
			))

			assert.Equal(t, tc.IsExpectedCookie, checkCookiesExist(cookies), fmt.Sprintf(
				"Expected code value %v. Code from response %v. Message in body %s.",
				tc.ExpectedCode, w.Code, w.Body,
			))

		})
	}

}

func checkCookiesExist(cookies []*http.Cookie) bool {
	var isAuthCookie bool = false
	var isRefreshCookie bool = false

	for _, cookie := range cookies {
		if cookie.Name == "refresh-token" {
			isRefreshCookie = true
		}
		if cookie.Name == "access-token" {
			isAuthCookie = true
		}
	}

	return isRefreshCookie && isAuthCookie
}

func toBase64(login string, password string) string {
	s := fmt.Sprintf("%s:%s", login, password)
	return base64.StdEncoding.EncodeToString([]byte(s))
}
