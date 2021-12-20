// gRPC Catalog Client in BFF

package clients

import (
	"context"
	"log"
	"time"

	pb "github.com/vincentandr/shopping-microservice/src/catalog/catalogpb"
	"google.golang.org/grpc"
)

const (
	catalogRpcPort = ":50051"
)

var (
	catalogClientConn *grpc.ClientConn
	catalogClient pb.CatalogServiceClient
)

func NewCatalogClient() error{
	catalogClientConn, err := grpc.Dial("localhost" + catalogRpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	catalogClient = pb.NewCatalogServiceClient(catalogClientConn)

	return err
}

func DisconnectCatalogClient() error{
	err := catalogClientConn.Close()

	if err != nil{
		log.Fatalln(err)
	}

	return err
}

func GetProducts() (*pb.GetProductsResponse , error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	products, err := catalogClient.GetProducts(ctx, &pb.EmptyRequest{})
	if err != nil{
		log.Fatalln(err)
		return nil, err
	}

	return products, nil

}

func GetProductsWithName(name string) (*pb.GetProductsResponse , error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	products, err := catalogClient.GetProductsWithName(ctx, &pb.GetProductsRequest{Name: name})
	if err != nil{
		log.Fatalln(err)
		return nil, err
	}

	return products, nil

}