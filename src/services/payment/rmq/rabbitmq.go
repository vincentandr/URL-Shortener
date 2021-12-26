package rmq

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vincentandr/shopping-microservice/src/model"
)

// RabbitMQ ...
var (
	amqConn *amqp.Connection
	amqChannel    *amqp.Channel
)

// NewRabbitMQ instantiates the RabbitMQ instances using configuration defined in environment variables.
func NewRabbitMQ() (error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("amqp.Dial %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel %w", err)
	}

	err = ch.ExchangeDeclare(
		"tasks", // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("ch.ExchangeDeclare %w", err)
	}

	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return fmt.Errorf("ch.Qos %w", err)
	}

	amqConn = conn
	amqChannel = ch

	return nil
}

func Close() error {
	if err := amqConn.Close(); err != nil{
		return fmt.Errorf("amqp close connection error: %v", err)
	}

	return nil
}

func PaymentSuccessfulEventPublish(order model.UserOrder) error{
	// Connecting bytes encoder and decoder
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(order); err != nil {
		return fmt.Errorf("byte encoding error: %v", err)
	}

	err := amqChannel.Publish(
		"tasks",
		"event.payment.success",
		false,
		false,
		amqp.Publishing{
                        ContentType: "text/plain",
                        Body: b.Bytes(),
        })
	if err != nil {
		return fmt.Errorf("publish payment successful event failed: %v", err)
	}

	return nil
}