//go:build wireinject
// +build wireinject

package main

import (
	"gateway/configs"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type App struct {
	Router *gin.Engine
}

func InitializeApp(appCfg *configs.AppConfig) (*App, error) {
	wire.Build(
		provideRouter,
		provideApp,
	)
	return nil, nil
}

func provideApp(router *gin.Engine) *App {
	return &App{
		Router: router,
	}
}

func provideRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return r
}
