package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

const queneName = "broker"

func main() {
	// connect to rabbit mq
	conn, err := connectToRabbit()
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
func connectToRabbit() (*amqp.Connection, error) {
	count := 1
	backoff := time.Second
	log.Println("Connecting to Rabbit...")
	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			count++
			backoff = time.Duration(count*count) * time.Second
			log.Println("Rabit is not ready yet, backing off...")
			time.Sleep(backoff)
		} else {
			return conn, nil
		}

		if count > 10 {
			return nil, err
		}
	}
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
