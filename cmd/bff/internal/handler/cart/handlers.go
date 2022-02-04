package cartHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	intf "github.com/vincentandr/shopping-microservice/cmd/bff/internal/interface/cart"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
)

// pb.CartServiceClient pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"

type GrpcClient struct {
	Client intf.ICartGrpcClient
}

func NewGrpcClient(client intf.ICartGrpcClient) *GrpcClient {
	return &GrpcClient{Client: client}
}

// Cart service
func (c *GrpcClient) GetCartItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := c.Client.Grpc_GetCartItems(ctx, &pb.GetCartItemsRequest{UserId: args["userId"]})
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

func (c *GrpcClient) AddOrUpdateCartQty(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	qty, err := strconv.Atoi(r.URL.Query().Get("qty"))
	// Check if input qty is non numeric and <= 0
	if err != nil || qty <= 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	items, err := c.Client.Grpc_AddOrUpdateCart(ctx, &pb.AddOrUpdateCartRequest{
		UserId: args["userId"],
		ProductId: args["productId"],
		NewQty: int32(qty),
	})
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

func (c *GrpcClient) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := c.Client.Grpc_RemoveItemFromCart(ctx, &pb.RemoveItemFromCartRequest{UserId: args["userId"], ProductId: args["productId"]})
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

func (c *GrpcClient) RemoveAllCartItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := c.Client.Grpc_RemoveAllCartItems(ctx, &pb.RemoveAllCartItemsRequest{UserId: args["userId"]})
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

func (c *GrpcClient) Checkout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	orderId, err := c.Client.Grpc_Checkout(ctx, &pb.CheckoutRequest{UserId: args["userId"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(orderId); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}