//go:build wireinject
// +build wireinject

package main

import (
	"auth-service/configs"
	"auth-service/internal/db"
	"auth-service/internal/handlers"
	"auth-service/internal/kafka/consumer"
	"auth-service/internal/kafka/producer"
	"auth-service/internal/middleware"
	"auth-service/internal/redis"
	"auth-service/internal/repositories"
	"auth-service/internal/routes"
	"auth-service/internal/services"
	redisutil "auth-service/internal/utils/redis"

	tokenmanager "auth-service/internal/services/TokenManager"
	"auth-service/internal/utils/jwt"

	"auth-service/internal/utils/twofa"

	authv1 "music-player/api/proto/auth/v1"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
)

type App struct {
	Router        *gin.Engine
	GRPCServer    *configs.GRPCServer
	KafkaProducer *producer.Producer
	KafkaConsumer *consumer.Consumer
}

func InitializeApp(appCfg *configs.AppConfig, dbCfg *configs.DBConfig, redisCfg *configs.RedisConfig, kafkaCfg *configs.KafkaConfig) (*App, error) {
	wire.Build(
		// Infrastructure
		db.NewGormDB,
		redis.NewRedisClient,
		producer.NewProducer,
		consumer.NewConsumer,

		// Repositories
		repositories.NewUserRepository,

		// Utilities
		provideTwoFAUtil,
		provideRedisUtil,
		provideJWTConfig,
		provideJWTService,
		provideTokenManager,

		// Services
		services.NewEventPublisher,
		services.NewUserService,
		services.NewTwoFAService,

		// Middleware
		middleware.NewAuthMiddleware,

		// Handlers
		handlers.NewUserHandler,
		handlers.NewTwoFAHandler,
		handlers.NewAuthGRPCHandler,
		handlers.NewJWKSHandler,

		// Server components
		provideRouter,
		provideGRPCServer,
		provideApp,
	)
	return nil, nil
}

func provideRouter(userHandler *handlers.UserHandler, twoFAHandler *handlers.TwoFAHandler, jwksHandler *handlers.JWKSHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	routes.RegisterUserRoutes(api, userHandler, authMiddleware)
	routes.RegisterTwoFARoutes(api, twoFAHandler, authMiddleware)
	routes.RegisterJWKSRoutes(api, jwksHandler)
	return r
}

func provideApp(router *gin.Engine, grpcServer *configs.GRPCServer, kafkaProducer *producer.Producer, kafkaConsumer *consumer.Consumer, authGRPCHandler *handlers.AuthGRPCHandler) *App {
	authv1.RegisterAuthServiceServer(grpcServer.GetServer(), authGRPCHandler)

	return &App{
		Router:        router,
		GRPCServer:    grpcServer,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
	}
}

func provideGRPCServer(appCfg *configs.AppConfig) (*configs.GRPCServer, error) {
	return configs.NewGRPCServer(appCfg.GRPCPort)
}

func provideTwoFAUtil() *twofa.TwoFAUtil {
	return twofa.NewTwoFAUtil("SupaGoodSongs")
}

func provideRedisUtil(client *goredis.Client) *redisutil.RedisUtil {
	return redisutil.NewRedisUtil(client)
}

func provideJWTConfig() *jwt.JWTConfig {
	cfg, err := jwt.LoadJWTConfig()
	if err != nil {
		panic("JWT config error: " + err.Error())
	}
	return cfg
}

func provideJWTService(cfg *jwt.JWTConfig) jwt.JWTService {
	return jwt.NewJWTService(cfg)
}

func provideTokenManager(jwtSvc jwt.JWTService, redisUtil *redisutil.RedisUtil) tokenmanager.TokenManager {
	return tokenmanager.NewTokenManager(jwtSvc, redisUtil)
}
