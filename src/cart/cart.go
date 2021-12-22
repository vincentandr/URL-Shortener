package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	db "github.com/vincentandr/shopping-microservice/src/cart/cartdb"
	pb "github.com/vincentandr/shopping-microservice/src/cart/cartpb"
	"github.com/vincentandr/shopping-microservice/src/cart/clients"
	"github.com/vincentandr/shopping-microservice/src/catalog/catalogpb"
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

	// Convert product ids in cart to int
	ids, err := convIdsStrToInt(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := clients.GetProductsByIds(ctx, ids)
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := appendItemToResponse(products, res)
	if err != nil{
		return nil, err
	}

	return items, nil
}

func (s *Server) AddOrUpdateCart(ctx context.Context, in *pb.AddOrUpdateCartRequest) (*pb.ItemsResponse, error) {
	res, err := db.AddOrUpdateCart(ctx, in.UserId, int(in.ProductId), int(in.NewQty))
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

	// Convert product ids in cart to int
	ids, err := convIdsStrToInt(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := clients.GetProductsByIds(ctx, ids)
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := appendItemToResponse(products, res)
	if err != nil{
		return nil, err
	}

    return items, nil
}

func (s *Server) RemoveItemFromCart(ctx context.Context, in *pb.RemoveItemFromCartRequest) (*pb.ItemsResponse, error) {
	res, err := db.RemoveItemFromCart(ctx, in.UserId, int(in.ProductId))
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

    // Convert product ids in cart to int
	ids, err := convIdsStrToInt(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := clients.GetProductsByIds(ctx, ids)
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := appendItemToResponse(products, res)
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

    // Convert product ids in cart to int
	ids, err := convIdsStrToInt(res)
	if err != nil{
		return nil, err
	}

	// RPC call catalog server to get cart products' names
	products, err := clients.GetProductsByIds(ctx, ids)
	if err != nil{
		return nil, err
	}

	// Return response in format product id, product name, and qty in cart
	items, err := appendItemToResponse(products, res)
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

	// RPC call to checkout to create order and return order id
	response, err := clients.PaymentCheckout(ctx, in.UserId, res)
	if err != nil{
		return nil, err
	}

	return &pb.CheckoutResponse{OrderId: response.OrderId}, nil
}

func convIdsStrToInt(hm map[string]string) ([]int, error) {
	ids := make([]int, len(hm))
	i := 0
	for k, _ := range hm {
		val, err := strconv.Atoi(k)
		if err != nil{
			return nil, fmt.Errorf("failed to convert id from string to int: %v", err)
		}

		ids[i] = val
		i++
	}

	return ids, nil
}

func appendItemToResponse(catalogRes *catalogpb.GetProductsByIdsResponse, hm map[string]string) (*pb.ItemsResponse, error){
	var items pb.ItemsResponse

	for _, prod := range catalogRes.Products {
		strProdId := strconv.Itoa(int(prod.ProductId))
		qty, err := strconv.Atoi(hm[strProdId])
		if err != nil{
			return nil, fmt.Errorf("failed to convert id from string to int: %v", err)
		}

		item := &pb.ItemResponse{ProductId: prod.ProductId, Name: prod.Name, Price: prod.Price, Qty: int32(qty)}
		items.Products = append(items.Products, item)
	}

	return &items, nil
}

func main() {
    // Init Redis Db
	db.NewDb()
    
    // gRPC
	err := clients.NewCatalogClient()
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectCatalogClient(); err != nil {
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