package configs

import (
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers      []string
	Version      sarama.KafkaVersion
	GroupID      string
	ClientID     string
	Debug        bool
	TopicPattern string
}

func LoadKafkaConfig() *KafkaConfig {
	brokersStr := viper.GetString("KAFKA_BROKERS")
	var brokers []string
	if brokersStr != "" {
		brokers = strings.Split(brokersStr, ",")
		for i := range brokers {
			brokers[i] = strings.TrimSpace(brokers[i])
		}
	}

	groupID := viper.GetString("KAFKA_GROUP_ID")
	clientID := viper.GetString("KAFKA_CLIENT_ID")
	pattern := viper.GetString("KAFKA_TOPIC_PATTERN")
	versionStr := viper.GetString("KAFKA_VERSION")

	version, err := sarama.ParseKafkaVersion(versionStr)
	if err != nil {
		log.Fatalf("[KafkaConfig] Invalid Kafka version '%s': %v", versionStr, err)
	}

	cfg := &KafkaConfig{
		Brokers:      brokers,
		Version:      version,
		GroupID:      groupID,
		ClientID:     clientID,
		Debug:        viper.GetString("APP_ENV") == "development",
		TopicPattern: pattern,
	}

	if len(cfg.Brokers) == 0 {
		log.Println("[WARN] KAFKA_BROKERS is empty. Kafka connections will fail if not set.")
	}
	if cfg.GroupID == "" {
		log.Println("[WARN] KAFKA_GROUP_ID is empty. Consumers will fail to join group unless set explicitly.")
	}
	if cfg.ClientID == "" {
		log.Println("[WARN] KAFKA_CLIENT_ID is empty. Consider setting it for easier traceability.")
	}
	if cfg.TopicPattern == "" {
		log.Println("[WARN] KAFKA_TOPIC_PATTERN is empty. Consider setting it explicitly if using regex consumer.")
	}

	return cfg
}
