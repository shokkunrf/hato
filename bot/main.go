package main

import (
	"hato/config"
	"hato/mqtt"
	"log"
)

func main() {
	brokerConfig, err := config.GetBrokerConfig()
	if err != nil {
		log.Fatalln(err)
	}
	publisher := mqtt.MakePublisher(brokerConfig)
	err = publisher.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		publisher.Disconnect()
	}()
}
