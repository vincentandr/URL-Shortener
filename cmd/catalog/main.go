package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	db "github.com/vincentandr/shopping-microservice/cmd/catalog/internal/db"
	rmqCatalog "github.com/vincentandr/shopping-microservice/cmd/catalog/internal/pubsub"
	"github.com/vincentandr/shopping-microservice/internal/model"
	"github.com/vincentandr/shopping-microservice/internal/mongodb"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	rbmq "github.com/vincentandr/shopping-microservice/internal/rabbitmq"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedCatalogServiceServer
	repo *db.Repository
	rmqConsumer *rmqCatalog.RbmqListener
}

func (s *Server) Grpc_GetProducts(ctx context.Context, in *pb.EmptyRequest) (*pb.GetProductsResponse, error) {
	cursor, err := s.repo.GetProducts(ctx)
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

		product := &pb.GetProductResponse{ProductId: res.Product_id.Hex(), Name: res.Name, Price: res.Price, Qty: int32(res.Qty)}

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

	cursor, err := s.repo.GetProductsByIds(ctx, productIds)
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
	cursor, err := s.repo.GetProductsByName(ctx, in.Name)
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

		product := &pb.GetProductResponse{ProductId: res.Product_id.Hex(), Name: res.Name, Price: res.Price, Qty: int32(res.Qty)}

		products.Products = append(products.Products, product)
	}


	return &products, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}

	// Establish connection to mysql db
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, err:= mongodb.NewDb(ctx, os.Getenv("MONGODB_CATALOG_DB_NAME"))
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = client.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err = client.Conn.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	action, err := db.NewRepository(client)
	if err != nil {
		fmt.Println(err)
	}

	// Seed and create index for collection
	err = action.InitCollection(context.Background())
	if err != nil {
		fmt.Println(err)
	}

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
	consumer, err := rmqCatalog.NewConsumer(rmqClient)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = rmqClient.CancelConsumerDelivery(consumer.Tag); err != nil {
			fmt.Println(err)
		}
	}()

	// Initialize server
	srv := &Server{repo: action, rmqConsumer: consumer}

	// Listen to rabbitmq events and handle them
	srv.rmqConsumer.EventHandler(action)

	// gRPC
	lis, err := net.Listen("tcp", os.Getenv("GRPC_CATALOG_PORT"))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, srv)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}