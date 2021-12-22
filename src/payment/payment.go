package main

import (
	"context"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	db "github.com/vincentandr/shopping-microservice/src/payment/paymentdb"
	pb "github.com/vincentandr/shopping-microservice/src/payment/paymentpb"
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
	// Change order status to pending

	// Check if product catalog has enough qty, if yes then reserve. If not reject

	// Fire event to product catalog reducing qty

	// Fire event to cart emptying user cart

	return &pb.PaymentResponse{}, nil
}

func main() {
	// Create mongodb database
	db.NewDb()

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