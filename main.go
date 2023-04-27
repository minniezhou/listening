package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

const queneName = "broker"

func main() {
	// connect to rabbit mq
	conn, err := amqp.Dial("amqp://guest:guest@localhost")
	if err != nil {
		log.Panic("failed to connect to rabbit mq")
	}
	defer conn.Close()

	ch, err := declareChannel(conn)
	if err != nil {
		log.Panic("failed to declare channel")
	}

	c := NewConfig(conn, ch)

	c.consume()
}

func declareChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(
		queneName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return ch, nil
}
