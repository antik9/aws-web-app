package queue

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitClient struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	messages <-chan amqp.Delivery
	queue    amqp.Queue
}

var Client *RabbitClient

func NewClient(clientType string) *RabbitClient {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s", err)
	}

	q, err := ch.QueueDeclare("events", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s", err)
	}
	var messages <-chan amqp.Delivery

	if clientType == "consumer" {
		msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
		if err != nil {
			log.Fatalf("%s", err)
		}
		messages = msgs
	}
	return &RabbitClient{
		conn:     conn,
		channel:  ch,
		queue:    q,
		messages: messages,
	}
}

func (c *RabbitClient) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (c *RabbitClient) SendMessage(message string) {
	err := c.channel.Publish(
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func (c *RabbitClient) ReadMessage() string {
	message := <-c.messages
	defer message.Ack(false)
	return string(message.Body)
}
