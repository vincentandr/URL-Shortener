package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	db "github.com/vincentandr/shopping-microservice/src/services/cart/cartdb"
	pb "github.com/vincentandr/shopping-microservice/src/services/cart/cartpb"
	"github.com/vincentandr/shopping-microservice/src/services/cart/clients"
	"github.com/vincentandr/shopping-microservice/src/services/cart/rmq"
	"github.com/vincentandr/shopping-microservice/src/services/catalog/catalogpb"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type Server struct {
	pb.UnimplementedCartServiceServer
}

func (s *Server) GetCartItems(ctx context.Context, in *pb.GetCartItemsRequest) (*pb.ItemsResponse, error) {
	// Get product ids and its quantity in cart by userId
	res, err := db.GetCartItems(ctx, in.UserId)
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
	products, err := clients.GetProductsByIds(ctx, ids)
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
	res, err := db.AddOrUpdateCart(ctx, in.UserId, in.ProductId, int(in.NewQty))
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
	products, err := clients.GetProductsByIds(ctx, ids)
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
	res, err := db.RemoveItemFromCart(ctx, in.UserId, in.ProductId)
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
	products, err := clients.GetProductsByIds(ctx, ids)
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
	res, err := db.RemoveAllCartItems(ctx, in.UserId)
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
	products, err := clients.GetProductsByIds(ctx, ids)
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
	response, err := clients.PaymentCheckout(ctx, in.UserId, res)
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
    // Init Redis Db
	db.NewDb()
	defer db.Disconnect()

	// RabbitMQ Publisher
	err := rmq.NewRabbitMQ()
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = rmq.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rmq.EventHandler()
    
    // gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	err = clients.NewCatalogClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectCatalogClient(); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err = clients.NewPaymentClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectPaymentClient(); err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCartServiceServer(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}