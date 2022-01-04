package cartMock

import (
	"context"

	"github.com/stretchr/testify/mock"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
	"google.golang.org/grpc"
)

type GrpcMock struct {
	mock.Mock
}

func (m *GrpcMock) Grpc_GetOrders(ctx context.Context, in *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error) {
	args := m.Called(ctx, in, opts)
	
	return args.Get(0).(*pb.GetOrdersResponse), args.Error(1)
}
func (m *GrpcMock) Grpc_PaymentCheckout(ctx context.Context, in *pb.CheckoutRequest, opts ...grpc.CallOption) (*pb.CheckoutResponse, error) {
	args := m.Called(ctx, in, opts)
	
	return args.Get(0).(*pb.CheckoutResponse), args.Error(1)
}
func (m *GrpcMock) Grpc_MakePayment(ctx context.Context, in *pb.PaymentRequest, opts ...grpc.CallOption) (*pb.PaymentResponse, error) {
	args := m.Called(ctx, in, opts)
	
	return args.Get(0).(*pb.PaymentResponse), args.Error(1)
}