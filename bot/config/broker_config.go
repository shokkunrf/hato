package config

import (
	"errors"
	"fmt"
	"os"
)

type BrokerConfig struct {
	Origin string
	Topic  string
}

func GetBrokerConfig() (*BrokerConfig, error) {
	host := os.Getenv("BROKER_HOST")
	port := os.Getenv("BROKER_PORT")
	topic := os.Getenv("BROKER_TOPIC")
	if host == "" || port == "" || topic == "" {
		return nil, errors.New("[Config] env is empty")
	}

	return &BrokerConfig{
		Origin: fmt.Sprintf("tcp://%s:%s", host, port),
		Topic:  topic,
	}, nil
}
