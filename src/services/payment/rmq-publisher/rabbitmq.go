package rmqPublisher

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vincentandr/shopping-microservice/src/model"
)

func PaymentSuccessfulEventPublish(ch *amqp.Channel, order model.UserOrder) error{
	// Connecting bytes encoder and decoder
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(order); err != nil {
		return fmt.Errorf("byte encoding error: %v", err)
	}

	err := ch.Publish(
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