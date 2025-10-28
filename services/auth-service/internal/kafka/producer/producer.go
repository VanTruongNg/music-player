package producer

import (
	"auth-service/configs"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	cl *kgo.Client
}

func NewProducer(cfg *configs.KafkaConfig) (*Producer, error) {
	return NewProducerWithProfile(cfg, ProfileBalanced)
}

func NewProducerWithProfile(cfg *configs.KafkaConfig, profile ProducerProfile) (*Producer, error) {
	opts := GetProducerOpts(cfg, profile)

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	// Verify connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("[ERROR] Failed to connect to Kafka brokers: %w", err)
	}

	log.Printf("[INFO] Kafka producer connected to %v (profile=%s)", cfg.Brokers, profile)
	return &Producer{cl: client}, nil
}

// Publish sends a message to the specified Kafka topic (synchronous)
func (p *Producer) Publish(ctx context.Context, topic string, key string, message []byte) error {
	if topic == "" {
		return fmt.Errorf("topic cannot be empty")
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: message,
		Headers: []kgo.RecordHeader{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "source", Value: []byte("auth-service")},
			{Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}

	var produceErr error

	p.cl.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			produceErr = fmt.Errorf("failed to produce message: %w", err)
			log.Printf("❌ [KAFKA] Failed to publish to topic '%s': %v", topic, err)
			return
		}

		log.Printf("✅ [KAFKA] Message published successfully")
		log.Printf("   Topic: %s", r.Topic)
		log.Printf("   Partition: %d", r.Partition)
		log.Printf("   Offset: %d", r.Offset)
		log.Printf("   Key: %s", string(r.Key))
		log.Printf("   Size: %d bytes", len(r.Value))
	})

	if err := p.cl.Flush(ctx); err != nil {
		return fmt.Errorf("failed to flush kafka producer: %w", err)
	}

	return produceErr
}

func (p *Producer) PublishAsync(ctx context.Context, topic string, key string, message []byte) {
	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: message,
		Headers: []kgo.RecordHeader{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "source", Value: []byte("auth-service")},
		},
	}

	bgCtx := context.Background()
	p.cl.Produce(bgCtx, record, func(r *kgo.Record, err error) {
		if err != nil {
			log.Printf("❌ [KAFKA] Async publish failed to topic '%s': %v", topic, err)
			return
		}
		log.Printf("✅ [KAFKA] Async message published to %s (partition=%d, offset=%d)",
			r.Topic, r.Partition, r.Offset)
	})
}

func (p *Producer) Close() {
	if p.cl != nil {
		log.Println("[INFO] Closing Kafka producer...")
		p.cl.Close()
		log.Println("[INFO] Kafka producer connection closed")
	}
}
