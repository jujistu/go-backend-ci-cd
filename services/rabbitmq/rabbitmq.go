package rabbitmq

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	amqp "github.com/rabbitmq/amqp091-go"
)

var br Broker

func InitializeBroker(logger log.Logger) (*amqp.Connection, *amqp.Channel, error) {
	level.Info(logger).Log("RabbitMQ ", "connecting")

	rmqHost := os.Getenv("RMQ_HOST")
	rmqUserName := os.Getenv("RMQ_USERNAME")
	rmqPassword := os.Getenv("RMQ_PASSWORD")
	rmqPort := os.Getenv("RMQ_PORT")
	dsn := "amqp://" + rmqUserName + ":" + rmqPassword + "@" + rmqHost + ":" + rmqPort + "/"

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	br.SetUp(conn, ch)
	level.Info(logger).Log("RabbitMQ ", "connected")
	return conn, ch, nil
}

func GetRabbitMQBroker() *Broker {
	return &br
}
