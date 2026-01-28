package config

import (
	"os"
)

type Config struct {
	Kafka         KafkaConfig
	EnableTraces  bool
	EnableMetrics bool
	EnableLogs    bool
}

type KafkaConfig struct {
	Broker string // comma-separated list of broker addresses
	LogsTopic	 string
	MetricsTopic string
	TracesTopic string
	GroupID string
}

func LoadConfig() *Config {
	return &Config{
		Kafka: KafkaConfig{
			Broker:       getEnv("KAFKA_BROKERS", "localhost:29092"),
			LogsTopic:    getEnv("KAFKA_LOGS_TOPIC", "logs"),
			MetricsTopic: getEnv("KAFKA_METRICS_TOPIC", "metrics"),
			TracesTopic:  getEnv("KAFKA_TRACES_TOPIC", "traces"),
			GroupID:      getEnv("KAFKA_GROUP_ID", "analytics-consumer-group"),
		},
		EnableTraces:  getEnvBool("ENABLE_TRACES", true),
		EnableMetrics: getEnvBool("ENABLE_METRICS", false),
		EnableLogs:    getEnvBool("ENABLE_LOGS", false),
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		return val == "true" || val == "1"
	}
	return defaultValue
}