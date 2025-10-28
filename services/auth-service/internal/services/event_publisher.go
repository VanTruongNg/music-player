package services

import (
	"auth-service/internal/domain"
	"auth-service/internal/kafka/envelope"
	"auth-service/internal/kafka/producer"
	"context"
	"log"
	"time"
)

// EventPublisher handles publishing domain events to Kafka
type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, user *domain.User) error
	// Future events:
	// PublishUserUpdated(ctx context.Context, user *domain.User) error
	// PublishUserDeleted(ctx context.Context, userID string) error
	// PublishPasswordChanged(ctx context.Context, userID string) error
}

type kafkaEventPublisher struct {
	producer *producer.Producer
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(producer *producer.Producer) EventPublisher {
	return &kafkaEventPublisher{
		producer: producer,
	}
}

// PublishUserRegistered publishes user.registered event
func (p *kafkaEventPublisher) PublishUserRegistered(ctx context.Context, user *domain.User) error {
	// Build event data
	eventData := map[string]interface{}{
		"user_id":    user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"full_name":  user.FullName,
		"created_at": user.CreatedAt.Format(time.RFC3339),
	}

	// Create envelope
	env, err := envelope.NewEnvelope(
		"auth-service",
		envelope.PriorityHigh,
		eventData,
	)
	if err != nil {
		return err
	}

	// Add metadata
	env.Metadata = &envelope.Metadata{
		UserID:      user.ID,
		Environment: "production",
		CustomFields: map[string]string{
			"event_type": "user_lifecycle",
			"action":     "registration",
		},
	}

	// Validate envelope
	if err := env.Validate(); err != nil {
		return err
	}

	// Marshal to bytes
	messageBytes, err := env.Marshal()
	if err != nil {
		return err
	}

	// Publish to Kafka
	topic := envelope.TopicUserRegistered.String()
	key := user.ID

	if err := p.producer.Publish(ctx, topic, key, messageBytes); err != nil {
		log.Printf("[WARN] Failed to publish user.registered event: %v", err)
		return err
	}

	log.Printf("[INFO] Published user.registered event for user %s", user.ID)
	return nil
}
