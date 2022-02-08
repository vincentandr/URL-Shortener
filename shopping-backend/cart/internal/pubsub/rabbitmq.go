package rmqCart

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/streadway/amqp"
	db "github.com/vincentandr/shopping-microservice/cart/internal/db"
	"github.com/vincentandr/shopping-microservice/internal/model"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
)

// RabbitMQ ...
type RbmqListener struct {
	Msgs <-chan amqp.Delivery
	Tag string
}

// NewRabbitMQ instantiates the RabbitMQ instances using configuration defined in environment variables.
func NewConsumer(r *rbmq.Rabbitmq) (*RbmqListener, error) {
	q, err := r.Channel.QueueDeclare(
                "cartQueue",    // name
                false, // durable
                false, // delete when unused
                true,  // exclusive
                false, // no-wait
                nil,   // arguments
    )
	if err != nil {
		return nil, fmt.Errorf("amqChannel.QueueDeclare %w", err)
	}

	err = r.Channel.QueueBind(
                        q.Name,
                        "event.payment.success",
                        "tasks",
                        false,
                        nil)
	if err != nil {
		return nil, fmt.Errorf("amqChannel.ExchangeDeclare %w", err)
	}

	tag := "cartConsumer"

	msgs, err := r.Channel.Consume(
                q.Name, // queue
                tag,     // consumer
                false,   // auto ack
                false,  // exclusive
                false,  // no local
                false,  // no wait
                nil,    // args
    )
	if err != nil {
		return nil, fmt.Errorf("r.Channel.Consume %w", err)
	}

	return &RbmqListener{Msgs: msgs, Tag: tag}, nil
}

func (l *RbmqListener) EventHandler(r *db.Repository) {
	go func(){
		for msg := range l.Msgs {
			var order model.UserOrder

			switch msg.RoutingKey {
			case "event.payment.success":
				gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&order)

				err := EventPaymentSuccessful(r, order)
				if err != nil{
					fmt.Println(err)
				}
				msg.Ack(false)
			}
		}
	}()
}

func EventPaymentSuccessful(r *db.Repository, order model.UserOrder) error {
	_, err := r.RemoveAllCartItems(context.Background(), order.User_id)
	if err != nil {
		return fmt.Errorf("failed to execute remove cart items event payment: %v", err)
	}
	return nil
}