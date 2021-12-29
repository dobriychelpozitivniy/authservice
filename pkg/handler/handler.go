package handler

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "mailer-auth/docs"
	"mailer-auth/pkg/logger"
	"mailer-auth/pkg/service"
	"net/http"
)

type ConfigHandler struct {
	Host            string
	AccessTokenTTL  int
	RefreshTokenTTL int
}

type Handler struct {
	services      *service.Service
	config        *ConfigHandler
	IsPprofEnable bool
}

func NewHandler(cfg *ConfigHandler, services *service.Service, isPprofEnable bool) *Handler {
	return &Handler{
		config:        cfg,
		services:      services,
		IsPprofEnable: isPprofEnable,
	}
}

func (h *Handler) InitRoutes(g *gin.Engine) *gin.Engine {
	g.POST("/login", h.login)
	g.GET("/logout", h.logout)

	pprofGroup := g.Group("/debug", func(c *gin.Context) {
		if !h.IsPprofEnable {
			c.AbortWithStatus(http.StatusForbidden)

			return
		}
		c.Next()
	})

	g.GET("/toggle_debug", func(c *gin.Context) {
		h.IsPprofEnable = !h.IsPprofEnable
		c.AbortWithStatusJSON(200, fmt.Sprintf("pprof working is - %v", h.IsPprofEnable))
	})

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	pprof.RouteRegister(pprofGroup, "pprof")

	authorized := g.Group("/")

	authorized.Use(logger.Logger())
	authorized.Use(h.isAuthorized())

	authorized.POST("/me", h.me)
	authorized.POST("/i", h.info)

	return g
}
