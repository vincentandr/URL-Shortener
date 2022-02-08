package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vincentandr/shopping-microservice/internal/mongodb"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
	db "github.com/vincentandr/shopping-microservice/payment/internal/db"
	"github.com/vincentandr/shopping-microservice/payment/internal/server"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}
	
	// Create mongodb database
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

	// Initialize server
	srv := &server.Server{Repo: repo, RmqClient: rmqClient}

	// gRPC
	lis, err := net.Listen("tcp", os.Getenv("GRPC_PAYMENT_PORT"))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, srv)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}