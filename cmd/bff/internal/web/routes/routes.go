package routes

import (
	"github.com/gorilla/mux"
	cartHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handlers/cart"
	catalogHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handlers/catalog"
	paymentHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handlers/payment"
)

type Router struct {
	*mux.Router
}

// Router constructor
func NewRouter() *Router{
	r := &Router{
		Router: mux.NewRouter(),
	}
	r.routes()
	return r
}

func (r *Router) routes() {
	// Product catalog
	r.HandleFunc("/products", catalogHandlers.GetProducts).Methods("GET")
	r.HandleFunc("/products/search", catalogHandlers.GetProductsByName).Methods("GET")

	// Cart
	r.HandleFunc("/cart/{userId}", cartHandlers.GetCartItems).Methods("GET")
	r.HandleFunc("/cart/{userId}/{productId}", cartHandlers.AddOrUpdateCartQty).Methods("PUT")
	r.HandleFunc("/cart/{userId}/{productId}", cartHandlers.RemoveCartItem).Methods("DELETE")
	r.HandleFunc("/cart/{userId}", cartHandlers.RemoveAllCartItems).Methods("DELETE")
	r.HandleFunc("/cart/checkout/{userId}", cartHandlers.Checkout).Methods("GET")

	// Payment
	r.HandleFunc("/payment", paymentHandlers.GetOrders).Methods("GET")
	r.HandleFunc("/payment/{userId}", paymentHandlers.GetOrders).Methods("GET")
	r.HandleFunc("/payment/{orderId}", paymentHandlers.MakePayment).Methods("PUT")
}