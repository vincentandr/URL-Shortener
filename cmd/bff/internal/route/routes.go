package routes

import (
	"github.com/gorilla/mux"
	cartHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/cart"
	catalogHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/catalog"
	paymentHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/payment"
)

type Router struct {
	*mux.Router
	CartHandlers *cartHandlers.GrpcClient
	CatalogHandlers *catalogHandlers.GrpcClient
	PaymentHandlers *paymentHandlers.GrpcClient
}

// Router constructor
func NewRouter(handlerMap map[string]interface{}) *Router{
	r := &Router{
		Router: mux.NewRouter(),
		CartHandlers: handlerMap["cart"].(*cartHandlers.GrpcClient),
		CatalogHandlers: handlerMap["catalog"].(*catalogHandlers.GrpcClient),
		PaymentHandlers: handlerMap["payment"].(*paymentHandlers.GrpcClient),
	}
	r.routes()
	return r
}

func (r *Router) routes() {
	// Product catalog
	r.HandleFunc("/products", r.CatalogHandlers.GetProducts).Methods("GET")
	r.HandleFunc("/products/search", r.CatalogHandlers.GetProductsByName).Methods("GET")

	// Cart
	r.HandleFunc("/cart/{userId}", r.CartHandlers.GetCartItems).Methods("GET")
	r.HandleFunc("/cart/{userId}/{productId}", r.CartHandlers.AddOrUpdateCartQty).Methods("PUT")
	r.HandleFunc("/cart/{userId}/{productId}", r.CartHandlers.RemoveCartItem).Methods("DELETE")
	r.HandleFunc("/cart/{userId}", r.CartHandlers.RemoveAllCartItems).Methods("DELETE")
	r.HandleFunc("/cart/checkout/{userId}", r.CartHandlers.Checkout).Methods("GET")

	// Payment
	r.HandleFunc("/payment", r.PaymentHandlers.GetOrders).Methods("GET")
	r.HandleFunc("/payment/{userId}", r.PaymentHandlers.GetOrders).Methods("GET")
	r.HandleFunc("/payment/{orderId}", r.PaymentHandlers.MakePayment).Methods("PUT")
}