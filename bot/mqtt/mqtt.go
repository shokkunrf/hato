package mqtt

import (
	"fmt"
	"hato/config"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	Client       paho.Client
	PrimaryTopic string
}

func MakePublisher(conf *config.BrokerConfig) *Publisher {
	opts := paho.NewClientOptions()
	opts.AddBroker(conf.Origin)
	c := paho.NewClient(opts)

	return &Publisher{
		Client:       c,
		PrimaryTopic: conf.Topic,
	}
}

func (p *Publisher) Connect() error {
	token := p.Client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("[Publisher] %s", token.Error())
	}
	return nil
}

func (p *Publisher) Disconnect() {
	p.Client.Disconnect(250)
}

func (p *Publisher) Publish(subTopic string, message string) error {
	topic := p.PrimaryTopic
	if subTopic != "" {
		topic += "/" + subTopic
	}
	token := p.Client.Publish(topic, 2, false, message)

	if token.Wait() {
		return fmt.Errorf("[Publisher] %s", token.Error())
	}
	return nil
}
