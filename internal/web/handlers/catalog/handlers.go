package catalogHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
)

var (
	Client pb.CatalogServiceClient
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	// Call function in catalog_client
	products, err := Client.GetProducts(ctx, &pb.EmptyRequest{})
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

	products, err := Client.GetProductsByName(ctx, &pb.GetProductsByNameRequest{Name: name})
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
