package main

import (
	"context"
	"log"
	"net"

	db "github.com/vincentandr/shopping-microservice/src/cart/cartdb"
	pb "github.com/vincentandr/shopping-microservice/src/cart/cartpb"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type Server struct {
	pb.UnimplementedCartServiceServer
}

func (s *Server) GetCartItems(ctx context.Context, in *pb.GetCartItemsRequest) (*pb.GetCartItemsResponse, error) {
	products, err := db.GetCartItems(in.UserId)
	if err != nil {
		log.Fatalln(err)
	}

	var res pb.GetCartItemsResponse

	for i, _ := range products {
		product := &pb.GetCartItemResponse{ProductId: int32(products[i].Product_id), Name: products[i].Name, Qty: int32(products[i].Qty)}
		res.Products = append(res.Products, product)
	}

	return &res, nil
}

func (s *Server) AddOrUpdateCart(ctx context.Context, in *pb.AddOrUpdateCartRequest) (*pb.AddOrUpdateCartResponse, error) {
    return &pb.AddOrUpdateCartResponse{}, nil
}

func (s *Server) RemoveItemFromCart(ctx context.Context, in *pb.RemoveItemFromCartRequest) (*pb.RemoveItemFromCartResponse, error) {
    return &pb.RemoveItemFromCartResponse{}, nil
}

func main() {
    // Init Redis Db
	db.NewDb()
    defer func(){
        if err := db.Disconnect(); err != nil{
            log.Fatalln("Failed to close redis database connection")
        }
    }()
    
    // gRPC
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCartServiceServer(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}