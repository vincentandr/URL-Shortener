package cartInterface

import (
	"context"

	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	"google.golang.org/grpc"
)

type ICartGrpcClient interface {
	Grpc_GetCartItems(ctx context.Context, in *pb.GetCartItemsRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error)
	Grpc_AddOrUpdateCart(ctx context.Context, in *pb.AddOrUpdateCartRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error)
	Grpc_RemoveItemFromCart(ctx context.Context, in *pb.RemoveItemFromCartRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error)
	Grpc_RemoveAllCartItems(ctx context.Context, in *pb.RemoveAllCartItemsRequest, opts ...grpc.CallOption) (*pb.ItemsResponse, error)
	Grpc_Checkout(ctx context.Context, in *pb.CheckoutRequest, opts ...grpc.CallOption) (*pb.CheckoutResponse, error)
}