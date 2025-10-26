//go:build wireinject
// +build wireinject

package main

import (
	"notification/configs"

	"notification/internal/kafka/consumer"
	"notification/internal/kafka/producer"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type App struct {
	Router        *gin.Engine
	KafkaProducer *producer.Producer
	KafkaConsumer *consumer.Consumer
}

func InitializeApp(app *configs.AppConfig, kafkaCfg *configs.KafkaConfig) (*App, error) {
	wire.Build(
		provideRouter,
		provideApp,
		producer.NewProducer,
		consumer.NewConsumer,
	)

	return nil, nil
}

func provideApp(router *gin.Engine, kafkaProducer *producer.Producer, kafkaConsumer *consumer.Consumer) *App {
	return &App{
		Router:        router,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
	}
}

func provideRouter() *gin.Engine {
	r := gin.Default()

	return r
}
