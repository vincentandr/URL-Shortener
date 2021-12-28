package paymentHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
)

var (
	Client pb.PaymentServiceClient
)

func MakePayment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3 * time.Second)
	defer cancel()

	args := mux.Vars(r)

	items, err := Client.MakePayment(ctx, &pb.PaymentRequest{OrderId: args["orderId"]})
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