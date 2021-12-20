package main

import (
	"context"
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
	products, err := db.GetProducts()
	if err != nil {
		log.Fatalln(err)
	}

	var res pb.GetProductsResponse

	for i, _ := range products {
		product := &pb.GetProductResponse{ProductId: int32(products[i].Product_id), Name: products[i].Name, Qty: int32(products[i].Qty)}
		res.Products = append(res.Products, product)
	}

	return &res, nil
}

func (s *Server) GetProductsWithName(ctx context.Context, in *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	log.Printf("Received: %v", in.GetName())

	products, err := db.GetProductsWithName(in.Name)
	if err != nil {
		log.Fatalln(err)
	}

	var res pb.GetProductsResponse

	for i, _ := range products {
		product := &pb.GetProductResponse{ProductId: int32(products[i].Product_id), Name: products[i].Name, Qty: int32(products[i].Qty)}
		res.Products = append(res.Products, product)
	}

	return &res, nil
}

func main() {
	// Establish connection to mysql db
	db.NewDb()
	defer func(){
		if err := db.Disconnect(); err != nil{
			log.Fatalln(err)
		}
	}()

	// Create new schema and table seeds
	// err := db.InitSchema()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// gRPC
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	
}