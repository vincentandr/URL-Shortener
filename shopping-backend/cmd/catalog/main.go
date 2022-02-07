package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	db "github.com/vincentandr/shopping-microservice/cmd/catalog/internal/db"
	rmqCatalog "github.com/vincentandr/shopping-microservice/cmd/catalog/internal/pubsub"
	"github.com/vincentandr/shopping-microservice/cmd/catalog/internal/server"
	"github.com/vincentandr/shopping-microservice/internal/mongodb"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}

	// Establish connection to mysql db
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, err:= mongodb.NewDb(ctx, os.Getenv("MONGODB_DB_NAME"))
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = client.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err = client.Conn.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	repo, err := db.NewRepository(client)
	if err != nil {
		fmt.Println(err)
	}

	// Seed and create index for collection
	err = repo.InitCollection(context.Background())
	if err != nil {
		fmt.Println(err)
	}

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
	consumer, err := rmqCatalog.NewConsumer(rmqClient)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = rmqClient.CancelConsumerDelivery(consumer.Tag); err != nil {
			fmt.Println(err)
		}
	}()

	// Initialize server
	srv := &server.Server{Repo: repo, RmqConsumer: consumer}

	// Listen to rabbitmq events and handle them
	srv.RmqConsumer.EventHandler(repo)

	// gRPC
	lis, err := net.Listen("tcp", os.Getenv("GRPC_CATALOG_PORT"))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, srv)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}