package kafka

import (
	"auth-service/configs"
	"log"

	"github.com/IBM/sarama"
)

// NewKafkaProducer creates a new Kafka producer using KafkaConfig
func NewKafkaProducer(cfg *configs.KafkaConfig) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = cfg.Version
	config.ClientID = cfg.ClientID

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	if cfg.Debug {
		sarama.Logger = log.Default()
	}

	return sarama.NewSyncProducer(cfg.Brokers, config)
}

// NewKafkaConsumer creates a new Kafka consumer using KafkaConfig
func NewKafkaConsumer(cfg *configs.KafkaConfig) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = cfg.Version
	config.ClientID = cfg.ClientID

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	if cfg.Debug {
		sarama.Logger = log.Default()
	}

	return sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
}
