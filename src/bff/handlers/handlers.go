package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/vincentandr/shopping-microservice/src/bff/clients"
)
	
func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	// Call function in catalog_client
	products, err := clients.GetProducts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return products as json response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(products); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetProductsByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")

	products, err := clients.GetProductsByName(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := clients.GetCartItems(ctx, args["userId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(items); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddOrUpdateCartQty(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	productId, err := strconv.Atoi(args["productId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qty, err := strconv.Atoi(r.URL.Query().Get("qty"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	

	items, err := clients.AddOrUpdateCartQty(ctx, args["userId"], productId, qty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(items); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	productId, err := strconv.Atoi(args["productId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items, err := clients.RemoveCartItem(ctx, args["userId"], productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(items); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RemoveAllCartItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := clients.RemoveAllCartItems(ctx, args["userId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(items); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}