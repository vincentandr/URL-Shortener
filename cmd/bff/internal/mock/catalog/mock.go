package cartMock

import (
	"context"

	"github.com/stretchr/testify/mock"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	"google.golang.org/grpc"
)

type GrpcMock struct {
	mock.Mock
}

func (m *GrpcMock) Grpc_GetProducts(ctx context.Context, in *pb.EmptyRequest, opts ...grpc.CallOption) (*pb.GetProductsResponse, error) {
	args := m.Called(ctx, in, opts)
	
	return args.Get(0).(*pb.GetProductsResponse), args.Error(1)
}
func (m *GrpcMock) Grpc_GetProductsByIds(ctx context.Context, in *pb.GetProductsByIdsRequest, opts ...grpc.CallOption) (*pb.GetProductsByIdsResponse, error) {
	args := m.Called(ctx, in, opts)
	
	return args.Get(0).(*pb.GetProductsByIdsResponse), args.Error(1)
}
func (m *GrpcMock) Grpc_GetProductsByName(ctx context.Context, in *pb.GetProductsByNameRequest, opts ...grpc.CallOption) (*pb.GetProductsResponse, error) {
	args := m.Called(ctx, in, opts)
	
	return args.Get(0).(*pb.GetProductsResponse), args.Error(1)
}
