package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Product struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Qty int `json:"quantity"`
}

// Frontend service

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// Product catalog service

func GetProducts(w http.ResponseWriter, r *http.Request) {
	prods := make([]Product, 5)

	for i, _ := range prods {
		prods[i].Id = strconv.Itoa(i)
		prods[i].Name = "product " + strconv.Itoa(i)
		prods[i].Qty = i
		i++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(prods); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

}

// Payment service

func Payment(w http.ResponseWriter, r *http.Request) {

}

// Cart service

func GetCartItems(w http.ResponseWriter, r *http.Request) {

}

func AddToCart(w http.ResponseWriter, r *http.Request) {

}