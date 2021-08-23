package main

import (
	"fmt"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type MQTTClientOptions struct {
	Address		string
	ClientID	string
	Username	string
	Password	string
}

type MQTTClient struct {
	client mqtt.Client
}

func NewMQTT(options MQTTClientOptions) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions().AddBroker(options.Address).SetClientID(options.ClientID).SetUsername(options.Username).SetPassword(options.Password).SetKeepAlive(2 * time.Second).SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	m := &MQTTClient {
		c,
	}

	return m, nil
}

func (m *MQTTClient) Publish(tag string, feed time.Time) error {
	token := m.client.Publish(tag, 0, false, fmt.Sprintf("%d", feed.Unix()))
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}
