// gRPC payment Client in cart

package clients

import (
	"context"
	"fmt"

	cpb "github.com/vincentandr/shopping-microservice/src/cart/cartpb"
	pb "github.com/vincentandr/shopping-microservice/src/payment/paymentpb"
	"google.golang.org/grpc"
)

const (
	paymentRpcPort = ":50053"
)

var (
	paymentClientConn *grpc.ClientConn
	paymentClient pb.PaymentServiceClient
)

func NewPaymentClient() error{
	paymentClientConn, err := grpc.Dial("localhost" + paymentRpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to RPC client: %v", err)
	}

	paymentClient = pb.NewPaymentServiceClient(paymentClientConn)

	return nil
}

func DisconnectPaymentClient() error{
	err := paymentClientConn.Close()

	if err != nil{
		return fmt.Errorf("failed to disconnect from RPC client: %v", err)
	}

	return err
}

func PaymentCheckout(ctx context.Context, userId string, items *cpb.ItemsResponse) (*pb.CheckoutResponse , error){
	itemsForOrder := make([]*pb.ItemResponse, len(items.Products))

	for i, item := range items.Products {
		itemsForOrder[i] = &pb.ItemResponse{ProductId: item.ProductId, Name: item.Name, Price: item.Price, Qty: item.Qty}
	}

	paymentRes, err := paymentClient.PaymentCheckout(ctx, &pb.CheckoutRequest{UserId: userId, Items: itemsForOrder})
	if err != nil{
		return nil, err
	}

	return &pb.CheckoutResponse{OrderId: paymentRes.OrderId}, nil

}