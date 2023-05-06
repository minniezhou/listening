package main

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

const (
	RABBITMQ_DEFAULT_PASS = "guest"
	RABBITMQ_DEFAULT_USER = "guest"
	RABBITMQ_URL          = "localhost"
)
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
		rabbit_pass := getEnv("RABBITMQ_DEFAULT_PASS", RABBITMQ_DEFAULT_PASS)
		rabbit_user := getEnv("RABBITMQ_DEFAULT_USER", RABBITMQ_DEFAULT_USER)
		rabbit_url := getEnv("RABBITMQ_URL", RABBITMQ_URL)
		rabbit_addr := "amqp://" + rabbit_pass + ":" + rabbit_user + "@" + rabbit_url + ":5672/"
		fmt.Println(rabbit_addr)
		conn, err := amqp.Dial(rabbit_addr)
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

func getEnv(key, default_value string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return default_value
	}
	return value
}
