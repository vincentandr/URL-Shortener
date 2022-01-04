package paymentHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	intf "github.com/vincentandr/shopping-microservice/cmd/bff/internal/interface/payment"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
)

type GrpcClient struct {
	Client intf.IPaymentGrpcClient
}

func NewGrpcClient(client intf.IPaymentGrpcClient) *GrpcClient {
	return &GrpcClient{Client: client}
}

var (
	Client pb.PaymentServiceClient
)

func (c *GrpcClient) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	orders, err := c.Client.Grpc_GetOrders(ctx, &pb.GetOrdersRequest{UserId: args["userId"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(orders); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *GrpcClient) MakePayment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := c.Client.Grpc_MakePayment(ctx, &pb.PaymentRequest{OrderId: args["orderId"]})
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