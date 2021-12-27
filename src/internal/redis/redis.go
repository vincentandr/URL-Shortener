package rds

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Redis struct{
    Conn *redis.Client
}

func NewDb() *Redis {
    name, _ := strconv.Atoi(os.Getenv("REDIS_DB_NAME"))

	conn := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
        Password: os.Getenv("REDIS_DB_PASSWORD"), // no password set
        DB:       name,  // use default DB
    })

    _, err := conn.Ping(context.Background()).Result()
    if err != nil{
        fmt.Printf("failed to connect to redis db: %v", err)
        return nil
    }

    return &Redis{
        Conn: conn,
    }
}

func (r *Redis) Close() error {
    // Defer close connection
    if err := r.Conn.Close(); err != nil {
        return fmt.Errorf("Redis close connection error: %v", err)
    }
    return nil
}