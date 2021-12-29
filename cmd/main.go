package main

import (
	"flag"
	"fmt"
	"mailer-auth/pkg/config"
	"mailer-auth/pkg/handler"
	"mailer-auth/pkg/repository"
	"mailer-auth/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title Auth Service
// @version 1.0
// @description API Server for Auth Service

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Cookie
func main() {
	zapLog, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(zapLog)

	logger := zapLog.Sugar()

	flag.String("config", "configs/local", "display colorized output")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		logger.Panicf("Error parse flag config path: %s", err)
	}

	cfg, err := config.Load(viper.GetString("config"))
	if err != nil {
		logger.Panicf("Error init config: %s", err.Error())
	}

	logger.Info("CONFIG: %+v", cfg)

	repos, err := repository.NewRepository(&repository.ConfigRepository{
		DBHost:       cfg.DBHost,
		DBUsername:   cfg.DBUsername,
		DBPassword:   cfg.DBPassword,
		DBPort:       cfg.DBPort,
		DBTimeout:    cfg.DBTimeout,
		DBName:       cfg.DBName,
		DBCollection: cfg.DBCollection,
	})
	if err != nil {
		logger.Panicf("Error create repository layout: %s", err)
	}

	err = repos.CreateTestUser()
	if err != nil {
		logger.Panicf("Error create test user: %s", err)
	}

	services := service.NewService(&service.ConfigService{
		SigningKey:      cfg.Key,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}, repos)

	handlers := handler.NewHandler(&handler.ConfigHandler{
		Host:            cfg.Host,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}, services, false)

	g := InitServer(handlers)
	if err := g.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		logger.Panicf("Error start server: %s", err.Error())
	}
}

func InitServer(routes *handler.Handler) *gin.Engine {
	g := gin.New()
	g.Use(gin.Logger())
	routes.InitRoutes(g)

	return g
}
