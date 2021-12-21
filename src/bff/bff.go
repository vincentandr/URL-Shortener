package main

import (
	"fmt"
	"net/http"

	"github.com/vincentandr/shopping-microservice/src/bff/clients"
	"github.com/vincentandr/shopping-microservice/src/bff/routes"
)

var port string = ":3000"

func main() {
	fmt.Println("Server starting at port " + port)

	r := routes.NewRouter()

	err := clients.NewCatalogClient()
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectCatalogClient(); err != nil {
			panic(err)
		}
	}()

	err = clients.NewCartClient()
	if err != nil {
		panic(err)
	}
	defer func(){
		if err = clients.DisconnectCartClient(); err != nil {
			panic(err)
		}
	}()

	http.ListenAndServe(port, r)
}