package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	db "github.com/vincentandr/shopping-microservice/src/services/payment/paymentdb"
	pb "github.com/vincentandr/shopping-microservice/src/services/payment/paymentpb"
	"github.com/vincentandr/shopping-microservice/src/services/payment/rmq"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

const (
	port = ":50053"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *Server) PaymentCheckout(ctx context.Context, in *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// Change order status to draft
	orderId, err := db.Checkout(ctx, in.UserId, in.Items)
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
	order, err := db.GetOrder(ctx, orderId)
	if err != nil{
		return nil, err
	}

	// Fire event to product catalog reducing qty and to cart emptying user cart
	err = rmq.PaymentSuccessfulEventPublish(order)
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{}, nil
}

func main() {
	// Create mongodb database
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.NewDb(ctx)
	defer db.Disconnect(ctx)

	// RabbitMQ Publisher
	err := rmq.NewRabbitMQ()
	if err != nil {
		fmt.Println("RabbitMQ initialization error")
	}
	defer func(){
		if err = rmq.Close(); err != nil {
			fmt.Println("RabbitMQ close connection error")
		}
	}()

	// gRPC
	lis, err := net.Listen("tcp", port)
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