package config

import (
	"log"

	"DistributeX/internal/common"
	"github.com/spf13/viper"
)

func LoadConfig() (*common.Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	config := &common.Config{
		RedisAddress:   viper.GetString("REDIS_ADDRESS"),
		GRPCPort:       viper.GetString("GRPC_PORT"),
		NodeTimeout:    viper.GetDuration("NODE_TIMEOUT"),
		KafkaBroker:    viper.GetString("KAFKA_BROKER"),
		KafkaTaskTopic: viper.GetString("KAFKA_TASK_TOPIC"),
		KafkaLogTopic:  viper.GetString("KAFKA_LOG_TOPIC"),
	}

	return config, nil
}
