package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	cartHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handlers/cart"
	catalogHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handlers/catalog"
	paymentHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handlers/payment"
	"github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/routes"
	cartGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/cart"
	catalogGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/catalog"
	paymentGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/payment"
)

const (
	port = ":3000"
	path = ""
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}
	
	fmt.Println("Server starting at port " + port)

	r := routes.NewRouter()

	// RPC
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	catalogRpc, err := catalogGrpc.NewGrpcClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = catalogRpc.Close(); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	cartRpc, err := cartGrpc.NewGrpcClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = cartRpc.Close(); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	paymentRpc, err := paymentGrpc.NewGrpcClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = paymentRpc.Close(); err != nil {
			panic(err)
		}
	}()

	catalogHandlers.Client = catalogRpc.Client
	cartHandlers.Client = cartRpc.Client
	paymentHandlers.Client = paymentRpc.Client

	http.ListenAndServe(port, r)
}