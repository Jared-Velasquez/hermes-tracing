package config

import (
	"os"
)

type Config struct {
	Kafka KafkaConfig
	Elasticsearch ElasticsearchConfig
}

type KafkaConfig struct {
	Broker       string // comma-separated list of broker addresses
	LogsTopic	 string
	MetricsTopic string
	TracesTopic  string
	GroupID      string
}

type ElasticsearchConfig struct {
	Addresses []string // list of Elasticsearch node addresses
	Index     string
}

func LoadConfig() *Config {
	return &Config {
		Kafka: KafkaConfig {
			Broker: getEnv("KAFKA_BROKERS", "localhost:29092"),
			LogsTopic: getEnv("KAFKA_LOGS_TOPIC", "logs"),
			MetricsTopic: getEnv("KAFKA_METRICS_TOPIC", "metrics"),
			TracesTopic: getEnv("KAFKA_TRACES_TOPIC", "traces"),
			GroupID: getEnv("KAFKA_GROUP_ID", "analytics-consumer-group"),
		},
		Elasticsearch: ElasticsearchConfig {
			Addresses: []string{getEnv("ES_ADDRESS", "http://localhost:9200")},
			Index:     getEnv("ES_INDEX", "traces"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}