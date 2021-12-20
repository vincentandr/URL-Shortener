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
	r.HandleFunc("/search", handlers.GetProductsWithName).Methods("GET")
	r.HandleFunc("/viewCart/{userId}", handlers.GetCartItems).Methods("GET")
	r.HandleFunc("/addOrUpdateCart", handlers.AddOrUpdateCart).Methods("POST")
	r.HandleFunc("/removeAllCart", handlers.RemoveAllCart).Methods("PUT")
	r.HandleFunc("/pay", handlers.Payment).Methods("POST")
}