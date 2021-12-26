package routes

import (
	"github.com/gorilla/mux"
	carthandlers "github.com/vincentandr/shopping-microservice/src/services/bff/handlers/cart"
	cataloghandlers "github.com/vincentandr/shopping-microservice/src/services/bff/handlers/catalog"
	paymenthandlers "github.com/vincentandr/shopping-microservice/src/services/bff/handlers/payment"
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
	r.HandleFunc("/products", cataloghandlers.GetProducts).Methods("GET")
	r.HandleFunc("/products/search", cataloghandlers.GetProductsByName).Methods("GET")

	// Cart
	r.HandleFunc("/cart/{userId}", carthandlers.GetCartItems).Methods("GET")
	r.HandleFunc("/cart/{userId}/{productId}", carthandlers.AddOrUpdateCartQty).Methods("PUT")
	r.HandleFunc("/cart/{userId}/{productId}", carthandlers.RemoveCartItem).Methods("DELETE")
	r.HandleFunc("/cart/{userId}", carthandlers.RemoveAllCartItems).Methods("DELETE")
	r.HandleFunc("/cart/checkout/{userId}", carthandlers.Checkout).Methods("GET")

	// Payment
	r.HandleFunc("/payment/{orderId}", paymenthandlers.MakePayment).Methods("PUT")
}