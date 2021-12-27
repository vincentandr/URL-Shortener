package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vincentandr/shopping-microservice/src/internal/mongodb"
	rbmq "github.com/vincentandr/shopping-microservice/src/internal/rabbitmq"
	db "github.com/vincentandr/shopping-microservice/src/services/payment/paymentdb"
	pb "github.com/vincentandr/shopping-microservice/src/services/payment/paymentpb"
	rmqPublisher "github.com/vincentandr/shopping-microservice/src/services/payment/rmq-publisher"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

var (
	rmqClient *rbmq.Rabbitmq
	action *db.Action
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *Server) PaymentCheckout(ctx context.Context, in *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// Change order status to draft
	orderId, err := action.Checkout(ctx, in.UserId, in.Items)
	if err != nil {
		return nil, err
	}

	return &pb.CheckoutResponse{OrderId: orderId}, nil
}

func (s *Server) MakePayment(ctx context.Context, in *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	// Get order document
	orderId, err := primitive.ObjectIDFromHex(in.OrderId)
	if err != nil{
		return nil, fmt.Errorf("failed to convert from hex to objectID: %v", err)
	}
	order, err := action.GetOrder(ctx, orderId)
	if err != nil{
		return nil, err
	}

	// Fire event to product catalog reducing qty and to cart emptying user cart
	err = rmqPublisher.PaymentSuccessfulEventPublish(rmqClient.Channel, order)
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v", err)
	}
	
	// Create mongodb database
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, err:= mongodb.NewDb(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	action, err = db.NewAction(client)
	if err != nil {
		fmt.Println(err)
	}

	// RabbitMQ client
	rmqClient, err = rbmq.NewRabbitMQ()
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = rmqClient.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// gRPC
	lis, err := net.Listen("tcp", os.Getenv("GRPC_PAYMENT_PORT"))
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}