package paymentGrpc

import (
	"context"
	"fmt"
	"os"

	pb "github.com/vincentandr/shopping-microservice/src/services/payment/paymentpb"
	"google.golang.org/grpc"
)

type PaymentGrpc struct {
	Conn *grpc.ClientConn
	Client pb.PaymentServiceClient
}

func NewGrpcClient(ctx context.Context) (*PaymentGrpc, error){
	target := os.Getenv("GRPC_PAYMENT_HOST") + os.Getenv("GRPC_PAYMENT_PORT")
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC client: %v", err)
	}

	paymentClient := pb.NewPaymentServiceClient(conn)

	return &PaymentGrpc{Conn: conn, Client: paymentClient}, nil
}

func (g *PaymentGrpc) Close() error{
	err := g.Conn.Close()
	if err != nil{
		return fmt.Errorf("failed to disconnect from RPC client: %v", err)
	}

	return err
}