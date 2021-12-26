package eventHandlers

import (
	"context"
	"fmt"

	"github.com/vincentandr/shopping-microservice/src/model"
	"github.com/vincentandr/shopping-microservice/src/services/catalog/catalogdb"
)

func EventPaymentSuccessful(order model.UserOrder) error {
	err := catalogdb.UpdateProducts(context.Background(), order.Items)
	if err != nil {
		return fmt.Errorf("failed to execute remove cart items event payment: %v", err)
	}

	return nil
}