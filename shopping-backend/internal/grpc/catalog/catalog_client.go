package catalogGrpc

import (
	"context"
	"fmt"
	"os"

	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	"google.golang.org/grpc"
)

type CatalogGrpc struct {
	Conn *grpc.ClientConn
	Client pb.CatalogServiceClient
}

func NewGrpcClient(ctx context.Context) (*CatalogGrpc, error){
	target := os.Getenv("GRPC_CATALOG_HOST") + os.Getenv("GRPC_CATALOG_PORT")
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC client: %v", err)
	}

	catalogClient := pb.NewCatalogServiceClient(conn)

	return &CatalogGrpc{Conn: conn, Client: catalogClient}, nil
}

func (g *CatalogGrpc) Close() error{
	err := g.Conn.Close()
	if err != nil{
		return fmt.Errorf("failed to disconnect from RPC client: %v", err)
	}

	return err
}