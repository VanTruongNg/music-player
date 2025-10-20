//go:build wireinject
// +build wireinject

package main

import (
	"notification/configs"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type App struct {
	Router *gin.Engine
}

func InitializeApp(app *configs.AppConfig) (*App, error) {
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

	return r
}
