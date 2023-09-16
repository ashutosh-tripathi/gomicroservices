package main

import (
	"fmt"
	"listener-service/cmd/api/event"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//try to connect to rabbitmq
	//start listening for messages
	//create consumer
	//watch queue and consume events
	connection, err := connect()
	if err != nil {
		fmt.Println("connection error")
		os.Exit(1)
	}
	fmt.Println("successfully established connection to RabbitMQ")
	defer connection.Close()
	fmt.Println("going to start listening for messages")
	consumer, err := event.NewConsumer(connection, "testqueueue")
	if err != nil {
		panic(err)
	}
	err = consumer.Listen([]string{"temp"})
	if err != nil {
		fmt.Println("error", err)
	}

}
func connect() (*amqp.Connection, error) {

	fmt.Println("Going to connect to rabbitmq")
	count := 0
	var connection *amqp.Connection
	var err error

	for count < 10 {
		connection, err = amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			time.Sleep(time.Second * 10)
			count++
		} else {
			break
		}

	}
	return connection, nil

}
