package cartMock

import (
	"context"

	"github.com/stretchr/testify/mock"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	"google.golang.org/grpc"
)

type GrpcMock struct {
	mock.Mock
}

func (m *GrpcMock) Grpc_GetCartItems(ctx context.Context, in *pb.GetCartItemsRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error) {
	args := m.Called(ctx, in, opts)

	var r0 *pb.ItemsResponse

	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.ItemsResponse)
	}
	
	return r0, args.Error(1)
}

func (m *GrpcMock) Grpc_AddOrUpdateCart(ctx context.Context, in *pb.AddOrUpdateCartRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error) {
	args := m.Called(ctx, in, opts)

	var r0 *pb.ItemsResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.ItemsResponse)
	}
	
	return r0, args.Error(1)
}

func (m *GrpcMock) Grpc_RemoveItemFromCart(ctx context.Context, in *pb.RemoveItemFromCartRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error) {
	args := m.Called(ctx, in, opts)

	var r0 *pb.ItemsResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.ItemsResponse)
	}
	
	return r0, args.Error(1)
}

func (m *GrpcMock) Grpc_RemoveAllCartItems(ctx context.Context, in *pb.RemoveAllCartItemsRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error) {
	args := m.Called(ctx, in, opts)
	
	var r0 *pb.ItemsResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.ItemsResponse)
	}
	
	return r0, args.Error(1)
}

func (m *GrpcMock) Grpc_Checkout(ctx context.Context, in *pb.CheckoutRequest, opts ...grpc.CallOption) (*pb.CheckoutResponse, error) {
	args := m.Called(ctx, in, opts)
	
	var r0 *pb.CheckoutResponse
	
	if args.Get(0) != nil {
		r0 = args.Get(0).(*pb.CheckoutResponse)
	}
	
	return r0, args.Error(1)
}
