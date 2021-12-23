// gRPC Cart Client in BFF

package clients

import (
	"context"
	"fmt"

	pb "github.com/vincentandr/shopping-microservice/src/cart/cartpb"
	"google.golang.org/grpc"
)

const (
	cartRpcPort = ":50052"
)

var (
	cartClientConn *grpc.ClientConn
	cartClient pb.CartServiceClient
)

func NewCartClient() error{
	cartClientConn, err := grpc.Dial("localhost" + cartRpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to RPC client: %v", err)
	}

	cartClient = pb.NewCartServiceClient(cartClientConn)

	return err
}

func DisconnectCartClient() error{
	err := cartClientConn.Close()

	if err != nil{
		return fmt.Errorf("failed to disconnect from RPC client: %v", err)
	}

	return err
}

func GetCartItems(ctx context.Context, userId string) (*pb.ItemsResponse , error){
	items, err := cartClient.GetCartItems(ctx, &pb.GetCartItemsRequest{UserId: userId})
	if err != nil{
		return nil, err
	}

	return items, nil

}

func AddOrUpdateCartQty(ctx context.Context, userId string, productId string, qty int) (*pb.ItemsResponse , error){
	items, err := cartClient.AddOrUpdateCart(ctx, &pb.AddOrUpdateCartRequest{UserId: userId, ProductId: productId, NewQty: int32(qty)})
	if err != nil{
		return nil, err
	}

	return items, nil

}

func RemoveCartItem(ctx context.Context, userId string, productId string) (*pb.ItemsResponse , error){
	items, err := cartClient.RemoveItemFromCart(ctx, &pb.RemoveItemFromCartRequest{UserId: userId, ProductId: productId})
	if err != nil{
		return nil, err
	}

	return items, nil

}

func RemoveAllCartItems(ctx context.Context, userId string) (*pb.ItemsResponse , error){
	items, err := cartClient.RemoveAllCartItems(ctx, &pb.RemoveAllCartItemsRequest{UserId: userId})
	if err != nil{
		return nil, err
	}

	return items, nil

}

func Checkout(ctx context.Context, userId string) (*pb.CheckoutResponse , error){
	res, err := cartClient.Checkout(ctx, &pb.CheckoutRequest{UserId: userId})
	if err != nil {
		return nil, err
	}

	return res, nil
}