package config

import (
	"os"
)

type Config struct {
	Kafka KafkaConfig
}

type KafkaConfig struct {
	Broker string // comma-separated list of broker addresses
	LogsTopic	 string
	MetricsTopic string
	TracesTopic string
	GroupID string
}

func LoadConfig() *Config {
	return &Config {
		Kafka: KafkaConfig {
			Broker: getEnv("KAFKA_BROKERS", "localhost:9092"),
			LogsTopic: getEnv("KAFKA_LOGS_TOPIC", "logs"),
			MetricsTopic: getEnv("KAFKA_METRICS_TOPIC", "metrics"),
			TracesTopic: getEnv("KAFKA_TRACES_TOPIC", "traces"),
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