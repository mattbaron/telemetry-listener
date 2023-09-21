package publisher

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQP struct {
	Brokers    []string
	Username   string
	Password   string
	Exchange   string
	VHost      string
	connection *amqp.Connection
	channel    *amqp.Channel
	auth       []amqp.Authentication
}

func NewAMQP() *AMQP {
	client := &AMQP{
		Brokers:  make([]string, 0),
		Username: "guest",
		Password: "guest",
		VHost:    "telegraf",
	}
	return client
}

func (client *AMQP) IsConnected() bool {
	return client.connection != nil && !client.connection.IsClosed()
}

func (client *AMQP) Connect() error {
	for _, broker := range client.Brokers {
		err, connection := client.connectBroker(broker)

		if err == nil {
			client.connection = connection
			break
		}
	}

	if client.connection == nil {
		return errors.New("could not connect to any broker")
	}

	channel, err := client.connection.Channel()
	if err != nil {
		return errors.New("could not create channel")
	}
	client.channel = channel

	return nil
}

func (client *AMQP) connectBroker(broker string) (error, *amqp.Connection) {
	if client.auth == nil {
		client.PlainAuthentication("guest", "guest")
	}

	connection, err := amqp.DialConfig(broker, amqp.Config{SASL: client.auth, Vhost: client.VHost})

	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	return err, connection
}

func (client *AMQP) PlainAuthentication(username string, password string) {
	client.auth = []amqp.Authentication{
		&amqp.PlainAuth{
			Username: string(client.Username),
			Password: string(client.Password),
		},
	}
}
