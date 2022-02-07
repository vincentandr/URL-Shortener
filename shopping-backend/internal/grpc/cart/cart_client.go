package cartGrpc

import (
	"context"
	"fmt"
	"os"

	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	"google.golang.org/grpc"
)

type CartGrpc struct {
	Conn *grpc.ClientConn
	Client pb.CartServiceClient
}

func NewGrpcClient(ctx context.Context) (*CartGrpc, error){
	target := os.Getenv("GRPC_CART_HOST") + os.Getenv("GRPC_CART_PORT")
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC client: %v", err)
	}

	cartClient := pb.NewCartServiceClient(conn)

	return &CartGrpc{Conn: conn, Client: cartClient}, nil
}

func (g *CartGrpc) Close() error{
	err := g.Conn.Close()
	if err != nil{
		return fmt.Errorf("failed to disconnect from RPC client: %v", err)
	}

	return err
}