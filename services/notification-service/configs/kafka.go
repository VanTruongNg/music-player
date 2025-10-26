package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers  []string
	GroupID  string
	ClientID string
	Debug    bool
	Topics   []string
}

func LoadKafkaConfig() *KafkaConfig {
	get := func(k, def string) string {
		if v := viper.GetString(k); v != "" {
			return v
		}
		return def
	}

	cfg := &KafkaConfig{
		Brokers:  splitTrim(get("KAFKA_BROKERS", "localhost:9092")),
		GroupID:  get("KAFKA_GROUP_ID", "default-group"),
		ClientID: get("KAFKA_CLIENT_ID", "default-client"),
		Topics:   splitTrim(get("KAFKA_TOPICS", "")),
		Debug:    get("APP_ENV", "") == "development",
	}

	if len(cfg.Brokers) == 0 {
		log.Println("[WARN] KAFKA_BROKERS not set, defaulting to localhost:9092")
	}
	return cfg
}

func splitTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
