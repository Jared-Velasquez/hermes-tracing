package config

import (
	"os"
)

type Config struct {
	Kafka 		  KafkaConfig
	Elasticsearch ElasticsearchConfig
	Server 		  AnalyticsConfig
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

type AnalyticsConfig struct {
	Address string // HTTP server address
	Port	string // HTTP server port
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
		Server: AnalyticsConfig {
			Address: getEnv("ANALYTICS_SERVER_ADDRESS", "0.0.0.0"),
			Port:    getEnv("ANALYTICS_SERVER_PORT", "8080"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}