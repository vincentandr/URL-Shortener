package main

import (
	"context"
	"fmt"
	"log"
	"net"

	db "github.com/vincentandr/shopping-microservice/src/catalog/catalogdb"
	pb "github.com/vincentandr/shopping-microservice/src/catalog/catalogpb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Server struct {
	pb.UnimplementedCatalogServiceServer
}

func (s *Server) GetProducts(ctx context.Context, in *pb.EmptyRequest) (*pb.GetProductsResponse, error) {
	products, err := db.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	var res pb.GetProductsResponse

	for i, _ := range products {
		product := &pb.GetProductResponse{ProductId: int32(products[i].Product_id), Name: products[i].Name, Price: products[i].Price, Qty: int32(products[i].Qty)}
		res.Products = append(res.Products, product)
	}

	return &res, nil
}

func (s *Server) GetProductsByIds(ctx context.Context, in *pb.GetProductsByIdsRequest) (*pb.GetProductsByIdsResponse, error) {
	ids := make([]int, len(in.ProductIds))
	i := 0
	for _, val := range in.ProductIds {
		ids[i] = int(val)
		i++
	}

	products, err := db.GetProductsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	var res pb.GetProductsByIdsResponse

	for i, _ := range products {
		product := &pb.GetProductByIdsResponse{ProductId: int32(products[i].Product_id), Name: products[i].Name, Price: products[i].Price,}
		res.Products = append(res.Products, product)
	}

	return &res, nil
}

func (s *Server) GetProductsByName(ctx context.Context, in *pb.GetProductsByNameRequest) (*pb.GetProductsResponse, error) {
	products, err := db.GetProductsByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	var res pb.GetProductsResponse

	for i, _ := range products {
		product := &pb.GetProductResponse{ProductId: int32(products[i].Product_id), Name: products[i].Name, Price: products[i].Price, Qty: int32(products[i].Qty)}
		res.Products = append(res.Products, product)
	}

	return &res, nil
}

func main() {
	// Establish connection to mysql db
	db.NewDb()

	// Create new schema and table seeds
	err := db.InitSchema()
	if err != nil {
		fmt.Println("failed to create schema")
	}

	// gRPC
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}

	
}