package cartdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	dbConn *redis.Client
)

func NewDb() {
	dbConn = redis.NewClient(&redis.Options{
        Addr:     "localhost:6000",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
}

func Disconnect() error{
	err := dbConn.Close()

	if err != nil {
		return fmt.Errorf("failed to disconnect db connection: %v", err)
	}

	return nil
}

func GetCartItems(ctx context.Context, userId string) (map[string]string, error) {
	res, err := dbConn.HGetAll(ctx, userId).Result()
	if err != nil{
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}

	return res, nil
}

func AddOrUpdateCart(ctx context.Context, userId string, productId int, newQty int) (map[string]string, error) {
    err := dbConn.HSet(ctx, userId, productId, newQty).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot set hash cart item to redis: %v", err)
    }

    res, err := GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}

func RemoveItemFromCart(ctx context.Context, userId string, productId int) (map[string]string, error) {
    err := dbConn.HDel(ctx, userId, strconv.Itoa(productId)).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot delete cart item from redis: %v", err)
    }

    res, err := GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}

func RemoveAllCartItems(ctx context.Context, userId string) (map[string]string, error) {
    err := dbConn.Del(ctx, userId).Err()
    if err != nil {
        return nil, fmt.Errorf("cannot delete cart items from redis: %v", err)
    }

    res, err := GetCartItems(ctx, userId)
    if err != nil {
        return nil, err
    }

	return res, nil
}