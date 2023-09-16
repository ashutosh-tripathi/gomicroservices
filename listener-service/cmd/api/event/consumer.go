package event

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(_conn *amqp.Connection, _queueName string) (*Consumer, error) {
	consumer := Consumer{conn: _conn, queueName: _queueName}
	err := consumer.Setup()
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (c *Consumer) Setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
	return nil
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	que, err := ch.QueueDeclare("testqueueue", false, false, true, false, nil)
	if err != nil {
		return err
	}
	for _, v := range topics {
		err := ch.QueueBind(que.Name, v, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}
	messages, err := ch.Consume(que.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range messages {
			fmt.Println("received message", d)
			fmt.Println("received message body", string(d.Body))
			var payLoad Payload
			_ = json.Unmarshal(d.Body, &payLoad)
			fmt.Println("payLoad:", payLoad)
		}
	}()
	<-forever
	return nil
}
