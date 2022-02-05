package paymentHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	intf "github.com/vincentandr/shopping-microservice/cmd/bff/internal/interface/payment"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
)

type GrpcClient struct {
	Client intf.IPaymentGrpcClient
}

func NewGrpcClient(client intf.IPaymentGrpcClient) *GrpcClient {
	return &GrpcClient{Client: client}
}

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

func (c *GrpcClient) GetDraftOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	order, err := c.Client.Grpc_GetDraftOrder(ctx, &pb.GetDraftOrderRequest{UserId: args["userId"]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client_secret, err := CreateStripePaymentIntent(order.Subtotal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct{
		Order *pb.GetOrderResponse `json:"order"`
		Client_secret string `json:"client_secret"`
	}{Order: order, Client_secret: client_secret}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(res); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *GrpcClient) MakePayment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req *pb.PaymentRequest

	err = json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.OrderId = args["orderId"]
	
	items, err := c.Client.Grpc_MakePayment(ctx, req)
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

func CreateStripePaymentIntent(subtotal float32) (string, error) {
	// Create a PaymentIntent with amount and currency
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(subtotal * 100)), // Stripe charges in cents so have to multiply by 100
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	pi, err := paymentintent.New(params)

	if err != nil {
		return "", fmt.Errorf("pi.New: %v", err)
	}


	return pi.ClientSecret, nil
}