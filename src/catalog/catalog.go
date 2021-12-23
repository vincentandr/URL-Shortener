package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	db "github.com/vincentandr/shopping-microservice/src/catalog/catalogdb"
	pb "github.com/vincentandr/shopping-microservice/src/catalog/catalogpb"
	"github.com/vincentandr/shopping-microservice/src/catalog/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Server struct {
	pb.UnimplementedCatalogServiceServer
}

func (s *Server) GetProducts(ctx context.Context, in *pb.EmptyRequest) (*pb.GetProductsResponse, error) {
	cursor, err := db.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	var products pb.GetProductsResponse

	// Must have capital letter and bson tag to be able to decode properly
	res := model.Product{}

	for cursor.Next(ctx) {
		// Convert document to above struct
		err := cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}

		product := &pb.GetProductResponse{ProductId: res.Product_id.Hex(), Name: res.Name, Price: res.Price, Qty: int32(res.Qty)}

		products.Products = append(products.Products, product)
	}

	return &products, nil
}

func (s *Server) GetProductsByIds(ctx context.Context, in *pb.GetProductsByIdsRequest) (*pb.GetProductsByIdsResponse, error) {
	// Convert string to ObjectID for collection filter
	productIds := make([]primitive.ObjectID, len(in.ProductIds))

	for i, id := range in.ProductIds {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil{
			return nil, fmt.Errorf("failed to convert from hex to object ID: %v", err)
		}

		productIds[i] = objectId
	}

	cursor, err := db.GetProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	var products pb.GetProductsByIdsResponse

	// Must have capital letter and bson tag to be able to decode properly
	res := model.Product{}

	for cursor.Next(ctx) {
		// Convert document to above struct
		err := cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}

		product := &pb.GetProductByIdsResponse{ProductId: res.Product_id.Hex(), Name: res.Name, Price: res.Price}

		products.Products = append(products.Products, product)
	}


	return &products, nil
}

func (s *Server) GetProductsByName(ctx context.Context, in *pb.GetProductsByNameRequest) (*pb.GetProductsResponse, error) {
	cursor, err := db.GetProductsByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	var products pb.GetProductsResponse

	// Must have capital letter and bson tag to be able to decode properly
	res := model.Product{}

	for cursor.Next(ctx) {
		// Convert document to above struct
		err := cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}

		product := &pb.GetProductResponse{ProductId: res.Product_id.Hex(), Name: res.Name, Price: res.Price, Qty: int32(res.Qty)}

		products.Products = append(products.Products, product)
	}


	return &products, nil
}

func main() {
	// Establish connection to mysql db
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	db.NewDb(ctx)
	defer db.Disconnect(ctx)

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