package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vincentandr/shopping-microservice/src/bff/clients"
)
	
func GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := clients.GetProducts()
	if err != nil {
		log.Fatalln("Could not get products")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(products); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetProductsWithName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	products, err := clients.GetProductsWithName(name)
	if err != nil {
		log.Fatalf("Could not get products with name %s", name)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(products); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


// Payment service

func Payment(w http.ResponseWriter, r *http.Request) {
}

// Cart service

func GetCartItems(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)

	items, err := clients.GetCartItems(args["userId"])
	if err != nil {
		log.Fatalln("Could not get cart items")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(items); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddOrUpdateCart(w http.ResponseWriter, r *http.Request) {

}

func RemoveAllCart(w http.ResponseWriter, r *http.Request) {

}