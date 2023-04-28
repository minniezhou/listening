package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConfig(conn *amqp.Connection, ch *amqp.Channel) Config {
	return Config{
		conn: conn,
		ch:   ch,
	}
}

func (c *Config) consume() {
	msgs, err := c.ch.Consume(
		queneName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		log.Panic(err)
	}
	var forever chan struct{}

	go func() {
		for d := range msgs {
			_ = c.handleEvents(d)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
