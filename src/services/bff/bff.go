package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vincentandr/shopping-microservice/src/services/bff/clients"
	"github.com/vincentandr/shopping-microservice/src/services/bff/routes"
)

var port string = ":3000"

func main() {
	fmt.Println("Server starting at port " + port)

	r := routes.NewRouter()

	// RPC

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err := clients.NewCatalogClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectCatalogClient(); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	err = clients.NewCartClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectCartClient(); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err = clients.NewpaymentClient(ctx)
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectPaymentClient(); err != nil {
			panic(err)
		}
	}()

	http.ListenAndServe(port, r)
}