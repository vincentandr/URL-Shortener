package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	db "github.com/vincentandr/shopping-microservice/cmd/payment/internal/db"
	rmqPayment "github.com/vincentandr/shopping-microservice/cmd/payment/internal/pubsub"
	"github.com/vincentandr/shopping-microservice/internal/model"
	"github.com/vincentandr/shopping-microservice/internal/mongodb"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
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

func (s *Server) GetOrders(ctx context.Context, in *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	// Get all orders from db
	cursor, err := action.GetOrders(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	orders := pb.GetOrdersResponse{}
	
	res := model.Order{}

	for cursor.Next(ctx) {
		// Convert document to above struct
		err := cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}

		var items []*pb.ItemResponse

		for _, item := range res.Items {
			var temp pb.ItemResponse
			
			// Convert model.Product to pb.ItemResponse
			temp.ProductId = item.Product_id.Hex()
			temp.Name = item.Name
			temp.Price = item.Price
			temp.Qty = int32(item.Qty)

			items = append(items, &temp)
		}

		order := &pb.GetOrderResponse{OrderId: res.Order_id.Hex(), UserId: res.User_id, Items: items, Status: res.Status}

		orders.Orders = append(orders.Orders, order)
	}
	
	return &orders, nil
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
	order, err := action.GetItemsFromOrder(ctx, orderId)
	if err != nil{
		return nil, err
	}

	// Fire event to product catalog reducing qty and to cart emptying user cart
	err = rmqPayment.PaymentSuccessfulEventPublish(rmqClient.Channel, order)
	if err != nil {
		return nil, err
	}

	// Change order status
	err = action.MakePayment(ctx, orderId)
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}
	
	// Create mongodb database
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, err:= mongodb.NewDb(ctx, os.Getenv("MONGODB_PAYMENT_DB_NAME"))
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
		if err = rmqClient.CloseChannel(); err != nil {
			fmt.Println(err)
		}
		if err = rmqClient.CloseConn(); err != nil {
			fmt.Println(err)
		}
	}()

	// gRPC
	lis, err := net.Listen("tcp", os.Getenv("GRPC_PAYMENT_PORT"))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}