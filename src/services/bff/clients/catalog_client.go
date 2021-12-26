// gRPC Catalog Client in BFF

package clients

import (
	"context"
	"fmt"

	pb "github.com/vincentandr/shopping-microservice/src/services/catalog/catalogpb"
	"google.golang.org/grpc"
)

const (
	catalogRpcPort = ":50051"
)

var (
	catalogClientConn *grpc.ClientConn
	catalogClient pb.CatalogServiceClient
)

func NewCatalogClient(ctx context.Context) error{
	catalogClientConn, err := grpc.DialContext(ctx, "localhost" + catalogRpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to RPC client: %v", err)
	}

	catalogClient = pb.NewCatalogServiceClient(catalogClientConn)

	return nil
}

func DisconnectCatalogClient() error{
	err := catalogClientConn.Close()

	if err != nil{
		return fmt.Errorf("failed to disconnect from RPC client: %v", err)
	}

	return nil
}

func GetProducts(ctx context.Context) (*pb.GetProductsResponse , error){
	products, err := catalogClient.GetProducts(ctx, &pb.EmptyRequest{})
	if err != nil{
		return nil, err
	}

	return products, nil

}

func GetProductsByName(ctx context.Context, name string) (*pb.GetProductsResponse , error){
	products, err := catalogClient.GetProductsByName(ctx, &pb.GetProductsByNameRequest{Name: name})
	if err != nil{
		return nil, err
	}

	return products, nil

}