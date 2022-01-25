package server

import (
	"context"
	"fmt"

	db "github.com/vincentandr/shopping-microservice/cmd/payment/internal/db"
	rmqPayment "github.com/vincentandr/shopping-microservice/cmd/payment/internal/pubsub"
	"github.com/vincentandr/shopping-microservice/internal/model"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
	Repo *db.Repository
	RmqClient *rbmq.Rabbitmq
}

func (s *Server) Grpc_GetOrders(ctx context.Context, in *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	// Get all orders from db
	cursor, err := s.Repo.GetOrders(ctx, in.UserId)
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

func (s *Server) Grpc_PaymentCheckout(ctx context.Context, in *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// Change order status to draft
	orderId, err := s.Repo.Checkout(ctx, in.UserId, in.Items)
	if err != nil {
		return nil, err
	}

	return &pb.CheckoutResponse{OrderId: orderId}, nil
}

func (s *Server) Grpc_MakePayment(ctx context.Context, in *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	// Get order document
	orderId, err := primitive.ObjectIDFromHex(in.OrderId)
	if err != nil{
		return nil, fmt.Errorf("failed to convert from hex to objectID: %v", err)
	}
	order, err := s.Repo.GetItemsFromOrder(ctx, orderId)
	if err != nil{
		return nil, err
	}

	// Fire event to product catalog reducing qty and to cart emptying user cart
	err = rmqPayment.PaymentSuccessfulEventPublish(s.RmqClient.Channel, order)
	if err != nil {
		return nil, err
	}

	// Change order status
	err = s.Repo.MakePayment(ctx, orderId)
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{}, nil
}