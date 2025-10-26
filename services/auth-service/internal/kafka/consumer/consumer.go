package consumer

import (
	"auth-service/configs"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	cl *kgo.Client
}

func NewConsumer(cfg *configs.KafkaConfig) (*Consumer, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.BlockRebalanceOnPoll(),
		kgo.SessionTimeout(45 * time.Second),
		kgo.HeartbeatInterval(3 * time.Second),
	}
	if cfg.Debug {
		opts = append(opts, kgo.WithLogger(kgo.BasicLogger(os.Stderr, kgo.LogLevelInfo, nil)))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("[ERROR] Failed to connect to Kafka brokers: %w", err)
	}

	log.Printf("[INFO] Kafka consumer connected to %v (group: %s)", cfg.Brokers, cfg.GroupID)
	return &Consumer{cl: client}, nil
}

func (c *Consumer) Close() {
	c.cl.Close()
	log.Println("[INFO] Kafka consumer connection closed")
}
