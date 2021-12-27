package rbmq

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

// RabbitMQ ...
type Rabbitmq struct {
	Conn *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQ instantiates the RabbitMQ instances using configuration defined in environment variables.
func NewRabbitMQ() (*Rabbitmq, error) {
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	hostname := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, hostname, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("amqp.Dial %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("conn.Channel %w", err)
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
		return nil, fmt.Errorf("ch.ExchangeDeclare %w", err)
	}

	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, fmt.Errorf("ch.Qos %w", err)
	}

	return &Rabbitmq{
		Conn: conn,
		Channel: ch,
	}, nil
}

func (r* Rabbitmq) CancelConsumerDelivery(tag string) error {
	if err := r.Channel.Cancel(tag, false); err != nil{
		return fmt.Errorf("Rabbitmq cancel channel error: %v", err)
	}

	return nil
}

func (r* Rabbitmq) CloseConn() error {
	if err := r.Conn.Close(); err != nil{
		return fmt.Errorf("Rabbitmq close connection error: %v", err)
	}

	return nil
}

func (r* Rabbitmq) CloseChannel() error {
	if err := r.Channel.Close(); err != nil{
		return fmt.Errorf("Rabbitmq close channel error: %v", err)
	}

	return nil
}
