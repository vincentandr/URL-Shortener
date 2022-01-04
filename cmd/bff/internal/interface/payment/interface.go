package paymentInterface

import (
	"context"

	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
	"google.golang.org/grpc"
)

type IPaymentGrpcClient interface {
	Grpc_GetOrders(ctx context.Context, in *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error)
	Grpc_PaymentCheckout(ctx context.Context, in *pb.CheckoutRequest, opts ...grpc.CallOption) (*pb.CheckoutResponse, error)
	Grpc_MakePayment(ctx context.Context, in *pb.PaymentRequest, opts ...grpc.CallOption) (*pb.PaymentResponse, error)
}