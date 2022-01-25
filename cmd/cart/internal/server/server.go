package server

import (
	"context"
	"fmt"
	"strconv"

	db "github.com/vincentandr/shopping-microservice/cmd/cart/internal/db"
	rmqCart "github.com/vincentandr/shopping-microservice/cmd/cart/internal/pubsub"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	catalogpb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	paymentpb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
)

type Server struct {
	pb.UnimplementedCartServiceServer
	CatalogClient catalogpb.CatalogServiceClient
	PaymentClient paymentpb.PaymentServiceClient
	Repo *db.Repository
	RmqConsumer *rmqCart.RbmqListener
}

func (s *Server) Grpc_GetCartItems(ctx context.Context, in *pb.GetCartItemsRequest) (*pb.ItemsResponse, error) {
	// Get product ids and its quantity in cart by userId
	res, err := s.Repo.GetCartItems(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

	// Get Product ID Keys from map
	ids := GetMapKeys(res)

	// RPC call catalog server to get cart products' names
	products, err := s.CatalogClient.Grpc_GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
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

func (s *Server) Grpc_AddOrUpdateCart(ctx context.Context, in *pb.AddOrUpdateCartRequest) (*pb.ItemsResponse, error) {
	res, err := s.Repo.AddOrUpdateCart(ctx, in.UserId, in.ProductId, int(in.NewQty))
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

	// Get Product ID Keys from map
	ids := GetMapKeys(res)

	// RPC call catalog server to get cart products' names
	products, err := s.CatalogClient.Grpc_GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
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

func (s *Server) Grpc_RemoveItemFromCart(ctx context.Context, in *pb.RemoveItemFromCartRequest) (*pb.ItemsResponse, error) {
	res, err := s.Repo.RemoveItemFromCart(ctx, in.UserId, in.ProductId)
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

    // Get Product ID Keys from map
	ids := GetMapKeys(res)

	// RPC call catalog server to get cart products' names
	products, err := s.CatalogClient.Grpc_GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
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

func (s *Server) Grpc_RemoveAllCartItems(ctx context.Context, in *pb.RemoveAllCartItemsRequest) (*pb.ItemsResponse, error) {
	res, err := s.Repo.RemoveAllCartItems(ctx, in.UserId)
	if err != nil{
		return nil, err
	}

	// Return empty response if there is no items in cart
	if len(res) == 0 {
		return &pb.ItemsResponse{}, nil
	}

    // Get Product ID Keys from map
	ids := GetMapKeys(res)

	// RPC call catalog server to get cart products' names
	products, err := s.CatalogClient.Grpc_GetProductsByIds(ctx, &catalogpb.GetProductsByIdsRequest{ProductIds: ids})
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

func (s *Server) Grpc_Checkout(ctx context.Context, in *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// Get user id's cart items
	res, err := s.Grpc_GetCartItems(ctx, &pb.GetCartItemsRequest{UserId: in.UserId})
	if err != nil{
		return nil, err
	}

	// RPC call to payment checkout to create order and return order id
	itemsForOrder := make([]*paymentpb.ItemResponse, len(res.Products))

	for i, item := range res.Products {
		itemsForOrder[i] = &paymentpb.ItemResponse{ProductId: item.ProductId, Name: item.Name, Price: item.Price, Qty: item.Qty}
	}

	response, err := s.PaymentClient.Grpc_PaymentCheckout(ctx, &paymentpb.CheckoutRequest{UserId: in.UserId, Items: itemsForOrder})
	if err != nil{
		return nil, err
	}

	return &pb.CheckoutResponse{OrderId: response.OrderId}, nil
}

func GetMapKeys(hm map[string]string) ([]string) {
	ids := make([]string, len(hm))
	i := 0
	for k := range hm {
		ids[i] = k
		i++
	}

	return ids
}

func AppendItemToResponse(catalogRes *catalogpb.GetProductsByIdsResponse, hm map[string]string) (*pb.ItemsResponse, error){
	items := pb.ItemsResponse{}

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