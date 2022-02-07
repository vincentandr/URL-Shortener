package cartdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Repository struct {
	Conn *redis.Client
}

func NewRepository(conn *redis.Client) (*Repository) {
    return &Repository{Conn: conn}
}

func (r *Repository) GetCartItems(ctx context.Context, userId string) (map[string]string, error) {
	res, err := r.Conn.HGetAll(ctx, userId).Result()
	if err != nil{
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}

	return res, nil
}

func (r *Repository) AddOrUpdateCart(ctx context.Context, duration time.Duration, userId string, productId string, newQty int) (map[string]string, error) {
    pipeline := r.Conn.Pipeline()

    pipeline.HSet(ctx, userId, productId, newQty)
    pipeline.Expire(ctx, userId, duration)

    _, err := pipeline.Exec(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to add cart items using pipeline: %v", err)
    }

    res, err := r.GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}

func (r *Repository) RemoveItemFromCart(ctx context.Context, userId string, productId string) (map[string]string, error) {
    err := r.Conn.HDel(ctx, userId, productId).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot delete cart item from redis: %v", err)
    }

    res, err := r.GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}

func (r *Repository) RemoveAllCartItems(ctx context.Context, userId string) (map[string]string, error) {
    err := r.Conn.Del(ctx, userId).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot delete cart items from redis: %v", err)
    }

    res, err := r.GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}