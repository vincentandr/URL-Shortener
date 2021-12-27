package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vincentandr/shopping-microservice/src/internal/mongodb"
	rbmq "github.com/vincentandr/shopping-microservice/src/internal/rabbitmq"
	"github.com/vincentandr/shopping-microservice/src/model"
	db "github.com/vincentandr/shopping-microservice/src/services/catalog/catalogdb"
	pb "github.com/vincentandr/shopping-microservice/src/services/catalog/catalogpb"
	rmqConsumer "github.com/vincentandr/shopping-microservice/src/services/catalog/rmq-consumer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

var (
	action *db.Action
)

type Server struct {
	pb.UnimplementedCatalogServiceServer
}

func (s *Server) GetProducts(ctx context.Context, in *pb.EmptyRequest) (*pb.GetProductsResponse, error) {
	cursor, err := action.GetProducts(ctx)
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

	cursor, err := action.GetProductsByIds(ctx, productIds)
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
	cursor, err := action.GetProductsByName(ctx, in.Name)
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
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v", err)
	}

	// Establish connection to mysql db
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, err:= mongodb.NewDb(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	err = client.Conn.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	action, err = db.NewAction(client)
	if err != nil {
		fmt.Println(err)
	}

	// Seed collection
	err = action.SeedCollection(context.Background())
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
	consumer, err := rmqConsumer.NewConsumer(rmqClient)
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
	lis, err := net.Listen("tcp", os.Getenv("GRPC_CATALOG_PORT"))
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