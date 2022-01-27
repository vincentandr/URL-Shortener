package server

import (
	"context"
	"fmt"

	db "github.com/vincentandr/shopping-microservice/cmd/catalog/internal/db"
	rmqCatalog "github.com/vincentandr/shopping-microservice/cmd/catalog/internal/pubsub"
	"github.com/vincentandr/shopping-microservice/internal/model"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	pb.UnimplementedCatalogServiceServer
	Repo *db.Repository
	RmqConsumer *rmqCatalog.RbmqListener
}

func (s *Server) Grpc_GetProducts(ctx context.Context, in *pb.EmptyRequest) (*pb.GetProductsResponse, error) {
	cursor, err := s.Repo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := pb.GetProductsResponse{}

	// Must have capital letter and bson tag to be able to decode properly
	res := model.Product{}

	for cursor.Next(ctx) {
		// Convert document to above struct
		err := cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}

		product := &pb.GetProductResponse{ProductId: res.Product_id.Hex(), Name: res.Name, Price: res.Price, Qty: int32(res.Qty), Desc: res.Desc, Image: res.Image}

		products.Products = append(products.Products, product)
	}

	return &products, nil
}

func (s *Server) Grpc_GetProductsByIds(ctx context.Context, in *pb.GetProductsByIdsRequest) (*pb.GetProductsByIdsResponse, error) {
	// Convert string to ObjectID for collection filter
	productIds := make([]primitive.ObjectID, len(in.ProductIds))

	for i, id := range in.ProductIds {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil{
			return nil, fmt.Errorf("failed to convert from hex to object ID: %v", err)
		}

		productIds[i] = objectId
	}

	cursor, err := s.Repo.GetProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	products := pb.GetProductsByIdsResponse{}

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

func (s *Server) Grpc_GetProductsByName(ctx context.Context, in *pb.GetProductsByNameRequest) (*pb.GetProductsResponse, error) {
	cursor, err := s.Repo.GetProductsByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	products := pb.GetProductsResponse{}

	// Must have capital letter and bson tag to be able to decode properly
	res := model.Product{}

	for cursor.Next(ctx) {
		// Convert document to above struct
		err := cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}

		product := &pb.GetProductResponse{ProductId: res.Product_id.Hex(), Name: res.Name, Price: res.Price, Qty: int32(res.Qty), Desc: res.Desc, Image: res.Image}

		products.Products = append(products.Products, product)
	}


	return &products, nil
}