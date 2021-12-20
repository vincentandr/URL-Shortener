// gRPC Cart Client in BFF

package clients

import (
	"context"
	"log"
	"time"

	pb "github.com/vincentandr/shopping-microservice/src/cart/cartpb"
	"google.golang.org/grpc"
)

const (
	cartRpcPort = ":50051"
)

var (
	cartClientConn *grpc.ClientConn
	cartClient pb.CartServiceClient
)

func NewCartClient() error{
	cartClientConn, err := grpc.Dial("localhost" + cartRpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cartClient = pb.NewCartServiceClient(cartClientConn)

	return err
}

func DisconnectCartClient() error{
	err := cartClientConn.Close()

	if err != nil{
		log.Fatalln(err)
	}

	return err
}

func GetCartItems(userId string) (*pb.GetCartItemsResponse , error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	items, err := cartClient.GetCartItems(ctx, &pb.GetCartItemsRequest{UserId: userId})
	if err != nil{
		log.Fatalln(err)
		return nil, err
	}

	return items, nil

}