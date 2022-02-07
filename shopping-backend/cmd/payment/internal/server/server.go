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

	for cursor.Next(ctx) {
		res := model.Order{}

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
			temp.Desc = item.Desc
			temp.Image = item.Image

			items = append(items, &temp)
		}

		customer := &pb.Customer{
			FirstName: res.Customer.First_name,
			LastName: res.Customer.Last_name,
			Address: res.Customer.Address,
			Email: res.Customer.Email,
			Area: res.Customer.Area,
			Postal: res.Customer.Postal,
			Phone: res.Customer.Phone,
		}

		order := &pb.GetOrderResponse{
			OrderId: res.Order_id.Hex(),
			UserId: res.User_id,
			Items: items,
			Status: res.Status,
			Subtotal: res.Subtotal,
			Customer: customer,
		}

		orders.Orders = append(orders.Orders, order)
	}
	
	return &orders, nil
}

func (s *Server) Grpc_GetDraftOrder(ctx context.Context, in *pb.GetDraftOrderRequest) (*pb.GetOrderResponse, error) {
	// Get all orders from db
	order, found, err := s.Repo.GetDraftOrder(ctx, in.UserId)
	if !found {
		// Not found
		return &pb.GetOrderResponse{}, nil
	} else if err != nil {
		// Found but decode error
		return nil, err
	}

	var items []*pb.ItemResponse

	for _, item := range order.Items {
		var temp pb.ItemResponse
		
		// Convert model.Product to pb.ItemResponse
		temp.ProductId = item.Product_id.Hex()
		temp.Name = item.Name
		temp.Price = item.Price
		temp.Qty = int32(item.Qty)
		temp.Desc = item.Desc
		temp.Image = item.Image

		items = append(items, &temp)
	}

	customer := &pb.Customer{
		FirstName: order.Customer.First_name,
		LastName: order.Customer.Last_name,
		Address: order.Customer.Address,
		Email: order.Customer.Email,
		Area: order.Customer.Area,
		Postal: order.Customer.Postal,
		Phone: order.Customer.Phone,
	}
	
	return &pb.GetOrderResponse{
		OrderId: order.Order_id.Hex(),
		UserId: order.User_id,
		Items: items,
		Status: order.Status,
		Subtotal: order.Subtotal,
		Customer: customer,
		}, nil
}

func (s *Server) Grpc_PaymentCheckout(ctx context.Context, in *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// Change order status to draft
	orderId, err := s.Repo.Checkout(ctx, in.UserId, in.Items, in.Subtotal)
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

	// Insertable struct with bson tag
	customer := model.Customer{
		First_name: in.Customer.FirstName,
		Last_name: in.Customer.LastName,
		Address: in.Customer.Address,
		Email: in.Customer.Email,
		Area: in.Customer.Area,
		Postal: in.Customer.Postal,
		Phone: in.Customer.Phone,
	}

	// Change order status
	err = s.Repo.MakePayment(ctx, orderId, customer, order.Subtotal)
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{}, nil
}