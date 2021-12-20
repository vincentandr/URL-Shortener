package main

import (
	"fmt"
	"log"
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
		panic("Cannot create catalog client")
	}
	defer func(){
		if err = clients.DisconnectCatalogClient(); err != nil {
			log.Fatalln("Could not disconnect catalog client")
		}
	}()

	http.ListenAndServe(port, r)
}