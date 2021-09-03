package mqtt

import (
	"fmt"
	"hato/config"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	client       paho.Client
	primaryTopic string
}

func MakePublisher(conf *config.BrokerConfig) *Publisher {
	opts := paho.NewClientOptions()
	opts.AddBroker(conf.Origin)
	c := paho.NewClient(opts)

	return &Publisher{
		client:       c,
		primaryTopic: conf.Topic,
	}
}

func (p *Publisher) Connect() error {
	token := p.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("[Publisher] %s", token.Error())
	}
	return nil
}

func (p *Publisher) Disconnect() {
	p.client.Disconnect(250)
}

func (p *Publisher) Publish(subTopic string, message string) error {
	topic := p.primaryTopic
	if subTopic != "" {
		topic += "/" + subTopic
	}
	token := p.client.Publish(topic, 2, false, message)

	if token.Wait() {
		return fmt.Errorf("[Publisher] %s", token.Error())
	}
	return nil
}
