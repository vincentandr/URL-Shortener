// gRPC payment Client in BFF

package clients

import (
	"context"
	"fmt"

	pb "github.com/vincentandr/shopping-microservice/src/services/payment/paymentpb"
	"google.golang.org/grpc"
)

const (
	paymentRpcPort = ":50053"
)

var (
	paymentClientConn *grpc.ClientConn
	paymentClient pb.PaymentServiceClient
)

func NewpaymentClient(ctx context.Context) error{
	paymentClientConn, err := grpc.DialContext(ctx, "localhost" + paymentRpcPort, grpc.WithInsecure(), grpc.WithBlock())
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

	return nil
}

func MakePayment(ctx context.Context, orderId string) (*pb.PaymentResponse , error){
	products, err := paymentClient.MakePayment(ctx, &pb.PaymentRequest{OrderId: orderId})
	if err != nil{
		return nil, err
	}

	return products, nil

}