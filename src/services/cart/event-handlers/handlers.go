package eventHandlers

import (
	"context"
	"fmt"

	"github.com/vincentandr/shopping-microservice/src/model"
	"github.com/vincentandr/shopping-microservice/src/services/cart/cartdb"
)

func EventPaymentSuccessful(order model.UserOrder) error {
	_, err := cartdb.RemoveAllCartItems(context.Background(), order.User_id)
	if err != nil {
		return fmt.Errorf("failed to execute remove cart items event payment: %v", err)
	}
	return nil
}