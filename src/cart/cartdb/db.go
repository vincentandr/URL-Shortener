package cartdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	dbConn *redis.Client
	ctx context.Context = context.Background()
)

// Fields must have capital letter to be exported and used in another package
type Product struct {
	Product_id int
	Name string
	Qty int
}

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
		log.Fatalln(err)
	}

	return err
}

func GetCartItems(userId string) ([]Product, error) {
	_, err := dbConn.Get(ctx, userId).Result()
	if err != nil{
		log.Fatalln(err)
		return nil, err
	}

	return []Product{}, nil
}

func AddOrUpdateCart(userId string, productId string, newQty int32) ([]Product, error) {
	err := dbConn.HSet(ctx, userId, productId, newQty).Err()
    if err != nil {
        log.Fatalln(err)
    }

	return []Product{}, nil
}

func ExampleClient() {
    err := dbConn.Set(ctx, "key", "value",  2 * time.Hour).Err()
    if err != nil {
        panic(err)
    }

    val, err := dbConn.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := dbConn.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
    // Output: key value
    // key2 does not exist
}