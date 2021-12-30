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
	catalogGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/catalog"
	paymentGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/payment"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	catalogpb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	paymentpb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
	rds "github.com/vincentandr/shopping-microservice/internal/redis"
	"google.golang.org/grpc"
)

var (
	catalogClient catalogpb.CatalogServiceClient
	paymentClient paymentpb.PaymentServiceClient
	action *db.Action
)

type Server struct {
	pb.UnimplementedCartServiceServer
}

func (s *Server) GetCartItems(ctx context.Context, in *pb.GetCartItemsRequest) (*pb.ItemsResponse, error) {
	// Get product ids and its quantity in cart by userId
	res, err := action.GetCartItems(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

	// Get Product ID Keys from map
	ids, err := GetMapKeys(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := catalogClient.GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := AppendItemToResponse(products, res)
	if err != nil{
		return nil, err
	}

	return items, nil
}

func (s *Server) AddOrUpdateCart(ctx context.Context, in *pb.AddOrUpdateCartRequest) (*pb.ItemsResponse, error) {
	res, err := action.AddOrUpdateCart(ctx, in.UserId, in.ProductId, int(in.NewQty))
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

	// Get Product ID Keys from map
	ids, err := GetMapKeys(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := catalogClient.GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := AppendItemToResponse(products, res)
	if err != nil{
		return nil, err
	}

    return items, nil
}

func (s *Server) RemoveItemFromCart(ctx context.Context, in *pb.RemoveItemFromCartRequest) (*pb.ItemsResponse, error) {
	res, err := action.RemoveItemFromCart(ctx, in.UserId, in.ProductId)
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

    // Get Product ID Keys from map
	ids, err := GetMapKeys(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := catalogClient.GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := AppendItemToResponse(products, res)
	if err != nil{
		return nil, err
	}

    return items, nil
}

func (s *Server) RemoveAllCartItems(ctx context.Context, in *pb.RemoveAllCartItemsRequest) (*pb.ItemsResponse, error) {
	res, err := action.RemoveAllCartItems(ctx, in.UserId)
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

    // Get Product ID Keys from map
	ids, err := GetMapKeys(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := catalogClient.GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := AppendItemToResponse(products, res)
	if err != nil{
		return nil, err
	}

    return items, nil
}

func (s *Server) Checkout(ctx context.Context, in *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// Get user id's cart items
	res, err := s.GetCartItems(ctx, &pb.GetCartItemsRequest{UserId: in.UserId})
	if err != nil{
		return nil, err
	}

	// RPC call to payment checkout to create order and return order id
	itemsForOrder := make([]*paymentpb.ItemResponse, len(res.Products))

	for i, item := range res.Products {
		itemsForOrder[i] = &paymentpb.ItemResponse{ProductId: item.ProductId, Name: item.Name, Price: item.Price, Qty: item.Qty}
	}

	response, err := paymentClient.PaymentCheckout(ctx, &paymentpb.CheckoutRequest{UserId: in.UserId, Items: itemsForOrder})
	if err != nil{
		return nil, err
	}

	return &pb.CheckoutResponse{OrderId: response.OrderId}, nil
}

func GetMapKeys(hm map[string]string) ([]string, error) {
	ids := make([]string, len(hm))
	i := 0
	for k := range hm {
		ids[i] = k
		i++
	}

	return ids, nil
}

func AppendItemToResponse(catalogRes *catalogpb.GetProductsByIdsResponse, hm map[string]string) (*pb.ItemsResponse, error){
	var items pb.ItemsResponse

	for _, prod := range catalogRes.Products {
		qty, err := strconv.Atoi(hm[prod.ProductId])
		if err != nil{
			return nil, fmt.Errorf("failed to convert qty from string to int: %v", err)
		}

		item := &pb.ItemResponse{ProductId: prod.ProductId, Name: prod.Name, Price: prod.Price, Qty: int32(qty)}
		items.Products = append(items.Products, item)
	}

	return &items, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}

    // Init Redis db
	rdb := rds.NewDb()
	defer func(){
		if err = rdb.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Database actions
	action = db.NewAction(rdb.Conn)

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

	// Listen to rabbitmq events and handle them
	consumer.EventHandler(action)
    
    // gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	catalogRpc, err := catalogGrpc.NewGrpcClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = catalogRpc.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
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

	catalogClient = catalogRpc.Client
	paymentClient = paymentRpc.Client

	lis, err := net.Listen("tcp", os.Getenv("GRPC_CART_PORT"))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCartServiceServer(s, &Server{})

	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}