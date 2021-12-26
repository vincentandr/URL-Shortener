package cataloghandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vincentandr/shopping-microservice/src/services/bff/clients"
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
