package producer

import (
	"context"
	"fmt"
	"log"
	"notification/configs"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	cl *kgo.Client
}

func NewProducer(cfg *configs.KafkaConfig) (*Producer, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),
		kgo.RequiredAcks(kgo.AllISRAcks()),
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

	log.Printf("[INFO] Kafka producer connected to %v", cfg.Brokers)
	return &Producer{cl: client}, nil
}

func (p *Producer) Close() {
	p.cl.Close()
	log.Println("[INFO] Kafka producer connection closed")
}
