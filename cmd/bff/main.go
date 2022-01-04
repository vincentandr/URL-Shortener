package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	cartHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/cart"
	catalogHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/catalog"
	paymentHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/payment"
	routes "github.com/vincentandr/shopping-microservice/cmd/bff/internal/route"
	cartGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/cart"
	catalogGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/catalog"
	paymentGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/payment"
)

type Server struct {
	Router *routes.Router
	Port string
}

func NewHTTPServer(r *routes.Router, port string) *Server{
	return &Server{Router: r, Port: port}
}

func (s *Server) HTTPListenAndServe() {
	http.ListenAndServe(s.Port, s.Router)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
	}

	// RPC
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	catalogRpc, err := catalogGrpc.NewGrpcClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = catalogRpc.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	cartRpc, err := cartGrpc.NewGrpcClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = cartRpc.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	paymentRpc, err := paymentGrpc.NewGrpcClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func(){
		if err = paymentRpc.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(cartRpc.Client),
		"catalog": catalogHandlers.NewGrpcClient(catalogRpc.Client),
		"payment": paymentHandlers.NewGrpcClient(paymentRpc.Client),
	}

	// Create router
	r := routes.NewRouter(handlerMap)
	
	// Create server
	srv := NewHTTPServer(r, os.Getenv("WEB_SERVER_PORT"))

	fmt.Println("Server starting at port " + srv.Port)

	srv.HTTPListenAndServe()
}