package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	db "github.com/vincentandr/shopping-microservice/cmd/cart/internal/db"
	rmqCart "github.com/vincentandr/shopping-microservice/cmd/cart/internal/pubsub"
	"github.com/vincentandr/shopping-microservice/cmd/cart/internal/server"
	paymentGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/payment"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
	rds "github.com/vincentandr/shopping-microservice/internal/redis"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}

    // Init Redis db
	idx, _ := strconv.Atoi(os.Getenv("REDIS_DB_INDEX"))
	rdb := rds.NewDb(idx)
	defer func(){
		if err = rdb.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Database repo
	repo := db.NewRepository(rdb.Conn)

	// RabbitMQ client
	rmqClient, err := rbmq.NewRabbitMQ()
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = rmqClient.CloseChannel(); err != nil {
			fmt.Println(err)
		}
		if err = rmqClient.CloseConn(); err != nil {
			fmt.Println(err)
		}
	}()

	// Instantiate a new rabbitmq consumer
	consumer, err := rmqCart.NewConsumer(rmqClient)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = rmqClient.CancelConsumerDelivery(consumer.Tag); err != nil {
			fmt.Println(err)
		}
	}()
    
    // gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	paymentRpc, err := paymentGrpc.NewGrpcClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = paymentRpc.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Initialize server
	srv := &server.Server{
		PaymentClient: paymentRpc.Client,
		Repo: repo,
		RmqConsumer: consumer,
	}

	// Listen to rabbitmq events and handle them
	srv.RmqConsumer.EventHandler(repo)

	lis, err := net.Listen("tcp", os.Getenv("GRPC_CART_PORT"))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCartServiceServer(s, srv)

	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}