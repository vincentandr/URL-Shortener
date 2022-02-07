package rds

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

type Redis struct{
    Conn *redis.Client
}

func NewDb(database int) *Redis {
    opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
    if err != nil {
        fmt.Printf("failed to parse URL: %v\n", err)
    }

    opt.DB = database
    opt.Password = os.Getenv("REDIS_DB_PASSWORD")

	conn := redis.NewClient(opt)

    _, err = conn.Ping(context.Background()).Result()
    if err != nil{
        fmt.Printf("failed to connect to redis db: %v\n", err)
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