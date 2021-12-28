package cartHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
)

var (
	Client pb.CartServiceClient
)

// Cart service

func GetCartItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := Client.GetCartItems(ctx, &pb.GetCartItemsRequest{UserId: args["userId"]})
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

	qty, err := strconv.Atoi(r.URL.Query().Get("qty"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	items, err := Client.AddOrUpdateCart(ctx, &pb.AddOrUpdateCartRequest{UserId: args["userId"], ProductId: args["productId"], NewQty: int32(qty)})
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

	items, err := Client.RemoveItemFromCart(ctx, &pb.RemoveItemFromCartRequest{UserId: args["userId"], ProductId: args["productId"]})
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

	items, err := Client.RemoveAllCartItems(ctx, &pb.RemoveAllCartItemsRequest{UserId: args["userId"]})
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

func Checkout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	orderId, err := Client.Checkout(ctx, &pb.CheckoutRequest{UserId: args["userId"]})
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