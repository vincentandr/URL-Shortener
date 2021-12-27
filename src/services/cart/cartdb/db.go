package cartdb

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Action struct {
	Conn *redis.Client
}

func NewAction(conn *redis.Client) (*Action) {
    return &Action{Conn: conn}
}

func (a *Action) GetCartItems(ctx context.Context, userId string) (map[string]string, error) {
	res, err := a.Conn.HGetAll(ctx, userId).Result()
	if err != nil{
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}

	return res, nil
}

func (a *Action) AddOrUpdateCart(ctx context.Context, userId string, productId string, newQty int) (map[string]string, error) {
    err := a.Conn.HSet(ctx, userId, productId, newQty).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot set hash cart item to redis: %v", err)
    }

    res, err := a.GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}

func (a *Action) RemoveItemFromCart(ctx context.Context, userId string, productId string) (map[string]string, error) {
    err := a.Conn.HDel(ctx, userId, productId).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot delete cart item from redis: %v", err)
    }

    res, err := a.GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}

func (a *Action) RemoveAllCartItems(ctx context.Context, userId string) (map[string]string, error) {
    err := a.Conn.Del(ctx, userId).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot delete cart items from redis: %v", err)
    }

    res, err := a.GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}