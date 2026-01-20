package config

import (
	"os"
)

type Config struct {
	Kafka KafkaConfig
}

type KafkaConfig struct {
	Brokers string // comma-separated list of broker addresses
	LogsTopic	 string
	MetricsTopic string
	TracesTopics string
	GroupID string
}

func LoadConfig() *Config {
	return &Config {
		Kafka: KafkaConfig {
			Brokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
			LogsTopic: getEnv("KAFKA_LOGS_TOPIC", "logs"),
			MetricsTopic: getEnv("KAFKA_METRICS_TOPIC", "metrics"),
			TracesTopics: getEnv("KAFKA_TRACES_TOPIC", "traces"),
			GroupID: getEnv("KAFKA_GROUP_ID", "analytics-consumer-group"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}