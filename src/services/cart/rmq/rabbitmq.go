package rmq

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vincentandr/shopping-microservice/src/model"
	events "github.com/vincentandr/shopping-microservice/src/services/cart/event-handlers"
)

// RabbitMQ ...
var (
	amqConn *amqp.Connection
	amqChannel    *amqp.Channel
	msgs <-chan amqp.Delivery
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

	q, err := ch.QueueDeclare(
                "cartQueue",    // name
                false, // durable
                false, // delete when unused
                true,  // exclusive
                false, // no-wait
                nil,   // arguments
    )
	if err != nil {
		return fmt.Errorf("amqChannel.QueueDeclare %w", err)
	}

	err = ch.QueueBind(
                        q.Name,
                        "event.payment.success",
                        "tasks",
                        false,
                        nil)
	if err != nil {
		return fmt.Errorf("amqChannel.ExchangeDeclare %w", err)
	}

	msgs, err = ch.Consume(
                q.Name, // queue
                "",     // consumer
                false,   // auto ack
                false,  // exclusive
                false,  // no local
                false,  // no wait
                nil,    // args
    )
	if err != nil {
		return fmt.Errorf("ch.Consume %w", err)
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

func EventHandler() {
	go func(){
		for msg := range msgs {
			var order model.UserOrder

			switch msg.RoutingKey {
			case "event.payment.success":
				gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&order)

				err := events.EventPaymentSuccessful(order)
				if err != nil{
					fmt.Println(err)
				}
			}
		}
	}()
}