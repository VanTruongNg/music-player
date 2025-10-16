//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"gateway/configs"
	"gateway/internal/handlers"
	"gateway/internal/middleware"
	"gateway/internal/routes"
	"gateway/internal/utils/jwt"

	"gateway/internal/redis"
	redisutil "gateway/internal/utils/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
)

type App struct {
	Router         *gin.Engine
	GRPCClients    *configs.GRPCClients
	AuthHandler    *handlers.AuthHandler
	TwoFAHandler   handlers.TwoFAHandler
	UserHandler    handlers.UserHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func InitializeApp(appCfg *configs.AppConfig, redisCfg *configs.RedisConfig) (*App, error) {
	wire.Build(
		// Infrastructure
		redis.NewRedisClient,

		// Handlers
		handlers.NewAuthHandler,
		handlers.NewTwoFAHandler,
		handlers.NewUserHandler,

		// JWT utilities and middleware
		provideJWKSClient,
		jwt.NewJWTVerifier,
		middleware.NewAuthMiddleware,

		// Router and App
		provideRouter,
		provideGRPCClients,
		provideApp,

		// Utilities
		provideRedisUtil,
	)
	return nil, nil
}

func provideApp(
	router *gin.Engine,
	grpcClients *configs.GRPCClients,
	authHandler *handlers.AuthHandler,
	twoFAHandler handlers.TwoFAHandler,
	userHandler handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *App {
	return &App{
		Router:         router,
		GRPCClients:    grpcClients,
		AuthHandler:    authHandler,
		TwoFAHandler:   twoFAHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
}

func provideRouter(
	authHandler *handlers.AuthHandler,
	twoFAHandler handlers.TwoFAHandler,
	userHandler handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	r := gin.Default()

	routes.SetupAuthRoutes(r, authHandler, twoFAHandler, userHandler, authMiddleware)

	return r
}

func provideGRPCClients(appCfg *configs.AppConfig) (*configs.GRPCClients, error) {
	ctx := context.Background()

	return configs.NewGRPCClients(ctx, appCfg.AuthServiceAddr)
}

func provideJWKSClient(appCfg *configs.AppConfig) *jwt.JWKSClient {
	return jwt.NewJWKSClient(appCfg.AuthServiceHTTPURL)
}

func provideRedisUtil(redisClient *goredis.Client) *redisutil.RedisUtil {
	return redisutil.NewRedisUtil(redisClient)
}
