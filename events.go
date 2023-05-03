package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

const (
	Authorization string = "authentication"
	Logging       string = "logging"
	Send          string = "send"
)

type requestType struct {
	Action string   `json:"action"`
	Auth   authType `json:"auth,omitempty"`
	Log    logType  `json:"log,omitempty"`
	Send   sendType `json:"send,omitempty"`
}

type authType struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type logType struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type sendType struct {
	From       string   `josn:"from,omitempty"`
	FromName   string   `josn:"from_name,omitempty"`
	To         string   `josn:"to"`
	Subject    string   `josn:"subject"`
	Body       string   `json:"body"`
	Attachment []string `josn:"attachments,omitempty"`
}

func (c *Config) handleEvents(delivery amqp.Delivery) error {
	fmt.Println("handing events...")
	var message requestType
	err := json.Unmarshal(delivery.Body, &message)
	if err != nil {
		return err
	}
	switch message.Action {
	case Authorization:
		{
			return c.handleAuthorization(message.Auth)
		}
	case Logging:
		{
			return c.handleLogging(message.Log)
		}
	default:
		return errors.New("Unknown events")
	}
}

func (c *Config) handleLogging(request logType) error {
	postBody, _ := json.Marshal(request)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:4321/log", "application/json", responseBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var payloadfromService jsonResponse
	err = json.NewDecoder(resp.Body).Decode(&payloadfromService)
	if err != nil {
		log.Println("log failed")
		return nil
	}
	log.Println("Log via RabitQP succeeded!")
	return nil
}

func (c *Config) handleAuthorization(request authType) error {
	postBody, _ := json.Marshal(request)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:80/auth", "application/json", responseBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var payloadfromService jsonResponse
	err = json.NewDecoder(resp.Body).Decode(&payloadfromService)
	if err != nil {
		log.Println("authentication failed")
		return nil
	}
	log.Println("authentication event handled successfully you are ready to go!")
	return nil
}
