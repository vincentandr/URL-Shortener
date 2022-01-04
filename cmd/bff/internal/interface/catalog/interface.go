package catalogInterface

import (
	"context"

	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	"google.golang.org/grpc"
)

type ICatalogGrpcClient interface {
	Grpc_GetProducts(ctx context.Context, in *pb.EmptyRequest, opts ...grpc.CallOption) (*pb.GetProductsResponse, error)
	Grpc_GetProductsByIds(ctx context.Context, in *pb.GetProductsByIdsRequest, opts ...grpc.CallOption) (*pb.GetProductsByIdsResponse, error)
	Grpc_GetProductsByName(ctx context.Context, in *pb.GetProductsByNameRequest, opts ...grpc.CallOption) (*pb.GetProductsResponse, error)
}