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

	var r0 *pb.GetOrdersResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.GetOrdersResponse)
	}
	
	return r0, args.Error(1)
}
func (m *GrpcMock) Grpc_GetDraftOrder(ctx context.Context, in *pb.GetDraftOrderRequest, opts ...grpc.CallOption) (*pb.GetOrderResponse, error) {
	args := m.Called(ctx, in, opts)

	var r0 *pb.GetOrderResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.GetOrderResponse)
	}
	
	return r0, args.Error(1)
}
func (m *GrpcMock) Grpc_PaymentCheckout(ctx context.Context, in *pb.CheckoutRequest, opts ...grpc.CallOption) (*pb.CheckoutResponse, error) {
	args := m.Called(ctx, in, opts)
	
	var r0 *pb.CheckoutResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.CheckoutResponse)
	}
	
	return r0, args.Error(1)
}
func (m *GrpcMock) Grpc_MakePayment(ctx context.Context, in *pb.PaymentRequest, opts ...grpc.CallOption) (*pb.PaymentResponse, error) {
	args := m.Called(ctx, in, opts)
	
	var r0 *pb.PaymentResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.PaymentResponse)
	}
	
	return r0, args.Error(1)
}