//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"gateway/configs"
	"gateway/internal/handlers"
	"gateway/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type App struct {
	Router       *gin.Engine
	GRPCClients  *configs.GRPCClients
	AuthHandler  *handlers.AuthHandler
	TwoFAHandler *handlers.TwoFAHandler
	UserHandler  *handlers.UserHandler
}

func InitializeApp(appCfg *configs.AppConfig) (*App, error) {
	wire.Build(
		// Handlers
		handlers.NewAuthHandler,
		handlers.NewTwoFAHandler,
		handlers.NewUserHandler,

		// Router and App
		provideRouter,
		provideGRPCClients,
		provideApp,
	)
	return nil, nil
}

func provideApp(
	router *gin.Engine,
	grpcClients *configs.GRPCClients,
	authHandler *handlers.AuthHandler,
	twoFAHandler *handlers.TwoFAHandler,
	userHandler *handlers.UserHandler,
) *App {
	return &App{
		Router:       router,
		GRPCClients:  grpcClients,
		AuthHandler:  authHandler,
		TwoFAHandler: twoFAHandler,
		UserHandler:  userHandler,
	}
}

func provideRouter(
	authHandler *handlers.AuthHandler,
	twoFAHandler *handlers.TwoFAHandler,
	userHandler *handlers.UserHandler,
) *gin.Engine {
	r := gin.Default()

	// Setup routes
	routes.SetupAuthRoutes(r, authHandler, twoFAHandler, userHandler)

	return r
}

func provideGRPCClients(appCfg *configs.AppConfig) (*configs.GRPCClients, error) {
	ctx := context.Background()

	return configs.NewGRPCClients(ctx, appCfg.AuthServiceAddr)
}
