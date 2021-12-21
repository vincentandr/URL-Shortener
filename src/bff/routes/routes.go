package routes

import (
	"github.com/gorilla/mux"
	"github.com/vincentandr/shopping-microservice/src/bff/handlers"
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
	r.HandleFunc("/", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/search", handlers.GetProductsByName).Methods("GET")
	r.HandleFunc("/viewCart/{userId}", handlers.GetCartItems).Methods("GET")
	r.HandleFunc("/addOrUpdateCartQty/{userId}/{productId}", handlers.AddOrUpdateCartQty).Methods("PUT")
	r.HandleFunc("/removeCartItem/{userId}/{productId}", handlers.RemoveCartItem).Methods("DELETE")
	r.HandleFunc("/removeAllCartItems/{userId}", handlers.RemoveAllCartItems).Methods("DELETE")
	r.HandleFunc("/pay", handlers.Payment).Methods("POST")
}