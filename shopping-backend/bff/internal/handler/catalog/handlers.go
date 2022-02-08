package catalogHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	intf "github.com/vincentandr/shopping-microservice/bff/internal/interface/catalog"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
)

type GrpcClient struct {
	Client intf.ICatalogGrpcClient
}

func NewGrpcClient(client intf.ICatalogGrpcClient) *GrpcClient {
	return &GrpcClient{Client: client}
}

func (c *GrpcClient) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	// Call function in catalog_client
	products, err := c.Client.Grpc_GetProducts(ctx, &pb.EmptyRequest{})
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

func (c *GrpcClient) GetProductsByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	name := strings.ToLower(r.URL.Query().Get("name"))

	products, err := c.Client.Grpc_GetProductsByName(ctx, &pb.GetProductsByNameRequest{Name: name})
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
